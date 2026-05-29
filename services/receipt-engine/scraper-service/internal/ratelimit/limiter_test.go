package ratelimit

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestLimiter_Allow(t *testing.T) {
	l := New(rate.Limit(100), 10)
	if !l.Allow("test") {
		t.Error("expected allow on first call")
	}
}

func TestLimiter_PerProvider(t *testing.T) {
	l := New(rate.Limit(1), 1)

	if !l.Allow("a") {
		t.Error("expected allow for a")
	}
	if !l.Allow("b") {
		t.Error("expected allow for b")
	}
	if l.Allow("a") {
		t.Error("expected deny for a (rate exceeded)")
	}
	if l.Allow("b") {
		t.Error("expected deny for b (rate exceeded)")
	}
}

func TestLimiter_Wait(t *testing.T) {
	l := New(rate.Limit(1000), 5)
	ctx := context.Background()
	if err := l.Wait(ctx, "test"); err != nil {
		t.Fatalf("Wait: %v", err)
	}
}

func TestLimiter_ConcurrentProviders(t *testing.T) {
	l := New(rate.Limit(10), 3)

	ctx := context.Background()
	if err := l.Wait(ctx, "p1"); err != nil {
		t.Fatal(err)
	}
	if err := l.Wait(ctx, "p2"); err != nil {
		t.Fatal(err)
	}
	if l.Get("p1") == l.Get("p2") {
		t.Error("expected different limiters for different providers")
	}
}

func TestLimiter_Timeout(t *testing.T) {
	l := New(rate.Limit(1), 1)
	l.Allow("x")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx, "x")
	if err == nil {
		t.Error("expected timeout error")
	}
}
