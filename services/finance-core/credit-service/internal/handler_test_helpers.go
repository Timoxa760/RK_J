package internal

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func testToken(t *testing.T, phone string) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": phone,
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	return token
}
