//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ledongthuc/pdf"
)

func main() {
	data, _ := os.ReadFile("/Users/denis/RK_J/Кредитный договор.pdf")
	r, _ := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	fmt.Println("pages:", r.NumPage())
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		text, err := page.GetPlainText(nil)
		if err != nil {
			fmt.Printf("page %d err: %v\n", i, err)
			continue
		}
		fmt.Printf("page %d chars: %d\n", i, len(text))
	}
}
