package advisor

import "strings"

var standardCategories = []string{
	"Продукты",
	"Кафе и рестораны",
	"Транспорт",
	"Доставка",
	"Подписки",
	"ЖКХ",
	"Развлечения",
	"Одежда",
	"Здоровье",
	"Прочие расходы",
}

var aliasToStandard = map[string]string{
	"прочее":           "Прочие расходы",
	"прочие расходы":   "Прочие расходы",
	"прочие":           "Прочие расходы",
	"связь":            "Прочие расходы",
	"продукты":         "Продукты",
	"кафе и рестораны": "Кафе и рестораны",
	"кафе":             "Кафе и рестораны",
	"рестораны":        "Кафе и рестораны",
	"транспорт":        "Транспорт",
	"доставка":         "Доставка",
	"подписки":         "Подписки",
	"подписка":         "Подписки",
	"жкх":              "ЖКХ",
	"коммунальные":     "ЖКХ",
	"развлечения":      "Развлечения",
	"одежда":           "Одежда",
	"здоровье":         "Здоровье",
	"кредиты":          "Прочие расходы",
	"доход":            "Прочие расходы",
}

func NormalizeCategory(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "Прочие расходы"
	}
	lower := strings.ToLower(trimmed)
	if std, ok := aliasToStandard[lower]; ok {
		return std
	}
	for _, standard := range standardCategories {
		stdLower := strings.ToLower(standard)
		if lower == stdLower || strings.Contains(lower, stdLower) || strings.Contains(stdLower, lower) {
			return standard
		}
	}
	return "Прочие расходы"
}
