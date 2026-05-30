package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const demoSMSCode = "0000"

// LoginRequest — тело POST /api/v1/auth/login (API_Contract).
type LoginRequest struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password,omitempty"` // legacy, не в контракте
}

// LoginUser — объект user в ответе login.
type LoginUser struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

// LoginResponse — 200 OK по API_Contract.
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	User         LoginUser `json:"user"`
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

	code := req.Code
	if code == "" && req.Password != "" {
		code = req.Password
	}

	if req.Phone == "" || code == "" {
		http.Error(w, `{"error":"phone and code required"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := users[req.Phone]
	mu.Unlock()

	if h.demoMode {
		if code != demoSMSCode {
			writeLoginUnauthorized(w)
			return
		}
		if !exists {
			mu.Lock()
			users[req.Phone] = User{Phone: req.Phone, Code: demoSMSCode}
			mu.Unlock()
		}
	} else {
		if !exists || user.Code != code {
			writeLoginUnauthorized(w)
			return
		}
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
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		User: LoginUser{
			ID:    userID,
			Phone: req.Phone,
			Role:  "user",
		},
	})
}

func writeLoginUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "invalid code"})
}

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "money-mind-demo-secret-key-2026"
}
