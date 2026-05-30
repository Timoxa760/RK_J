package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend_project/internal/onlysq"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

func BuildChatReply(snap Snapshot, req ChatRequest, llm *onlysq.Client) string {
	if llm != nil && llm.Enabled() {
		ctxJSON, _ := json.Marshal(snap)
		hist, _ := json.Marshal(req.History)
		userPrompt := fmt.Sprintf("Snapshot:\n%s\n\nHistory:\n%s\n\nUser message:\n%s", ctxJSON, hist, req.Message)
		if raw, err := llm.Complete(context.Background(), onlysq.AdvisorSystemPrompt, userPrompt); err == nil && strings.TrimSpace(raw) != "" {
			return strings.TrimSpace(raw)
		}
	}
	return heuristicChat(snap, req.Message)
}

func heuristicChat(snap Snapshot, message string) string {
	q := strings.ToLower(strings.TrimSpace(message))
	plan := BuildPlanResponse(snap, nil)

	if strings.Contains(q, "план") || strings.Contains(q, "шаг") {
		lines := []string{"Краткий план:"}
		for i, s := range plan.Plan.Steps {
			lines = append(lines, fmt.Sprintf("%d. %s — %s", i+1, s.Title, s.Description))
		}
		return strings.Join(lines, "\n")
	}
	if strings.Contains(q, "урез") || strings.Contains(q, "сократ") || strings.Contains(q, "эконом") {
		m := plan.Diagnosis.MainAction
		return fmt.Sprintf("%s: %s", m.Title, m.Description)
	}
	if strings.Contains(q, "цел") || strings.Contains(q, "когда") || strings.Contains(q, "дойду") {
		return plan.Plan.GoalProgress
	}
	if strings.Contains(q, "ставк") || strings.Contains(q, "кредит") {
		if len(snap.Credits.Credits) == 0 {
			return "Загрузите PDF договора в разделе кредитов — тогда смогу сравнить ставку с рынком."
		}
		c := snap.Credits.Credits[0]
		if c.BenchmarkRate > 0 {
			return fmt.Sprintf("Ставка %.1f%%, ориентир рынка %.1f%% (%s).", c.Rate, c.BenchmarkRate, c.RateVsMarket)
		}
	}
	return fmt.Sprintf("%s %s", plan.Diagnosis.MainAction.Title, plan.Diagnosis.MainAction.Description)
}
