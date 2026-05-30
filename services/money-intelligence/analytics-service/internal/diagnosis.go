package internal

import (
	"net/http"
)

type diagnosisIndicator struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Norm   string  `json:"norm"`
	Status string  `json:"status"`
}

type diagnosisAction struct {
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	PotentialSavings   float64 `json:"potential_savings"`
	Difficulty         string  `json:"difficulty"`
}

type diagnosisResponse struct {
	Score          int                  `json:"score"`
	Grade          string               `json:"grade"`
	Indicators     []diagnosisIndicator `json:"indicators"`
	MainAction     diagnosisAction      `json:"main_action"`
	NextCheckDays  int                  `json:"next_check_days"`
}

func (h *Handler) diagnosis(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, diagnosisResponse{
		Score: 72,
		Grade: "B",
		Indicators: []diagnosisIndicator{
			{Name: "Долговая нагрузка", Value: 28, Norm: "<30", Status: "good"},
			{Name: "Подушка безопасности", Value: 4.2, Norm: ">3", Status: "good"},
			{Name: "Накопления от дохода", Value: 15, Norm: ">20", Status: "warning"},
			{Name: "Импульсивные траты", Value: 32, Norm: "<25", Status: "critical"},
			{Name: "Стабильность доходов", Value: 85, Norm: ">70", Status: "good"},
		},
		MainAction: diagnosisAction{
			Title:            "Сократите доставку еды",
			Description:      "Вы тратите 9 000 ₽ в месяц на доставку. Готовьте дома 3 раза в неделю — это сэкономит 4 500 ₽ в месяц.",
			PotentialSavings: 4500,
			Difficulty:       "easy",
		},
		NextCheckDays: 30,
	})
}
