package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDashboard_Contract(t *testing.T) {
	h := New(true)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/credits/dashboard", nil)
	w := httptest.NewRecorder()
	h.dashboard(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp dashboardResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.DTI != 0.28 || len(resp.Credits) != 1 {
		t.Fatalf("unexpected: %+v", resp)
	}
}
