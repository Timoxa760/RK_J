package whisper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Transcribe_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method %s", r.Method)
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
