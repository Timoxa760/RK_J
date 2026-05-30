package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend_project/internal/otp"
	"backend_project/internal/userstore"
)

func testDeps() *Deps {
	return NewDeps(userstore.NewMemory(), otp.NewMemoryStore())
}

func TestRegister_Success(t *testing.T) {
	h := NewRegisterHandler(testDeps())
	body := `{"phone":"+79994444444","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_Conflict(t *testing.T) {
	deps := testDeps()
	h := NewRegisterHandler(deps)
	body := `{"phone":"+79995555555","password":"secret123"}`
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if i == 0 && w.Code != http.StatusOK {
			t.Fatalf("first register expected 200, got %d", w.Code)
		}
		if i == 1 && w.Code != http.StatusConflict {
			t.Fatalf("second register expected 409, got %d", w.Code)
		}
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	h := NewRegisterHandler(testDeps())
	body := `{"phone":"+79996666666","password":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_MissingPhone(t *testing.T) {
	h := NewRegisterHandler(testDeps())
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
