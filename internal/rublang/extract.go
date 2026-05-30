package rublang

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var cleanNumRe = regexp.MustCompile(`[\s.,]`)

var (
	thousandWords = regexp.MustCompile(`(?i)(\d+(?:[\s.,]\d*)?)\s*(?:` +
		`тыс(?:яч(?:и|ей|а)?|и|ч)?|` +
		`тыщ(?:а|и|ей|у)?|` +
		`т(?:р|р\.|\.р\.?)|` +
		`штук(?:а|и|у|ами)?` +
		`)`)

	thousandLetterK    = regexp.MustCompile(`(?i)(\d+(?:[\s.,]\d*)?)\s+k(?:\s|$|[^a-z0-9])`)
	thousandLetterCyrK = regexp.MustCompile(`(?i)(\d+(?:[\s.,]\d*)?)\s+к(?:\s|$|[^0-9а-яё])`)

	kosarSuffix = regexp.MustCompile(`(?i)(\d+)\s*(?:` +
		`кос(?:ар(?:ь|я|ей|ями|ик)?)?|` +
		`к(?:ес(?:ов|а|ик)?|ос(?:ов|а)?)` +
		`)`)

	kosarWord = regexp.MustCompile(`(?i)(?:^|\s)(?:один\s+)?кос(?:ар(?:ь|я|ик)?|арь)(?:\s|$|[^а-яё])`)

	hundredSuffix = regexp.MustCompile(`(?i)(\d+)\s*сот(?:ка|ку|очка|ки|ок|ен)?`)
	hundredWord   = regexp.MustCompile(`(?i)(?:^|\s)(?:одна\s+)?сот(?:ка|очка|ку)(?:\s|$|[^а-яё])`)

	attachedK    = regexp.MustCompile(`(?i)(\d+(?:[\s.,]\d*)?)k(?:\s|$|[^a-z0-9])`)
	attachedCyrK = regexp.MustCompile(`(?i)(\d+(?:[\s.,]\d*)?)к(?:\s|$|[^0-9а-яё])`)

	rublesSuffix = regexp.MustCompile(`(?i)(\d+(?:[\s.,]\d*)?)\s*(?:₽|руб(?:лей|ля|\.?)?|р(?:\s|$|[^а-яё]))`)
	spacedThousands = regexp.MustCompile(`(?:^|\s)((?:\d{1,3})(?: \d{3})+)(?:\s|$|[^0-9])`)
	bareLarge       = regexp.MustCompile(`(?:^|\s)(\d{4,7})(?:\s|$|[^0-9])`)
)

// Normalize приводит текст к нижнему регистру с одиночными пробелами.
func Normalize(text string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(text))), " ")
}

// ExtractAll возвращает все найденные суммы в рублях (порядок появления, без дубликатов).
func ExtractAll(text string) []int {
	normalized := Normalize(text)
	if normalized == "" {
		return nil
	}

	seen := make(map[int]struct{})
	var out []int

	add := func(v int) {
		if v <= 0 {
			return
		}
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}

	parseBase := func(raw string) (float64, bool) {
		raw = cleanNumRe.ReplaceAllString(strings.TrimSpace(raw), "")
		if raw == "" {
			return 0, false
		}
		v, err := strconv.ParseFloat(raw, 64)
		return v, err == nil && v > 0
	}

	collect := func(re *regexp.Regexp, mult float64) {
		for _, m := range re.FindAllStringSubmatch(normalized, -1) {
			if v, ok := parseBase(m[1]); ok {
				add(int(v * mult))
			}
		}
	}

	collect(thousandWords, 1000)
	collect(thousandLetterK, 1000)
	collect(thousandLetterCyrK, 1000)
	collect(kosarSuffix, 1000)
	if kosarWord.MatchString(normalized) {
		add(1000)
	}
	collect(hundredSuffix, 100)
	if hundredWord.MatchString(normalized) {
		add(100)
	}
	collect(attachedK, 1000)
	collect(attachedCyrK, 1000)
	collect(rublesSuffix, 1)

	for _, m := range spacedThousands.FindAllStringSubmatch(normalized, -1) {
		if v, ok := parseBase(m[1]); ok {
			add(int(v))
		}
	}
	for _, m := range bareLarge.FindAllStringSubmatch(normalized, -1) {
		if v, ok := parseBase(m[1]); ok {
			add(int(v))
		}
	}

	for _, word := range strings.Fields(normalized) {
		if v, ok := parseAttachedSuffixWord(word); ok {
			add(v)
		}
	}

	return out
}

func parseAttachedSuffixWord(word string) (int, bool) {
	word = strings.TrimSpace(word)
	if word == "" {
		return 0, false
	}
	last, size := utf8.DecodeLastRuneInString(word)
	mult := 0
	switch last {
	case 'k', 'K', 'к', 'К':
		mult = 1000
	default:
		return 0, false
	}
	base := strings.TrimSpace(word[:len(word)-size])
	if base == "" {
		return 0, false
	}
	base = cleanNumRe.ReplaceAllString(base, "")
	v, err := strconv.ParseFloat(base, 64)
	if err != nil || v <= 0 {
		return 0, false
	}
	return int(v * float64(mult)), true
}

// ExtractPrimary возвращает наиболее вероятную сумму (максимальную из найденных).
func ExtractPrimary(text string) (float64, bool) {
	all := ExtractAll(text)
	if len(all) == 0 {
		return 0, false
	}
	max := all[0]
	for _, v := range all[1:] {
		if v > max {
			max = v
		}
	}
	return float64(max), true
}
