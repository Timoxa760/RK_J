package internal

import (
	"context"
	"testing"
	"time"
)

type memDedupRepo struct {
	hashes map[string]time.Time
}

func (m *memDedupRepo) Exists(_ context.Context, hash string) (bool, error) {
	_, ok := m.hashes[hash]
	return ok, nil
}
func (m *memDedupRepo) Save(_ context.Context, hash string, expiresAt time.Time) error {
	if m.hashes == nil {
		m.hashes = make(map[string]time.Time)
	}
	m.hashes[hash] = expiresAt
	return nil
}

func TestDedup_Hash(t *testing.T) {
	d := NewDedup(&memDedupRepo{}, time.Hour)
	h1 := d.Hash("r1", "x5", "user1")
	h2 := d.Hash("r1", "x5", "user1")
	if h1 != h2 {
		t.Error("same input should produce same hash")
	}
}

func TestDedup_FirstNotDuplicate(t *testing.T) {
	d := NewDedup(&memDedupRepo{}, time.Hour)
	dup, err := d.IsDuplicate(context.Background(), "r1", "x5", "user1")
	if err != nil {
		t.Fatalf("IsDuplicate: %v", err)
	}
	if dup {
		t.Error("first call should not be duplicate")
	}
}

func TestDedup_SecondIsDuplicate(t *testing.T) {
	d := NewDedup(&memDedupRepo{}, time.Hour)
	d.IsDuplicate(context.Background(), "r1", "x5", "user1")
	dup, _ := d.IsDuplicate(context.Background(), "r1", "x5", "user1")
	if !dup {
		t.Error("second call should be duplicate")
	}
}

func TestDedup_DifferentUsers(t *testing.T) {
	d := NewDedup(&memDedupRepo{}, time.Hour)
	d1, _ := d.IsDuplicate(context.Background(), "r1", "x5", "user1")
	d2, _ := d.IsDuplicate(context.Background(), "r1", "x5", "user2")
	if d1 {
		t.Error("user1 should not be duplicate")
	}
	if d2 {
		t.Error("user2 should not be duplicate")
	}
}
