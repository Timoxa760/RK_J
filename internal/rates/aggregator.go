package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Benchmark — этalonная ставка для сравнения со scan.
type Benchmark struct {
	BenchmarkRate float64   `json:"benchmark_rate"`
	Source        string    `json:"source"`
	FetchedAt     time.Time `json:"fetched_at"`
}

// Client — обёртка над внешним агрегатором (MVP: mock + env override).
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: os.Getenv("RATES_AGGREGATOR_URL"),
		apiKey:  os.Getenv("RATES_AGGREGATOR_API_KEY"),
		http:    &http.Client{Timeout: 8 * time.Second},
	}
}

func (c *Client) Fetch(ctx context.Context, product string, amount float64, termMonths int) (Benchmark, error) {
	if c.baseURL != "" {
		url := fmt.Sprintf("%s/rates?product=%s&amount=%.0f&term=%d", c.baseURL, product, amount, termMonths)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return Benchmark{}, err
		}
		if c.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+c.apiKey)
		}
		resp, err := c.http.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var b Benchmark
			if json.NewDecoder(resp.Body).Decode(&b) == nil && b.BenchmarkRate > 0 {
				b.Source = "aggregator"
				if b.FetchedAt.IsZero() {
					b.FetchedAt = time.Now().UTC()
				}
				return b, nil
			}
			io.Copy(io.Discard, resp.Body)
		}
	}
	return c.mockBenchmark(product, amount, termMonths), nil
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
