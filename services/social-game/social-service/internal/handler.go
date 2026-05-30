package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handler — челленджи и лидерборды (демо in-memory).
type Handler struct {
	mu         sync.RWMutex
	challenges map[string]challengeRecord
	seq        int
}

type challengeRecord struct {
	ID              string
	Type            string
	Status          string
	InviteToken     string
	Participants    int
	CreatedAt       time.Time
}

// New создаёт обработчик social-service.
func New() *Handler {
	return &Handler{challenges: make(map[string]challengeRecord)}
}

// Register монтирует маршруты /api/v1/challenges.
func (h *Handler) Register(r chi.Router) {
	r.Post("/api/v1/challenges", h.create)
	r.Get("/api/v1/challenges/{id}/leaderboard", h.leaderboard)
}

type createRequest struct {
	Type            string `json:"type"`
	Title           string `json:"title"`
	DurationDays    int    `json:"duration_days"`
	MaxParticipants int    `json:"max_participants"`
}

type createResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	InviteToken  string `json:"invite_token"`
	Participants int    `json:"participants"`
	CreatedAt    string `json:"created_at"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		http.Error(w, `{"error":"type required"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.seq++
	id := fmt.Sprintf("challenge-%d", h.seq)
	rec := challengeRecord{
		ID:           id,
		Type:         req.Type,
		Status:       "active",
		InviteToken:  fmt.Sprintf("invite-%d", h.seq),
		Participants: 1,
		CreatedAt:    time.Now().UTC(),
	}
	h.challenges[id] = rec
	h.mu.Unlock()

	writeJSON(w, createResponse{
		ID:           rec.ID,
		Type:         rec.Type,
		Status:       rec.Status,
		InviteToken:  rec.InviteToken,
		Participants: rec.Participants,
		CreatedAt:    rec.CreatedAt.Format(time.RFC3339),
	})
}

type leaderboardEntry struct {
	Position      int     `json:"position"`
	Username      string  `json:"username"`
	Avatar        *string `json:"avatar"`
	RelativeScore float64 `json:"relative_score"`
}

type myPosition struct {
	Position          int `json:"position"`
	TotalParticipants int `json:"total_participants"`
}

type leaderboardResponse struct {
	ChallengeID  string             `json:"challenge_id"`
	Type         string             `json:"type"`
	Leaderboard  []leaderboardEntry `json:"leaderboard"`
	MyPosition   myPosition         `json:"my_position"`
}

func (h *Handler) leaderboard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.mu.RLock()
	rec, ok := h.challenges[id]
	h.mu.RUnlock()
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	writeJSON(w, leaderboardResponse{
		ChallengeID: rec.ID,
		Type:        rec.Type,
		Leaderboard: []leaderboardEntry{
			{Position: 1, Username: "Анна", Avatar: nil, RelativeScore: 0},
			{Position: 2, Username: "Иван", Avatar: nil, RelativeScore: 0.35},
		},
		MyPosition: myPosition{Position: 2, TotalParticipants: 3},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
