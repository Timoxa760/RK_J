package whisper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Transcribe_ASR(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/asr") {
			t.Fatalf("path %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.MultipartForm.File["audio_file"] == nil {
			t.Fatal("expected audio_file field")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"text":"продукты 500 рублей"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "whisper-1")
	text, err := c.Transcribe(context.Background(), "voice.webm", []byte("fake-audio"))
	if err != nil {
		t.Fatal(err)
	}
	if text != "продукты 500 рублей" {
		t.Fatalf("got %q", text)
	}
}

func TestClient_Transcribe_OpenAI(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		if r.MultipartForm.File["file"] == nil {
			t.Fatal("expected file field")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"text":"тест"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL+"/v1/audio/transcriptions", "", "whisper-1")
	text, err := c.Transcribe(context.Background(), "voice.webm", []byte("fake-audio"))
	if err != nil {
		t.Fatal(err)
	}
	if text != "тест" {
		t.Fatalf("got %q", text)
	}
}

func TestClient_NotConfigured(t *testing.T) {
	c := NewClient("", "", "")
	if c.Enabled() {
		t.Fatal("expected disabled")
	}
	_, err := c.Transcribe(context.Background(), "a.webm", []byte("x"))
	if err == nil {
		t.Fatal("expected error")
	}
}
