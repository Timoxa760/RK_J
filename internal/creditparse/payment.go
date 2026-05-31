package creditparse

import "math"

// EstimateMonthlyPayment — аннуитетный платёж по сумме, годовой ставке и сроку в месяцах.
func EstimateMonthlyPayment(amount, rate float64, termMonths int) float64 {
	if amount <= 0 || termMonths <= 0 {
		return 0
	}
	if rate <= 0 {
		return math.Round(amount / float64(termMonths))
	}
	r := rate / 100 / 12
	pow := math.Pow(1+r, float64(termMonths))
	if pow <= 1 {
		return math.Round(amount / float64(termMonths))
	}
	payment := amount * r * pow / (pow - 1)
	return math.Round(payment)
}

func ensureMonthlyPayment(amount, rate float64, termMonths int, payment float64) (float64, bool) {
	if payment > 0 {
		return payment, false
	}
	estimated := EstimateMonthlyPayment(amount, rate, termMonths)
	if estimated <= 0 {
		return 0, false
	}
	return estimated, true
}
