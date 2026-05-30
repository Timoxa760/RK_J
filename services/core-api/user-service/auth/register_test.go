package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	mu.Lock()
	delete(users, "+79994444444")
	mu.Unlock()

	h := NewRegisterHandler(true)
	body := `{"phone":"+79994444444"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp RegisterResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Message != "SMS sent" || resp.ExpiresIn != 300 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestRegister_Conflict(t *testing.T) {
	mu.Lock()
	users["+79995555555"] = User{Phone: "+79995555555", Code: demoSMSCode}
	mu.Unlock()

	h := NewRegisterHandler(true)
	body := `{"phone":"+79995555555"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestRegister_MissingPhone(t *testing.T) {
	h := NewRegisterHandler(true)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_WrongMethod(t *testing.T) {
	h := NewRegisterHandler(true)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/register", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
