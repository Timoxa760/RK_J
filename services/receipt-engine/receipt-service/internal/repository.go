package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceiptRepo struct {
	pool *pgxpool.Pool
}

func NewReceiptRepo(pool *pgxpool.Pool) *ReceiptRepo {
	return &ReceiptRepo{pool: pool}
}

func (r *ReceiptRepo) Insert(ctx context.Context, receipt *RawReceipt) error {
	itemsJSON, err := json.Marshal(receipt.Items)
	if err != nil {
		return fmt.Errorf("repo: marshal items: %w", err)
	}

	query := `
		INSERT INTO receipts (id, user_id, provider, store_name, purchased_at, total_amount, items)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING`

	_, err = r.pool.Exec(ctx, query,
		receipt.ID,
		receipt.UserID,
		receipt.Provider,
		receipt.Store,
		receipt.Date,
		receipt.Total,
		itemsJSON,
	)
	if err != nil {
		return fmt.Errorf("repo: insert receipt: %w", err)
	}
	return nil
}

func (r *ReceiptRepo) GetByID(ctx context.Context, id string) (*RawReceipt, error) {
	query := `SELECT id, user_id, provider, store_name, purchased_at, total_amount, items FROM receipts WHERE id = $1`

	var receipt RawReceipt
	var itemsJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&receipt.ID,
		&receipt.UserID,
		&receipt.Provider,
		&receipt.Store,
		&receipt.Date,
		&receipt.Total,
		&itemsJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("repo: get receipt: %w", err)
	}

	if err := json.Unmarshal(itemsJSON, &receipt.Items); err != nil {
		return nil, fmt.Errorf("repo: unmarshal items: %w", err)
	}
	return &receipt, nil
}

func (r *ReceiptRepo) ListByUser(ctx context.Context, userID string, limit, offset int) ([]RawReceipt, error) {
	query := `SELECT id, user_id, provider, store_name, purchased_at, total_amount, items
		FROM receipts WHERE user_id = $1
		ORDER BY purchased_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repo: list receipts: %w", err)
	}
	defer rows.Close()

	var receipts []RawReceipt
	for rows.Next() {
		var r RawReceipt
		var itemsJSON []byte
		if err := rows.Scan(&r.ID, &r.UserID, &r.Provider, &r.Store, &r.Date, &r.Total, &itemsJSON); err != nil {
			return nil, fmt.Errorf("repo: scan receipt: %w", err)
		}
		if err := json.Unmarshal(itemsJSON, &r.Items); err != nil {
			return nil, fmt.Errorf("repo: unmarshal items: %w", err)
		}
		receipts = append(receipts, r)
	}
	return receipts, nil
}

func (r *ReceiptRepo) Exists(ctx context.Context, hash string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM receipt_dedup WHERE hash = $1 AND expires_at > NOW())`
	var exists bool
	if err := r.pool.QueryRow(ctx, query, hash).Scan(&exists); err != nil {
		return false, fmt.Errorf("repo: check dedup: %w", err)
	}
	return exists, nil
}

func (r *ReceiptRepo) Save(ctx context.Context, hash string, expiresAt time.Time) error {
	query := `INSERT INTO receipt_dedup (hash, expires_at) VALUES ($1, $2) ON CONFLICT (hash) DO NOTHING`
	_, err := r.pool.Exec(ctx, query, hash, expiresAt)
	if err != nil {
		return fmt.Errorf("repo: save dedup: %w", err)
	}
	return nil
}
