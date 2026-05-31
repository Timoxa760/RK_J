//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

func main() {
	data, _ := os.ReadFile("/Users/denis/RK_J/Кредитный договор.pdf")
	r, _ := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	for i := 1; i <= r.NumPage(); i++ {
		text, _ := r.Page(i).GetPlainText(nil)
		fmt.Printf("\n===== PAGE %d =====\n%s\n", i, text)
		// show lines with only digits/spaces
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			digits := 0
			for _, ch := range line {
				if ch >= '0' && ch <= '9' {
					digits++
				}
			}
			if digits >= 3 && len(line) < 40 {
				fmt.Printf("  NUMLINE p%d: %q\n", i, line)
			}
		}
	}
}
