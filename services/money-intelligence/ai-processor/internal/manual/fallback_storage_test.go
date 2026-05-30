package manual

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"backend_project/internal/expensestore"
	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

func TestFallbackStorage_FileOnly(t *testing.T) {
	dir := t.TempDir()
	fileStore, err := expensestore.NewFileStore(filepath.Join(dir, "expenses.json"))
	if err != nil {
		t.Fatal(err)
	}
	store := NewFallbackStorage(nil, fileStore)
	date := time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC)

	out, err := store.SaveExpenses(context.Background(), "+7999", "manual", date, []expense.Item{{
		Amount: 1500, Category: "Продукты", Description: "Пятёрочка",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Amount != 1500 {
		t.Fatalf("out: %+v", out)
	}
}
