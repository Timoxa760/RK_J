package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	deps := testDeps()
	reg := NewRegisterHandler(deps)
	regBody := `{"phone":"+79991111111","password":"secret123"}`
	regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	reg.ServeHTTP(regW, regReq)
	if regW.Code != http.StatusOK {
		t.Fatalf("register failed: %s", regW.Body.String())
	}

	h := NewLoginHandler(deps)
	body := `{"phone":"+79991111111","password":"secret123"}`
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
	if resp.User.Phone != "+79991111111" {
		t.Errorf("unexpected user: %+v", resp.User)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	deps := testDeps()
	reg := NewRegisterHandler(deps)
	regBody := `{"phone":"+79992222222","password":"secret123"}`
	regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	reg.ServeHTTP(regW, regReq)

	h := NewLoginHandler(deps)
	body := `{"phone":"+79992222222","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestPasswordReset_StubCode(t *testing.T) {
	deps := testDeps()
	reg := NewRegisterHandler(deps)
	regBody := `{"phone":"+79993333333","password":"oldpass12"}`
	regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	reg.ServeHTTP(regW, regReq)

	forgot := NewForgotPasswordHandler(deps)
	fReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password/forgot", strings.NewReader(`{"phone":"+79993333333"}`))
	fReq.Header.Set("Content-Type", "application/json")
	fW := httptest.NewRecorder()
	forgot.ServeHTTP(fW, fReq)
	if fW.Code != http.StatusOK {
		t.Fatalf("forgot expected 200, got %d", fW.Code)
	}

	reset := NewResetPasswordHandler(deps)
	rBody := `{"phone":"+79993333333","code":"0000","new_password":"newpass123"}`
	rReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password/reset", strings.NewReader(rBody))
	rReq.Header.Set("Content-Type", "application/json")
	rW := httptest.NewRecorder()
	reset.ServeHTTP(rW, rReq)
	if rW.Code != http.StatusOK {
		t.Fatalf("reset expected 200, got %d: %s", rW.Code, rW.Body.String())
	}

	login := NewLoginHandler(deps)
	lBody := `{"phone":"+79993333333","password":"newpass123"}`
	lReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(lBody))
	lReq.Header.Set("Content-Type", "application/json")
	lW := httptest.NewRecorder()
	login.ServeHTTP(lW, lReq)
	if lW.Code != http.StatusOK {
		t.Fatalf("login after reset expected 200, got %d", lW.Code)
	}
}
