package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSankey_Demo(t *testing.T) {
	h := New(true, "test-secret")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/sankey", nil)
	w := httptest.NewRecorder()
	h.sankey(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp sankeyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Nodes) < 3 || len(resp.Links) < 2 {
		t.Fatalf("expected sankey data, got %+v", resp)
	}
}

func TestListReceipts_Demo(t *testing.T) {
	h := New(true, "test-secret")
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "+79951239340",
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/receipts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	h.listReceipts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp receiptsListResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Receipts) < 1 {
		t.Fatalf("expected demo receipts, got %+v", resp)
	}
}

func TestTimemachine_ContractShape(t *testing.T) {
	h := New(true, "test-secret")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/timemachine", nil)
	w := httptest.NewRecorder()
	h.timemachine(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp timemachineResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Months) != 60 || len(resp.RealSavings) != 60 {
		t.Fatalf("expected 60 months, got %d", len(resp.Months))
	}
}
