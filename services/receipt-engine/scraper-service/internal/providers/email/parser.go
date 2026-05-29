package email

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseReceiptHTML(htmlContent string) *scrap.RawReceipt {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	receipt := &scrap.RawReceipt{
		Date: time.Now(),
	}

	parseTable(doc, receipt)
	parseDivs(doc, receipt)

	return receipt
}

func parseTable(doc *goquery.Document, r *scrap.RawReceipt) {
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(_ int, row *goquery.Selection) {
			cells := row.Find("td, th")
			if cells.Length() < 2 {
				return
			}

			key := strings.TrimSpace(cells.First().Text())
			val := strings.TrimSpace(cells.Last().Text())

			switch {
			case containsAny(key, "магазин", "store", "seller", "продавец"):
				r.Store = val

			case containsAny(key, "дата", "date", "время", "time"):
				if t, err := parseDate(val); err == nil {
					r.Date = t
				}

			case containsAny(key, "сумма", "total", "итог", "amount", "к оплате"):
				r.Total = parsePrice(val)

			case containsAny(key, "фискальный", "фн", "фд", "фп", "fn", "fd", "fp", "чек", "номер"):
			}
		})

		table.Find("tr").Each(func(_ int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() < 2 {
				return
			}

			name := strings.TrimSpace(cells.First().Text())
			if !isItemName(name) || isLabel(name) {
				return
			}

			item := scrap.RawItem{Name: name, Quantity: 1}

			if cells.Length() >= 3 {
				item.Price = parsePrice(strings.TrimSpace(cells.Eq(1).Text()))
				if qty := parseInt(strings.TrimSpace(cells.Eq(2).Text())); qty > 0 {
					item.Quantity = qty
				}
			} else {
				item.Price = parsePrice(strings.TrimSpace(cells.Eq(1).Text()))
			}

			r.Items = append(r.Items, item)
		})
	})
}

func parseDivs(doc *goquery.Document, r *scrap.RawReceipt) {
	doc.Find("div, span, p").Each(func(_ int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		lower := strings.ToLower(text)

		if r.Store == "" && (containsAny(lower, "магазин:", "пятёрочка", "магнит", "вкусвилл", "лента")) {
			for _, name := range []string{"пятёрочка", "магнит", "вкусвилл", "лента", "ozon", "wildberries"} {
				if strings.Contains(lower, name) {
					r.Store = capitalize(name)
					break
				}
			}
		}

		if r.Total == 0 && (containsAny(lower, "итог:", "сумма:", "total:", "к оплате:")) {
			r.Total = parsePrice(text)
		}
	})
}

func isItemName(s string) bool {
	if len(s) < 2 || len(s) > 100 {
		return false
	}
	skip := []string{"магазин", "дата", "сумма", "итог", "наименование", "цена", "кол-во",
		"количество", "товар", "название", "фискальный", "чек", "кассовый", "итого",
		"подытог", "скидка", "налог", "ндс", "к оплате", "всего", "store", "date",
		"total", "name", "price", "qty", "quantity", "item", "description", "amount"}
	lower := strings.ToLower(strings.TrimSpace(s))
	for _, skipWord := range skip {
		if lower == skipWord {
			return false
		}
	}
	return true
}

func isLabel(s string) bool {
	labels := []string{"магазин", "дата", "сумма", "итог", "наименование", "товар", "цена", "кол-во"}
	lower := strings.ToLower(strings.TrimSpace(s))
	for _, l := range labels {
		if lower == l {
			return true
		}
	}
	return false
}

func containsAny(s string, substrs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range substrs {
		if strings.Contains(lower, sub) {
			return true
		}
	}
	return false
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(parseDigits(s))
	if err != nil {
		return 0
	}
	return n
}

func parsePrice(s string) float64 {
	s = strings.NewReplacer(" ", "", "\u00a0", "", "₽", "", "руб", "", ",", ".").Replace(s)
	s = strings.TrimSpace(s)
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseDigits(s string) string {
	var d strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			d.WriteRune(r)
		}
	}
	return d.String()
}

func parseDate(s string) (time.Time, error) {
	formats := []string{
		"02.01.2006",
		"2006-01-02",
		"02/01/2006",
		"02.01.06",
		"2006/01/02",
	}
	s = strings.TrimSpace(s)
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}

func capitalize(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}
