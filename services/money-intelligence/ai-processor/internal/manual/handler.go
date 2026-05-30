package manual

import (
	"encoding/json"
	"net/http"

	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

// Handler обрабатывает POST /expenses/manual с PostgreSQL.
type Handler struct {
	proc *Processor
}

// NewHandler создаёт HTTP-обработчик manual expenses и Processor для voice.
func NewHandler(repo *Repo, parser *expense.Parser) (*Handler, *Processor) {
	return NewHandlerWithStorage(newRepoStorage(repo), parser)
}

// NewHandlerWithStorage создаёт обработчик с произвольным Storage (PG, file, fallback).
func NewHandlerWithStorage(store Storage, parser *expense.Parser) (*Handler, *Processor) {
	proc := NewProcessor(parser, store)
	return &Handler{proc: proc}, proc
}

// Create принимает JSON и сохраняет траты.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	resp, code, err := h.proc.Create(r.Context(), req)
	if err != nil {
		respondError(w, code, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
