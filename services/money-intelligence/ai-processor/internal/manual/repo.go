package manual

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pg *pgxpool.Pool
	ch driver.Conn
}

func NewRepo(pg *pgxpool.Pool, ch driver.Conn) *Repo {
	return &Repo{pg: pg, ch: ch}
}

type Expense struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	RawText     string    `json:"raw_text"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"expense_date"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *Repo) Insert(ctx context.Context, e *Expense) error {
	query := `
		INSERT INTO manual_expenses (user_id, raw_text, amount, category, description, expense_date, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`
	err := r.pg.QueryRow(ctx, query,
		e.UserID, e.RawText, e.Amount, e.Category, e.Description, e.Date, e.Source,
	).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return fmt.Errorf("pg insert: %w", err)
	}

	if r.ch != nil {
		batch, err := r.ch.PrepareBatch(ctx, `
			INSERT INTO receipt_items (user_id, store, category, item_name, price, quantity, purchased_at, is_impulsive)
		`)
		if err != nil {
			return fmt.Errorf("ch prepare: %w", err)
		}
		if err := batch.Append(
			e.UserID, "manual_input", e.Category, e.Description,
			e.Amount, uint32(1), e.Date, uint8(0),
		); err != nil {
			return fmt.Errorf("ch append: %w", err)
		}
		if err := batch.Send(); err != nil {
			return fmt.Errorf("ch send: %w", err)
		}
	}

	return nil
}
