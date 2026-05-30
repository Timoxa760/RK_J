package auth

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	iroot "backend_project/internal/auth"
)

const defaultAccessTokenTTL = 4 * time.Hour

// LoginUser — объект user в ответе login.
type LoginUser struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

// LoginResponse — 200 OK по API_Contract.
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	User         LoginUser `json:"user"`
}

func issueTokens(phone, role string) (LoginResponse, error) {
	secret := iroot.JWTSecret()
	accessTTL := accessTokenTTL()
	userID := phone

	accessClaims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(accessTTL).Unix(),
		"typ": "access",
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secret))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"typ": "refresh",
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("refresh token: %w", err)
	}

	return LoginResponse{
		AccessToken:  accessToken,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
		User: LoginUser{
			ID:    userID,
			Phone: phone,
			Role:  role,
		},
	}, nil
}

func accessTokenTTL() time.Duration {
	if raw := strings.TrimSpace(os.Getenv("JWT_ACCESS_TTL")); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil && d > 0 {
			return d
		}
	}
	if mins := strings.TrimSpace(os.Getenv("JWT_ACCESS_TTL_MINUTES")); mins != "" {
		if n, err := strconv.Atoi(mins); err == nil && n > 0 {
			return time.Duration(n) * time.Minute
		}
	}
	return defaultAccessTokenTTL
}
