package creditparse

import (
	"regexp"

	"backend_project/internal/rublang"
)

var transferAmountPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(?:передает|выдает|выдаёт|предоставляет)` + rxWord + `?\s+(?:заемщику\s+)?(\d[\d\s.,]*)\s*(?:₽|руб` + rxWord + `|р\b|\()`),
	regexp.MustCompile(`(?i)(?:передает|выдает|выдаёт|предоставляет)` + rxWord + `?\s+(?:заемщику\s+)?(\d[\d\s.,]*)(?:\s*(?:₽|руб` + rxWord + `))`),
}

func pickLoanAmount(text, normalized string) float64 {
	if v := firstAmount(amountPatterns, normalized); v > 0 {
		return v
	}
	if v := firstAmount(transferAmountPatterns, normalized); v > 0 {
		return v
	}

	best := 0
	bestScore := -1
	for _, candidate := range rublang.ExtractAll(text) {
		score := scoreLoanAmount(candidate)
		if score < 0 {
			continue
		}
		if score > bestScore || (score == bestScore && candidate > best) {
			best = candidate
			bestScore = score
		}
	}
	return float64(best)
}

func scoreLoanAmount(amount int) int {
	if amount < 1_000 || amount > 30_000_000 {
		return -1
	}
	score := 0
	switch {
	case amount >= 10_000 && amount <= 5_000_000:
		score += 3
	case amount >= 1_000 && amount < 10_000:
		score += 1
	}
	if amount == 1_000 {
		score -= 2
	}
	return score
}

func inferTermMonths(normalized string, amount, rate, payment float64) (int, bool) {
	if n, ok := firstInt(termMonthPatterns, normalized); ok {
		return n, true
	}
	if years, ok := firstInt(termYearPatterns, normalized); ok && years <= 50 {
		return years * 12, true
	}
	if n, ok := termMonthsFromDates(normalized); ok {
		return n, true
	}
	if days, ok := firstInt(termDayPatterns, normalized); ok && days <= 3650 {
		months := days / 30
		if months < 1 {
			months = 1
		}
		return months, true
	}
	if payment > 0 && amount > 0 && rate > 0 {
		if n, ok := termMonthsFromPayment(amount, rate, payment); ok {
			return n, true
		}
	}
	return 0, false
}

func termMonthsFromPayment(amount, rate, payment float64) (int, bool) {
	monthlyRate := rate / 100 / 12
	if payment <= 0 || amount <= 0 || monthlyRate <= 0 {
		return 0, false
	}
	// payment = amount * r * (1+r)^n / ((1+r)^n - 1)
	// Решаем перебором — срок в договорах обычно до 50 лет.
	for months := 1; months <= 600; months++ {
		pow := pow1p(monthlyRate, months)
		expected := amount * monthlyRate * pow / (pow - 1)
		if nearlyEqual(expected, payment, 0.03) {
			return months, true
		}
	}
	return 0, false
}

func pow1p(base float64, exp int) float64 {
	out := 1.0
	for i := 0; i < exp; i++ {
		out *= 1 + base
	}
	return out
}

func nearlyEqual(a, b, relTol float64) bool {
	if a == 0 || b == 0 {
		return false
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	base := a
	if b > base {
		base = b
	}
	return diff/base <= relTol
}
