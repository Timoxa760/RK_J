package x5club

import (
	"testing"
)

func TestMapper_ToRawReceipts(t *testing.T) {
	m := NewMapper()
	items := []HistoryItem{
		{
			ID: "test-1", StoreName: "Пятёрочка", Date: "2026-05-28",
			Total: 100.50,
			Items: []struct {
				Name     string  `json:"name"`
				Price    float64 `json:"price"`
				Quantity int     `json:"quantity"`
				Category string  `json:"category,omitempty"`
			}{
				{Name: "Хлеб", Price: 45.00, Quantity: 1},
			},
		},
	}

	receipts := m.ToRawReceipts(items)

	if len(receipts) != 1 {
		t.Fatalf("expected 1 receipt, got %d", len(receipts))
	}

	r := receipts[0]
	if r.ID != "test-1" {
		t.Errorf("ID: expected test-1, got %s", r.ID)
	}
	if r.Provider != "x5club" {
		t.Errorf("Provider: expected x5club, got %s", r.Provider)
	}
	if r.Store != "Пятёрочка" {
		t.Errorf("Store: expected Пятёрочка, got %s", r.Store)
	}
	if r.Total != 100.50 {
		t.Errorf("Total: expected 100.50, got %.2f", r.Total)
	}
	if len(r.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(r.Items))
	}
	if r.Items[0].Name != "Хлеб" {
		t.Errorf("Item name: expected Хлеб, got %s", r.Items[0].Name)
	}
}

func TestMapper_EmptyList(t *testing.T) {
	m := NewMapper()
	receipts := m.ToRawReceipts(nil)
	if len(receipts) != 0 {
		t.Errorf("expected 0 receipts for nil input, got %d", len(receipts))
	}

	receipts = m.ToRawReceipts([]HistoryItem{})
	if len(receipts) != 0 {
		t.Errorf("expected 0 receipts for empty input, got %d", len(receipts))
	}
}

func TestMapper_DateParsing(t *testing.T) {
	m := NewMapper()
	items := []HistoryItem{
		{ID: "t1", StoreName: "S", Date: "invalid-date", Total: 1},
	}

	receipts := m.ToRawReceipts(items)
	if len(receipts) != 1 {
		t.Fatalf("expected 1 receipt, got %d", len(receipts))
	}

	if !receipts[0].Date.IsZero() {
		t.Errorf("expected zero time for invalid date, got %v", receipts[0].Date)
	}
}
