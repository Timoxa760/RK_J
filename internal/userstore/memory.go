package userstore

import (
	"context"
	"sync"
)

type Memory struct {
	mu    sync.Mutex
	users map[string]*User
	next  int
}

func NewMemory() *Memory {
	return &Memory{users: make(map[string]*User), next: 1}
}

func (m *Memory) Create(_ context.Context, phone, passwordHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[phone]; ok {
		return ErrExists
	}
	m.users[phone] = &User{
		ID:           m.next,
		Phone:        phone,
		Role:         "user",
		PasswordHash: passwordHash,
	}
	m.next++
	return nil
}

func (m *Memory) GetByPhone(_ context.Context, phone string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[phone]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *u
	return &copy, nil
}

func (m *Memory) UpdatePassword(_ context.Context, phone, passwordHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.users[phone]
	if !ok {
		return ErrNotFound
	}
	u.PasswordHash = passwordHash
	return nil
}
