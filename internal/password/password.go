package password

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const MinLength = 8

var (
	ErrTooShort = errors.New("password must be at least 8 characters")
	ErrMismatch = errors.New("password mismatch")
)

func Validate(plain string) error {
	if utf8.RuneCountInString(plain) < MinLength {
		return ErrTooShort
	}
	return nil
}

func Hash(plain string) (string, error) {
	if err := Validate(plain); err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

func Compare(hash, plain string) error {
	if hash == "" {
		return bcrypt.ErrMismatchedHashAndPassword
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
