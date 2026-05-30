package onboarding

import (
	"strings"

	"backend_project/internal/rublang"
)

type Step string

const (
	StepIncome   Step = "income"
	StepCushion  Step = "cushion"
	StepGoal     Step = "goal"
	StepExpenses Step = "expenses"
)

type FixedExpense struct {
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
}

type EmergencyBreakdownPatch struct {
	Cash        float64 `json:"cash"`
	Deposit     float64 `json:"deposit"`
	Investments float64 `json:"investments"`
}

type Patch struct {
	ActiveIncome         *float64                 `json:"active_income,omitempty"`
	PassiveIncome        *float64                 `json:"passive_income,omitempty"`
	EmergencyFund        *float64                 `json:"emergency_fund,omitempty"`
	EmergencyBreakdown   *EmergencyBreakdownPatch `json:"emergency_breakdown,omitempty"`
	GoalKind             *string                  `json:"goal_kind,omitempty"`
	GoalTitle            *string                  `json:"goal_title,omitempty"`
	GoalAmount           *float64                 `json:"goal_amount,omitempty"`
	FixedExpenses        []FixedExpense           `json:"fixed_expenses,omitempty"`
	SkippedExpenses      *bool                    `json:"skipped_expenses,omitempty"`
}

type ParseResult struct {
	Parsed bool  `json:"parsed"`
	Step   Step  `json:"step"`
	Patch  Patch `json:"patch"`
}

func Parse(step Step, rawText string) ParseResult {
	text := strings.TrimSpace(rawText)
	res := ParseResult{Step: step}
	if text == "" {
		return res
	}

	switch step {
	case StepIncome:
		res.Patch = parseIncome(text)
	case StepCushion:
		res.Patch = parseCushion(text)
	case StepGoal:
		res.Patch = parseGoal(text)
	case StepExpenses:
		res.Patch = parseExpenses(text)
	}

	res.Parsed = patchMeaningful(step, res.Patch)
	return res
}

func parseIncome(text string) Patch {
	nums := rublang.ExtractAll(text)
	passive := passiveIncomeHint(text)
	var p Patch

	switch {
	case len(nums) >= 2:
		a, b := float64(nums[0]), float64(nums[1])
		p.ActiveIncome, p.PassiveIncome = &a, &b
	case len(nums) == 1 && passive:
		v := float64(nums[0])
		zero := float64(0)
		p.ActiveIncome, p.PassiveIncome = &zero, &v
	case len(nums) == 1:
		v := float64(nums[0])
		zero := float64(0)
		p.ActiveIncome, p.PassiveIncome = &v, &zero
	}
	return p
}

func parseCushion(text string) Patch {
	if b := parseEmergencyBreakdown(text); b.hasAny() {
		total := b.total()
		return Patch{
			EmergencyFund: &total,
			EmergencyBreakdown: &EmergencyBreakdownPatch{
				Cash:        b.Cash,
				Deposit:     b.Deposit,
				Investments: b.Investments,
			},
		}
	}
	nums := rublang.ExtractAll(text)
	if len(nums) == 0 {
		return Patch{}
	}
	v := float64(nums[0])
	return Patch{EmergencyFund: &v}
}

func parseGoal(text string) Patch {
	normalized := rublang.Normalize(text)
	nums := rublang.ExtractAll(text)
	kind, title := detectGoalKind(normalized)

	var amount float64
	if len(nums) > 0 {
		amount = float64(nums[0])
	}

	return Patch{
		GoalKind:   &kind,
		GoalTitle:  &title,
		GoalAmount: &amount,
	}
}

func parseExpenses(text string) Patch {
	if isSkipAnswer(text) {
		skipped := true
		return Patch{SkippedExpenses: &skipped, FixedExpenses: []FixedExpense{}}
	}

	nums := rublang.ExtractAll(text)
	normalized := rublang.Normalize(text)
	items := []FixedExpense{}

	if strings.Contains(normalized, "аренд") || strings.Contains(normalized, "квартир") {
		if len(nums) > 0 {
			items = append(items, FixedExpense{Title: "Аренда", Amount: float64(nums[0])})
		}
	}
	if strings.Contains(normalized, "кредит") || strings.Contains(normalized, "ипотек") {
		idx := len(items)
		amount := float64(0)
		if idx < len(nums) {
			amount = float64(nums[idx])
		} else if len(nums) > 0 {
			amount = float64(nums[0])
		}
		if amount > 0 {
			items = append(items, FixedExpense{Title: "Кредит", Amount: amount})
		}
	}
	if strings.Contains(normalized, "парковк") {
		idx := len(items)
		if idx < len(nums) {
			items = append(items, FixedExpense{Title: "Парковка", Amount: float64(nums[idx])})
		}
	}
	if len(items) == 0 && len(nums) > 0 {
		items = append(items, FixedExpense{Title: "Обязательный платёж", Amount: float64(nums[0])})
	}

	return Patch{FixedExpenses: items}
}

func passiveIncomeHint(text string) bool {
	n := rublang.Normalize(text)
	return strings.Contains(n, "пассив") ||
		strings.Contains(n, "аренд") ||
		strings.Contains(n, "дивиденд") ||
		strings.Contains(n, "процент") ||
		strings.Contains(n, "подработк")
}

func isSkipAnswer(text string) bool {
	n := rublang.Normalize(text)
	keywords := []string{"пропуст", "нет обязательн", "не знаю", "позже", "ничего", "без обязательн", "пока нет"}
	for _, k := range keywords {
		if strings.Contains(n, k) {
			return true
		}
	}
	return false
}

func detectGoalKind(normalized string) (kind, title string) {
	kind, title = "save", "Накопления"
	switch {
	case strings.Contains(normalized, "подушк") || strings.Contains(normalized, "резерв") || strings.Contains(normalized, "запас"):
		return "cushion", "Подушка безопасности"
	case strings.Contains(normalized, "отпуск") || strings.Contains(normalized, "путешеств"):
		return "save", "Отпуск"
	case strings.Contains(normalized, "квартир") || strings.Contains(normalized, "машин") ||
		strings.Contains(normalized, "авто") || strings.Contains(normalized, "ремонт") ||
		strings.Contains(normalized, "ипотек") || strings.Contains(normalized, "покупк"):
		return "purchase", "Крупная покупка"
	case strings.Contains(normalized, "накоп"):
		return "save", "Накопления"
	default:
		return kind, title
	}
}

func patchMeaningful(step Step, p Patch) bool {
	switch step {
	case StepIncome:
		return (p.ActiveIncome != nil && *p.ActiveIncome > 0) ||
			(p.PassiveIncome != nil && *p.PassiveIncome > 0)
	case StepCushion:
		return p.EmergencyFund != nil && *p.EmergencyFund > 0
	case StepGoal:
		return p.GoalAmount != nil && *p.GoalAmount >= 1000
	case StepExpenses:
		if p.SkippedExpenses != nil && *p.SkippedExpenses {
			return true
		}
		return len(p.FixedExpenses) > 0
	default:
		return false
	}
}
