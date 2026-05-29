package x5club

import (
	"context"
	"testing"
)

func TestPool_FetchAll(t *testing.T) {
	client := NewClient(true)
	pool := NewPool(client, 3)

	items, err := pool.FetchAll(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}

	if len(items) == 0 {
		t.Fatal("expected items")
	}
}

func TestPool_SinglePage(t *testing.T) {
	client := NewClient(true)
	pool := NewPool(client, 3)

	items, err := pool.FetchAll(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("expected 3 items for large limit, got %d", len(items))
	}
}

func TestPool_ContextCancel(t *testing.T) {
	client := NewClient(true)
	pool := NewPool(client, 2)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := pool.FetchAll(ctx, 1, 1)
	if err != nil {
		t.Fatalf("FetchAll with cancelled ctx: %v", err)
	}
}

func TestPool_ZeroWorkers(t *testing.T) {
	pool := NewPool(NewClient(true), 0)
	if pool.workers != 1 {
		t.Errorf("expected 1 worker for zero input, got %d", pool.workers)
	}
}
