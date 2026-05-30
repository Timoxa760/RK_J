package manual

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

func TestDemoHandler_Create_Voice(t *testing.T) {
	h, _ := NewDemoHandler(expense.NewParser(nil))
	body := `{"user_id":"u1","raw_text":"продукты 5000","source":"voice"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses/manual", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp CreateResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Success || resp.Amount <= 0 || resp.ID == "" {
		t.Fatalf("unexpected: %+v", resp)
	}
	if len(resp.Expenses) != 1 {
		t.Fatalf("expenses: %+v", resp.Expenses)
	}
	if resp.ParsedBy != "regex" {
		t.Fatalf("parsed_by: %q", resp.ParsedBy)
	}
}

func TestDemoHandler_MissingUser(t *testing.T) {
	h, _ := NewDemoHandler(expense.NewParser(nil))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses/manual", strings.NewReader(`{"raw_text":"500"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d", w.Code)
	}
}
