//go:build ignore

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"backend_project/internal/pdfextract"
)

func main() {
	data, _ := os.ReadFile("/Users/denis/RK_J/general_conditions_nal_28052020.pdf")
	text, _ := pdfextract.TextFromPDF(data)
	n := strings.ReplaceAll(text, "\n", " ")
	n = regexp.MustCompile(`\s+`).ReplaceAllString(n, " ")
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile(`(?i)\d+[.,]\d+\s*%`),
		regexp.MustCompile(`(?i)\d+\s*%\s*годовых`),
		regexp.MustCompile(`(?i)срок.{0,30}\d+.{0,20}мес`),
		regexp.MustCompile(`(?i)сумм.{0,30}\d[\d\s]{3,}`),
	} {
		m := re.FindAllString(n, 8)
		if len(m) > 0 {
			fmt.Println(re.String(), "=>", m)
		}
	}
	fmt.Println("contains общие условия:", strings.Contains(strings.ToLower(text), "общие условия"))
	fmt.Println("contains индивидуальн:", strings.Contains(strings.ToLower(text), "индивидуальн"))
}
