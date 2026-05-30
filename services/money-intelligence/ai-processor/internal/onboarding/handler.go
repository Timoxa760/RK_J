package onboarding

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ParseRequest struct {
	Step    string `json:"step"`
	RawText string `json:"raw_text"`
	Locale  string `json:"locale,omitempty"`
}

type ParseResponse struct {
	Parsed bool        `json:"parsed"`
	Step   Step        `json:"step"`
	Patch  Patch       `json:"patch"`
	Message string     `json:"message,omitempty"`
}

// Handler — POST /onboarding/parse (regex + будущий OnlySQ).
type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

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
