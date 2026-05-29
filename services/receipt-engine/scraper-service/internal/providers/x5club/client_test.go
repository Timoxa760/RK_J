package x5club

import (
	"testing"
)

func TestSendCode_DemoMode(t *testing.T) {
	c := NewClient(true)
	if err := c.SendCode("+79991234567"); err != nil {
		t.Fatalf("SendCode in demo mode: %v", err)
	}
}

func TestLogin_DemoMode(t *testing.T) {
	c := NewClient(true)
	if err := c.Login("+79991234567", "0000"); err != nil {
		t.Fatalf("Login in demo mode: %v", err)
	}
}

func TestGetHistory_DemoMode(t *testing.T) {
	c := NewClient(true)
	receipts, pages, err := c.GetHistory(1, 10)
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if len(receipts) == 0 {
		t.Fatal("expected at least 1 receipt")
	}
	if pages < 1 {
		t.Errorf("expected at least 1 page, got %d", pages)
	}
}

func TestGetHistory_Pagination(t *testing.T) {
	c := NewClient(true)

	receipts, pages, err := c.GetHistory(1, 2)
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if len(receipts) != 2 {
		t.Errorf("expected 2 receipts on page 1, got %d", len(receipts))
	}

	receipts, _, err = c.GetHistory(2, 2)
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if len(receipts) == 0 {
		t.Error("expected receipts on page 2")
	}

	if pages != 2 {
		t.Errorf("expected 2 total pages with limit=2, got %d", pages)
	}
}

func TestGetHistory_OutOfRange(t *testing.T) {
	c := NewClient(true)
	receipts, _, err := c.GetHistory(999, 10)
	if err != nil {
		t.Fatalf("GetHistory out of range: %v", err)
	}
	if len(receipts) != 0 {
		t.Errorf("expected 0 receipts for out-of-range page, got %d", len(receipts))
	}
}

func TestToRawReceipts(t *testing.T) {
	m := NewMapper()
	c := NewClient(true)
	items, _, _ := c.GetHistory(1, 3)
	receipts := m.ToRawReceipts(items)

	if len(receipts) != len(items) {
		t.Fatalf("expected %d receipts, got %d", len(items), len(receipts))
	}

	for i, r := range receipts {
		if r.Provider != "x5club" {
			t.Errorf("receipt %d: expected provider x5club, got %s", i, r.Provider)
		}
		if r.Store == "" {
			t.Errorf("receipt %d: empty store name", i)
		}
		if len(r.Items) == 0 {
			t.Errorf("receipt %d: no items", i)
		}
		if r.Total <= 0 {
			t.Errorf("receipt %d: total <= 0", i)
		}
	}
}
