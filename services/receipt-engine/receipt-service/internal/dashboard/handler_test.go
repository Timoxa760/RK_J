package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSankey_Demo(t *testing.T) {
	h := New(true)
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

func TestTimemachine_ContractShape(t *testing.T) {
	h := New(true)
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
