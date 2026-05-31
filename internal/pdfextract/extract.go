package pdfextract

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

var (
	reHyphenBreak   = regexp.MustCompile(`-\s*\n\s*`)
	reCyrWordBreak  = regexp.MustCompile(`([а-яё])\s*\n\s*([а-яё]{1,6})([\s,.!?;:»)\]"'—\-]|$)`)
	reLatinWordBreak = regexp.MustCompile(`(?i)([a-z])\s*\n\s*([a-z]{1,6})([\s,.!?;:]|$)`)
	reGluedPrep     = regexp.MustCompile(`(?i)([а-яё]{3,})(на|по|за|от|до|из|при|для|без|под|над|про|об|со)([а-яё]{2,})`)
)

var lineBreakGlueStop = map[string]struct{}{
	"на": {}, "по": {}, "за": {}, "от": {}, "до": {}, "из": {}, "об": {}, "со": {},
	"при": {}, "для": {}, "или": {}, "не": {}, "ни": {}, "без": {}, "под": {}, "над": {},
	"про": {}, "и": {}, "в": {}, "к": {}, "с": {}, "у": {}, "о": {}, "а": {}, "но": {},
}

// TextFromPDF извлекает plain text из PDF для передачи в LLM.
func TextFromPDF(data []byte) (string, error) {
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		return "", fmt.Errorf("invalid pdf")
	}
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("read pdf: %w", err)
	}
	var buf strings.Builder
	total := r.NumPage()
	for i := 1; i <= total; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		buf.WriteByte('\n')
	}
	out := cleanExtractedText(buf.String())
	if len(out) < 20 {
		return "", fmt.Errorf("no extractable text in pdf")
	}
	return out, nil
}

func shouldJoinLineBreak(_ string, cont string) bool {
	_, stop := lineBreakGlueStop[strings.ToLower(cont)]
	return !stop
}

func joinCyrLineBreak(full, stem, cont, tail string) string {
	if !shouldJoinLineBreak(stem, cont) {
		return full
	}
	return stem + cont + tail
}

// cleanExtractedText склеивает слова, разорванные переносами строк в PDF.
func cleanExtractedText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = reHyphenBreak.ReplaceAllString(s, "-")
	s = reCyrWordBreak.ReplaceAllStringFunc(s, func(full string) string {
		parts := reCyrWordBreak.FindStringSubmatch(full)
		if len(parts) < 4 {
			return full
		}
		return joinCyrLineBreak(full, parts[1], parts[2], parts[3])
	})
	s = reLatinWordBreak.ReplaceAllString(s, "$1$2$3")
	s = reGluedPrep.ReplaceAllString(s, "$1 $2 $3")
	return strings.TrimSpace(s)
}
