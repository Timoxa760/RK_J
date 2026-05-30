package profile

import "time"

// FinancialProfile — финансовый профиль пользователя (API + snapshot).
type FinancialProfile struct {
	ActiveIncome        float64              `json:"active_income"`
	PassiveIncome       float64              `json:"passive_income"`
	EmergencyFund       float64              `json:"emergency_fund"`
	EmergencyBreakdown  EmergencyBreakdown   `json:"emergency_breakdown"`
	FixedExpenses       []FixedExpense       `json:"fixed_expenses"`
	GoalKind            string               `json:"goal_kind"`
	GoalTitle           string               `json:"goal_title"`
	GoalAmount          float64              `json:"goal_amount"`
	SkippedIncome       bool                 `json:"skipped_income"`
	SkippedCushion      bool                 `json:"skipped_cushion"`
	SkippedGoal         bool                 `json:"skipped_goal"`
	SkippedExpenses     bool                 `json:"skipped_expenses"`
	SurveyInputMode     string               `json:"survey_input_mode,omitempty"`
	OnboardingCompleted bool                 `json:"onboarding_completed"`
	UpdatedAt           time.Time            `json:"updated_at,omitempty"`
}

type EmergencyBreakdown struct {
	Cash         float64 `json:"cash"`
	Deposit      float64 `json:"deposit"`
	Investments  float64 `json:"investments"`
}

type FixedExpense struct {
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
}

func Default() FinancialProfile {
	return FinancialProfile{
		EmergencyBreakdown: EmergencyBreakdown{},
		FixedExpenses:      []FixedExpense{},
		GoalKind:           "save",
	}
}
