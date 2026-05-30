package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	iroot "backend_project/internal/auth"
	"backend_project/internal/password"
	"backend_project/internal/userstore"
)

// RegisterRequest — POST /api/v1/auth/register.
type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// RegisterResponse — 200 OK.
type RegisterResponse struct {
	Message string `json:"message"`
}

type RegisterHandler struct {
	deps *Deps
}

func NewRegisterHandler(deps *Deps) *RegisterHandler {
	return &RegisterHandler{deps: deps}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
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
		writeJSONError(w, http.StatusBadRequest, "phone required")
		return
	}

	if err := password.Validate(req.Password); err != nil {
		writeJSONError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid password")
		return
	}

	if h.deps.Users == nil {
		writeJSONError(w, http.StatusInternalServerError, "user storage unavailable")
		return
	}

	if err := h.deps.Users.Create(r.Context(), phone, hash); err != nil {
		if errors.Is(err, userstore.ErrExists) {
			writeJSONError(w, http.StatusConflict, "user already exists")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "registration failed")
		return
	}

	writeJSON(w, http.StatusOK, RegisterResponse{Message: "registered"})
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
