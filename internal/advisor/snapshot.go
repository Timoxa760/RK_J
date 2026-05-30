package advisor

import (
	"backend_project/internal/creditstore"
	"backend_project/internal/profile"
)

// Snapshot — единый снимок для ИИ.
type Snapshot struct {
	Profile          profile.FinancialProfile `json:"profile"`
	Credits          creditstore.Dashboard    `json:"credits"`
	DataCompleteness map[string]string        `json:"data_completeness"`
}

func BuildSnapshot(profiles *profile.FileStore, credits *creditstore.FileStore, userID string) Snapshot {
	p, _ := profiles.Get(userID)
	income := p.ActiveIncome + p.PassiveIncome
	if p.SkippedIncome {
		income = 0
	}
	dash := credits.Dashboard(userID, income)
	dc := map[string]string{
		"income":   completenessIncome(p),
		"cushion":  completenessCushion(p),
		"goal":     completenessGoal(p),
		"expenses": completenessExpenses(p),
		"credits":  completenessCredits(dash),
	}
	return Snapshot{Profile: p, Credits: dash, DataCompleteness: dc}
}

func completenessIncome(p profile.FinancialProfile) string {
	if p.SkippedIncome {
		return "skipped"
	}
	if p.ActiveIncome+p.PassiveIncome > 0 {
		return "known"
	}
	return "unknown"
}

func completenessCushion(p profile.FinancialProfile) string {
	if p.SkippedCushion {
		return "skipped"
	}
	if p.EmergencyFund > 0 {
		return "known"
	}
	return "unknown"
}

func completenessGoal(p profile.FinancialProfile) string {
	if p.SkippedGoal {
		return "skipped"
	}
	if p.GoalAmount >= 1000 {
		return "known"
	}
	return "unknown"
}

func completenessExpenses(p profile.FinancialProfile) string {
	if p.SkippedExpenses {
		return "skipped"
	}
	for _, e := range p.FixedExpenses {
		if e.Amount > 0 {
			return "known"
		}
	}
	return "unknown"
}

func completenessCredits(d creditstore.Dashboard) string {
	if len(d.Credits) > 0 {
		return "known"
	}
	return "unknown"
}

func ProfileIncome(p profile.FinancialProfile) float64 {
	if p.SkippedIncome {
		return 0
	}
	return p.ActiveIncome + p.PassiveIncome
}

func ProfileExpenses(p profile.FinancialProfile) float64 {
	if p.SkippedExpenses {
		return 0
	}
	var sum float64
	for _, e := range p.FixedExpenses {
		sum += e.Amount
	}
	return sum
}

func RunwayMonths(p profile.FinancialProfile) *float64 {
	if p.SkippedCushion || p.EmergencyFund <= 0 {
		return nil
	}
	exp := ProfileExpenses(p)
	if exp <= 0 {
		return nil
	}
	m := p.EmergencyFund / exp
	return &m
}

func FreeCashflow(p profile.FinancialProfile) float64 {
	return ProfileIncome(p) - ProfileExpenses(p)
}
