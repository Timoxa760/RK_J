package voice

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"backend_project/services/money-intelligence/ai-processor/internal/whisper"
)

func testToken(t *testing.T, sub string) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestHandler_Transcribe_OK(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	asr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"text":"180 тысяч в месяц"}`))
	}))
	defer asr.Close()

	wh := whisper.NewClient(asr.URL, "", "whisper-1")
	h := NewHandler(wh, nil)

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("audio", "recording.webm")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("fake-audio")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/voice/transcribe", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+testToken(t, "+79991234567"))

	rr := httptest.NewRecorder()
	h.Transcribe(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	if !bytes.Contains(rr.Body.Bytes(), []byte("180 тысяч")) {
		t.Fatalf("unexpected body %s", rr.Body.String())
	}
}

func TestHandler_Transcribe_Unauthorized(t *testing.T) {
	h := NewHandler(whisper.NewClient("http://127.0.0.1:1", "", ""), nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/voice/transcribe", nil)
	rr := httptest.NewRecorder()
	h.Transcribe(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", rr.Code)
	}
}
