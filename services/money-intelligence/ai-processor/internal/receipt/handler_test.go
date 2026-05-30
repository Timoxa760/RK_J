package receipt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"backend_project/services/money-intelligence/ai-processor/internal/expense"
	"backend_project/services/money-intelligence/ai-processor/internal/manual"
)

func testToken(t *testing.T, sub string) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestManualCreate_MapsFrontContract(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	_, proc := manual.NewDemoHandler(expense.NewParser(nil))
	h := NewHandler(nil, proc)

	body := `{"store":"Пятёрочка","amount":1032.5,"category":"Продукты","date":"2026-05-30T14:32:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/receipt/manual", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, "+79991234567"))
	w := httptest.NewRecorder()

	h.ManualCreate(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}

	var resp ManualResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ReceiptID == "" || resp.Status != "saved" {
		t.Fatalf("unexpected: %+v", resp)
	}
	if resp.Store != "Пятёрочка" || resp.Amount != 1032.5 || resp.Category == "" {
		t.Fatalf("fields: %+v", resp)
	}
}

func TestManualCreate_Unauthorized(t *testing.T) {
	_, proc := manual.NewDemoHandler(expense.NewParser(nil))
	h := NewHandler(nil, proc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/receipt/manual", strings.NewReader(`{"amount":100}`))
	w := httptest.NewRecorder()
	h.ManualCreate(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", w.Code)
	}
}
