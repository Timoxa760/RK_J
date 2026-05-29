package provider

import (
	"context"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type Provider interface {
	Name() string
	SendCode(ctx context.Context, phone string) error
	Login(ctx context.Context, creds map[string]string) error
	Sync(ctx context.Context) ([]scrap.RawReceipt, error)
}

type Type int

const (
	TypeAPI Type = iota
	TypeScraping
)

func Interval(t Type) int {
	if t == TypeAPI {
		return 6
	}
	return 24
}
