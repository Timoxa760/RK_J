package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	mu.Lock()
	users["+79991111111"] = User{Phone: "+79991111111", Code: demoSMSCode}
	mu.Unlock()
	h := NewLoginHandler(false)

	body := `{"phone":"+79991111111","code":"0000"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp LoginResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.AccessToken == "" {
		t.Error("expected access_token")
	}
	if resp.Token != resp.AccessToken {
		t.Error("expected token alias equals access_token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected refresh_token")
	}
	if resp.ExpiresIn != 900 {
		t.Errorf("expected 900s, got %d", resp.ExpiresIn)
	}
	if resp.User.Phone != "+79991111111" || resp.User.Role != "user" {
		t.Errorf("unexpected user: %+v", resp.User)
	}
}

func TestLogin_DemoCode(t *testing.T) {
	h := NewLoginHandler(true)
	body := `{"phone":"+79992222222","code":"0000"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 in demo mode, got %d", w.Code)
	}
}

func TestLogin_InvalidCode(t *testing.T) {
	h := NewLoginHandler(true)
	body := `{"phone":"+79993333333","code":"1234"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	h := NewLoginHandler(true)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
