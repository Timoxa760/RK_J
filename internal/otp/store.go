package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"
)

// Store хранит одноразовые коды.
type Store interface {
	Set(ctx context.Context, phone, code string, ttl time.Duration) error
	Get(ctx context.Context, phone string) (string, error)
	Delete(ctx context.Context, phone string) error
	RateCount(ctx context.Context, phone string, window time.Duration) (int, error)
	IncrementRate(ctx context.Context, phone string, window time.Duration) error
}

func Generate(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	max := big.NewInt(1)
	for i := 0; i < length; i++ {
		max.Mul(max, big.NewInt(10))
	}
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n.Int64()), nil
}

func TTLFromEnv() time.Duration {
	if v := os.Getenv("OTP_TTL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return 5 * time.Minute
}

func LengthFromEnv() int {
	if v := os.Getenv("OTP_LENGTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 4 && n <= 8 {
			return n
		}
	}
	return 6
}

func RateLimitFromEnv() int {
	if v := os.Getenv("OTP_RATE_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 5
}

const rateWindow = 15 * time.Minute

func RateWindow() time.Duration {
	return rateWindow
}
