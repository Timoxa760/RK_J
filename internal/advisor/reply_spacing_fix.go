package advisor

import (
	"regexp"
	"strings"
)

var (
	reExplicitSplitFixes = []*regexp.Regexp{
		regexp.MustCompile(`(?i)за\s+по\s+лнении`),
		regexp.MustCompile(`(?i)за\s+по\s+лнен`),
		regexp.MustCompile(`(?i)по\s+это\s+му`),
		regexp.MustCompile(`(?i)по\s+это\s+м([^а-яё]|$)`),
		regexp.MustCompile(`(?i)де\s+по\s+зит`),
		regexp.MustCompile(`(?i)го\s+до\s+вых`),
		regexp.MustCompile(`(?i)кр\s+из\s+ис`),
		regexp.MustCompile(`(?i)пред\s+по\s+чтительнее`),
		regexp.MustCompile(`(?i)рабо\s+тают`),
		regexp.MustCompile(`(?i)не\s+рабают`),
		regexp.MustCompile(`(?i)без\s+укания`),
		regexp.MustCompile(`(?i)ставки\s*цб`),
		regexp.MustCompile(`(?i)ис\s+по\s+льз`),
		regexp.MustCompile(`(?i)не\s+до\s+стающ`),
		regexp.MustCompile(`(?i)единственныйфи\s+на\s+нсовый`),
		regexp.MustCompile(`(?i)фи\s+на\s+нсовый`),
		regexp.MustCompile(`(?i)месяцона`),
		regexp.MustCompile(`(?i)увасуже`),
		regexp.MustCompile(`(?i)заменын[а-яёa-z]+`),
	}

	explicitSplitReplacements = []string{
		"заполнении",
		"заполнен",
		"поэтому",
		"поэтом$1",
		"депозит",
		"годовых",
		"кризис",
		"предпочтительнее",
		"работают",
		"не работают",
		"без уточнения",
		"ставки ЦБ",
		"использ",
		"недостающ",
		"единственный финансовый",
		"финансовый",
		"месяц она",
		"у вас уже",
		"замены нет",
	}

	// «ис по льзовать», «не до стающие» — короткие слоги + длинный хвост.
	reMixedFragmentChain = regexp.MustCompile(`(?i)(^|[\s(«"'—-])((?:[а-яё]{1,3}\s+){2,}[а-яё]{4,12})([\s,.!?;:»)\]"'—-]|$)`)

	// «фи на нсовый» — «на» внутри слова, не предлог (короткий левый фрагмент).
	rePrepInsideWordNa = regexp.MustCompile(`(?i)([а-яё]{1,3}) на ([б-джзклмнпрстфхцчшщ][а-яё]{3,})`)

	// «де по зит», «го до вых», «кр из ис» — 3+ коротких слога подряд (RE2 без lookahead).
	reShortFragmentChain = regexp.MustCompile(`(?i)(^|[\s(«"'—-])([а-яё]{1,3}(?: [а-яё]{1,3}){2,})([\s,.!?;:»)\]"'—-]|$)`)

	// «пред по чтительнее» — «по» внутри слова, не предлог.
	rePrepInsideWord = regexp.MustCompile(`(?i)([а-яё]{3,6}) по ([б-джзклмнпрстфхцчшщ][а-яё]{4,})`)

	// «ука за ть», «транс по рт» — короткий хвост после за|по|от|или внутри слова.
	reEmbeddedPrepSplit = regexp.MustCompile(`(?i)([а-яё]{3,}) (?:за|по|от|или) ([а-яё]{1,3})([\s,.!?;:«»\)\]"'\-—]|$)`)

	// «доход ы», «расход ов» — короткое окончание, не отдельное слово.
	reSuffixSplit = regexp.MustCompile(`(?i)([а-яё]{4,}) ([а-яё]{1,3})([\s,.!?;:«»\)\]"'\-—]|$)`)
)

var splitMergeStop = map[string]struct{}{
	"за": {}, "по": {}, "от": {}, "до": {}, "на": {}, "в": {}, "и": {}, "но": {},
	"или": {}, "при": {}, "для": {}, "из": {}, "с": {}, "к": {}, "у": {}, "о": {}, "об": {},
	"не": {}, "ни": {}, "вы": {}, "мы": {}, "он": {}, "я": {}, "ты": {}, "её": {}, "ее": {}, "их": {},
	"то": {}, "как": {}, "что": {}, "это": {}, "там": {}, "тут": {}, "уже": {}, "ещё": {}, "еще": {},
	"все": {}, "всё": {}, "без": {}, "про": {}, "под": {}, "над": {}, "мне": {}, "вас": {}, "нас": {},
	"ему": {}, "ей": {}, "им": {}, "бы": {}, "ли": {}, "же": {},
	"она": {}, "они": {}, "нет": {}, "есть": {},
}

// RepairSplitRussianWords убирает лишние пробелы внутри русских слов (артефакт LLM).
func RepairSplitRussianWords(s string) string {
	if s == "" {
		return s
	}
	for i, re := range reExplicitSplitFixes {
		s = re.ReplaceAllString(s, explicitSplitReplacements[i])
	}
	s = repairShortFragmentChains(s)
	s = repairMixedFragmentChains(s)
	s = replaceWithSubmatches(rePrepInsideWord, s, mergePrepInsideWord)
	s = replaceWithSubmatches(rePrepInsideWordNa, s, mergePrepInsideWordNa)
	s = replaceWithSubmatches(reEmbeddedPrepSplit, s, mergeIfNotStopword)
	for i := 0; i < 6; i++ {
		next := replaceWithSubmatches(reSuffixSplit, s, mergeIfNotStopword)
		if next == s {
			break
		}
		s = next
	}
	return s
}

func repairShortFragmentChains(s string) string {
	return replaceWithSubmatches(reShortFragmentChain, s, mergeShortFragmentChain)
}

func mergeShortFragmentChain(m []string) string {
	if len(m) < 4 {
		return m[0]
	}
	prefix := m[1]
	chain := m[2]
	suffix := m[3]
	parts := strings.Fields(chain)
	if len(parts) < 3 {
		return m[0]
	}
	if allSplitMergeStop(parts) {
		return m[0]
	}
	for _, p := range parts {
		if len([]rune(p)) > 3 {
			return m[0]
		}
	}
	merged := strings.Join(parts, "")
	runes := len([]rune(merged))
	if runes < 5 || runes > 24 {
		return m[0]
	}
	return prefix + merged + suffix
}

func repairMixedFragmentChains(s string) string {
	return replaceWithSubmatches(reMixedFragmentChain, s, mergeMixedFragmentChain)
}

func mergeMixedFragmentChain(m []string) string {
	if len(m) < 4 {
		return m[0]
	}
	prefix := m[1]
	chain := m[2]
	suffix := m[3]
	parts := strings.Fields(chain)
	if len(parts) < 3 {
		return m[0]
	}
	for _, p := range parts[:len(parts)-1] {
		if len([]rune(p)) > 3 {
			return m[0]
		}
		if len([]rune(p)) < 2 {
			return m[0]
		}
	}
	last := parts[len(parts)-1]
	if lr := len([]rune(last)); lr < 4 || lr > 12 {
		return m[0]
	}
	merged := strings.Join(parts, "")
	runes := len([]rune(merged))
	if runes < 6 || runes > 32 {
		return m[0]
	}
	return prefix + merged + suffix
}

func mergePrepInsideWordNa(m []string) string {
	if len(m) < 3 {
		return m[0]
	}
	stem := strings.ToLower(m[1])
	if _, skip := prepInsideWordNaFalsePositives[stem]; skip {
		return m[0]
	}
	return m[1] + "на" + m[2]
}

func mergePrepInsideWord(m []string) string {
	if len(m) < 3 {
		return m[0]
	}
	stem := strings.ToLower(m[1])
	if _, skip := prepInsideWordFalsePositives[stem]; skip {
		return m[0]
	}
	return m[1] + "по" + m[2]
}

var prepInsideWordFalsePositives = map[string]struct{}{
	"идти": {}, "идём": {}, "идем": {}, "иду": {}, "шли": {}, "шёл": {}, "шел": {},
	"смотреть": {}, "смотрите": {}, "смотрим": {}, "ориентир": {}, "опираться": {},
}

var prepInsideWordNaFalsePositives = map[string]struct{}{
	"он": {}, "она": {}, "они": {}, "мы": {}, "вы": {}, "то": {}, "ту": {}, "ему": {}, "ей": {}, "им": {},
}

func replaceWithSubmatches(re *regexp.Regexp, s string, merge func([]string) string) string {
	return re.ReplaceAllStringFunc(s, func(full string) string {
		return merge(re.FindStringSubmatch(full))
	})
}

func allSplitMergeStop(parts []string) bool {
	for _, p := range parts {
		if _, stop := splitMergeStop[strings.ToLower(p)]; !stop {
			return false
		}
	}
	return true
}

func mergeIfNotStopword(m []string) string {
	if len(m) < 4 {
		return m[0]
	}
	frag := strings.ToLower(m[2])
	stem := m[1]
	if !shouldMergeSplitFragment(stem, frag) {
		return m[0]
	}
	return m[1] + m[2] + m[3]
}

func shouldMergeSplitFragment(stem, frag string) bool {
	frag = strings.ToLower(frag)
	if _, stop := splitMergeStop[frag]; stop {
		return false
	}
	if len(frag) == 1 {
		r := []rune(frag)[0]
		if r == 'и' {
			return false
		}
		if isRussianVowel(r) && endsWithRussianVowel(stem) {
			return false
		}
	}
	return true
}

func isRussianVowel(r rune) bool {
	switch r {
	case 'а', 'е', 'ё', 'и', 'о', 'у', 'ы', 'э', 'ю', 'я',
		'А', 'Е', 'Ё', 'И', 'О', 'У', 'Ы', 'Э', 'Ю', 'Я':
		return true
	default:
		return false
	}
}

func endsWithRussianVowel(s string) bool {
	runes := []rune(s)
	if len(runes) == 0 {
		return false
	}
	return isRussianVowel(runes[len(runes)-1])
}

func repairStructuredText(s string) string {
	out := RepairSplitRussianWords(s)
	// Второй проход — цепочки могли образоваться после склейки хвостов.
	return RepairSplitRussianWords(out)
}
