package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	Message      string `json:"message,omitempty"`
}

type LoginHandler struct {
	demoMode bool
}

func NewLoginHandler(demoMode bool) *LoginHandler {
	return &LoginHandler{demoMode: demoMode}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if req.Phone == "" || req.Password == "" {
		http.Error(w, `{"error":"phone and password required"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := users[req.Phone]
	mu.Unlock()

	if h.demoMode {
		if !exists {
			mu.Lock()
			users[req.Phone] = User{Phone: req.Phone, Password: req.Password}
			user = users[req.Phone]
			mu.Unlock()
		}
	} else if !exists || user.Password != req.Password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Success: false, Message: "invalid credentials"})
		return
	}

	secret := getJWTSecret()
	userID := req.Phone

	accessClaims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"typ": "access",
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secret))
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"token: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"typ": "refresh",
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"refresh token: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	})
}

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "money-mind-demo-secret-key-2026"
}
