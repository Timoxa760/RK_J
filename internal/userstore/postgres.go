package userstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (p *Postgres) Create(ctx context.Context, phone, passwordHash string) error {
	tag, err := p.pool.Exec(ctx, `
		INSERT INTO users (phone, role, password_hash) VALUES ($1, 'user', $2)
		ON CONFLICT (phone) DO NOTHING
	`, phone, passwordHash)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrExists
	}
	return nil
}

func (p *Postgres) GetByPhone(ctx context.Context, phone string) (*User, error) {
	var u User
	var hash *string
	err := p.pool.QueryRow(ctx, `
		SELECT id, phone, role, password_hash FROM users WHERE phone = $1
	`, phone).Scan(&u.ID, &u.Phone, &u.Role, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if hash != nil {
		u.PasswordHash = *hash
	}
	return &u, nil
}

func (p *Postgres) UpdatePassword(ctx context.Context, phone, passwordHash string) error {
	tag, err := p.pool.Exec(ctx, `
		UPDATE users SET password_hash = $2 WHERE phone = $1
	`, phone, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
