package parser

import (
	"regexp"
	"strconv"
	"strings"
)

type ParsedExpense struct {
	Amount      float64
	Category    string
	Description string
}

var (
	amountRe      = regexp.MustCompile(`(\d+(?:[\s.,]?\d*)?)\s*(?:тыщ|тысяч|тысячи|руб|₽|рублей|рубля)?`)
	amountOnlyRe   = regexp.MustCompile(`^\d+$`)
	amountCleanRe = regexp.MustCompile(`[\s.,]`)
)

var categoryWords = []struct {
	keyword string
	cat     string
}{
	{"зарплата", "Доход"},
	{"аванс", "Доход"},
	{"зп", "Доход"},
	{"получил", "Доход"},
	{"продукт", "Продукты"},
	{"еда", "Продукты"},
	{"продукты", "Продукты"},
	{"кофе", "Кафе и рестораны"},
	{"кафе", "Кафе и рестораны"},
	{"ресторан", "Кафе и рестораны"},
	{"столовая", "Кафе и рестораны"},
	{"доставка", "Доставка"},
	{"такси", "Транспорт"},
	{"бензин", "Транспорт"},
	{"транспорт", "Транспорт"},
	{"метро", "Транспорт"},
	{"автобус", "Транспорт"},
	{"проезд", "Транспорт"},
	{"подписк", "Подписки"},
	{"инет", "Подписки"},
	{"интернет", "Подписки"},
	{"связь", "Подписки"},
	{"телефон", "Подписки"},
	{"коммуналк", "ЖКХ"},
	{"жкх", "ЖКХ"},
	{"квартир", "ЖКХ"},
	{"электричеств", "ЖКХ"},
	{"развлечен", "Развлечения"},
	{"кино", "Развлечения"},
	{"игр", "Развлечения"},
	{"кредит", "Кредиты"},
	{"ипотек", "Кредиты"},
	{"долг", "Кредиты"},
	{"здоровь", "Здоровье"},
	{"аптек", "Здоровье"},
	{"врач", "Здоровье"},
	{"лекарств", "Здоровье"},
	{"спорт", "Здоровье"},
	{"одежд", "Одежда"},
	{"обувь", "Одежда"},
	{"проебал", "Прочие расходы"},
	{"потерял", "Прочие расходы"},
	{"сломал", "Прочие расходы"},
	{"ремонт", "Прочие расходы"},
	{"подарк", "Прочие расходы"},
	{"купил", "Прочие расходы"},
	{"прочее", "Прочие расходы"},
}

func Parse(rawText string) *ParsedExpense {
	text := strings.ToLower(strings.TrimSpace(rawText))
	if text == "" {
		return nil
	}

	if amountOnlyRe.MatchString(text) {
		amt, _ := strconv.ParseFloat(text, 64)
		return &ParsedExpense{Amount: amt, Category: "Прочие расходы", Description: rawText}
	}

	m := amountRe.FindStringSubmatch(text)
	if m == nil {
		return nil
	}

	rawAmt := strings.TrimSpace(m[1])
	rawAmt = amountCleanRe.ReplaceAllString(rawAmt, "")
	amount, _ := strconv.ParseFloat(rawAmt, 64)
	if amount == 0 {
		return nil
	}

	category := "Прочие расходы"
	for _, cw := range categoryWords {
		if strings.Contains(text, cw.keyword) {
			category = cw.cat
			break
		}
	}

	desc := strings.TrimSpace(rawText)

	return &ParsedExpense{
		Amount:      amount,
		Category:    category,
		Description: desc,
	}
}
