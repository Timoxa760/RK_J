package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"backend_project/services/receipt-engine/receipt-service/internal"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			BatchTimeout: 100 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) SendParsed(ctx context.Context, receipt internal.RawReceipt) error {
	data, err := json.Marshal(receipt)
	if err != nil {
		return fmt.Errorf("kafka: marshal: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(receipt.ID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event", Value: []byte("receipt.parsed")},
			{Key: "provider", Value: []byte(receipt.Provider)},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("kafka: write: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
