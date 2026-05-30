package parser

import (
	"regexp"
	"strings"

	"backend_project/internal/rublang"
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
	{"купил", "Прочие расходы"},
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
	for _, cw := range categoryWords {
		if strings.Contains(text, cw.keyword) {
			category = cw.cat
			break
		}
	}

	store := detectStore(text)
	desc := shortDescription(store, category)

	return &ParsedExpense{
		Amount:      amount,
		Category:    category,
		Description: desc,
		Store:       store,
	}
}

func shortDescription(store, category string) string {
	if store != "" {
		return store
	}
	if category != "" && category != "Прочие расходы" {
		return category
	}
	return "Покупка"
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
