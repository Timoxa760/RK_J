package sms

import (
	"context"
	"testing"
)

func TestSmsRuSender_MissingAPIKey(t *testing.T) {
	s := &SmsRuSender{apiID: ""}
	err := s.SendOTP(context.Background(), "+79991234567", "482913")
	if err == nil {
		t.Fatal("expected error without api id")
	}
}

func TestOTPMessageLength(t *testing.T) {
	msg := OTPMessage("123456")
	if len([]rune(msg)) > 70 {
		t.Errorf("message too long for free SMS tier: %d chars", len([]rune(msg)))
	}
}

func TestPhoneDigits(t *testing.T) {
	if phoneDigits("+79991234567") != "79991234567" {
		t.Error("expected digits without plus")
	}
	if phoneDigits("bad") != "" {
		t.Error("expected empty for bad phone")
	}
}
