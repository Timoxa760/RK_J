package sms

import (
	"context"
	"log"
)

// ConsoleSender логирует OTP (только для DEMO_MODE и unit-тестов).
type ConsoleSender struct{}

func (s *ConsoleSender) SendOTP(_ context.Context, phone, code string) error {
	log.Printf("[sms:console] OTP for %s: %s", phone, code)
	return nil
}
