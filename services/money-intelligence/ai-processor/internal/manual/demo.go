package manual

import (
	"encoding/json"
	"net/http"

	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

// DemoHandler обрабатывает POST /expenses/manual без БД (DEMO_MODE).
type DemoHandler struct {
	proc *Processor
}

// NewDemoHandler создаёт in-memory обработчик и общий Processor для voice.
func NewDemoHandler(parser *expense.Parser) (*DemoHandler, *Processor) {
	store := newDemoStorage()
	proc := NewProcessor(parser, store)
	return &DemoHandler{proc: proc}, proc
}

// Create парсит текст и сохраняет расход в памяти.
func (h *DemoHandler) Create(w http.ResponseWriter, r *http.Request) {
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
