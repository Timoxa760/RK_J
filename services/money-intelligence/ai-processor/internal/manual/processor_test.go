package manual

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"backend_project/internal/expensestore"
	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

func TestProcessor_Create_MultipleExpenses_FileStore(t *testing.T) {
	dir := t.TempDir()
	fileStore, err := expensestore.NewFileStore(filepath.Join(dir, "expenses.json"))
	if err != nil {
		t.Fatal(err)
	}
	store := NewFallbackStorage(nil, fileStore)
	proc := NewProcessor(expense.NewParser(nil), store)

	resp, code, err := proc.Create(context.Background(), CreateRequest{
		UserID:  "+79991234567",
		RawText: "кофе 300 рублей и такси 500 рублей",
		Source:  "voice",
	})
	if err != nil {
		t.Fatalf("code=%d err=%v", code, err)
	}
	if len(resp.Expenses) != 2 {
		t.Fatalf("expected 2 expenses, got %+v", resp.Expenses)
	}

	recs, err := fileStore.ListSince("+79991234567", time.Now().AddDate(0, 0, -1))
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 2 {
		t.Fatalf("file store len=%d want 2: %+v", len(recs), recs)
	}
}
