package manual

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend_project/services/money-intelligence/ai-processor/internal/parser"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

type CreateRequest struct {
	UserID      string  `json:"user_id"`
	RawText     string  `json:"raw_text"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Source      string  `json:"source"`
}

type CreateResponse struct {
	Success  bool    `json:"success"`
	ID       string  `json:"id,omitempty"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Parsed   bool    `json:"parsed"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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
		UserID: req.UserID,
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
		p := parser.Parse(req.RawText)
		if p != nil {
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

	if err := h.repo.Insert(r.Context(), expense); err != nil {
		log.Printf("manual: insert: %v", err)
		http.Error(w, `{"error":"save failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateResponse{
		Success:  true,
		ID:       expense.ID,
		Amount:   expense.Amount,
		Category: expense.Category,
		Parsed:   parsed,
	})
}
