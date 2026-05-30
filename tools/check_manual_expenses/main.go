package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://postgres:postgres@127.0.0.1:5432/moneymind?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		fmt.Println("connect:", err)
		os.Exit(1)
	}
	defer pool.Close()

	var exists bool
	err = pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'manual_expenses'
		)`).Scan(&exists)
	if err != nil {
		fmt.Println("query:", err)
		os.Exit(1)
	}
	fmt.Println("manual_expenses exists:", exists)
}
