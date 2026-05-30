package expensestore

import (
	"os"
	"testing"
	"time"
)

func TestFileStore_AppendAndList(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/expenses.json"
	store, err := NewFileStore(path)
	if err != nil {
		t.Fatal(err)
	}

	date := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	saved, err := store.Append("+7999", "manual", date, []Item{{
		Amount: 1500, Category: "Продукты", Description: "Пятёрочка",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(saved) != 1 || saved[0].Amount != 1500 {
		t.Fatalf("saved: %+v", saved)
	}

	list, err := store.ListSince("+7999", date)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %+v err=%v", list, err)
	}
}

func TestDefaultPath_Env(t *testing.T) {
	t.Setenv("EXPENSE_STORE_PATH", "/tmp/custom.json")
	if got := DefaultPath(); got != "/tmp/custom.json" {
		t.Fatalf("got %q", got)
	}
	_ = os.Unsetenv("EXPENSE_STORE_PATH")
}
