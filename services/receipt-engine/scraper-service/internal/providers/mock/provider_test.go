package mock

import (
	"context"
	"testing"
)

func TestProvider_Name(t *testing.T) {
	p := New("test-provider", "Test Store")
	if p.Name() != "test-provider" {
		t.Errorf("expected test-provider, got %s", p.Name())
	}
}

func TestProvider_Login(t *testing.T) {
	p := New("mock", "Store")
	if err := p.Login(context.Background(), nil); err != nil {
		t.Fatalf("Login: %v", err)
	}
}

func TestProvider_Sync(t *testing.T) {
	p := New("mock", "Store")
	receipts, err := p.Sync(context.Background())
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if len(receipts) != 2 {
		t.Errorf("expected 2 receipts, got %d", len(receipts))
	}
}

func TestProvider_ReceiptFields(t *testing.T) {
	p := New("mock", "MyStore")
	receipts, _ := p.Sync(context.Background())

	for i, r := range receipts {
		if r.Provider != "mock" {
			t.Errorf("receipt %d: provider mismatch", i)
		}
		if r.Store != "MyStore" {
			t.Errorf("receipt %d: store mismatch", i)
		}
		if r.Total <= 0 {
			t.Errorf("receipt %d: total <= 0", i)
		}
	}
}
