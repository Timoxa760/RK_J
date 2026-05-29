package internal

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func testPool(t interface{ Skip(...interface{}) }) *pgxpool.Pool {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		return nil
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil
	}
	return pool
}
