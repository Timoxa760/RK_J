package creditparse

import (
	"regexp"
	"time"
)

var (
	reDateDMY = `(\d{2})\.(\d{2})\.(\d{4})`
	reCreditUntilSpaced = regexp.MustCompile(`(?i)(?:возвращает|возврат|погаш|уплат` + rxWord + `)\s+(?:кредит|займ|за[её]m)` + rxWord + `?\s+до\s+` + reDateDMY)
	reCreditUntilGlued  = regexp.MustCompile(`(?i)(?:кредит|займ|за[её]m)до\s+` + reDateDMY)
	reCreditUntilInline = regexp.MustCompile(`(?i)(?:кредит|займ|за[её]m)\s+до\s+` + reDateDMY)
	reDisburseUntil     = regexp.MustCompile(`(?i)(?:передает|выдает|предоставляет|выдаёт)` + rxWord + `?\s+до\s+` + reDateDMY)
)

func parseDMYParts(day, month, year string) (time.Time, bool) {
	t, err := time.Parse("02.01.2006", day+"."+month+"."+year)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func firstDate(re *regexp.Regexp, text string) (time.Time, bool) {
	m := re.FindStringSubmatch(text)
	if len(m) < 4 {
		return time.Time{}, false
	}
	return parseDMYParts(m[1], m[2], m[3])
}

func firstDateLiteral(text string) (time.Time, bool) {
	m := regexp.MustCompile(reDateDMY).FindStringSubmatch(text)
	if len(m) < 4 {
		return time.Time{}, false
	}
	return parseDMYParts(m[1], m[2], m[3])
}

func monthsBetween(start, end time.Time) int {
	if !end.After(start) {
		return 0
	}
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	total := years*12 + months
	if end.Day() < start.Day() {
		total--
	}
	if total < 1 {
		return 1
	}
	return total
}

func creditEndDate(text string) (time.Time, bool) {
	for _, re := range []*regexp.Regexp{reCreditUntilSpaced, reCreditUntilGlued, reCreditUntilInline} {
		if t, ok := firstDate(re, text); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

func termMonthsFromDates(text string) (int, bool) {
	end, ok := creditEndDate(text)
	if !ok {
		return 0, false
	}

	start, ok := firstDate(reDisburseUntil, text)
	if !ok {
		start, ok = firstDateLiteral(text)
	}
	if !ok {
		return 0, false
	}

	months := monthsBetween(start, end)
	if months < 1 || months > 600 {
		return 0, false
	}
	return months, true
}
