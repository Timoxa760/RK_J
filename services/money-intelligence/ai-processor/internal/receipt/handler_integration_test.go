package receipt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"backend_project/internal/expensestore"
	"backend_project/services/money-intelligence/ai-processor/internal/expense"
	"backend_project/services/money-intelligence/ai-processor/internal/manual"
)

func TestManualCreate_FileFallbackIntegration(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	dir := t.TempDir()
	fileStore, err := expensestore.NewFileStore(filepath.Join(dir, "expenses.json"))
	if err != nil {
		t.Fatal(err)
	}
	store := manual.NewFallbackStorage(nil, fileStore)
	proc := manual.NewProcessor(expense.NewParser(nil), store)
	h := NewHandler(nil, proc)

	body := `{"store":"Пятёрочка","amount":1500,"category":"Продукты","date":"2026-05-30T12:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/receipt/manual", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "+79991234567",
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("test-secret"))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	h.ManualCreate(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}

	recs, err := fileStore.ListSince("+79991234567", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	if err != nil || len(recs) != 1 {
		t.Fatalf("file records: %+v err=%v", recs, err)
	}

	var resp ManualResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != "saved" || resp.Amount != 1500 {
		t.Fatalf("resp: %+v", resp)
	}
}
