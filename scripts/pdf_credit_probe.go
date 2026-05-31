//go:build ignore

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"backend_project/internal/creditparse"
	"backend_project/internal/pdfextract"
	"backend_project/internal/rublang"
)

func main() {
	path := "/Users/denis/RK_J/Кредитный договор.pdf"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}
	text, err := pdfextract.TextFromPDF(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "extract: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("=== extracted %d chars ===\n", len(text))
	if len(text) > 2500 {
		fmt.Println(text[:2500])
		fmt.Println("\n... [truncated] ...\n")
		fmt.Println(text[len(text)-800:])
	} else {
		fmt.Println(text)
	}
	fmt.Printf("\n=== general conditions? %v ===\n", creditparse.IsGeneralConditionsTemplate(text, path))
	n := rublang.Normalize(text)
	fields, ok := creditparse.ParseFromText(text)
	fmt.Printf("=== parse ok=%v fields=%+v ===\n", ok, fields)

	snippet := "заемщик выплачивает кредитору проценты за пользование кредитом — 18% годовых\nсрок кредита 12 месяцев\nпередает заемщику 500 000 ₽"
	fs, osn := creditparse.ParseFromText(snippet)
	fmt.Printf("=== snippet parse ok=%v fields=%+v ===\n", osn, fs)

	if idx := strings.Index(n, "18%"); idx >= 0 {
		start := idx - 80
		if start < 0 {
			start = 0
		}
		end := idx + 120
		if end > len(n) {
			end = len(n)
		}
		snippet2 := n[start:end] + "\nсрок кредита 12 месяцев\nсумма кредита 500 000 рублей"
		f3, ok3 := creditparse.ParseFromText(snippet2)
		fmt.Printf("=== pdf-rate-chunk parse ok=%v fields=%+v ===\n", ok3, f3)
	}
	fmt.Printf("=== amounts: %v ===\n", rublang.ExtractAll(text))
	_ = os.WriteFile("/tmp/credit_norm.txt", []byte(n), 0o644)
	for _, m := range regexp.MustCompile(`\d+\s*(?:мес|лет|год|дн|%)`).FindAllString(n, -1) {
		fmt.Println("NUM:", m)
	}
	for _, pat := range []string{`возврат[^.]{0,160}`, `погаш[^.]{0,160}`, `график[^.]{0,160}`, `ежем[^.]{0,100}`, `передает заемщику[^.]{0,80}`} {
		for _, m := range regexp.MustCompile(pat).FindAllString(n, -1) {
			fmt.Println(">>", m)
		}
	}
	for _, kw := range []string{"срок", "месяц", "лет", "год", "плат", "18%", "500", "погаш", "возврат", "аннуит", "график"} {
		idx := strings.Index(n, kw)
		if idx >= 0 {
			start := idx - 50
			if start < 0 {
				start = 0
			}
			end := idx + 150
			if end > len(n) {
				end = len(n)
			}
			fmt.Printf("\n--- %q ---\n%s\n", kw, n[start:end])
		}
	}
	if len(n) > 1200 {
		fmt.Printf("\n=== normalized sample ===\n%s\n", n[:1200])
	}
	for _, marker := range []string{"3. ", "4. ", "5. ", "6. ", "7. ", "8. ", "9. ", "10. ", "11. ", "12. "} {
		idx := strings.Index(n, marker)
		if idx >= 0 {
			end := idx + 400
			if end > len(n) {
				end = len(n)
			}
			fmt.Printf("\n=== section %s ===\n%s\n", strings.TrimSpace(marker), n[idx:end])
		}
	}
}
