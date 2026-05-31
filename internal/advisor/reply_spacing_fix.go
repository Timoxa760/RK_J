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
	}

	explicitSplitReplacements = []string{
		"заполнении",
		"заполнен",
		"поэтому",
		"поэтом$1",
	}

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
}

// RepairSplitRussianWords убирает лишние пробелы внутри русских слов (артефакт LLM).
func RepairSplitRussianWords(s string) string {
	if s == "" {
		return s
	}
	for i, re := range reExplicitSplitFixes {
		s = re.ReplaceAllString(s, explicitSplitReplacements[i])
	}
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

func replaceWithSubmatches(re *regexp.Regexp, s string, merge func([]string) string) string {
	return re.ReplaceAllStringFunc(s, func(full string) string {
		return merge(re.FindStringSubmatch(full))
	})
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
	return RepairSplitRussianWords(s)
}
