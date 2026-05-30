package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCreateAndLeaderboard(t *testing.T) {
	r := chi.NewRouter()
	New().Register(r)

	body := bytes.NewBufferString(`{"type":"least_spend","title":"test","duration_days":7}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", body)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}

	var created createResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/challenges/"+created.ID+"/leaderboard", nil)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("leaderboard status=%d body=%s", rec2.Code, rec2.Body.String())
	}
}
