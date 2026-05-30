package internal

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler — отчёты и дайджест (demo).
type Handler struct{}

// New создаёт обработчик reporting-service.
func New() *Handler {
	return &Handler{}
}

// Register монтирует маршруты reporting API.
func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/digest/latest", h.digestLatest)
}

type periodRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type categoryDigest struct {
	Name    string  `json:"name"`
	Total   float64 `json:"total"`
	Percent float64 `json:"percent"`
	Trend   string  `json:"trend"`
}

type storeDigest struct {
	Name   string  `json:"name"`
	Total  float64 `json:"total"`
	Visits int     `json:"visits"`
}

type digestResponse struct {
	Period            periodRange      `json:"period"`
	TotalSpent        float64          `json:"total_spent"`
	TotalIncome       float64          `json:"total_income"`
	Saved             float64          `json:"saved"`
	ByCategory        []categoryDigest `json:"by_category"`
	WordCloud         []string         `json:"word_cloud"`
	TopStores         []storeDigest    `json:"top_stores"`
	MindfulnessRating int              `json:"mindfulness_rating"`
	AIAdvice          string           `json:"ai_advice"`
	InsightsSummary   string           `json:"insights_summary"`
}

func (h *Handler) digestLatest(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, digestResponse{
		Period:      periodRange{From: "2026-04-01", To: "2026-04-30"},
		TotalSpent:  145000,
		TotalIncome: 180000,
		Saved:       35000,
		ByCategory: []categoryDigest{
			{Name: "Продукты", Total: 52000, Percent: 35.9, Trend: "+8.3%"},
			{Name: "Кафе", Total: 28000, Percent: 19.3, Trend: "+2.1%"},
		},
		WordCloud: []string{"молоко", "латте", "хлеб", "сыр", "такси"},
		TopStores: []storeDigest{
			{Name: "Пятёрочка", Total: 9100, Visits: 14},
		},
		MindfulnessRating: 72,
		AIAdvice:          "Попробуйте сократить доставку — это 9 000 ₽ в месяц",
		InsightsSummary:   "Найдено 2 скрытые подписки и 3 переплаты",
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
