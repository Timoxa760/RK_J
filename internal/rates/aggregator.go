package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Benchmark — этalonная ставка для сравнения со scan.
type Benchmark struct {
	BenchmarkRate float64   `json:"benchmark_rate"`
	Source        string    `json:"source"`
	FetchedAt     time.Time `json:"fetched_at"`
}

type cacheEntry struct {
	value     Benchmark
	expiresAt time.Time
}

// Client — обёртка над внешним агрегатором (MVP: mock + env override).
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
	ttl     time.Duration

	mu    sync.RWMutex
	cache map[string]cacheEntry
}

func NewClient() *Client {
	return &Client{
		baseURL: os.Getenv("RATES_AGGREGATOR_URL"),
		apiKey:  os.Getenv("RATES_AGGREGATOR_API_KEY"),
		http:    &http.Client{Timeout: 8 * time.Second},
		ttl:     24 * time.Hour,
		cache:   make(map[string]cacheEntry),
	}
}

func (c *Client) cacheKey(product string, amount float64, termMonths int) string {
	return fmt.Sprintf("%s:%.0f:%d", product, amount, termMonths)
}

func (c *Client) getCached(key string) (Benchmark, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return Benchmark{}, false
	}
	return entry.value, true
}

func (c *Client) setCached(key string, b Benchmark) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{value: b, expiresAt: time.Now().Add(c.ttl)}
}

func (c *Client) Fetch(ctx context.Context, product string, amount float64, termMonths int) (Benchmark, error) {
	key := c.cacheKey(product, amount, termMonths)
	if b, ok := c.getCached(key); ok {
		return b, nil
	}

	var b Benchmark
	var err error
	if c.baseURL != "" {
		b, err = c.fetchRemote(ctx, product, amount, termMonths)
	} else {
		b = c.mockBenchmark(product, amount, termMonths)
	}
	if err != nil {
		return Benchmark{}, err
	}
	c.setCached(key, b)
	return b, nil
}

func (c *Client) fetchRemote(ctx context.Context, product string, amount float64, termMonths int) (Benchmark, error) {
	url := fmt.Sprintf("%s/rates?product=%s&amount=%.0f&term=%d", c.baseURL, product, amount, termMonths)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Benchmark{}, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return c.mockBenchmark(product, amount, termMonths), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, resp.Body)
		return c.mockBenchmark(product, amount, termMonths), nil
	}
	var b Benchmark
	if json.NewDecoder(resp.Body).Decode(&b) != nil || b.BenchmarkRate <= 0 {
		return c.mockBenchmark(product, amount, termMonths), nil
	}
	b.Source = "aggregator"
	if b.FetchedAt.IsZero() {
		b.FetchedAt = time.Now().UTC()
	}
	return b, nil
}

func (c *Client) mockBenchmark(product string, amount float64, termMonths int) Benchmark {
	base := 12.5
	if v := os.Getenv("RATES_MOCK_BENCHMARK"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			base = f
		}
	}
	if product == "mortgage" {
		base = 8.5
	}
	if termMonths > 60 {
		base += 1.2
	}
	if amount > 2_000_000 {
		base -= 0.8
	}
	return Benchmark{
		BenchmarkRate: base,
		Source:        "mock_cbr_adjacent",
		FetchedAt:     time.Now().UTC(),
	}
}
