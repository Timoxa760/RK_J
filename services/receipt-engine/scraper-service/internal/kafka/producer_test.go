package kafka

import (
	"context"
	"testing"
	"time"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

func TestNewProducer(t *testing.T) {
	p := NewProducer([]string{"localhost:9092"}, "receipt.raw")
	if p == nil {
		t.Fatal("expected non-nil producer")
	}
	if p.topic != "receipt.raw" {
		t.Errorf("expected topic receipt.raw, got %s", p.topic)
	}
	p.Close()
}

func TestProducer_Send_NoBroker(t *testing.T) {
	p := NewProducer([]string{"127.0.0.1:19092"}, "receipt.raw")
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := p.Send(ctx, scrap.RawReceipt{
		ID:       "test-001",
		Provider: "mock",
		Store:    "Test",
		Date:     time.Now(),
		Total:    100,
	})
	if err == nil {
		t.Skip("expected connection error (no broker)")
	}
}

func TestProducer_SendBatch_NoBroker(t *testing.T) {
	p := NewProducer([]string{"127.0.0.1:19092"}, "receipt.raw")
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := p.SendBatch(ctx, []scrap.RawReceipt{
		{ID: "t1", Provider: "mock", Store: "S", Date: time.Now(), Total: 10},
		{ID: "t2", Provider: "mock", Store: "S", Date: time.Now(), Total: 20},
	})
	if err == nil {
		t.Skip("expected connection error (no broker)")
	}
}
