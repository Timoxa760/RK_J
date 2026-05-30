package manual

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"backend_project/services/money-intelligence/ai-processor/internal/parser"
)

// DemoHandler обрабатывает POST /expenses/manual без БД (DEMO_MODE).
type DemoHandler struct {
	mu      sync.Mutex
	byUser  map[string][]*Expense
}

// NewDemoHandler создаёт in-memory обработчик ручного ввода.
func NewDemoHandler() *DemoHandler {
	return &DemoHandler{byUser: make(map[string][]*Expense)}
}

// Create парсит текст и сохраняет расход в памяти.
func (h *DemoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, `{"error":"user_id required"}`, http.StatusBadRequest)
		return
	}
	if req.Source == "" {
		req.Source = "manual"
	}

	expense := &Expense{
		UserID:  req.UserID,
		RawText: req.RawText,
		Source:  req.Source,
		Date:    time.Now(),
	}
	if req.Date != "" {
		if t, err := time.Parse("2006-01-02", req.Date); err == nil {
			expense.Date = t
		}
	}

	parsed := false
	if req.RawText != "" {
		if p := parser.Parse(req.RawText); p != nil {
			expense.Amount = p.Amount
			expense.Category = p.Category
			expense.Description = p.Description
			parsed = true
		}
	}
	if req.Amount > 0 {
		expense.Amount = req.Amount
	}
	if req.Category != "" {
		expense.Category = req.Category
	}
	if req.Description != "" {
		expense.Description = req.Description
	}
	if expense.Amount <= 0 {
		http.Error(w, `{"error":"amount required"}`, http.StatusBadRequest)
		return
	}
	if expense.Category == "" {
		expense.Category = "Прочие расходы"
	}

	expense.ID = fmt.Sprintf("demo-%d", time.Now().UnixNano())
	expense.CreatedAt = time.Now()

	h.mu.Lock()
	h.byUser[req.UserID] = append(h.byUser[req.UserID], expense)
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateResponse{
		Success:  true,
		ID:       expense.ID,
		Amount:   expense.Amount,
		Category: expense.Category,
		Parsed:   parsed,
	})
}
