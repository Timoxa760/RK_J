package x5club

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	baseURL   = "https://api.x5club.ru/api/v2"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type HistoryItem struct {
	ID         string `json:"id"`
	StoreName  string `json:"store_name"`
	StoreAddr  string `json:"store_address"`
	Date       string `json:"date"`
	Total      float64 `json:"total"`
	Currency   string `json:"currency"`
	ItemsCount int    `json:"items_count"`
	Items      []struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
		Category string  `json:"category,omitempty"`
	} `json:"items,omitempty"`
}

type HistoryResponse struct {
	Success bool          `json:"success"`
	Data    struct {
		Receipts []HistoryItem `json:"receipts"`
		Page     int           `json:"page"`
		Pages    int           `json:"pages"`
		Total    int           `json:"total"`
	} `json:"data,omitempty"`
}

type Client struct {
	httpCli  *http.Client
	demoMode bool
}

func NewClient(demoMode bool) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		httpCli: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		demoMode: demoMode,
	}
}

func (c *Client) Login(phone, password string) error {
	if c.demoMode {
		return nil
	}

	body := LoginRequest{Phone: phone, Password: password}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", baseURL+"/auth/login", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("x5: create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("x5: login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("x5: login failed (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) GetHistory(page, limit int) ([]HistoryItem, int, error) {
	if c.demoMode {
		return c.demoHistory(page, limit)
	}

	u := fmt.Sprintf("%s/history?page=%d&limit=%d", baseURL, page, limit)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("x5: create history request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("x5: history request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("x5: history failed (%d): %s", resp.StatusCode, string(body))
	}

	var histResp HistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&histResp); err != nil {
		return nil, 0, fmt.Errorf("x5: decode history: %w", err)
	}

	return histResp.Data.Receipts, histResp.Data.Pages, nil
}

func (c *Client) demoHistory(page, limit int) ([]HistoryItem, int, error) {
	receipts := []HistoryItem{
		{
			ID: "x5-demo-001", StoreName: "Пятёрочка", Date: "2026-05-28",
			Total: 1032.50, ItemsCount: 5,
			Items: []struct {
				Name     string  `json:"name"`
				Price    float64 `json:"price"`
				Quantity int     `json:"quantity"`
				Category string  `json:"category,omitempty"`
			}{
				{Name: "Молоко 3.2%", Price: 78.00, Quantity: 2},
				{Name: "Хлеб белый", Price: 45.00, Quantity: 1},
				{Name: "Сыр Российский", Price: 189.00, Quantity: 1},
				{Name: "Масло сливочное", Price: 150.00, Quantity: 1},
				{Name: "Колбаса докторская", Price: 320.00, Quantity: 1},
			},
		},
		{
			ID: "x5-demo-002", StoreName: "Пятёрочка", Date: "2026-05-26",
			Total: 654.00, ItemsCount: 3,
			Items: []struct {
				Name     string  `json:"name"`
				Price    float64 `json:"price"`
				Quantity int     `json:"quantity"`
				Category string  `json:"category,omitempty"`
			}{
				{Name: "Яйца С1", Price: 120.00, Quantity: 1},
				{Name: "Кефир 3.2%", Price: 85.00, Quantity: 2},
				{Name: "Печенье", Price: 89.00, Quantity: 1},
			},
		},
		{
			ID: "x5-demo-003", StoreName: "Перекрёсток", Date: "2026-05-24",
			Total: 2340.00, ItemsCount: 7, Items: []struct {
				Name     string  `json:"name"`
				Price    float64 `json:"price"`
				Quantity int     `json:"quantity"`
				Category string  `json:"category,omitempty"`
			}{
				{Name: "Стейк говяжий", Price: 850.00, Quantity: 1},
				{Name: "Вино красное", Price: 680.00, Quantity: 1},
				{Name: "Салат Цезарь", Price: 320.00, Quantity: 1},
				{Name: "Сыр пармезан", Price: 490.00, Quantity: 1},
			},
		},
	}

	start := (page - 1) * limit
	if start > len(receipts) {
		return nil, 0, nil
	}
	end := start + limit
	if end > len(receipts) {
		end = len(receipts)
	}

	totalPages := (len(receipts) + limit - 1) / limit
	return receipts[start:end], totalPages, nil
}


