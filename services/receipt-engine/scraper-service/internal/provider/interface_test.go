package provider

import (
	"context"
	"testing"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type testProvider struct{}

func (p *testProvider) Name() string { return "test" }
func (p *testProvider) Login(_ context.Context, _ map[string]string) error { return nil }
func (p *testProvider) Sync(_ context.Context) ([]scrap.RawReceipt, error) { return nil, nil }

func TestProviderInterface(t *testing.T) {
	var p Provider = &testProvider{}
	if p.Name() != "test" {
		t.Errorf("expected test, got %s", p.Name())
	}
}

func TestInterval(t *testing.T) {
	if Interval(TypeAPI) != 6 {
		t.Errorf("expected 6h for API, got %d", Interval(TypeAPI))
	}
	if Interval(TypeScraping) != 24 {
		t.Errorf("expected 24h for scraping, got %d", Interval(TypeScraping))
	}
}
