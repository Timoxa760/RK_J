package expense

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"backend_project/internal/llm"
	"backend_project/services/money-intelligence/ai-processor/internal/parser"
)

var jsonFenceRe = regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")

type llmPayload struct {
	Expenses []Item `json:"expenses"`
	Advice   string `json:"advice"`
}

// Parser извлекает траты из текста через Gemini с fallback на regex.
type Parser struct {
	llm *llm.Client
}

// NewParser создаёт парсер расходов.
func NewParser(llmClient *llm.Client) *Parser {
	return &Parser{llm: llmClient}
}

// ParseInput — входные переопределения из HTTP-запроса.
type ParseInput struct {
	RawText     string
	Amount      float64
	Category    string
	Description string
}

// Parse возвращает траты и совет из текста пользователя.
func (p *Parser) Parse(ctx context.Context, in ParseInput) Result {
	text := strings.TrimSpace(in.RawText)
	if text != "" && p.llm != nil && p.llm.Enabled() {
		if res, ok := p.parseLLM(ctx, text); ok {
			p.applyOverrides(&res, in)
			return res
		}
	}
	res := p.parseRegex(text, in)
	p.applyOverrides(&res, in)
	return res
}

func (p *Parser) parseLLM(ctx context.Context, text string) (Result, bool) {
	raw, err := p.llm.Complete(ctx, llm.ExpenseSystemPrompt, text)
	if err != nil {
		return Result{}, false
	}
	raw = extractJSON(raw)
	var payload llmPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return Result{}, false
	}
	items := normalizeItems(payload.Expenses)
	if len(items) == 0 {
		return Result{}, false
	}
	return Result{
		Expenses: items,
		Advice:   strings.TrimSpace(payload.Advice),
		ParsedBy: "gemini",
		Parsed:   true,
	}, true
}

func (p *Parser) parseRegex(text string, in ParseInput) Result {
	if text == "" && in.Amount <= 0 {
		return Result{}
	}
	if parsed := parser.ParseAll(text); len(parsed) > 0 {
		items := make([]Item, 0, len(parsed))
		for _, pe := range parsed {
			items = append(items, Item{
				Amount:      pe.Amount,
				Category:    pe.Category,
				Description: pe.Description,
			})
		}
		return Result{
			Expenses: items,
			ParsedBy: "regex",
			Parsed:   true,
		}
	}
	if in.Amount > 0 {
		cat := in.Category
		if cat == "" {
			cat = "Прочие расходы"
		}
		desc := in.Description
		if desc == "" {
			desc = text
		}
		return Result{
			Expenses: []Item{{
				Amount:      in.Amount,
				Category:    cat,
				Description: desc,
			}},
			ParsedBy: "regex",
			Parsed:   true,
		}
	}
	return Result{}
}

func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if m := jsonFenceRe.FindStringSubmatch(s); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return s
}

func normalizeItems(items []Item) []Item {
	out := make([]Item, 0, len(items))
	for _, it := range items {
		if it.Amount <= 0 {
			continue
		}
		if strings.TrimSpace(it.Category) == "" {
			it.Category = "Прочие расходы"
		}
		if strings.TrimSpace(it.Description) == "" {
			it.Description = it.Category
		}
		out = append(out, it)
	}
	return out
}

func (p *Parser) applyOverrides(res *Result, in ParseInput) {
	if len(res.Expenses) == 0 {
		return
	}
	if in.Amount > 0 {
		res.Expenses[0].Amount = in.Amount
	}
	if in.Category != "" {
		res.Expenses[0].Category = in.Category
	}
	if in.Description != "" {
		res.Expenses[0].Description = in.Description
	}
}
