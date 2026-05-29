package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"unicode/utf8"
)

var (
	ErrEmptyID       = errors.New("receipt: empty id")
	ErrEmptyUserID   = errors.New("receipt: empty user_id")
	ErrEmptyProvider = errors.New("receipt: empty provider")
	ErrEmptyStore    = errors.New("receipt: empty store_name")
	ErrNegativeTotal = errors.New("receipt: negative total_amount")
	ErrNoItems       = errors.New("receipt: no items")
	ErrTooLarge      = errors.New("receipt: exceeds 10MB")
)

const maxRawSize = 10 * 1024 * 1024

func ValidateRaw(data []byte) error {
	if len(data) == 0 {
		return errors.New("receipt: empty body")
	}
	if len(data) > maxRawSize {
		return ErrTooLarge
	}
	if !utf8.Valid(data) {
		return errors.New("receipt: invalid utf-8")
	}
	if !json.Valid(data) {
		return errors.New("receipt: invalid json")
	}
	return nil
}

func ValidateReceipt(r *RawReceipt) error {
	if r.ID == "" {
		return ErrEmptyID
	}
	if r.UserID == "" {
		return ErrEmptyUserID
	}
	if r.Provider == "" {
		return ErrEmptyProvider
	}
	if r.Store == "" {
		return ErrEmptyStore
	}
	if r.Total < 0 {
		return ErrNegativeTotal
	}
	if len(r.Items) == 0 {
		return ErrNoItems
	}
	for i, item := range r.Items {
		if item.Name == "" {
			return fmt.Errorf("receipt: item[%d] empty name", i)
		}
		if item.Price < 0 {
			return fmt.Errorf("receipt: item[%d] negative price", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("receipt: item[%d] invalid quantity", i)
		}
	}
	return nil
}
