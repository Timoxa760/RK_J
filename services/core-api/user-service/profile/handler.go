package profile

import (
	"encoding/json"
	"net/http"

	iroot "backend_project/internal/auth"
	iprofile "backend_project/internal/profile"

	"github.com/go-chi/chi/v5"
)

// Handler — GET/PATCH profile, POST onboarding complete.
type Handler struct {
	store *iprofile.FileStore
}

func NewHandler(store *iprofile.FileStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/users/me/profile", h.get)
	r.Patch("/api/v1/users/me/profile", h.patch)
	r.Post("/api/v1/users/me/onboarding/complete", h.complete)
}

func (h *Handler) userID(r *http.Request) (string, bool) {
	id, err := iroot.UserIDFromRequest(r, "")
	return id, err == nil && id != ""
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	uid, ok := h.userID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	p, err := h.store.Get(uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) patch(w http.ResponseWriter, r *http.Request) {
	uid, ok := h.userID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	cur, err := h.store.Get(uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	next, err := applyPatchMap(cur, raw)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.store.Save(uid, next); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, next)
}

func (h *Handler) complete(w http.ResponseWriter, r *http.Request) {
	uid, ok := h.userID(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	cur, err := h.store.Get(uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	cur.OnboardingCompleted = true
	if err := h.store.Save(uid, cur); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"onboarding_completed": true})
}

func applyPatchMap(cur iprofile.FinancialProfile, raw map[string]json.RawMessage) (iprofile.FinancialProfile, error) {
	data, err := json.Marshal(cur)
	if err != nil {
		return cur, err
	}
	var base map[string]json.RawMessage
	if err := json.Unmarshal(data, &base); err != nil {
		return cur, err
	}
	for k, v := range raw {
		base[k] = v
	}
	merged, err := json.Marshal(base)
	if err != nil {
		return cur, err
	}
	var out iprofile.FinancialProfile
	if err := json.Unmarshal(merged, &out); err != nil {
		return cur, err
	}
	return out, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
