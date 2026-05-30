package parser

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"backend_project/internal/rublang"
	"backend_project/services/money-intelligence/ai-processor/internal/categorizer"
)

type ParsedExpense struct {
	Amount      float64
	Category    string
	Description string
	Store       string
}

var nonLettersRe = regexp.MustCompile(`[^a-zа-яё0-9]+`)

var storePatterns = []struct {
	match func(string) bool
	name  string
}{
	{
		match: func(n string) bool {
			return strings.Contains(n, "пятероч") ||
				strings.Contains(n, "терочк") ||
				strings.Contains(n, "spiteroc") ||
				strings.Contains(n, "piteroc")
		},
		name: "Пятёрочка",
	},
	{match: func(n string) bool { return strings.Contains(n, "перекрест") }, name: "Перекрёсток"},
	{match: func(n string) bool { return strings.Contains(n, "магнит") }, name: "Магнит"},
	{match: func(n string) bool { return strings.Contains(n, "lenta") || strings.Contains(n, "лента") }, name: "Лента"},
	{match: func(n string) bool { return strings.Contains(n, "ашан") || strings.Contains(n, "auchan") }, name: "Ашан"},
	{match: func(n string) bool { return strings.Contains(n, "dixi") || strings.Contains(n, "diksi") || strings.Contains(n, "дикси") }, name: "Дикси"},
}

var categoryWords = []struct {
	keyword string
	cat     string
}{
	{"пятёроч", "Продукты"},
	{"пятероч", "Продукты"},
	{"терочк", "Продукты"},
	{"перекрест", "Продукты"},
	{"магнит", "Продукты"},
	{"лента", "Продукты"},
	{"ашан", "Продукты"},
	{"дикси", "Продукты"},
	{"продукт", "Продукты"},
	{"еда", "Продукты"},
	{"продукты", "Продукты"},
	{"кофе", "Кафе и рестораны"},
	{"кафе", "Кафе и рестораны"},
	{"ресторан", "Кафе и рестораны"},
	{"такси", "Транспорт"},
	{"бензин", "Транспорт"},
	{"транспорт", "Транспорт"},
}

var storeWords = []struct {
	keyword string
	name    string
}{
	{"пятёроч", "Пятёрочка"},
	{"пятероч", "Пятёрочка"},
	{"перекрест", "Перекрёсток"},
	{"магнит", "Магнит"},
	{"лента", "Лента"},
	{"ашан", "Ашан"},
	{"дикси", "Дикси"},
}

var stripAmountRe = regexp.MustCompile(`(?i)\d+(?:[\s.,]\d*)?\s*(?:₽|руб(?:лей|ля|\.?)?|р(?:\.|\s|$)|тыс(?:\.|яч(?:и|ей|а)?|и|ч)?|тыщ(?:а|и|ей|у)?|k|к\b)`)

// normalizeFuzzy приводит текст к виду для нечёткого поиска магазинов (Whisper-опечатки).
func normalizeFuzzy(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "ё", "е")
	text = nonLettersRe.ReplaceAllString(text, "")
	return text
}

// Parse извлекает сумму, категорию и короткое название из текста (regex fallback для voice/manual).
func Parse(rawText string) *ParsedExpense {
	text := rublang.Normalize(rawText)
	if text == "" {
		return nil
	}

	amount, ok := rublang.ExtractPrimary(rawText)
	if !ok || amount == 0 {
		return nil
	}

	category := "Прочие расходы"
	if cat := categorizer.CategoryForText(text); cat != "Прочее" {
		category = cat
	} else {
		for _, cw := range categoryWords {
			if strings.Contains(text, cw.keyword) {
				category = cw.cat
				break
			}
		}
	}

	store := detectStore(text)
	desc := shortDescription(store, category, text, amount)

	return &ParsedExpense{
		Amount:      amount,
		Category:    category,
		Description: desc,
		Store:       store,
	}
}

func shortDescription(store, category, rawText string, amount float64) string {
	if store != "" {
		return store
	}
	if label := categorizer.ProductLabelFromText(rawText); label != "" {
		return label
	}
	if label := extractSpendLabel(rawText, amount); label != "" {
		return label
	}
	if category != "" && category != "Прочие расходы" {
		return category
	}
	return "Покупка"
}

func extractSpendLabel(rawText string, amount float64) string {
	_ = amount
	text := rublang.Normalize(rawText)
	text = stripAmountRe.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	for _, prefix := range []string{"купил", "купила", "купили", "потратил", "потратила", "заплатил", "заплатила", "оплатил", "оплатила", "на"} {
		text = strings.TrimPrefix(text, prefix+" ")
	}
	text = strings.Trim(text, " ,.-—")
	if text == "" {
		return ""
	}
	return titlePhrase(text)
}

func titlePhrase(text string) string {
	words := strings.Fields(text)
	for i, w := range words {
		words[i] = titleToken(w)
	}
	return strings.Join(words, " ")
}

func titleToken(word string) string {
	if word == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(word)
	return string(unicode.ToUpper(r)) + word[size:]
}

func detectStore(text string) string {
	n := normalizeFuzzy(text)
	for _, sp := range storePatterns {
		if sp.match(n) {
			return sp.name
		}
	}
	for _, sw := range storeWords {
		if strings.Contains(text, sw.keyword) {
			return sw.name
		}
	}
	return ""
}
