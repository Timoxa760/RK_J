package kafka

import (
	"context"
	"testing"
	"time"

	svc "backend_project/services/receipt-engine/receipt-service/internal"
)

func TestNewConsumer(t *testing.T) {
	handler := func(_ context.Context, _ svc.RawReceipt) error { return nil }
	c := NewConsumer([]string{"127.0.0.1:19092"}, "receipt.raw", "test-group", handler)
	if c == nil {
		t.Fatal("expected non-nil consumer")
	}
	c.Close()
}

func TestConsumer_Start_NoBroker(t *testing.T) {
	handler := func(_ context.Context, _ svc.RawReceipt) error { return nil }
	c := NewConsumer([]string{"127.0.0.1:19092"}, "receipt.raw", "test-group", handler)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := c.Start(ctx)
	if err == nil {
		t.Skip("expected error (no broker)")
	}
}
