package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConnect_Success(t *testing.T) {
	h := NewConnectHandler(true)
	body := `{"credentials":{"phone":"+79991111111","password":"secret"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/connect?provider=x5club", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp ConnectResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Status != "active" || resp.Provider != "x5club" {
		t.Errorf("unexpected response: %+v", resp)
	}

	key := "+79991111111:x5club"
	mu.Lock()
	cp, ok := providers[key]
	mu.Unlock()
	if !ok {
		t.Fatal("provider not stored")
	}
	if cp.Credentials == "" {
		t.Error("expected encrypted credentials")
	}
}

func TestConnect_WrongMethod(t *testing.T) {
	h := NewConnectHandler(true)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/connect", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestConnect_MissingFields(t *testing.T) {
	h := NewConnectHandler(true)

	tests := []string{
		`{"credentials":{"p":"s"}}`,
		`{"credentials":null}`,
	}

	for _, body := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/connect", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for body %s, got %d", body, w.Code)
		}
	}
}
