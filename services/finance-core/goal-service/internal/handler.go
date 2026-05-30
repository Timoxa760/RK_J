package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// Goal — финансовая цель пользователя.
type Goal struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id,omitempty"`
	Title           string    `json:"title"`
	TargetAmount    float64   `json:"target_amount"`
	CurrentAmount   float64   `json:"current_amount"`
	ProgressPercent float64   `json:"progress_percent"`
	TargetDate      string    `json:"target_date,omitempty"`
	AutoSavePercent float64   `json:"auto_save_percent,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

// Handler — CRUD целей (demo in-memory).
type Handler struct {
	mu    sync.RWMutex
	goals map[string]Goal
	seq   int
}

// New создаёт обработчик goal-service.
func New() *Handler {
	return &Handler{goals: make(map[string]Goal)}
}

// Register монтирует маршруты /api/v1/goals.
func (h *Handler) Register(r chi.Router) {
	r.Post("/api/v1/goals", h.create)
	r.Get("/api/v1/goals/{id}", h.getByID)
	r.Get("/api/v1/goals", h.list)
}

type createRequest struct {
	UserID          string  `json:"user_id"`
	Title           string  `json:"title"`
	TargetAmount    float64 `json:"target_amount"`
	TargetDate      string  `json:"target_date"`
	AutoSavePercent float64 `json:"auto_save_percent"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.Title == "" || req.TargetAmount <= 0 {
		http.Error(w, `{"error":"title and target_amount required"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.seq++
	id := fmt.Sprintf("goal-%d", h.seq)
	g := Goal{
		ID:              id,
		UserID:          req.UserID,
		Title:           req.Title,
		TargetAmount:    req.TargetAmount,
		CurrentAmount:   0,
		ProgressPercent: 0,
		TargetDate:      req.TargetDate,
		AutoSavePercent: req.AutoSavePercent,
		CreatedAt:       time.Now().UTC(),
	}
	h.goals[id] = g
	h.mu.Unlock()

	writeJSON(w, http.StatusOK, g)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.mu.RLock()
	g, ok := h.goals[id]
	h.mu.RUnlock()
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	items := make([]Goal, 0, len(h.goals))
	for _, g := range h.goals {
		items = append(items, g)
	}
	h.mu.RUnlock()
	writeJSON(w, http.StatusOK, map[string][]Goal{"goals": items})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
