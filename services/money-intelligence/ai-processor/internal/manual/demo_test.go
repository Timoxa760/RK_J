package manual

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDemoHandler_Create_Voice(t *testing.T) {
	h := NewDemoHandler()
	body := `{"user_id":"u1","raw_text":"продукты 5000","source":"voice"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses/manual", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp CreateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success || resp.Amount <= 0 || resp.ID == "" {
		t.Fatalf("unexpected: %+v", resp)
	}
}
