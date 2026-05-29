package fns

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_DemoMode(t *testing.T) {
	h := NewHandler(true)

	body := `{"fn":"test","fd":"test","fp":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fns/ticket", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var receipt map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &receipt); err != nil {
		t.Fatalf("json decode: %v", err)
	}

	if receipt["store_name"] != "Пятёрочка" {
		t.Errorf("expected Пятёрочка, got %v", receipt["store_name"])
	}
}

func TestHandler_MissingFields(t *testing.T) {
	h := NewHandler(true)

	tests := []string{
		`{"fn":"","fd":"test","fp":"test"}`,
		`{"fn":"test","fd":"","fp":"test"}`,
		`{"fn":"test","fd":"test","fp":""}`,
		`{}`,
	}

	for _, body := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/fns/ticket", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for body %s, got %d", body, w.Code)
		}
	}
}

func TestHandler_WrongMethod(t *testing.T) {
	h := NewHandler(true)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/fns/ticket", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandler_InvalidJSON(t *testing.T) {
	h := NewHandler(true)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fns/ticket", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
