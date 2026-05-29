package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type DedupRepo interface {
	Exists(ctx context.Context, hash string) (bool, error)
	Save(ctx context.Context, hash string, expiresAt time.Time) error
}

type Dedup struct {
	cache   sync.Map
	repo    DedupRepo
	ttl     time.Duration
	done    chan struct{}
	cleanup *time.Ticker
}

func NewDedup(repo DedupRepo, ttl time.Duration) *Dedup {
	d := &Dedup{
		repo:    repo,
		ttl:     ttl,
		done:    make(chan struct{}),
		cleanup: time.NewTicker(5 * time.Minute),
	}
	go d.evictLoop()
	return d
}

func (d *Dedup) Hash(receiptID, provider, userID string) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", receiptID, provider, userID)))
	return hex.EncodeToString(h[:])
}

func (d *Dedup) IsDuplicate(ctx context.Context, receiptID, provider, userID string) (bool, error) {
	hash := d.Hash(receiptID, provider, userID)

	if _, ok := d.cache.Load(hash); ok {
		return true, nil
	}

	exists, err := d.repo.Exists(ctx, hash)
	if err != nil {
		return false, err
	}
	if exists {
		d.cache.Store(hash, time.Now())
		return true, nil
	}

	d.cache.Store(hash, time.Now())
	if err := d.repo.Save(ctx, hash, time.Now().Add(d.ttl)); err != nil {
		return false, err
	}

	return false, nil
}

func (d *Dedup) evictLoop() {
	for range d.cleanup.C {
		d.cache.Range(func(key, value interface{}) bool {
			ts, ok := value.(time.Time)
			if ok && time.Since(ts) > d.ttl {
				d.cache.Delete(key)
			}
			return true
		})
	}
}

func (d *Dedup) Stop() {
	d.cleanup.Stop()
	close(d.done)
}
