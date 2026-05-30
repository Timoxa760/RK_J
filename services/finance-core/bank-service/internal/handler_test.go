package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAccounts(t *testing.T) {
	r := chi.NewRouter()
	New().Register(r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/banks/accounts", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestTransactions(t *testing.T) {
	r := chi.NewRouter()
	New().Register(r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/banks/transactions", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
