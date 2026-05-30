package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Complete_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, ":generateContent") {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.URL.Query().Get("key") != "test-key" {
			t.Fatalf("key %q", r.URL.Query().Get("key"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"{\"expenses\":[]}"}]}}]}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-key", "gemini-2.0-flash")
	out, err := c.Complete(context.Background(), "sys", "user text")
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"expenses":[]}` {
		t.Fatalf("got %q", out)
	}
}

func TestClient_DisabledWithoutKey(t *testing.T) {
	c := NewClient("http://localhost", "", "gemini-2.0-flash")
	if c.Enabled() {
		t.Fatal("expected disabled")
	}
	_, err := c.Complete(context.Background(), "s", "u")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewFromEnv(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "abc")
	t.Setenv("GEMINI_MODEL", "gemini-2.0-flash")
	c := NewFromEnv()
	if !c.Enabled() || c.model != "gemini-2.0-flash" {
		t.Fatalf("unexpected client %+v", c)
	}
}
