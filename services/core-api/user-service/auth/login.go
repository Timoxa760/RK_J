package auth

import (
	"encoding/json"
	"net/http"

	iroot "backend_project/internal/auth"
	"backend_project/internal/password"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest — POST /api/v1/auth/login.
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginHandler struct {
	deps *Deps
}

func NewLoginHandler(deps *Deps) *LoginHandler {
	return &LoginHandler{deps: deps}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Phone == "" || req.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "phone and password required")
		return
	}

	phone := iroot.NormalizePhone(req.Phone)
	if phone == "" {
		writeJSONError(w, http.StatusBadRequest, "phone and password required")
		return
	}

	if h.deps.Users == nil {
		writeJSONError(w, http.StatusInternalServerError, "user storage unavailable")
		return
	}

	user, err := h.deps.Users.GetByPhone(r.Context(), phone)
	if err != nil {
		writeLoginUnauthorized(w)
		return
	}

	if err := password.Compare(user.PasswordHash, req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			writeLoginUnauthorized(w)
			return
		}
		writeLoginUnauthorized(w)
		return
	}

	resp, err := issueTokens(phone, user.Role)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "token issue failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeLoginUnauthorized(w http.ResponseWriter) {
	writeJSONError(w, http.StatusUnauthorized, "invalid phone or password")
}
