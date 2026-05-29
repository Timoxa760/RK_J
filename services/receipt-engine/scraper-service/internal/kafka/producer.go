package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			BatchTimeout: 100 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
		topic: topic,
	}
}

func (p *Producer) Send(ctx context.Context, receipt scrap.RawReceipt) error {
	data, err := json.Marshal(receipt)
	if err != nil {
		return fmt.Errorf("kafka: marshal receipt: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(receipt.ID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "provider", Value: []byte(receipt.Provider)},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("kafka: write message: %w", err)
	}
	return nil
}

func (p *Producer) SendBatch(ctx context.Context, receipts []scrap.RawReceipt) error {
	messages := make([]kafka.Message, 0, len(receipts))
	for _, r := range receipts {
		data, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("kafka: marshal: %w", err)
		}
		messages = append(messages, kafka.Message{
			Key:   []byte(r.ID),
			Value: data,
			Headers: []kafka.Header{
				{Key: "provider", Value: []byte(r.Provider)},
			},
		})
	}

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("kafka: write batch: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
