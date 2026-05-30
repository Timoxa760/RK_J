package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInsights_Contract(t *testing.T) {
	h := New(true)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights", nil)
	w := httptest.NewRecorder()
	h.insights(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp insightsResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Insights) < 2 {
		t.Fatalf("expected insights, got %+v", resp)
	}
}

func TestSimulate_ReduceDelivery(t *testing.T) {
	h := New(true)
	body := `{"scenario":"reduce_delivery","reduction_percent":50,"months":3}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scenarios/simulate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.simulate(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp simulateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Months) != 3 {
		t.Fatalf("expected 3 months, got %d", len(resp.Months))
	}
}

func TestSimulate_InvalidScenario(t *testing.T) {
	h := New(true)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scenarios/simulate", strings.NewReader(`{"scenario":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.simulate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
