package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handler — HTTP API кредитного дашборда (демо-данные по API_Contract).
type Handler struct {
	demoMode bool
}

// New создаёт обработчик credit-service.
func New(demoMode bool) *Handler {
	return &Handler{demoMode: demoMode}
}

// Register монтирует маршруты /api/v1/credits/*.
func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/credits/dashboard", h.dashboard)
	r.Post("/api/v1/credits/scan", h.scan)
}

type creditItem struct {
	ID             string  `json:"id"`
	Bank           string  `json:"bank"`
	Amount         float64 `json:"amount"`
	Rate           float64 `json:"rate"`
	TermMonths     int     `json:"term_months"`
	Remaining      float64 `json:"remaining"`
	MonthlyPayment float64 `json:"monthly_payment"`
	NextPayment    string  `json:"next_payment"`
}

type dashboardResponse struct {
	DTI               float64      `json:"dti"`
	StressTestMonths  float64      `json:"stress_test_months"`
	Savings           float64      `json:"savings"`
	TotalDebt         float64      `json:"total_debt"`
	MonthlyPayments   float64      `json:"monthly_payments"`
	MonthlyIncome     float64      `json:"monthly_income"`
	Credits           []creditItem `json:"credits"`
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	resp := dashboardResponse{
		DTI:              0.28,
		StressTestMonths: 4.2,
		Savings:          340000,
		TotalDebt:        1200000,
		MonthlyPayments:  42000,
		MonthlyIncome:    180000,
		Credits: []creditItem{
			{
				ID:             "demo-credit-1",
				Bank:           "Т-Банк",
				Amount:         1200000,
				Rate:           14.5,
				TermMonths:     36,
				Remaining:      980000,
				MonthlyPayment: 42000,
				NextPayment:    time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
			},
		},
	}
	writeJSON(w, resp)
}

type scanParsed struct {
	Amount         float64 `json:"amount"`
	Rate           float64 `json:"rate"`
	TermMonths     int     `json:"term_months"`
	MonthlyPayment float64 `json:"monthly_payment"`
	Bank           string  `json:"bank"`
}

type scanResponse struct {
	Parsed     scanParsed `json:"parsed"`
	Confidence float64    `json:"confidence"`
}

func (h *Handler) scan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		http.Error(w, `{"error":"invalid multipart"}`, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"file required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()
	_ = header
	writeJSON(w, scanResponse{
		Parsed: scanParsed{
			Amount:         1200000,
			Rate:           14.5,
			TermMonths:     36,
			MonthlyPayment: 42000,
			Bank:           "Т-Банк",
		},
		Confidence: 0.87,
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
