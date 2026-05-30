package otp

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type memEntry struct {
	code      string
	expiresAt time.Time
}

type rateEntry struct {
	count     int
	expiresAt time.Time
}

// MemoryStore — in-process OTP store для unit-тестов.
type MemoryStore struct {
	mu    sync.Mutex
	otps  map[string]memEntry
	rates map[string]rateEntry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		otps:  make(map[string]memEntry),
		rates: make(map[string]rateEntry),
	}
}

func (s *MemoryStore) Set(_ context.Context, phone, code string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.otps[phone] = memEntry{code: code, expiresAt: time.Now().Add(ttl)}
	return nil
}

func (s *MemoryStore) Get(_ context.Context, phone string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.otps[phone]
	if !ok || time.Now().After(e.expiresAt) {
		return "", fmt.Errorf("otp not found")
	}
	return e.code, nil
}

func (s *MemoryStore) Delete(_ context.Context, phone string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.otps, phone)
	return nil
}

func (s *MemoryStore) RateCount(_ context.Context, phone string, window time.Duration) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.rates[phone]
	if !ok || time.Now().After(e.expiresAt) {
		return 0, nil
	}
	return e.count, nil
}

func (s *MemoryStore) IncrementRate(_ context.Context, phone string, window time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.rates[phone]
	if time.Now().After(e.expiresAt) {
		e = rateEntry{expiresAt: time.Now().Add(window)}
	}
	e.count++
	s.rates[phone] = e
	return nil
}
