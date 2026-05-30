package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handler — бюджеты пользователя (демо-набор для UI).
type Handler struct{}

// New создаёт обработчик budget-service.
func New() *Handler {
	return &Handler{}
}

// Register монтирует маршруты /api/v1/budgets.
func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/budgets", h.list)
}

type budgetItem struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Spent     float64 `json:"spent"`
	Period    string  `json:"period"`
	StartDate string  `json:"start_date"`
}

type listResponse struct {
	Budgets []budgetItem `json:"budgets"`
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UTC().Format("2006-01-01")
	writeJSON(w, listResponse{
		Budgets: []budgetItem{
			{
				ID:        "budget-food",
				Name:      "Продукты",
				Amount:    60000,
				Spent:     0,
				Period:    "monthly",
				StartDate: start,
			},
		},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
