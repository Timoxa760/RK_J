package fns

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type stubProducer struct {
	sent bool
}

func (s *stubProducer) Send(_ context.Context, _ scrap.RawReceipt) error {
	s.sent = true
	return nil
}

func TestScanHandler_OK(t *testing.T) {
	fnsH := NewHandler(true)
	prod := &stubProducer{}
	h := NewScanHandler(fnsH, prod)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/receipt/fns/scan", bytes.NewBufferString(
		`{"fn":"9289000100123456","fd":"12345","fp":"1234567890"}`,
	))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if !prod.sent {
		t.Fatal("expected kafka send")
	}
}
