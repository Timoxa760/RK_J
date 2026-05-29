package magnit

import (
	"testing"
)

func TestMapper_ToRawReceipts(t *testing.T) {
	m := NewMapper(NewINNDetector())
	receipts := []Receipt{
		{
			ID: "mg-001", StoreName: "Магнит", Date: "2026-05-28",
			Total: 100.50,
			Items: []ReceiptItem{
				{Name: "Хлеб", Price: 45.00, Quantity: 1},
			},
		},
	}

	raw := m.ToRawReceipts(receipts)
	if len(raw) != 1 {
		t.Fatalf("expected 1, got %d", len(raw))
	}

	r := raw[0]
	if r.Provider != "magnit" {
		t.Errorf("provider: expected magnit, got %s", r.Provider)
	}
	if r.Store != "Магнит" {
		t.Errorf("store: expected Магнит, got %s", r.Store)
	}
	if r.Total != 100.50 {
		t.Errorf("total: expected 100.50, got %.2f", r.Total)
	}
	if len(r.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(r.Items))
	}
}

func TestMapper_ToRawReceiptsWithINN(t *testing.T) {
	m := NewMapper(NewINNDetector())
	receipts := []Receipt{
		{ID: "mg-001", StoreName: "Магнит", Date: "2026-05-28", Total: 50},
	}

	raw := m.ToRawReceiptsWithINN(receipts)
	if len(raw) != 1 {
		t.Fatalf("expected 1, got %d", len(raw))
	}

	if raw[0].ID != "2309085638_mg-001" {
		t.Errorf("ID: expected 2309085638_mg-001, got %s", raw[0].ID)
	}
}

func TestMapper_EmptyInput(t *testing.T) {
	m := NewMapper(NewINNDetector())
	raw := m.ToRawReceipts(nil)
	if len(raw) != 0 {
		t.Errorf("expected 0 for nil, got %d", len(raw))
	}
	raw = m.ToRawReceipts([]Receipt{})
	if len(raw) != 0 {
		t.Errorf("expected 0 for empty, got %d", len(raw))
	}
}
