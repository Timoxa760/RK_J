package creditparse

import (
	"regexp"
	"strings"
)

var (
	reGeneralTitle    = regexp.MustCompile(`(?i)общие\s+условия\s+договора`)
	reIndividualTitle = regexp.MustCompile(`(?i)индивидуальные\s+условия\s+договора`)
	reGeneralFilename = regexp.MustCompile(`(?i)(general[_\-\s]?conditions|общие[_\-\s]?услов)`)
	reGeneralProduct  = regexp.MustCompile(`(?i)общие\s+условия\s+(?:договора\s+)?(?:потребительского\s+)?(?:кредит|займ|за[её]м|микрозайм|микрофинанс)`)
)

// IsGeneralConditionsTemplate определяет типовой документ банка (общие условия продукта),
// в котором нет суммы, ставки и срока конкретного кредита клиента.
func IsGeneralConditionsTemplate(text, filename string) bool {
	if reGeneralFilename.MatchString(strings.TrimSpace(filename)) {
		return true
	}
	prep := strings.TrimSpace(text)
	if len(prep) < 40 {
		return false
	}

	intro := prep
	if len(intro) > 4000 {
		intro = intro[:4000]
	}
	lowerIntro := strings.ToLower(intro)

	if !strings.Contains(lowerIntro, "общие условия") {
		return false
	}

	if reGeneralProduct.MatchString(lowerIntro) {
		return true
	}

	lead := intro
	if len(lead) > 1200 {
		lead = lead[:1200]
	}
	if reIndividualTitle.MatchString(lead) {
		leadHead := lead
		if len(leadHead) > 500 {
			leadHead = leadHead[:500]
		}
		if !reGeneralTitle.MatchString(leadHead) {
			return false
		}
	}

	if reGeneralTitle.MatchString(lead) {
		return true
	}

	// Общие условия в шапке, а ставка/сумма отсылают к «индивидуальным условиям».
	lowerAll := strings.ToLower(prep)
	if strings.Contains(lowerAll, "указан") && strings.Contains(lowerAll, "индивидуальн") {
		return true
	}

	return strings.Contains(lowerIntro, "общие условия выдачи")
}
