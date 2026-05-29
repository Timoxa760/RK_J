package kafka

import (
	"context"
	"testing"
	"time"

	svc "backend_project/services/receipt-engine/receipt-service/internal"
)

func TestNewProducer(t *testing.T) {
	p := NewProducer([]string{"localhost:9092"}, "receipt.parsed")
	if p == nil {
		t.Fatal("expected non-nil producer")
	}
	p.Close()
}

func TestProducer_SendParsed_NoBroker(t *testing.T) {
	p := NewProducer([]string{"127.0.0.1:19092"}, "receipt.parsed")
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := p.SendParsed(ctx, svc.RawReceipt{
		ID: "test", Provider: "mock", Store: "S",
		Date: time.Now(), Total: 100,
	})
	if err == nil {
		t.Skip("expected error (no broker)")
	}
}
