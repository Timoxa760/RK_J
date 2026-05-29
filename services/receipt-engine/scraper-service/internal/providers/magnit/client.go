package magnit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL   = "https://api.magnit.ru/api/v2"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

type SendCodeRequest struct {
	Phone string `json:"phone"`
}

type VerifyCodeRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

type ReceiptItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category,omitempty"`
}

type Receipt struct {
	ID        string        `json:"id"`
	StoreName string        `json:"store_name"`
	StoreAddr string        `json:"store_address"`
	Date      string        `json:"date"`
	Total     float64       `json:"total"`
	Currency  string        `json:"currency"`
	ItemsCnt  int           `json:"items_count"`
	Items     []ReceiptItem `json:"items,omitempty"`
}

type ReceiptsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Receipts []Receipt `json:"receipts"`
		Page     int       `json:"page"`
		Pages    int       `json:"pages"`
		Total    int       `json:"total"`
	} `json:"data,omitempty"`
}

type Client struct {
	httpCli  *http.Client
	token    string
	demoMode bool
}

func NewClient(demoMode bool) *Client {
	return &Client{
		httpCli: &http.Client{
			Timeout: 30 * time.Second,
		},
		demoMode: demoMode,
	}
}

func (c *Client) SendCode(phone string) error {
	if c.demoMode {
		return nil
	}

	body := SendCodeRequest{Phone: phone}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", baseURL+"/auth/send-code", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("magnit: create send-code request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("magnit: send-code request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("magnit: send-code failed (%d): %s", resp.StatusCode, string(b))
	}

	return nil
}

func (c *Client) Login(phone, code string) error {
	if c.demoMode {
		c.token = "demo-token"
		return nil
	}

	body := VerifyCodeRequest{Phone: phone, Code: code}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", baseURL+"/auth/verify-code", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("magnit: create verify request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("magnit: verify request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("magnit: verify failed (%d): %s", resp.StatusCode, string(b))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("magnit: decode verify response: %w", err)
	}

	c.token = loginResp.Token
	return nil
}

func (c *Client) GetReceipts(page int) ([]Receipt, int, error) {
	if c.demoMode {
		return c.demoReceipts(page)
	}

	u := fmt.Sprintf("%s/receipts?page=%d", baseURL, page)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("magnit: create receipts request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("magnit: receipts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("magnit: receipts failed (%d): %s", resp.StatusCode, string(body))
	}

	var recResp ReceiptsResponse
	if err := json.NewDecoder(resp.Body).Decode(&recResp); err != nil {
		return nil, 0, fmt.Errorf("magnit: decode receipts: %w", err)
	}

	return recResp.Data.Receipts, recResp.Data.Pages, nil
}

func (c *Client) demoReceipts(page int) ([]Receipt, int, error) {
	allReceipts := []Receipt{
		{
			ID: "mg-demo-001", StoreName: "Магнит", StoreAddr: "ул. Ленина, 1",
			Date: "2026-05-28", Total: 845.30, ItemsCnt: 4,
			Items: []ReceiptItem{
				{Name: "Молоко Простоквашино 3.2%", Price: 89.99, Quantity: 1},
				{Name: "Батон нарезной", Price: 52.00, Quantity: 2},
				{Name: "Сосиски молочные", Price: 320.50, Quantity: 1},
				{Name: "Масло подсолнечное", Price: 149.99, Quantity: 1},
			},
		},
		{
			ID: "mg-demo-002", StoreName: "Магнит Косметик", StoreAddr: "пр. Мира, 15",
			Date: "2026-05-26", Total: 1230.00, ItemsCnt: 5,
			Items: []ReceiptItem{
				{Name: "Шампунь Head&Shoulders", Price: 450.00, Quantity: 1},
				{Name: "Зубная паста Colgate", Price: 280.00, Quantity: 1},
				{Name: "Крем для рук", Price: 350.00, Quantity: 1},
				{Name: "Мыло жидкое", Price: 150.00, Quantity: 1},
			},
		},
		{
			ID: "mg-demo-003", StoreName: "Магнит Семейный", StoreAddr: "ул. Победы, 50",
			Date: "2026-05-24", Total: 2765.80, ItemsCnt: 6,
			Items: []ReceiptItem{
				{Name: "Курица охлаждённая", Price: 380.90, Quantity: 1},
				{Name: "Рис круглозёрный", Price: 120.50, Quantity: 2},
				{Name: "Гречка", Price: 95.00, Quantity: 1},
				{Name: "Тушёнка говяжья", Price: 450.00, Quantity: 1},
				{Name: "Макароны", Price: 89.90, Quantity: 2},
			},
		},
		{
			ID: "mg-demo-004", StoreName: "Магнит", StoreAddr: "ул. Гагарина, 23",
			Date: "2026-05-22", Total: 567.40, ItemsCnt: 3,
			Items: []ReceiptItem{
				{Name: "Кефир 1%", Price: 75.50, Quantity: 2},
				{Name: "Творог обезжиренный", Price: 135.00, Quantity: 1},
				{Name: "Сметана 15%", Price: 89.90, Quantity: 1},
			},
		},
		{
			ID: "mg-demo-005", StoreName: "Магнит Косметик", StoreAddr: "ул. Кирова, 8",
			Date: "2026-05-20", Total: 1980.00, ItemsCnt: 4,
			Items: []ReceiptItem{
				{Name: "Духи женские", Price: 1200.00, Quantity: 1},
				{Name: "Лак для ногтей", Price: 280.00, Quantity: 2},
				{Name: "Помада", Price: 450.00, Quantity: 1},
			},
		},
	}

	limit := 20
	start := (page - 1) * limit
	if start > len(allReceipts) {
		return nil, 0, nil
	}
	end := start + limit
	if end > len(allReceipts) {
		end = len(allReceipts)
	}

	totalPages := (len(allReceipts) + limit - 1) / limit
	return allReceipts[start:end], totalPages, nil
}
