package pdfextract

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

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
	out := strings.TrimSpace(buf.String())
	if len(out) < 20 {
		return "", fmt.Errorf("no extractable text in pdf")
	}
	return out, nil
}
