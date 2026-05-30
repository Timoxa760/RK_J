package fns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

// ScanRequest — тело POST /receipt/fns/scan (контракт front).
type ScanRequest struct {
	FN string `json:"fn"`
	FD string `json:"fd"`
	FP string `json:"fp"`
}

// ScanResponse — ответ POST /receipt/fns/scan для Nuxt front.
type ScanResponse struct {
	ReceiptID string           `json:"receipt_id"`
	Store     string           `json:"store"`
	INN       string           `json:"inn,omitempty"`
	Date      string           `json:"date,omitempty"`
	Total     float64          `json:"total"`
	Items     []scrap.RawItem  `json:"items,omitempty"`
	Category  string           `json:"category"`
}

type ScanHandler struct {
	fns      *Handler
	producer receiptProducer
}

type receiptProducer interface {
	Send(ctx context.Context, receipt scrap.RawReceipt) error
}

// NewScanHandler создаёт обработчик front-контракта /receipt/fns/scan.
func NewScanHandler(fns *Handler, producer receiptProducer) *ScanHandler {
	return &ScanHandler{fns: fns, producer: producer}
}

func (h *ScanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.FN == "" || req.FD == "" || req.FP == "" {
		http.Error(w, `{"error":"fn, fd, fp are required"}`, http.StatusBadRequest)
		return
	}

	raw, err := h.fns.CheckTicket(r.Context(), req.FN, req.FD, req.FP)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	if h.producer != nil {
		if err := h.producer.Send(r.Context(), *raw); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"kafka: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	}

	resp := ScanResponse{
		ReceiptID: firstNonEmpty(raw.ID, uuid.NewString()),
		Store:     raw.Store,
		Total:     raw.Total,
		Items:     raw.Items,
		Category:  categorizeReceipt(raw),
		Date:      raw.Date.UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func categorizeReceipt(r *scrap.RawReceipt) string {
	if r == nil || len(r.Items) == 0 {
		return "Продукты"
	}
	return "Продукты"
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
