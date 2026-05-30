package auth

import (
	"encoding/json"
	"net/http"
	"sync"
)

// User — учётная запись в in-memory хранилище (демо / MVP).
type User struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

var (
	mu    sync.Mutex
	users = make(map[string]User)
)

// RegisterRequest — POST /api/v1/auth/register (только phone).
type RegisterRequest struct {
	Phone string `json:"phone"`
}

// RegisterResponse — 200 OK по API_Contract.
type RegisterResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

type RegisterHandler struct {
	demoMode bool
}

func NewRegisterHandler(demoMode bool) *RegisterHandler {
	return &RegisterHandler{demoMode: demoMode}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		http.Error(w, `{"error":"phone required"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[req.Phone]; exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "user already exists"})
		return
	}

	users[req.Phone] = User{
		Phone: req.Phone,
		Code:  demoSMSCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		Message:   "SMS sent",
		ExpiresIn: 300,
	})
}
