package creditparse

import (
	"regexp"
	"strconv"
	"strings"

	"backend_project/internal/rublang"
)

// rxWord — суффикс русского/латинского слова (Go \w не включает кириллицу).
const rxWord = `[а-яёa-z0-9_]*`

var (
	ratePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:ставк` + rxWord + `|процент` + rxWord + `|полн` + rxWord + `\s+стоимост` + rxWord + `|psk|пск)[^\d]{0,32}(\d+[.,]?\d*)\s*%`),
		regexp.MustCompile(`(?i)(\d+[.,]?\d*)\s*%\s*(?:годовых|г\.?\s*г\.?|p\.?\s*a\.?|в\s+год)`),
		regexp.MustCompile(`(?i)(\d+[.,]\d+)\s*%`),
	}
	termMonthPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:срок` + rxWord + `)[^\d]{0,32}(\d+)\s*(?:мес(?:\.|яцев|яца|)?|month)`),
		regexp.MustCompile(`(?i)(?:срок` + rxWord + `|на)\s*(\d+)\s*(?:мес(?:\.|яцев|яца|)?|month)`),
		regexp.MustCompile(`(?i)(\d+)\s*(?:мес(?:\.|яцев|яца|)?)\s*(?:кредит|займ|за[её]м|договор|ссуд)`),
		regexp.MustCompile(`(?i)(?:срок` + rxWord + `\s+(?:кредит|займ|за[её]м|договор|ссуд))[^\d]{0,16}(\d+)\s*(?:мес(?:\.|яцев|яца|)?|month)`),
	}
	termYearPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:срок` + rxWord + `)[^\d]{0,32}(\d+)\s*(?:лет|года|год(?:\s|$|[^а-яё]))`),
	}
	termDayPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:срок` + rxWord + `)[^\d]{0,32}(\d+)\s*(?:дн(?:\.|ей|я|)?|day)`),
		regexp.MustCompile(`(?i)(?:срок` + rxWord + `|на)\s*(\d+)\s*(?:дн(?:\.|ей|я|)?|day)`),
	}
	paymentPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:ежемесяч` + rxWord + `|аннуитет` + rxWord + `|плат[её]ж` + rxWord + `|раз\s+в\s+(?:месяц|мес))[^\d]{0,32}(\d[\d\s.,]*)`),
		regexp.MustCompile(`(?i)(?:сумм` + rxWord + `\s+плат[её]ж` + rxWord + `)[^\d]{0,16}(\d[\d\s.,]*)`),
	}
	amountPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:сумм` + rxWord + `\s+(?:кредит|займ|за[её]м|ссуд|наличн)` + rxWord + `|размер` + rxWord + `\s+(?:кредит|займ|за[её]м|ссуд)` + rxWord + `|лимит` + rxWord + `\s+(?:кредит|займ|за[её]м)` + rxWord + `|основн` + rxWord + `\s+долг` + rxWord + `)[^\d]{0,24}(\d[\d\s.,]*)`),
	}
)

var bankHints = []struct {
	key  string
	name string
}{
	{"т-банк", "Т-Банк"},
	{"«тбанк»", "Т-Банк"},
	{"tinkoff", "Т-Банк"},
	{"тинькофф", "Т-Банк"},
	{"сбер", "Сбербанк"},
	{"sber", "Сбербанк"},
	{"втб", "ВТБ"},
	{"vtb", "ВТБ"},
	{"альфа", "Альфа-Банк"},
	{"alfa", "Альфа-Банк"},
	{"райфф", "Райффайзен"},
	{"газпром", "Газпромбанк"},
	{"росбанк", "Росбанк"},
	{"открытие", "Открытие"},
	{"совком", "Совкомбанк"},
	{"почта банк", "Почта Банк"},
	{"мкб", "МКБ"},
	{"займер", "Займер"},
	{"money man", "MoneyMan"},
	{"moneyman", "MoneyMan"},
	{"webbankir", "Webbankir"},
	{"веббанкир", "Webbankir"},
	{"lime", "Lime"},
	{"лайм", "Lime"},
	{"home credit", "Home Credit"},
	{"хоум кредит", "Home Credit"},
}

// Fields — распознанные поля договора.
type Fields struct {
	Amount           float64
	Rate             float64
	TermMonths       int
	MonthlyPayment   float64
	Bank             string
	PaymentEstimated bool
}

// ParseFromText извлекает условия кредита из текста PDF (regex fallback без LLM).
func ParseFromText(text string) (Fields, bool) {
	normalized := rublang.Normalize(text)
	if len(normalized) < 20 {
		return Fields{}, false
	}

	rate, okRate := firstFloat(ratePatterns, normalized)
	if !okRate {
		return Fields{}, false
	}

	payment, _ := firstAmountFloat(paymentPatterns, normalized)
	amount := pickLoanAmount(text, normalized)
	if amount <= 0 {
		return Fields{}, false
	}

	term, okTerm := inferTermMonths(normalized, amount, rate, payment)
	if !okTerm {
		return Fields{}, false
	}

	bank := detectProductLabel(normalized)
	payment, paymentEstimated := ensureMonthlyPayment(amount, rate, term, payment)
	return Fields{
		Amount:           amount,
		Rate:             rate,
		TermMonths:       term,
		MonthlyPayment:   payment,
		Bank:             bank,
		PaymentEstimated: paymentEstimated,
	}, true
}

func detectBankName(text string) string {
	for _, hint := range bankHints {
		if strings.Contains(text, hint.key) {
			return hint.name
		}
	}
	return ""
}

func detectProductLabel(text string) string {
	if name := detectBankName(text); name != "" {
		return name
	}
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "микрозайм"), strings.Contains(lower, "микрофинанс"), strings.Contains(lower, " мфо"):
		return "Займ"
	case strings.Contains(lower, "договор займа"), strings.Contains(lower, "договор микрозайма"), strings.Contains(lower, " займ "), strings.Contains(lower, " займа "):
		return "Займ"
	case strings.Contains(lower, "наличн"):
		return "Кредит наличными"
	default:
		return "Кредит"
	}
}

func firstFloat(patterns []*regexp.Regexp, text string) (float64, bool) {
	for _, re := range patterns {
		m := re.FindStringSubmatch(text)
		if len(m) < 2 {
			continue
		}
		if v, ok := parseNumber(m[1]); ok && v > 0 && v <= 876 {
			return v, true
		}
	}
	return 0, false
}

func firstInt(patterns []*regexp.Regexp, text string) (int, bool) {
	for _, re := range patterns {
		m := re.FindStringSubmatch(text)
		if len(m) < 2 {
			continue
		}
		raw := strings.ReplaceAll(strings.TrimSpace(m[1]), " ", "")
		n, err := strconv.Atoi(raw)
		if err == nil && n > 0 && n <= 600 {
			return n, true
		}
	}
	return 0, false
}

func firstTermMonths(monthPatterns, yearPatterns, dayPatterns []*regexp.Regexp, text string) (int, bool) {
	if n, ok := firstInt(monthPatterns, text); ok {
		return n, true
	}
	if years, ok := firstInt(yearPatterns, text); ok && years <= 50 {
		return years * 12, true
	}
	if days, ok := firstInt(dayPatterns, text); ok && days <= 3650 {
		months := days / 30
		if months < 1 {
			months = 1
		}
		return months, true
	}
	return 0, false
}

func firstAmount(patterns []*regexp.Regexp, text string) float64 {
	v, ok := firstAmountFloat(patterns, text)
	if ok {
		return v
	}
	return 0
}

func firstAmountFloat(patterns []*regexp.Regexp, text string) (float64, bool) {
	for _, re := range patterns {
		m := re.FindStringSubmatch(text)
		if len(m) < 2 {
			continue
		}
		if v, ok := parseNumber(m[1]); ok && v >= 100 {
			return v, true
		}
	}
	return 0, false
}

func parseNumber(raw string) (float64, bool) {
	raw = strings.TrimSpace(raw)
	raw = strings.ReplaceAll(raw, " ", "")
	raw = strings.ReplaceAll(raw, ",", ".")
	if raw == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(raw, 64)
	return v, err == nil && v > 0
}

func maxInt(values []int) int {
	max := 0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
