package fns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type TicketRequest struct {
	FN string `json:"fn"`
	FD string `json:"fd"`
	FP string `json:"fp"`
}

type fnsResponse struct {
	Document struct {
		Receipt struct {
			RetailPlace string `json:"retailPlace"`
			DateTime    int64  `json:"dateTime"`
			TotalSum    int64  `json:"totalSum"`
			Items       []struct {
				Name     string `json:"name"`
				Sum      int64  `json:"sum"`
				Quantity int    `json:"quantity"`
				Price    int64  `json:"price"`
			} `json:"items"`
		} `json:"receipt"`
	} `json:"document"`
}

type Handler struct {
	httpCli  *http.Client
	demoMode bool
}

func NewHandler(demoMode bool) *Handler {
	return &Handler{
		httpCli: &http.Client{
			Timeout: 15 * time.Second,
		},
		demoMode: demoMode,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if req.FN == "" || req.FD == "" || req.FP == "" {
		http.Error(w, `{"error":"fn, fd, fp are required"}`, http.StatusBadRequest)
		return
	}

	receipt, err := h.CheckTicket(r.Context(), req.FN, req.FD, req.FP)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func (h *Handler) CheckTicket(ctx interface{}, fn, fd, fp string) (*scrap.RawReceipt, error) {
	if h.demoMode {
		return h.mockReceipt(), nil
	}

	url := fmt.Sprintf("https://proverkacheka.nalog.ru:8888/v1/inns/*/kkts/*/fss/%s/tickets/%s?fiscalSign=%s&sendToEmail=no", fn, fd, fp)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fns: create request: %w", err)
	}

	req.Header.Set("Device-Id", fmt.Sprintf("web-%d", time.Now().UnixNano()%1000000))
	req.Header.Set("Device-OS", "Windows")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := h.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fns: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fns: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fns: API returned %d: %s", resp.StatusCode, string(body))
	}

	var fnsResp fnsResponse
	if err := json.Unmarshal(body, &fnsResp); err != nil {
		return nil, fmt.Errorf("fns: parse response: %w", err)
	}

	return mapToReceipt(&fnsResp), nil
}

func mapToReceipt(fr *fnsResponse) *scrap.RawReceipt {
	r := fr.Document.Receipt
	rp := &scrap.RawReceipt{
		Provider: "fns",
		Store:    r.RetailPlace,
		Date:     time.UnixMilli(r.DateTime),
		Total:    float64(r.TotalSum) / 100,
	}

	for _, item := range r.Items {
		rp.Items = append(rp.Items, scrap.RawItem{
			Name:     strings.TrimSpace(item.Name),
			Price:    float64(item.Price) / 100,
			Quantity: item.Quantity,
		})
	}

	return rp
}

func (h *Handler) mockReceipt() *scrap.RawReceipt {
	return &scrap.RawReceipt{
		Provider: "fns",
		Store:    "Пятёрочка",
		Date:     time.Now().Add(-2 * time.Hour),
		Total:    1032.50,
		Items: []scrap.RawItem{
			{Name: "Молоко 3.2%", Price: 78.00, Quantity: 2},
			{Name: "Хлеб белый", Price: 45.00, Quantity: 1},
			{Name: "Сыр Российский", Price: 189.00, Quantity: 1},
			{Name: "Масло сливочное", Price: 150.00, Quantity: 1},
		},
	}
}
