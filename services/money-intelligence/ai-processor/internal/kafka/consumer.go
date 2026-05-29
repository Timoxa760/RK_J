package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"backend_project/services/money-intelligence/ai-processor/internal"
)

type Handler func(ctx context.Context, receipt internal.RawReceipt) error

type Consumer struct {
	reader  *kafka.Reader
	handler Handler
}

func NewConsumer(brokers []string, topic, groupID string, handler Handler) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       1,
			MaxBytes:       10 * 1024 * 1024,
			CommitInterval: time.Second,
			StartOffset:    kafka.LastOffset,
		}),
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("kafka consumer: starting on %s (group=%s)", c.reader.Config().Topic, c.reader.Config().GroupID)
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("kafka consumer: read: %w", err)
		}
		var receipt internal.RawReceipt
		if err := json.Unmarshal(msg.Value, &receipt); err != nil {
			log.Printf("kafka consumer: skip invalid message (key=%s): %v", string(msg.Key), err)
			continue
		}
		if err := c.handler(ctx, receipt); err != nil {
			log.Printf("kafka consumer: handler error (key=%s): %v", string(msg.Key), err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
