package creditparse

import "testing"

func TestEstimateMonthlyPayment_TBankSample(t *testing.T) {
	got := EstimateMonthlyPayment(500_000, 18, 12)
	if got < 45_000 || got > 47_000 {
		t.Fatalf("payment=%v want ~45800", got)
	}
}

func TestParseFromText_TBankBusinessLoan_EstimatesPayment(t *testing.T) {
	text := `
	Кредитный договор № 4
	09.09.2025
	АО «ТБанк» (далее — Кредитор)
	1.1. Кредитор до 13.09.2025 передает Заемщику 500 000 ₽ (далее — Кредит)
	1.2. Заемщик выплачивает проценты за пользование Кредитом — 18% годовых.
	3.1. Заемщик возвращает кредит до 12.09.2026 Кредитору.
	`
	got, ok := ParseFromText(text)
	if !ok {
		t.Fatal("expected parse ok")
	}
	if got.MonthlyPayment <= 0 {
		t.Fatalf("expected estimated payment, got %+v", got)
	}
	if !got.PaymentEstimated {
		t.Fatalf("expected payment_estimated=true")
	}
}
