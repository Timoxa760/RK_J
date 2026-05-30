package onlysq

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Complete_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("auth %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"expenses\":[]}"}}]}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL+"/v1", "test-key", "gpt-4o-mini")
	out, err := c.Complete(context.Background(), "sys", "user text")
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"expenses":[]}` {
		t.Fatalf("got %q", out)
	}
}

func TestClient_DisabledWithoutKey(t *testing.T) {
	c := NewClient("http://localhost/v1", "", "m")
	if c.Enabled() {
		t.Fatal("expected disabled")
	}
	_, err := c.Complete(context.Background(), "s", "u")
	if err == nil {
		t.Fatal("expected error")
	}
}
