package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCreateGoal(t *testing.T) {
	h := New()
	body := `{"title":"Отпуск","target_amount":150000,"target_date":"2026-12-01","auto_save_percent":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/goals", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.create(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var g Goal
	json.NewDecoder(w.Body).Decode(&g)
	if g.ID == "" || g.Title != "Отпуск" || g.TargetAmount != 150000 {
		t.Fatalf("unexpected: %+v", g)
	}
}

func TestGetGoal(t *testing.T) {
	h := New()
	createBody := `{"title":"Машина","target_amount":500000}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/goals", strings.NewReader(createBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	h.create(wCreate, reqCreate)
	var created Goal
	json.NewDecoder(wCreate.Body).Decode(&created)

	r := chi.NewRouter()
	h.Register(r)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/goals/"+created.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}
