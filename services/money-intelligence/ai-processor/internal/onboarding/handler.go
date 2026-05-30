package onboarding

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend_project/internal/llm"
)

type ParseRequest struct {
	Step    string `json:"step"`
	RawText string `json:"raw_text"`
	Locale  string `json:"locale,omitempty"`
}

type ParseResponse struct {
	Parsed  bool   `json:"parsed"`
	Step    Step   `json:"step"`
	Patch   Patch  `json:"patch"`
	Message string `json:"message,omitempty"`
}

// Handler — POST /onboarding/parse (Gemini + regex fallback).
type Handler struct {
	llm *llm.Client
}

func NewHandler(llm *llm.Client) *Handler {
	return &Handler{llm: llm}
}

func (h *Handler) Parse(w http.ResponseWriter, r *http.Request) {
	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	step := Step(strings.TrimSpace(req.Step))
	if step == "" {
		http.Error(w, `{"error":"step required"}`, http.StatusBadRequest)
		return
	}

	result := Parse(step, req.RawText)
	if !result.Parsed && h.llm != nil && h.llm.Enabled() {
		if llmResult, ok := h.parseLLM(r, step, req.RawText); ok {
			result = llmResult
		}
	}

	resp := ParseResponse{
		Parsed: result.Parsed,
		Step:   result.Step,
		Patch:  result.Patch,
	}
	if !result.Parsed {
		resp.Message = "Не удалось разобрать ответ — попробуйте иначе или введите вручную"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) parseLLM(r *http.Request, step Step, text string) (ParseResult, bool) {
	prompt := llm.OnboardingParsePrompt + "\nШаг: " + string(step)
	raw, err := h.llm.Complete(r.Context(), prompt, text)
	if err != nil {
		return ParseResult{}, false
	}
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end <= start {
		return ParseResult{}, false
	}
	var patch Patch
	if err := json.Unmarshal([]byte(raw[start:end+1]), &patch); err != nil {
		return ParseResult{}, false
	}
	return ParseResult{Parsed: true, Step: step, Patch: patch}, true
}
