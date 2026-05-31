package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Complete_GoogleQueryKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-key" {
			t.Fatalf("key %q", r.URL.Query().Get("key"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"{\"expenses\":[]}"}]}}]}`))
	}))
	defer srv.Close()

	// isGoogleDirect() матчит подстроку generativelanguage.googleapis.com в baseURL.
	base := srv.URL + "/generativelanguage.googleapis.com/v1beta"
	c := NewClient(base, "test-key", "gemini-2.5-flash")
	out, err := c.Complete(context.Background(), "sys", "user text")
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"expenses":[]}` {
		t.Fatalf("got %q", out)
	}
}

func TestClient_Complete_AntigravityOpenAI(t *testing.T) {
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path %s", r.URL.Path)
		}
		auth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL+"/v1", "antigravity-key", "claude-sonnet-4-6")
	out, err := c.Complete(context.Background(), "sys", "user")
	if err != nil {
		t.Fatal(err)
	}
	if out != "ok" {
		t.Fatalf("got %q", out)
	}
	if auth != "Bearer antigravity-key" {
		t.Fatalf("auth %q", auth)
	}
}

func TestNormalizeBaseURL_Antigravity(t *testing.T) {
	got := normalizeBaseURL("http://127.0.0.1:8045/v1beta")
	if got != "http://127.0.0.1:8045/v1" {
		t.Fatalf("got %q", got)
	}
}

func TestClient_DisabledWithoutKey(t *testing.T) {
	c := NewClient("http://localhost/v1", "", "claude-sonnet-4-6")
	if c.Enabled() {
		t.Fatal("expected disabled")
	}
	_, err := c.Complete(context.Background(), "s", "u")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewFromEnv_AntigravityProvider(t *testing.T) {
	t.Setenv("GEMINI_PROVIDER", "antigravity")
	t.Setenv("GEMINI_API_KEY", "proxy-key")
	t.Setenv("GEMINI_BASE_URL", "")
	t.Setenv("GEMINI_MODEL", "")
	c := NewFromEnv()
	if c.baseURL != defaultAntigravityBaseURL {
		t.Fatalf("baseURL %q", c.baseURL)
	}
	if c.model != defaultAntigravityModel {
		t.Fatalf("model %q", c.model)
	}
	if !c.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestExtractChatDelta_PreservesSpaces(t *testing.T) {
	var chunk chatCompletionResponse
	if err := json.Unmarshal([]byte(`{"choices":[{"delta":{"content":" на"}}]}`), &chunk); err != nil {
		t.Fatal(err)
	}
	got := extractChatDelta(chunk)
	if got != " на" {
		t.Fatalf("got %q", got)
	}
}

func TestEndpointURL_GoogleUsesQueryKey(t *testing.T) {
	c := NewClient("https://generativelanguage.googleapis.com/v1beta", "k", "m")
	u := c.geminiEndpointURL("generateContent", false)
	if !strings.Contains(u, "key=k") {
		t.Fatalf("url %q", u)
	}
}
