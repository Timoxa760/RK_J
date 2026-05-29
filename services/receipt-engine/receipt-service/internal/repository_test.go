package internal

import (
	"context"
	"testing"
	"time"
)

func TestReceiptRepo_InsertAndGetByID(t *testing.T) {
	pool := testPool(t)
	if pool == nil {
		t.Skip("no postgres available")
	}

	repo := NewReceiptRepo(pool)

	receipt := &RawReceipt{
		ID: "test-insert-001", UserID: "u1", Provider: "x5", Store: "Пятёрочка",
		Date: time.Now(), Total: 100.50,
		Items: []Item{{Name: "Хлеб", Price: 50.25, Quantity: 2}},
	}

	if err := repo.Insert(context.Background(), receipt); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	got, err := repo.GetByID(context.Background(), "test-insert-001")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}

	if got.Store != receipt.Store {
		t.Errorf("store: expected %s, got %s", receipt.Store, got.Store)
	}
	if len(got.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(got.Items))
	}
}

func TestReceiptRepo_ListByUser(t *testing.T) {
	pool := testPool(t)
	if pool == nil {
		t.Skip("no postgres available")
	}

	repo := NewReceiptRepo(pool)

	receipts, err := repo.ListByUser(context.Background(), "u1", 10, 0)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}
	if len(receipts) == 0 {
		t.Log("no receipts found for u1")
	}
}

func TestReceiptRepo_Dedup(t *testing.T) {
	pool := testPool(t)
	if pool == nil {
		t.Skip("no postgres available")
	}

	repo := NewReceiptRepo(pool)

	exists, err := repo.Exists(context.Background(), "test-hash-001")
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}

	if err := repo.Save(context.Background(), "test-hash-001", time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("Save: %v", err)
	}

	exists, err = repo.Exists(context.Background(), "test-hash-001")
	if err != nil {
		t.Fatalf("Exists after save: %v", err)
	}
	if !exists {
		t.Error("expected hash to exist after save")
	}
}
