package ratelimit

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type Limiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

func New(r rate.Limit, burst int) *Limiter {
	return &Limiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

func (l *Limiter) Get(provider string) *rate.Limiter {
	l.mu.RLock()
	lim, ok := l.limiters[provider]
	l.mu.RUnlock()
	if ok {
		return lim
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if lim, ok = l.limiters[provider]; ok {
		return lim
	}

	lim = rate.NewLimiter(l.rate, l.burst)
	l.limiters[provider] = lim
	return lim
}

func (l *Limiter) Allow(provider string) bool {
	return l.Get(provider).Allow()
}

func (l *Limiter) Wait(ctx context.Context, provider string) error {
	return l.Get(provider).Wait(ctx)
}
