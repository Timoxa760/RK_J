package auth

import (
	"regexp"
	"strings"
)

var phoneNonDigit = regexp.MustCompile(`\D`)

// NormalizePhone приводит номер телефона к каноническому виду +7XXXXXXXXXX.
func NormalizePhone(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	digits := phoneNonDigit.ReplaceAllString(raw, "")
	if len(digits) == 11 && (digits[0] == '8' || digits[0] == '7') {
		digits = digits[1:]
	}
	if len(digits) == 10 {
		return "+7" + digits
	}
	if strings.HasPrefix(raw, "+") && len(digits) >= 10 {
		return "+" + digits
	}
	return raw
}
