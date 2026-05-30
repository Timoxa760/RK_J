package userstore

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("user not found")
var ErrExists = errors.New("user already exists")

// User — учётная запись.
type User struct {
	ID           int
	Phone        string
	Role         string
	PasswordHash string
}

// Store — persistent users.
type Store interface {
	Create(ctx context.Context, phone, passwordHash string) error
	GetByPhone(ctx context.Context, phone string) (*User, error)
	UpdatePassword(ctx context.Context, phone, passwordHash string) error
}
