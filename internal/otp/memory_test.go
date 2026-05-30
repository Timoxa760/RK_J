package otp

import (
	"context"
	"testing"
	"time"
)

func TestGenerate_Length(t *testing.T) {
	code, err := Generate(6)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 6 {
		t.Errorf("expected 6 digits, got %q", code)
	}
}

func TestMemoryStore_SetGetDelete(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()
	if err := s.Set(ctx, "+79990000001", "123456", time.Minute); err != nil {
		t.Fatal(err)
	}
	got, err := s.Get(ctx, "+79990000001")
	if err != nil || got != "123456" {
		t.Fatalf("got=%q err=%v", got, err)
	}
	if err := s.Delete(ctx, "+79990000001"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(ctx, "+79990000001"); err == nil {
		t.Error("expected error after delete")
	}
}

func TestMemoryStore_RateLimit(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()
	window := time.Minute
	_ = s.IncrementRate(ctx, "+79990000002", window)
	n, err := s.RateCount(ctx, "+79990000002", window)
	if err != nil || n != 1 {
		t.Fatalf("count=%d err=%v", n, err)
	}
}
