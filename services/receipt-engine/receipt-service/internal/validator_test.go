package internal

import (
	"testing"
	"time"
)

func TestValidateRaw_Empty(t *testing.T) {
	if err := ValidateRaw(nil); err == nil {
		t.Error("expected error for nil")
	}
}

func TestValidateRaw_TooLarge(t *testing.T) {
	data := make([]byte, 11*1024*1024)
	if err := ValidateRaw(data); err != ErrTooLarge {
		t.Errorf("expected ErrTooLarge, got %v", err)
	}
}

func TestValidateRaw_InvalidUTF8(t *testing.T) {
	data := []byte{0xff, 0xfe, 0x00, 0x01}
	if err := ValidateRaw(data); err == nil {
		t.Error("expected error for invalid utf-8")
	}
}

func TestValidateRaw_InvalidJSON(t *testing.T) {
	data := []byte(`{invalid}`)
	if err := ValidateRaw(data); err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestValidateRaw_Valid(t *testing.T) {
	data := []byte(`{"id":"test"}`)
	if err := ValidateRaw(data); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateReceipt_OK(t *testing.T) {
	r := &RawReceipt{
		ID: "r1", UserID: "u1", Provider: "x5", Store: "Store",
		Date: time.Now(), Total: 100,
		Items: []Item{{Name: "Item", Price: 50, Quantity: 2}},
	}
	if err := ValidateReceipt(r); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateReceipt_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		r    RawReceipt
		err  error
	}{
		{"empty id", RawReceipt{UserID: "u", Provider: "x", Store: "s", Total: 1, Items: []Item{{Name: "a", Price: 1, Quantity: 1}}}, ErrEmptyID},
		{"empty user_id", RawReceipt{ID: "r", Provider: "x", Store: "s", Total: 1, Items: []Item{{Name: "a", Price: 1, Quantity: 1}}}, ErrEmptyUserID},
		{"empty provider", RawReceipt{ID: "r", UserID: "u", Store: "s", Total: 1, Items: []Item{{Name: "a", Price: 1, Quantity: 1}}}, ErrEmptyProvider},
		{"empty store", RawReceipt{ID: "r", UserID: "u", Provider: "x", Total: 1, Items: []Item{{Name: "a", Price: 1, Quantity: 1}}}, ErrEmptyStore},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateReceipt(&tt.r); err != tt.err {
				t.Errorf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestValidateReceipt_NegativeTotal(t *testing.T) {
	r := &RawReceipt{
		ID: "r", UserID: "u", Provider: "x", Store: "s",
		Total: -1, Items: []Item{{Name: "a", Price: 1, Quantity: 1}},
	}
	if err := ValidateReceipt(r); err != ErrNegativeTotal {
		t.Errorf("expected ErrNegativeTotal, got %v", err)
	}
}

func TestValidateReceipt_NoItems(t *testing.T) {
	r := &RawReceipt{ID: "r", UserID: "u", Provider: "x", Store: "s", Total: 0}
	if err := ValidateReceipt(r); err != ErrNoItems {
		t.Errorf("expected ErrNoItems, got %v", err)
	}
}

func TestValidateReceipt_ItemErrors(t *testing.T) {
	tests := []struct {
		name string
		r    RawReceipt
	}{
		{"empty item name", RawReceipt{ID: "r", UserID: "u", Provider: "x", Store: "s", Total: 10, Items: []Item{{Name: "", Price: 1, Quantity: 1}}}},
		{"negative price", RawReceipt{ID: "r", UserID: "u", Provider: "x", Store: "s", Total: 10, Items: []Item{{Name: "a", Price: -1, Quantity: 1}}}},
		{"zero quantity", RawReceipt{ID: "r", UserID: "u", Provider: "x", Store: "s", Total: 10, Items: []Item{{Name: "a", Price: 1, Quantity: 0}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateReceipt(&tt.r); err == nil {
				t.Error("expected error")
			}
		})
	}
}
