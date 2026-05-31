package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	iroot "backend_project/internal/auth"
)

const smsRuSendURL = "https://sms.ru/sms/send"

// SmsRuSender отправляет SMS через SMS.ru API.
type SmsRuSender struct {
	apiID      string
	from       string
	testMode   bool
	httpClient *http.Client
}

func NewSmsRuSender() *SmsRuSender {
	test := os.Getenv("SMSRU_TEST") == "1"
	return &SmsRuSender{
		apiID:    strings.TrimSpace(os.Getenv("SMSRU_API_ID")),
		from:     strings.TrimSpace(os.Getenv("SMSRU_FROM")),
		testMode: test,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type smsRuResponse struct {
	Status     string            `json:"status"`
	StatusCode int               `json:"status_code"`
	StatusText string            `json:"status_text"`
	SMS        map[string]smsEntry `json:"sms"`
}

type smsEntry struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
}

func (s *SmsRuSender) SendOTP(ctx context.Context, phone, code string) error {
	if s.apiID == "" {
		return fmt.Errorf("SMSRU_API_ID not configured")
	}
	digits := phoneDigits(phone)
	if digits == "" {
		return fmt.Errorf("invalid phone")
	}
	msg := OTPMessage(code)

	form := url.Values{}
	form.Set("api_id", s.apiID)
	form.Set("to", digits)
	form.Set("msg", msg)
	form.Set("json", "1")
	if s.from != "" {
		form.Set("from", s.from)
	}
	if s.testMode {
		form.Set("test", "1")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, smsRuSendURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sms.ru request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var parsed smsRuResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("sms.ru response parse: %w (body=%s)", err, truncate(string(body), 200))
	}
	if parsed.StatusCode != 100 {
		return fmt.Errorf("sms.ru error %d: %s", parsed.StatusCode, parsed.StatusText)
	}
	if entry, ok := parsed.SMS[digits]; ok && entry.StatusCode != 100 {
		return fmt.Errorf("sms.ru delivery error %d: %s (register sender at https://sms.ru/?panel=senders)", entry.StatusCode, entry.StatusText)
	}
	return nil
}

func phoneDigits(phone string) string {
	normalized := iroot.NormalizePhone(phone)
	digits := strings.TrimPrefix(normalized, "+")
	if len(digits) < 10 {
		return ""
	}
	return digits
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
