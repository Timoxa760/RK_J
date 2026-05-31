package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend_project/internal/llm"
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
	Reply   string       `json:"reply"`
	Title   string       `json:"title,omitempty"`
	Blocks  []ReplyBlock `json:"blocks,omitempty"`
	Actions []ChatAction `json:"actions,omitempty"`
	Source  string       `json:"source"`
	ID      string       `json:"id,omitempty"`
}

type ChatResult struct {
	Reply      string
	Title      string
	Blocks     []ReplyBlock
	Structured StructuredReply
	Actions    []ChatAction
	Source     string
}

func BuildChatReply(snap Snapshot, req ChatRequest, client *llm.Client) ChatResult {
	source := "heuristic"
	var structured StructuredReply
	var plain string

	if client != nil && client.Enabled() {
		ctxJSON, _ := json.Marshal(snap)
		hist, _ := json.Marshal(req.History)
		userPrompt := fmt.Sprintf("Snapshot:\n%s\n\nHistory:\n%s\n\nUser message:\n%s", ctxJSON, hist, req.Message)
		if raw, err := client.Complete(context.Background(), llm.AdvisorSystemPrompt, userPrompt); err == nil && strings.TrimSpace(raw) != "" {
			structured, plain = ParseStructuredReply(raw)
			source = "gemini"
		}
	}
	if len(structured.Blocks) == 0 {
		plain = heuristicChat(snap, req.Message)
		structured = structuredFromPlain(plain)
		source = "heuristic"
	}

	actions := BuildChatActions(snap, req.Message, plain)
	return chatResultFromStructured(structured, plain, actions, source)
}

func BuildChatReplyStream(
	ctx context.Context,
	snap Snapshot,
	req ChatRequest,
	client *llm.Client,
	onDelta func(string) error,
) ChatResult {
	source := "heuristic"
	var structured StructuredReply
	var plain string

	if client != nil && client.Enabled() {
		ctxJSON, _ := json.Marshal(snap)
		hist, _ := json.Marshal(req.History)
		userPrompt := fmt.Sprintf("Snapshot:\n%s\n\nHistory:\n%s\n\nUser message:\n%s", ctxJSON, hist, req.Message)
		full, err := client.StreamComplete(ctx, llm.AdvisorSystemPrompt, userPrompt, onDelta)
		if err == nil && strings.TrimSpace(full) != "" {
			structured, plain = ParseStructuredReply(full)
			source = "gemini"
		}
	}
	if len(structured.Blocks) == 0 {
		plain = heuristicChat(snap, req.Message)
		structured = structuredFromPlain(plain)
		source = "heuristic"
	}

	actions := BuildChatActions(snap, req.Message, plain)
	result := chatResultFromStructured(structured, plain, actions, source)
	if onDelta != nil && source == "heuristic" {
		_ = onDelta(plain)
	}
	return result
}

func chatResultFromStructured(structured StructuredReply, plain string, actions []ChatAction, source string) ChatResult {
	return ChatResult{
		Reply:      plain,
		Title:      structured.Title,
		Blocks:     structured.Blocks,
		Structured: structured,
		Actions:    actions,
		Source:     source,
	}
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
		if top := TopCategory(snap.Spending); top != nil {
			saving := top.Amount * 0.15
			return fmt.Sprintf(
				"Больше всего уходит на «%s» — %.0f ₽ за месяц (%.0f%%). Попробуйте сократить на 10–15%% — это ~%.0f ₽/мес.",
				top.Name, top.Amount, top.Share*100, saving,
			)
		}
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
	if top := TopCategory(snap.Spending); top != nil {
		return fmt.Sprintf(
			"Главная статья расходов — «%s» (%.0f ₽/мес). %s",
			top.Name, top.Amount, plan.Diagnosis.MainAction.Description,
		)
	}
	return fmt.Sprintf("%s %s", plan.Diagnosis.MainAction.Title, plan.Diagnosis.MainAction.Description)
}
