package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend_project/internal/creditstore"
	"backend_project/internal/profile"
)

func TestDashboard_EmptyWithoutCredits(t *testing.T) {
	credits, err := creditstore.NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	profiles, err := profile.NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	h := NewHandler(credits, profiles, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/credits/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+testToken(t, "+79991234567"))
	w := httptest.NewRecorder()
	h.dashboard(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	var resp creditstore.Dashboard
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Credits) != 0 {
		t.Fatalf("expected empty credits, got %+v", resp)
	}
}

func TestDashboard_Unauthorized(t *testing.T) {
	credits, _ := creditstore.NewFileStore(t.TempDir())
	profiles, _ := profile.NewFileStore(t.TempDir())
	h := NewHandler(credits, profiles, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/credits/dashboard", nil)
	w := httptest.NewRecorder()
	h.dashboard(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", w.Code)
	}
}
