package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handler — HTTP API банковских счетов (демо по API_Contract).
type Handler struct{}

// New создаёт обработчик bank-service.
func New() *Handler {
	return &Handler{}
}

// Register монтирует маршруты /api/v1/banks/*.
func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/banks/accounts", h.accounts)
	r.Get("/api/v1/banks/transactions", h.transactions)
}

type accountItem struct {
	ID       string  `json:"id"`
	Bank     string  `json:"bank"`
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type accountsResponse struct {
	Accounts []accountItem `json:"accounts"`
}

func (h *Handler) accounts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, accountsResponse{
		Accounts: []accountItem{
			{
				ID:       "demo-account-1",
				Bank:     "Т-Банк",
				Name:     "Дебетовая Tinkoff Black",
				Balance:  340000,
				Currency: "RUB",
			},
		},
	})
}

type transactionItem struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

type transactionsResponse struct {
	Transactions []transactionItem `json:"transactions"`
}

func (h *Handler) transactions(w http.ResponseWriter, r *http.Request) {
	today := time.Now().UTC().Format("2006-01-02")
	writeJSON(w, transactionsResponse{
		Transactions: []transactionItem{
			{
				ID:          "demo-tx-1",
				Date:        today,
				Amount:      180000,
				Description: "Зарплата",
				Category:    "income",
			},
			{
				ID:          "demo-tx-2",
				Date:        today,
				Amount:      -1200,
				Description: "Пятёрочка",
				Category:    "food",
			},
		},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
