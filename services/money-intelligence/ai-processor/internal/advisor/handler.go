package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	iadvisor "backend_project/internal/advisor"
	iroot "backend_project/internal/auth"
	"backend_project/internal/creditstore"
	"backend_project/internal/llm"
	"backend_project/internal/profile"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	profiles profile.Store
	credits  *creditstore.FileStore
	spending iadvisor.SpendingProvider
	chatStore *iadvisor.ChatStore
	llm      *llm.Client
}

func NewHandler(
	profiles profile.Store,
	credits *creditstore.FileStore,
	spending iadvisor.SpendingProvider,
	chatStore *iadvisor.ChatStore,
	llmClient *llm.Client,
) *Handler {
	return &Handler{
		profiles:  profiles,
		credits:   credits,
		spending:  spending,
		chatStore: chatStore,
		llm:       llmClient,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/ai/plan", h.plan)
	r.Get("/api/v1/ai/diagnosis", h.diagnosis)
	r.Get("/api/v1/ai/chat/history", h.chatHistory)
	r.Delete("/api/v1/ai/chat/history", h.chatClear)
	r.Post("/api/v1/ai/chat", h.chat)
	r.Post("/api/v1/ai/chat/stream", h.chatStream)
}

func (h *Handler) userID(r *http.Request) (string, error) {
	return iroot.UserIDFromRequest(r, "")
}

func (h *Handler) snap(r *http.Request) (iadvisor.Snapshot, error) {
	uid, err := h.userID(r)
	if err != nil {
		return iadvisor.Snapshot{}, err
	}
	return iadvisor.BuildSnapshot(h.profiles, h.credits, h.spending, uid), nil
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

func (h *Handler) chatHistory(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			limit = n
		}
	}
	msgs, err := h.chatMessages(r.Context(), uid, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "history unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, iadvisor.ChatHistoryResponse{Messages: msgs})
}

func (h *Handler) chatClear(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if h.chatStore != nil && h.chatStore.Enabled() {
		if err := h.chatStore.ClearHistory(r.Context(), uid); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "clear failed"})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) chat(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
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

	history, _ := h.resolveHistory(r.Context(), uid, req.History, 10)
	req.History = history

	if h.chatStore != nil && h.chatStore.Enabled() {
		if _, err := h.chatStore.Append(r.Context(), uid, "user", strings.TrimSpace(req.Message), "", nil); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist failed"})
			return
		}
	}

	result := iadvisor.BuildChatReply(snap, req, h.llm)
	stored := iadvisor.EncodeStoredContent(result.Structured, result.Reply)
	resp := iadvisor.ChatResponse{
		Reply:   result.Reply,
		Title:   result.Title,
		Blocks:  result.Blocks,
		Actions: result.Actions,
		Source:  result.Source,
	}

	if h.chatStore != nil && h.chatStore.Enabled() {
		id, err := h.chatStore.Append(r.Context(), uid, "assistant", stored, result.Source, result.Actions)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist failed"})
			return
		}
		resp.ID = id
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) chatStream(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
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

	history, _ := h.resolveHistory(r.Context(), uid, req.History, 10)
	req.History = history

	if h.chatStore != nil && h.chatStore.Enabled() {
		if _, err := h.chatStore.Append(r.Context(), uid, "user", strings.TrimSpace(req.Message), "", nil); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "persist failed"})
			return
		}
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "streaming unsupported"})
		return
	}

	writeSSE(w, flusher, "ready", `{"status":"ok"}`)

	result := iadvisor.BuildChatReplyStream(r.Context(), snap, req, h.llm, func(delta string) error {
		payload, _ := json.Marshal(map[string]string{"text": delta})
		writeSSE(w, flusher, "delta", string(payload))
		return nil
	})

	resp := iadvisor.ChatResponse{
		Reply:   result.Reply,
		Title:   result.Title,
		Blocks:  result.Blocks,
		Actions: result.Actions,
		Source:  result.Source,
	}
	stored := iadvisor.EncodeStoredContent(result.Structured, result.Reply)
	if h.chatStore != nil && h.chatStore.Enabled() {
		id, err := h.chatStore.Append(r.Context(), uid, "assistant", stored, result.Source, result.Actions)
		if err == nil {
			resp.ID = id
		}
	}

	donePayload, _ := json.Marshal(resp)
	writeSSE(w, flusher, "done", string(donePayload))
}

func (h *Handler) resolveHistory(ctx context.Context, uid string, client []iadvisor.ChatMessage, limit int) ([]iadvisor.ChatMessage, error) {
	if h.chatStore != nil && h.chatStore.Enabled() {
		server, err := h.chatStore.RecentTurns(ctx, uid, limit)
		if err == nil && len(server) > 0 {
			return server, nil
		}
	}
	return client, nil
}

func (h *Handler) chatMessages(ctx context.Context, uid string, limit int) ([]iadvisor.StoredChatMessage, error) {
	if h.chatStore != nil && h.chatStore.Enabled() {
		return h.chatStore.ListHistory(ctx, uid, limit)
	}
	return nil, nil
}

func writeSSE(w http.ResponseWriter, flusher http.Flusher, event, data string) {
	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
