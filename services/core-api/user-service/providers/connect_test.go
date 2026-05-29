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
	body := `{"user_id":"+79991111111","provider":"x5club","credentials":{"phone":"+79991111111","password":"secret"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/connect", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp ConnectResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Error("expected success")
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
		`{"user_id":"","provider":"x5","credentials":{"p":"s"}}`,
		`{"user_id":"u1","provider":"","credentials":{"p":"s"}}`,
		`{"user_id":"u1","provider":"x5","credentials":null}`,
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
