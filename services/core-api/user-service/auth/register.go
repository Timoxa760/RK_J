package auth

import (
	"encoding/json"
	"net/http"
	"sync"
)

type User struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

var (
	mu     sync.Mutex
	users  = make(map[string]User)
)

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	UserID  string `json:"user_id,omitempty"`
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

	if req.Phone == "" || req.Password == "" {
		http.Error(w, `{"error":"phone and password required"}`, http.StatusBadRequest)
		return
	}

	userID := req.Phone

	mu.Lock()
	if _, exists := users[req.Phone]; !exists {
		users[req.Phone] = User{
			Phone:    req.Phone,
			Password: req.Password,
		}
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		Success: true,
		Message: "user registered",
		UserID:  userID,
	})
}
