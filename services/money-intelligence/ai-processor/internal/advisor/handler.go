package advisor

import (
	"encoding/json"
	"net/http"
	"strings"

	iadvisor "backend_project/internal/advisor"
	iroot "backend_project/internal/auth"
	"backend_project/internal/creditstore"
	"backend_project/internal/onlysq"
	"backend_project/internal/profile"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	profiles *profile.FileStore
	credits  *creditstore.FileStore
	llm      *onlysq.Client
}

func NewHandler(profiles *profile.FileStore, credits *creditstore.FileStore, llm *onlysq.Client) *Handler {
	return &Handler{profiles: profiles, credits: credits, llm: llm}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/ai/plan", h.plan)
	r.Get("/api/v1/ai/diagnosis", h.diagnosis)
	r.Post("/api/v1/ai/chat", h.chat)
}

func (h *Handler) userID(r *http.Request) (string, error) {
	return iroot.UserIDFromRequest(r, "")
}

func (h *Handler) snap(r *http.Request) (iadvisor.Snapshot, error) {
	uid, err := h.userID(r)
	if err != nil {
		return iadvisor.Snapshot{}, err
	}
	return iadvisor.BuildSnapshot(h.profiles, h.credits, uid), nil
}

func (h *Handler) plan(w http.ResponseWriter, r *http.Request) {
	snap, err := h.snap(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, iadvisor.BuildPlanResponse(snap, h.llm))
}

func (h *Handler) diagnosis(w http.ResponseWriter, r *http.Request) {
	snap, err := h.snap(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, iadvisor.BuildDiagnosis(snap, h.llm))
}

func (h *Handler) chat(w http.ResponseWriter, r *http.Request) {
	snap, err := h.snap(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	var req iadvisor.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "message required"})
		return
	}
	writeJSON(w, http.StatusOK, iadvisor.ChatResponse{Reply: iadvisor.BuildChatReply(snap, req, h.llm)})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
