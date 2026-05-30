package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema_manual_expenses.sql
var manualExpensesDDL string

// EnsureManualExpenses создаёт таблицу manual_expenses, если PostgreSQL доступен.
func EnsureManualExpenses(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("postgres pool is nil")
	}
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("postgres ping: %w", err)
	}
	if _, err := pool.Exec(ctx, manualExpensesDDL); err != nil {
		return fmt.Errorf("ensure manual_expenses: %w", err)
	}
	return nil
}

// Ping проверяет доступность PostgreSQL.
func Ping(ctx context.Context, pool *pgxpool.Pool) bool {
	if pool == nil {
		return false
	}
	return pool.Ping(ctx) == nil
}
