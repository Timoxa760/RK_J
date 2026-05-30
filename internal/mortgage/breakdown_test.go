package mortgage

import "testing"

func TestBuildBreakdown_returnsBanksAndApproval(t *testing.T) {
	b := BuildBreakdown(AnalyzeInput{
		MortgageAmount:   12_000_000,
		MonthlyIncome:    180_000,
		Savings:          340_000,
		ExistingDTI:      28,
		BenchmarkRate:    18,
		ExistingPayments: 50_000,
	})
	if len(b.Banks) < 3 {
		t.Fatalf("expected banks, got %d", len(b.Banks))
	}
	if b.ApprovalLevel == "" || b.OptimalBankID == "" {
		t.Fatalf("incomplete breakdown: %+v", b)
	}
	if b.ComfortablePayment <= 0 {
		t.Fatalf("expected comfortable payment")
	}
}
