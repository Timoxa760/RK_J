package mortgage

import (
	"fmt"
	"math"
	"sort"
)

type ApprovalLevel string

const (
	ApprovalHigh   ApprovalLevel = "high"
	ApprovalMedium ApprovalLevel = "medium"
	ApprovalLow    ApprovalLevel = "low"
)

type BankOffer struct {
	ID               string  `json:"id"`
	Bank             string  `json:"bank"`
	Rate             float64 `json:"rate"`
	MonthlyPayment   float64 `json:"monthly_payment"`
	TotalOverpayment float64 `json:"total_overpayment"`
	TermMonths       int     `json:"term_months"`
}

type Breakdown struct {
	ApprovalLevel      ApprovalLevel `json:"approval_level"`
	ApprovalReason     string        `json:"approval_reason"`
	SafeMortgageAmount float64       `json:"safe_mortgage_amount"`
	ComfortablePayment float64       `json:"comfortable_payment"`
	LoadRisk           string        `json:"load_risk"`
	ScenarioNow        string        `json:"scenario_now"`
	ScenarioWait       string        `json:"scenario_wait"`
	WaitMonths         int           `json:"wait_months"`
	Banks              []BankOffer   `json:"banks"`
	OptimalBankID      string        `json:"optimal_bank_id"`
}

type AnalyzeInput struct {
	MortgageAmount    float64
	MonthlyIncome     float64
	Savings           float64
	ExistingDTI       float64
	StressTestMonths  float64
	BenchmarkRate     float64
	ExistingPayments  float64
}

func approvalFromDTI(dti float64) ApprovalLevel {
	switch {
	case dti < 30:
		return ApprovalHigh
	case dti < 45:
		return ApprovalMedium
	default:
		return ApprovalLow
	}
}

func annuityPayment(principal, annualRate float64, termMonths int) float64 {
	if principal <= 0 || termMonths <= 0 {
		return 0
	}
	r := annualRate / 100 / 12
	if r <= 0 {
		return principal / float64(termMonths)
	}
	pow := math.Pow(1+r, float64(termMonths))
	return principal * r * pow / (pow - 1)
}

func totalOverpayment(principal, monthly float64, termMonths int) float64 {
	return math.Max(0, monthly*float64(termMonths)-principal)
}

func fmtRub(v float64) string {
	return fmt.Sprintf("%.0f", math.Round(v))
}

// BuildBreakdown — ипотечный разбор на profile + credits dashboard.
func BuildBreakdown(in AnalyzeInput) Breakdown {
	amount := in.MortgageAmount
	if amount <= 0 {
		amount = 12_000_000
	}
	income := in.MonthlyIncome
	if income <= 0 {
		income = 180_000
	}
	savings := in.Savings
	dti := in.ExistingDTI
	if dti < 0 {
		dti = 0
	}

	termMonths := 240
	baseRate := in.BenchmarkRate
	if baseRate <= 0 {
		baseRate = 18.0
	}

	bankDefs := []struct {
		id   string
		name string
		adj  float64
	}{
		{"tinkoff", "Т-Банк", -0.3},
		{"alfa", "Альфа-Банк", 0.0},
		{"sber", "Сбер", 0.2},
		{"vtb", "ВТБ", 0.5},
		{"gazprom", "Газпромбанк", 0.4},
	}

	banks := make([]BankOffer, 0, len(bankDefs))
	for _, b := range bankDefs {
		rate := baseRate + b.adj
		payment := annuityPayment(amount, rate, termMonths)
		banks = append(banks, BankOffer{
			ID:               b.id,
			Bank:             b.name,
			Rate:             math.Round(rate*10) / 10,
			MonthlyPayment:   math.Round(payment),
			TotalOverpayment: math.Round(totalOverpayment(amount, payment, termMonths)),
			TermMonths:       termMonths,
		})
	}

	sort.Slice(banks, func(i, j int) bool {
		return banks[i].MonthlyPayment < banks[j].MonthlyPayment
	})
	optimal := banks[0]

	comfortable := math.Min(income*0.35, optimal.MonthlyPayment)
	if comfortable <= 0 {
		comfortable = optimal.MonthlyPayment
	}

	projectedPayments := in.ExistingPayments + optimal.MonthlyPayment
	newDTI := dti
	if income > 0 {
		newDTI = projectedPayments / income * 100
	}
	level := approvalFromDTI(newDTI)

	reasons := map[ApprovalLevel]string{
		ApprovalHigh: fmt.Sprintf("Доход %s ₽/мес, на кредиты уйдёт ~%.0f%% — шанс одобрения выглядит хорошим.", fmtRub(income), newDTI),
		ApprovalMedium: fmt.Sprintf("На кредиты уйдёт ~%.0f%% дохода, запас %s ₽ — банк может попросить подтверждения.", newDTI, fmtRub(savings)),
		ApprovalLow: fmt.Sprintf("На кредиты уйдёт ~%.0f%% — ипотека на %s ₽ выглядит рискованной.", newDTI, fmtRub(amount)),
	}

	stressAfter := float64(0)
	if optimal.MonthlyPayment > 0 && savings > 0 {
		stressAfter = savings / optimal.MonthlyPayment
	}
	loadRisk := "После платежа запас станет маленьким — лучше подумать ещё раз."
	if stressAfter >= 4 {
		loadRisk = fmt.Sprintf("После платежа запаса хватит ~%.1f мес. — тесновато, но терпимо.", stressAfter)
	} else if stressAfter >= 2 {
		loadRisk = fmt.Sprintf("После платежа подушка сократится до ~%.1f мес. — ситуация станет более хрупкой.", stressAfter)
	}

	waitMonths := 8
	if newDTI > 45 {
		waitMonths = 12
	} else if newDTI < 30 {
		waitMonths = 4
	}

	return Breakdown{
		ApprovalLevel:      level,
		ApprovalReason:     reasons[level],
		SafeMortgageAmount: math.Round(income * 50),
		ComfortablePayment: math.Round(comfortable),
		LoadRisk:           loadRisk,
		ScenarioNow:        fmt.Sprintf("Платёж ~%s ₽/мес — до цели будете идти медленнее.", fmtRub(optimal.MonthlyPayment)),
		ScenarioWait:       fmt.Sprintf("Через %d мес. при накоплении ~%s ₽ условия могут стать мягче.", waitMonths, fmtRub(optimal.MonthlyPayment*3)),
		WaitMonths:         waitMonths,
		Banks:              banks,
		OptimalBankID:      optimal.ID,
	}
}
