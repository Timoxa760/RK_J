package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestUserIDFromRequest(t *testing.T) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "+79991234567",
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	got, err := UserIDFromRequest(req, "test-secret")
	if err != nil {
		t.Fatal(err)
	}
	if got != "+79991234567" {
		t.Fatalf("got %q", got)
	}
}
