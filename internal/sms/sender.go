package sms

import (
	"context"
	"fmt"
	"os"
)

// Sender доставляет OTP на телефон.
type Sender interface {
	SendOTP(ctx context.Context, phone, code string) error
}

func NewFromEnv() Sender {
	provider := os.Getenv("SMS_PROVIDER")
	if provider == "" {
		if os.Getenv("DEMO_MODE") == "true" {
			provider = "console"
		} else {
			provider = "smsru"
		}
	}
	switch provider {
	case "console":
		return &ConsoleSender{}
	case "smsru":
		return NewSmsRuSender()
	default:
		return NewSmsRuSender()
	}
}

func OTPMessage(code string) string {
	return fmt.Sprintf("Поток: код %s. Действует 5 мин.", code)
}
