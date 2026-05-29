package mock

import (
	"context"
	"time"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type Provider struct {
	name  string
	store string
}

func New(name, store string) *Provider {
	return &Provider{name: name, store: store}
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) SendCode(_ context.Context, _ string) error {
	return nil
}

func (p *Provider) Login(_ context.Context, _ map[string]string) error {
	return nil
}

func (p *Provider) Sync(_ context.Context) ([]scrap.RawReceipt, error) {
	return []scrap.RawReceipt{
		{
			ID:       "mock-001",
			UserID:   "mock-user",
			Provider: p.name,
			Store:    p.store,
			Date:     time.Now().Add(-2 * time.Hour),
			Total:    567.80,
			Items: []scrap.RawItem{
				{Name: "Mock Item 1", Price: 200.00, Quantity: 1},
				{Name: "Mock Item 2", Price: 367.80, Quantity: 2},
			},
		},
		{
			ID:       "mock-002",
			UserID:   "mock-user",
			Provider: p.name,
			Store:    p.store,
			Date:     time.Now().Add(-26 * time.Hour),
			Total:    1200.00,
			Items: []scrap.RawItem{
				{Name: "Mock Item 3", Price: 600.00, Quantity: 2},
			},
		},
	}, nil
}
