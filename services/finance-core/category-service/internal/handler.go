package internal

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler — справочник категорий расходов (демо-набор для UI).
type Handler struct{}

// New создаёт обработчик category-service.
func New() *Handler {
	return &Handler{}
}

// Register монтирует маршруты /api/v1/categories.
func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/categories", h.list)
}

type categoryItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
}

type listResponse struct {
	Categories []categoryItem `json:"categories"`
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, listResponse{
		Categories: []categoryItem{
			{ID: "cat-food", Name: "Продукты", Color: "#4CAF50"},
			{ID: "cat-cafe", Name: "Кафе и рестораны", Color: "#FF9800"},
			{ID: "cat-transport", Name: "Транспорт", Color: "#2196F3"},
			{ID: "cat-fun", Name: "Развлечения", Color: "#9C27B0"},
			{ID: "cat-other", Name: "Прочее", Color: "#607D8B"},
		},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
