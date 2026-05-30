package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var ErrUnauthorized = errors.New("unauthorized")

// UserIDFromRequest извлекает claim sub из Bearer JWT.
func UserIDFromRequest(r *http.Request, secret string) (string, error) {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", ErrUnauthorized
	}
	tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	if tokenStr == "" {
		return "", ErrUnauthorized
	}
	if secret == "" {
		secret = JWTSecret()
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnauthorized
	}
	sub, _ := claims["sub"].(string)
	sub = NormalizePhone(strings.TrimSpace(sub))
	if sub == "" {
		return "", ErrUnauthorized
	}
	return sub, nil
}

// JWTSecret читает секрет подписи JWT из окружения.
func JWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "test-secret"
}
