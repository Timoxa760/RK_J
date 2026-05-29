package scheduler

import (
	"context"
	"testing"
	"time"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
	"backend_project/services/receipt-engine/scraper-service/internal/provider"
)

type mockSchedProvider struct {
	name    string
	synced  int
}

func (m *mockSchedProvider) Name() string { return m.name }
func (m *mockSchedProvider) Login(_ context.Context, _ map[string]string) error { return nil }
func (m *mockSchedProvider) Sync(_ context.Context) ([]scrap.RawReceipt, error) {
	m.synced++
	return []scrap.RawReceipt{{ID: "sched-test"}}, nil
}

func TestScheduler_Add(t *testing.T) {
	s := New()
	p := &mockSchedProvider{name: "test"}
	s.Add(p, provider.TypeAPI)

	if len(s.tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(s.tasks))
	}
}

func TestScheduler_Interval(t *testing.T) {
	if provider.Interval(provider.TypeAPI) != 6 {
		t.Errorf("expected 6h for API, got %d", provider.Interval(provider.TypeAPI))
	}
	if provider.Interval(provider.TypeScraping) != 24 {
		t.Errorf("expected 24h for scraping, got %d", provider.Interval(provider.TypeScraping))
	}
}

func TestScheduler_StartStop(t *testing.T) {
	s := New()
	p := &mockSchedProvider{name: "fast"}
	s.Add(p, provider.TypeAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.Start(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not stop after context cancel")
	}
}
