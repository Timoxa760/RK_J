package pdfextract

import "testing"

func TestTextFromPDF_Invalid(t *testing.T) {
	_, err := TextFromPDF([]byte("not a pdf"))
	if err == nil {
		t.Fatal("expected error for non-pdf")
	}
}

func TestTextFromPDF_EmptyPDFHeader(t *testing.T) {
	_, err := TextFromPDF([]byte("%PDF-1.4\n"))
	if err == nil {
		t.Fatal("expected error for minimal pdf without text")
	}
}
