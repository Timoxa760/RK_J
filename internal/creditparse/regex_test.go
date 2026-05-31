package creditparse

import (
	"testing"

	"backend_project/internal/rublang"
)

func TestParseFromText_ConsumerLoan(t *testing.T) {
	text := `
	Кредитный договор №123
	Банк: Т-Банк
	Сумма кредита 1 200 000 рублей
	Срок 36 месяцев
	Процентная ставка 14.5% годовых
	Ежемесячный платёж 42 000 руб
	`
	got, ok := ParseFromText(text)
	if !ok {
		t.Fatal("expected parse ok")
	}
	if got.Rate != 14.5 || got.TermMonths != 36 {
		t.Fatalf("rate/term: %+v", got)
	}
	if got.Amount != 1_200_000 {
		t.Fatalf("amount=%v", got.Amount)
	}
	if got.MonthlyPayment != 42_000 {
		t.Fatalf("payment=%v", got.MonthlyPayment)
	}
	if got.Bank != "Т-Банк" {
		t.Fatalf("bank=%q", got.Bank)
	}
}

func TestParseFromText_TermYears(t *testing.T) {
	text := `
	Сумма кредита 800 000 руб
	Срок кредита 5 лет
	Процентная ставка 13.5% годовых
	`
	got, ok := ParseFromText(text)
	if !ok {
		t.Fatal("expected parse ok")
	}
	if got.TermMonths != 60 {
		t.Fatalf("term_months=%v want 60", got.TermMonths)
	}
}

func TestParseFromText_CashLoan(t *testing.T) {
	text := `
	ИНДИВИДУАЛЬНЫЕ УСЛОВИЯ ДОГОВОРА ПОТРЕБИТЕЛЬСКОГО КРЕДИТА НАЛИЧНЫМИ
	Сумма кредита наличными 350 000 рублей
	Срок кредита 24 месяца
	Процентная ставка 18.9% годовых
	`
	got, ok := ParseFromText(text)
	if !ok {
		t.Fatal("expected parse ok")
	}
	if got.Amount != 350_000 || got.TermMonths != 24 || got.Rate != 18.9 {
		t.Fatalf("got %+v", got)
	}
	if got.Bank != "Кредит наличными" {
		t.Fatalf("bank=%q", got.Bank)
	}
}

func TestParseFromText_Microloan(t *testing.T) {
	text := `
	Договор займа №77
	Сумма займа 30 000 руб
	Срок займа 30 дней
	Полная стоимость займа 365% годовых
	`
	n := rublang.Normalize(text)
	rate, okR := firstFloat(ratePatterns, n)
	term, okT := firstTermMonths(termMonthPatterns, termYearPatterns, termDayPatterns, n)
	amt := firstAmount(amountPatterns, n)
	if !okR || !okT || amt <= 0 {
		t.Fatalf("rate=%v okR=%v term=%v okT=%v amt=%v norm=%q", rate, okR, term, okT, amt, n[:120])
	}
	got, ok := ParseFromText(text)
	if !ok {
		t.Fatal("expected parse ok")
	}
	if got.Amount != 30_000 {
		t.Fatalf("amount=%v", got.Amount)
	}
	if got.Rate != 365 {
		t.Fatalf("rate=%v", got.Rate)
	}
	if got.TermMonths != 1 {
		t.Fatalf("term_months=%v want 1", got.TermMonths)
	}
	if got.Bank != "Займ" {
		t.Fatalf("bank=%q", got.Bank)
	}
}

func TestTermMonthsFromDates_TBankSample(t *testing.T) {
	text := rublang.Normalize(`
		09.09.2025
		Кредитор до 13.09.2025 передает Заемщику 500 000 ₽
		процентная ставка 18% годовых
		Заемщик возвращает кредит до 12.09.2026
	`)
	months, ok := termMonthsFromDates(text)
	if !ok {
		t.Fatal("expected date-based term")
	}
	if months < 11 || months > 12 {
		t.Fatalf("term_months=%v want 11-12", months)
	}
}

func TestParseFromText_TBankBusinessLoan(t *testing.T) {
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
	if got.Amount != 500_000 || got.Rate != 18 {
		t.Fatalf("amount/rate: %+v", got)
	}
	if got.TermMonths < 11 || got.TermMonths > 12 {
		t.Fatalf("term_months=%v", got.TermMonths)
	}
	if got.Bank != "Т-Банк" {
		t.Fatalf("bank=%q", got.Bank)
	}
}
