package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDigestLatest(t *testing.T) {
	h := New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/digest/latest", nil)
	w := httptest.NewRecorder()
	h.digestLatest(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp digestResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.TotalSpent <= 0 || resp.MindfulnessRating == 0 {
		t.Fatalf("unexpected: %+v", resp)
	}
}
