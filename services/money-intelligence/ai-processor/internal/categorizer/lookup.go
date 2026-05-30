package categorizer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var textDict = NewDict()

// CategoryForText подбирает категорию по ключевым словам в свободном тексте.
func CategoryForText(text string) string {
	t := strings.ToLower(strings.TrimSpace(text))
	if t == "" {
		return "Прочее"
	}
	return textDict.findCategory(t)
}

// ProductLabelFromText возвращает найденное слово-товар для короткого описания.
func ProductLabelFromText(text string) string {
	t := strings.ToLower(strings.TrimSpace(text))
	if t == "" {
		return ""
	}
	if _, ok := textDict.categories[t]; ok {
		return titleWord(t)
	}
	var best string
	for prefix := range textDict.categories {
		if prefix == "химия" || prefix == "корм" {
			continue
		}
		if !strings.Contains(t, prefix) {
			continue
		}
		if len(prefix) > len(best) {
			best = prefix
		}
	}
	if best == "" {
		return ""
	}
	return titleWord(best)
}

func titleWord(word string) string {
	word = strings.TrimSpace(word)
	if word == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(word)
	return string(unicode.ToUpper(r)) + word[size:]
}
