package pdfextract

import (
	"strings"
	"testing"
)

func TestCleanExtractedText_JoinsBrokenWords(t *testing.T) {
	raw := "получение Кре\nдита наличными\nАльфа-\nБАНК"
	got := cleanExtractedText(raw)
	if got != "получение Кредита наличными\nАльфа-БАНК" {
		t.Fatalf("got %q", got)
	}
}

func TestCleanExtractedText_KeepsPrepositions(t *testing.T) {
	raw := "действующего\n\nна\n\nосновании\n\nУстава"
	got := cleanExtractedText(raw)
	if strings.Contains(got, "действующегона") {
		t.Fatalf("glued preposition: %q", got)
	}
	if strings.Contains(got, "действующегона") || strings.Contains(got, "основании") == false {
		t.Fatalf("expected spaced preposition, got %q", got)
	}
}
