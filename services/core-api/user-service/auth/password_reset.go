package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	iroot "backend_project/internal/auth"
	"backend_project/internal/otp"
	"backend_project/internal/password"
	"backend_project/internal/userstore"
)

// ForgotPasswordRequest — POST /api/v1/auth/password/forgot.
type ForgotPasswordRequest struct {
	Phone string `json:"phone"`
}

type ForgotPasswordResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

type ForgotPasswordHandler struct {
	deps *Deps
}

func NewForgotPasswordHandler(deps *Deps) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{deps: deps}
}

func (h *ForgotPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	phone := iroot.NormalizePhone(req.Phone)
	if phone == "" {
		writeJSONError(w, http.StatusBadRequest, "phone required")
		return
	}

	ttl := otp.TTLFromEnv()
	ctx := r.Context()

	// Заглушка: «отправка» кода — всегда кладём stub в Redis, если пользователь есть.
	if h.deps.Users != nil {
		if _, err := h.deps.Users.GetByPhone(ctx, phone); err == nil {
			_ = h.deps.OTP.Set(ctx, resetKey(phone), resetStubCode, ttl)
		} else if !errors.Is(err, userstore.ErrNotFound) {
			writeJSONError(w, http.StatusInternalServerError, "request failed")
			return
		}
	}

	writeJSON(w, http.StatusOK, ForgotPasswordResponse{
		Message:   "If the account exists, a reset code has been sent",
		ExpiresIn: int(ttl.Seconds()),
	})
}

// ResetPasswordRequest — POST /api/v1/auth/password/reset.
type ResetPasswordRequest struct {
	Phone       string `json:"phone"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordHandler struct {
	deps *Deps
}

func NewResetPasswordHandler(deps *Deps) *ResetPasswordHandler {
	return &ResetPasswordHandler{deps: deps}
}

func (h *ResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	phone := iroot.NormalizePhone(req.Phone)
	if phone == "" || req.Code == "" || req.NewPassword == "" {
		writeJSONError(w, http.StatusBadRequest, "phone, code and new_password required")
		return
	}

	if err := password.Validate(req.NewPassword); err != nil {
		writeJSONError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	if !h.verifyResetCode(r, phone, req.Code) {
		writeJSONError(w, http.StatusUnauthorized, "invalid or expired code")
		return
	}

	hash, err := password.Hash(req.NewPassword)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid password")
		return
	}

	if h.deps.Users == nil {
		writeJSONError(w, http.StatusInternalServerError, "user storage unavailable")
		return
	}

	if err := h.deps.Users.UpdatePassword(r.Context(), phone, hash); err != nil {
		if errors.Is(err, userstore.ErrNotFound) {
			writeJSONError(w, http.StatusUnauthorized, "invalid or expired code")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "reset failed")
		return
	}

	_ = h.deps.OTP.Delete(r.Context(), resetKey(phone))

	writeJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (h *ResetPasswordHandler) verifyResetCode(r *http.Request, phone, code string) bool {
	ctx := r.Context()
	stored, err := h.deps.OTP.Get(ctx, resetKey(phone))
	if err == nil && stored == code {
		return true
	}
	// Stub: внутренний код (не показываем на фронте).
	return code == resetStubCode
}
