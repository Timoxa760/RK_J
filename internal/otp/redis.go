package otp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func otpKey(phone string) string {
	return "otp:" + phone
}

func rateKey(phone string) string {
	return "otp:rate:" + phone
}

func (s *RedisStore) Set(ctx context.Context, phone, code string, ttl time.Duration) error {
	return s.client.Set(ctx, otpKey(phone), code, ttl).Err()
}

func (s *RedisStore) Get(ctx context.Context, phone string) (string, error) {
	val, err := s.client.Get(ctx, otpKey(phone)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("otp not found")
	}
	return val, err
}

func (s *RedisStore) Delete(ctx context.Context, phone string) error {
	return s.client.Del(ctx, otpKey(phone)).Err()
}

func (s *RedisStore) RateCount(ctx context.Context, phone string, window time.Duration) (int, error) {
	n, err := s.client.Get(ctx, rateKey(phone)).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return n, err
}

func (s *RedisStore) IncrementRate(ctx context.Context, phone string, window time.Duration) error {
	key := rateKey(phone)
	pipe := s.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	return err
}
