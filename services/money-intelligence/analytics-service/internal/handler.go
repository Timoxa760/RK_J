package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handler — insights, forecast и сценарии (демо по API_Contract).
type Handler struct {
	demoMode bool
}

// New создаёт обработчик analytics-service.
func New(demoMode bool) *Handler {
	return &Handler{demoMode: demoMode}
}

// Register монтирует маршруты money-intelligence API.
func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/insights", h.insights)
	r.Get("/api/v1/forecast", h.forecast)
	r.Post("/api/v1/scenarios/simulate", h.simulate)
}

type insightItem struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount,omitempty"`
	Merchant    string  `json:"merchant,omitempty"`
	Store       string  `json:"store,omitempty"`
}

type insightsResponse struct {
	Insights []insightItem `json:"insights"`
}

func (h *Handler) insights(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, insightsResponse{
		Insights: []insightItem{
			{
				Type:        "subscription",
				Severity:    "warning",
				Title:       "Найдена скрытая подписка",
				Description: "Списывается 299 ₽ каждый месяц",
				Amount:      299,
				Merchant:    "Яндекс.Плюс",
			},
			{
				Type:        "duplicate",
				Severity:    "info",
				Title:       "Дублирование в чеке",
				Description: "Товар «Молоко 3.2%» пробит дважды",
				Amount:      156,
			},
			{
				Type:        "overprice",
				Severity:    "warning",
				Title:       "Переплата за товар",
				Description: "Молоко 3.2% куплено за 95 ₽, средняя — 78 ₽",
				Amount:      17,
				Store:       "Пятёрочка",
			},
		},
	})
}

type forecastPoint struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type forecastResponse struct {
	Days    int             `json:"days"`
	Total   float64         `json:"total"`
	Points  []forecastPoint `json:"points"`
}

func (h *Handler) forecast(w http.ResponseWriter, r *http.Request) {
	start := time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC)
	points := make([]forecastPoint, 7)
	var total float64
	for i := 0; i < 7; i++ {
		amt := 9800 + float64(i*120)
		points[i] = forecastPoint{
			Date:   start.AddDate(0, 0, i).Format("2006-01-02"),
			Amount: amt,
		}
		total += amt
	}
	writeJSON(w, forecastResponse{Days: 7, Total: total, Points: points})
}

type simulateRequest struct {
	Scenario          string  `json:"scenario"`
	ReductionPercent  float64 `json:"reduction_percent"`
	Months            int     `json:"months"`
}

type scenarioMeta struct {
	Name          string  `json:"name"`
	MonthlySaving float64 `json:"monthly_saving"`
	AnnualSaving  float64 `json:"annual_saving"`
}

type simulateResponse struct {
	Months            []string     `json:"months"`
	RealSavings       []float64    `json:"real_savings"`
	OptimizedSavings  []float64    `json:"optimized_savings"`
	DifferenceFinal   float64      `json:"difference_final"`
	Scenario          scenarioMeta `json:"scenario"`
}

var allowedScenarios = map[string]float64{
	"reduce_delivery":       4500,
	"reduce_cafe":         3200,
	"reduce_entertainment":  2800,
	"custom":                2000,
}

func (h *Handler) simulate(w http.ResponseWriter, r *http.Request) {
	var req simulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	baseSaving, ok := allowedScenarios[req.Scenario]
	if !ok {
		http.Error(w, `{"error":"invalid scenario"}`, http.StatusBadRequest)
		return
	}
	if req.Months <= 0 || req.Months > 60 {
		req.Months = 12
	}
	if req.ReductionPercent > 0 {
		baseSaving *= req.ReductionPercent / 100
	}

	months := make([]string, req.Months)
	real := make([]float64, req.Months)
	opt := make([]float64, req.Months)
	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	curReal := 500000.0
	curOpt := 500000.0
	for i := 0; i < req.Months; i++ {
		months[i] = start.AddDate(0, i, 0).Format("2006-01")
		curReal += 12000
		curOpt += 12000 + baseSaving
		real[i] = curReal
		opt[i] = curOpt
	}

	writeJSON(w, simulateResponse{
		Months:           months,
		RealSavings:      real,
		OptimizedSavings: opt,
		DifferenceFinal:  opt[req.Months-1] - real[req.Months-1],
		Scenario: scenarioMeta{
			Name:          req.Scenario,
			MonthlySaving: baseSaving,
			AnnualSaving:  baseSaving * 12,
		},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
