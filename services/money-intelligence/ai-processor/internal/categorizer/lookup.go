package categorizer

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var textDict = NewDict()

type productMatch struct {
	prefix   string
	category string
	index    int
}

// CategoryForText подбирает категорию по ключевым словам в свободном тексте.
func CategoryForText(text string) string {
	matches := findProductMatches(text)
	if len(matches) >= 2 && allSameCategory(matches, "Продукты") {
		return "Продукты"
	}
	if len(matches) == 1 {
		return matches[0].category
	}

	t := strings.ToLower(strings.TrimSpace(text))
	if t == "" {
		return "Прочее"
	}
	return textDict.findCategory(t)
}

// ProductLabelFromText возвращает найденные товары для короткого описания.
func ProductLabelFromText(text string) string {
	matches := findProductMatches(text)
	if len(matches) >= 2 && allSameCategory(matches, "Продукты") {
		labels := make([]string, 0, len(matches))
		for _, m := range matches {
			labels = append(labels, titleWord(m.prefix))
		}
		return strings.Join(labels, ", ")
	}
	if len(matches) == 1 {
		return titleWord(matches[0].prefix)
	}

	t := strings.ToLower(strings.TrimSpace(text))
	if t == "" {
		return ""
	}
	if _, ok := textDict.categories[t]; ok {
		return titleWord(t)
	}
	return ""
}

func findProductMatches(text string) []productMatch {
	t := strings.ToLower(strings.TrimSpace(text))
	if t == "" {
		return nil
	}

	var raw []productMatch
	for prefix, cat := range textDict.categories {
		if prefix == "химия" || prefix == "корм" {
			continue
		}
		idx := strings.Index(t, prefix)
		if idx < 0 {
			continue
		}
		raw = append(raw, productMatch{prefix: prefix, category: cat, index: idx})
	}
	if len(raw) == 0 {
		return nil
	}

	sort.Slice(raw, func(i, j int) bool {
		return len(raw[i].prefix) > len(raw[j].prefix)
	})

	var picked []productMatch
	for _, m := range raw {
		if spansOverlap(picked, m) {
			continue
		}
		picked = append(picked, m)
	}

	sort.Slice(picked, func(i, j int) bool {
		if picked[i].index != picked[j].index {
			return picked[i].index < picked[j].index
		}
		return len(picked[i].prefix) > len(picked[j].prefix)
	})
	return picked
}

func spansOverlap(picked []productMatch, candidate productMatch) bool {
	start := candidate.index
	end := candidate.index + len(candidate.prefix)
	for _, p := range picked {
		pStart := p.index
		pEnd := p.index + len(p.prefix)
		if start < pEnd && pStart < end {
			return true
		}
	}
	return false
}

func allSameCategory(matches []productMatch, category string) bool {
	for _, m := range matches {
		if m.category != category {
			return false
		}
	}
	return true
}

func titleWord(word string) string {
	word = strings.TrimSpace(word)
	if word == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(word)
	return string(unicode.ToUpper(r)) + word[size:]
}
