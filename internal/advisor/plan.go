package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"backend_project/internal/llm"
	"backend_project/internal/profile"
)

type PlanStep struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Plan struct {
	GoalTitle        string     `json:"goalTitle"`
	GoalProgress     string     `json:"goalProgress"`
	Steps            []PlanStep `json:"steps"`
	RunwayText       *string    `json:"runwayText"`
	FreeCashflowText *string    `json:"freeCashflowText"`
	UpdatedAt        int64      `json:"updatedAt"`
}

type MainAction struct {
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	PotentialSavings  float64 `json:"potential_savings"`
	Difficulty        string  `json:"difficulty"`
}

type Indicator struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Norm   string  `json:"norm"`
	Status string  `json:"status"`
}

type Diagnosis struct {
	Score         int         `json:"score"`
	Grade         string      `json:"grade"`
	Indicators    []Indicator `json:"indicators"`
	MainAction    MainAction  `json:"main_action"`
	NextCheckDays int         `json:"next_check_days"`
}

type PlanResponse struct {
	Plan      Plan      `json:"plan"`
	Diagnosis Diagnosis `json:"diagnosis"`
}

func BuildPlanResponse(snap Snapshot, client *llm.Client) PlanResponse {
	if client != nil && client.Enabled() {
		if resp, ok := tryLLMPlan(snap, client); ok {
			return resp
		}
	}
	return heuristicPlan(snap)
}

func tryLLMPlan(snap Snapshot, client *llm.Client) (PlanResponse, bool) {
	ctxJSON, err := json.Marshal(snap)
	if err != nil {
		return PlanResponse{}, false
	}
	raw, err := client.Complete(context.Background(), llm.PlanGenerationPrompt, string(ctxJSON))
	if err != nil {
		return PlanResponse{}, false
	}
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end <= start {
		return PlanResponse{}, false
	}
	var resp PlanResponse
	if err := json.Unmarshal([]byte(raw[start:end+1]), &resp); err != nil {
		return PlanResponse{}, false
	}
	if len(resp.Plan.Steps) < 3 {
		return PlanResponse{}, false
	}
	if resp.Plan.UpdatedAt == 0 {
		resp.Plan.UpdatedAt = time.Now().UnixMilli()
	}
	return resp, true
}

func heuristicPlan(snap Snapshot) PlanResponse {
	p := snap.Profile
	goalTitle := p.GoalTitle
	if goalTitle == "" {
		goalTitle = "Финансовая цель"
	}
	goalProgress := buildGoalProgress(p)
	runway := RunwayMonths(p)
	var runwayText *string
	if runway != nil {
		t := fmt.Sprintf("Запас примерно на %.0f мес.", math.Floor(*runway))
		runwayText = &t
	}
	fcf := FreeCashflow(p)
	var fcfText *string
	if fcf != 0 && snap.DataCompleteness["income"] != "skipped" {
		t := fmt.Sprintf("После расходов остаётся %.0f ₽/мес.", fcf)
		fcfText = &t
	}

	main := MainAction{
		Title:            "Записывайте траты голосом",
		Description:      "Так картина денег и советы обновляются быстрее.",
		PotentialSavings: 0,
		Difficulty:       "easy",
	}
	if len(snap.Credits.Credits) > 0 {
		c := snap.Credits.Credits[0]
		if c.RateVsMarket == "above" && c.BenchmarkRate > 0 {
			main = MainAction{
				Title:            "Сравните ставку по кредиту",
				Description:      fmt.Sprintf("Ваша ставка %.1f%%, рынок около %.1f%% — обсудите рефинанс.", c.Rate, c.BenchmarkRate),
				PotentialSavings: c.MonthlyPayment * 0.1,
				Difficulty:       "medium",
			}
		}
	}
	if snap.DataCompleteness["income"] == "skipped" {
		main = MainAction{
			Title:       "Укажите доход в профиле",
			Description: "Тогда прогноз и советы станут точнее.",
			Difficulty:  "easy",
		}
	}

	steps := []PlanStep{
		{Title: main.Title, Description: main.Description},
		{Title: "Проверьте прогресс цели", Description: goalProgress},
		{Title: "Добавляйте покупки", Description: "Записывайте траты — советник учитывает их в плане."},
	}

	score := 72
	grade := "B"
	dti := snap.Credits.DTI
	indicators := []Indicator{}
	if dti > 0 {
		st := "good"
		if dti >= 50 {
			st = "critical"
		} else if dti >= 35 {
			st = "warning"
		}
		indicators = append(indicators, Indicator{Name: "Долговая нагрузка", Value: dti, Norm: "<35", Status: st})
	}

	return PlanResponse{
		Plan: Plan{
			GoalTitle:        goalTitle,
			GoalProgress:     goalProgress,
			Steps:            steps,
			RunwayText:       runwayText,
			FreeCashflowText: fcfText,
			UpdatedAt:        time.Now().UnixMilli(),
		},
		Diagnosis: Diagnosis{
			Score:         score,
			Grade:         grade,
			Indicators:    indicators,
			MainAction:    main,
			NextCheckDays: 30,
		},
	}
}

func buildGoalProgress(p profile.FinancialProfile) string {
	if p.SkippedGoal || p.GoalAmount < 1000 {
		return "Цель не задана — можно указать в профиле."
	}
	return fmt.Sprintf("Цель «%s» — %.0f ₽.", p.GoalTitle, p.GoalAmount)
}

func BuildDiagnosis(snap Snapshot, client *llm.Client) Diagnosis {
	return BuildPlanResponse(snap, client).Diagnosis
}
