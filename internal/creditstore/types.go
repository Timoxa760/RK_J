package creditstore

import "time"

type Credit struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Bank           string    `json:"bank"`
	Amount         float64   `json:"amount"`
	Rate           float64   `json:"rate"`
	TermMonths     int       `json:"term_months"`
	Remaining      float64   `json:"remaining"`
	MonthlyPayment float64   `json:"monthly_payment"`
	BenchmarkRate  float64   `json:"benchmark_rate,omitempty"`
	RateVsMarket   string    `json:"rate_vs_market,omitempty"`
	NextPayment    string    `json:"next_payment,omitempty"`
	ScannedAt      time.Time `json:"scanned_at"`
	SourceFileHash string    `json:"source_file_hash,omitempty"`
}

type Dashboard struct {
	DTI              float64  `json:"dti"`
	StressTestMonths float64  `json:"stress_test_months,omitempty"`
	Savings          float64  `json:"savings"`
	TotalDebt        float64  `json:"total_debt"`
	MonthlyPayments  float64  `json:"monthly_payments"`
	MonthlyIncome    float64  `json:"monthly_income"`
	Credits          []Credit `json:"credits"`
}
