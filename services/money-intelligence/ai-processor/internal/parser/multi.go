package parser

import (
	"strings"

	"backend_project/internal/rublang"
)

// ParseAll извлекает несколько трат из одной фразы (regex fallback без LLM).
func ParseAll(rawText string) []ParsedExpense {
	text := strings.TrimSpace(rawText)
	if text == "" {
		return nil
	}

	segments := splitExpenseSegments(text)
	out := make([]ParsedExpense, 0, len(segments))
	for _, seg := range segments {
		if p := Parse(seg); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

func splitExpenseSegments(rawText string) []string {
	normalized := rublang.Normalize(rawText)
	if normalized == "" {
		return nil
	}

	for _, sep := range []string{",", ";"} {
		if strings.Contains(normalized, sep) {
			if parts := splitListSegments(normalized, sep); len(parts) > 1 {
				return parts
			}
		}
	}

	if strings.Contains(normalized, " и ") {
		if parts := splitListSegments(normalized, " и "); len(parts) > 1 && eachSegmentHasSingleAmount(parts) {
			return parts
		}
	}

	if ends := rublang.AmountPhraseEnds(normalized); len(ends) > 1 {
		return splitByAmountEnds(normalized, ends)
	}

	return []string{strings.TrimSpace(rawText)}
}

func splitListSegments(text, sep string) []string {
	parts := strings.Split(text, sep)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if len(rublang.ExtractAll(part)) == 0 {
			continue
		}
		out = append(out, part)
	}
	return out
}

func eachSegmentHasSingleAmount(parts []string) bool {
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if len(rublang.ExtractAll(part)) != 1 {
			return false
		}
	}
	return true
}

func splitByAmountEnds(normalized string, ends []int) []string {
	out := make([]string, 0, len(ends))
	start := 0
	for _, end := range ends {
		seg := strings.TrimSpace(normalized[start:end])
		if seg != "" {
			out = append(out, seg)
		}
		start = end
		for strings.HasPrefix(strings.TrimSpace(normalized[start:]), "и ") {
			start += 2
		}
	}
	trailing := strings.TrimSpace(normalized[start:])
	if trailing != "" {
		if len(out) == 0 {
			out = append(out, trailing)
		} else if len(rublang.ExtractAll(trailing)) > 0 {
			out = append(out, trailing)
		} else {
			out[len(out)-1] = strings.TrimSpace(out[len(out)-1] + " " + trailing)
		}
	}
	return out
}
