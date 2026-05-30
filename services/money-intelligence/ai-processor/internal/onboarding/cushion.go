package onboarding

import (
	"regexp"
	"strconv"
	"strings"

	"backend_project/internal/rublang"
)

var segmentAmountRe = regexp.MustCompile(`(\d+(?:[\s.,]\d*)?)`)

// EmergencyBreakdown — разбивка «запаса» по местам хранения.
type EmergencyBreakdown struct {
	Cash        float64 `json:"cash"`
	Deposit     float64 `json:"deposit"`
	Investments float64 `json:"investments"`
}

func (b EmergencyBreakdown) hasAny() bool {
	return b.Cash > 0 || b.Deposit > 0 || b.Investments > 0
}

func (b EmergencyBreakdown) total() float64 {
	return b.Cash + b.Deposit + b.Investments
}

var (
	cashKeywords = []string{
		"налич", "налик", "наличк", "наличкой", "кэш", "cash", "в кармане", "на руках", "дома",
	}
	depositKeywords = []string{
		"вклад", "вкладах", "депозит", "счет", "счёт", "банк", "накопит", "сбер",
	}
	investKeywords = []string{
		"инвест", "акци", "облига", "брокер", "etf", "фонд", "крипт", "валют",
	}
)

// parseEmergencyBreakdown извлекает суммы по сегментам фразы («наличными 50k, на вкладе 15k»).
func parseEmergencyBreakdown(text string) EmergencyBreakdown {
	var out EmergencyBreakdown
	for _, seg := range splitVoiceSegments(text) {
		amount := amountFromSegment(seg)
		if amount <= 0 {
			continue
		}
		n := rublang.Normalize(seg)
		switch {
		case containsAny(n, cashKeywords):
			out.Cash += amount
		case containsAny(n, depositKeywords):
			out.Deposit += amount
		case containsAny(n, investKeywords):
			out.Investments += amount
		}
	}
	return out
}

func splitVoiceSegments(text string) []string {
	n := rublang.Normalize(text)
	n = strings.NewReplacer(";", ",", " и ", ",").Replace(n)
	parts := strings.Split(n, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) <= 1 {
		return []string{n}
	}
	return out
}

func containsAny(text string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func amountFromSegment(seg string) float64 {
	nums := rublang.ExtractAll(seg)
	if len(nums) > 0 {
		return float64(nums[0])
	}
	n := rublang.Normalize(seg)
	m := segmentAmountRe.FindStringSubmatch(n)
	if len(m) < 2 {
		return 0
	}
	raw := strings.ReplaceAll(strings.ReplaceAll(m[1], " ", ""), ",", ".")
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v <= 0 {
		return 0
	}
	return v
}
