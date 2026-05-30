package internal

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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
	if resp.DTI != 0 {
		t.Fatalf("expected dti 0, got %v", resp.DTI)
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

func TestScan_Unauthorized(t *testing.T) {
	credits, _ := creditstore.NewFileStore(t.TempDir())
	profiles, _ := profile.NewFileStore(t.TempDir())
	h := NewHandler(credits, profiles, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/credits/scan", nil)
	w := httptest.NewRecorder()
	h.scan(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", w.Code)
	}
}

func TestScan_MissingFile(t *testing.T) {
	credits, _ := creditstore.NewFileStore(t.TempDir())
	profiles, _ := profile.NewFileStore(t.TempDir())
	h := NewHandler(credits, profiles, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/credits/scan", nil)
	req.Header.Set("Authorization", "Bearer "+testToken(t, "+79991234567"))
	w := httptest.NewRecorder()
	h.scan(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestScan_DemoModeFallback(t *testing.T) {
	t.Setenv("DEMO_MODE", "true")
	credits, _ := creditstore.NewFileStore(t.TempDir())
	profiles, _ := profile.NewFileStore(t.TempDir())
	h := NewHandler(credits, profiles, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "contract.pdf")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("%PDF-1.4\n"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/credits/scan", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+testToken(t, "+79991234567"))
	w := httptest.NewRecorder()
	h.scan(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestScan_RejectWithoutDemoMode(t *testing.T) {
	os.Unsetenv("DEMO_MODE")
	credits, _ := creditstore.NewFileStore(t.TempDir())
	profiles, _ := profile.NewFileStore(t.TempDir())
	h := NewHandler(credits, profiles, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "contract.pdf")
	part.Write([]byte("%PDF-1.4\n"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/credits/scan", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+testToken(t, "+79991234567"))
	w := httptest.NewRecorder()
	h.scan(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}
