package receipt

import (
	"errors"
	"net/http"

	"backend_project/internal/auth"
)

var errUnauthorized = errors.New("unauthorized")

func userIDFromRequest(r *http.Request) (string, error) {
	id, err := auth.UserIDFromRequest(r, auth.JWTSecret())
	if err != nil {
		return "", errUnauthorized
	}
	return id, nil
}
