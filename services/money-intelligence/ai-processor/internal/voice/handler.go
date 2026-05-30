package voice

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"backend_project/internal/auth"
	"backend_project/services/money-intelligence/ai-processor/internal/manual"
	"backend_project/services/money-intelligence/ai-processor/internal/whisper"
)

const maxFormBytes = 11 << 20

// Handler обрабатывает POST /expenses/voice.
type Handler struct {
	whisper   *whisper.Client
	processor *manual.Processor
}

// NewHandler создаёт HTTP-обработчик голосового ввода.
func NewHandler(w *whisper.Client, p *manual.Processor) *Handler {
	return &Handler{whisper: w, processor: p}
}

// Create принимает multipart audio и возвращает JSON с transcript и advice.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.whisper.Enabled() {
		http.Error(w, `{"error":"whisper unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	if err := r.ParseMultipartForm(maxFormBytes); err != nil {
		http.Error(w, `{"error":"invalid multipart"}`, http.StatusBadRequest)
		return
	}
	userID := strings.TrimSpace(r.FormValue("user_id"))
	if userID == "" {
		http.Error(w, `{"error":"user_id required"}`, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"file required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	audio, err := io.ReadAll(io.LimitReader(file, maxFormBytes))
	if err != nil {
		http.Error(w, `{"error":"read file failed"}`, http.StatusBadRequest)
		return
	}

	filename := "voice.webm"
	if header != nil && header.Filename != "" {
		filename = header.Filename
	}

	transcript, err := h.whisper.Transcribe(r.Context(), filename, audio)
	if err != nil {
		http.Error(w, `{"error":"transcription failed"}`, http.StatusServiceUnavailable)
		return
	}

	resp, code, err := h.processor.CreateFromTranscript(r.Context(), userID, transcript, r.FormValue("date"))
	if err != nil {
		status := http.StatusBadRequest
		if code == 500 {
			status = http.StatusInternalServerError
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, status)
		return
	}
	resp.Transcript = transcript

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// TranscribeResponse — ответ POST /voice/transcribe (только текст, без сохранения расхода).
type TranscribeResponse struct {
	Text string `json:"text"`
}

// Transcribe принимает multipart audio и возвращает распознанный текст для онбординга.
func (h *Handler) Transcribe(w http.ResponseWriter, r *http.Request) {
	if _, err := auth.UserIDFromRequest(r, auth.JWTSecret()); err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	if !h.whisper.Enabled() {
		http.Error(w, `{"error":"whisper unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	if err := r.ParseMultipartForm(maxFormBytes); err != nil {
		http.Error(w, `{"error":"invalid multipart"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		file, header, err = r.FormFile("file")
	}
	if err != nil {
		http.Error(w, `{"error":"audio required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	audio, err := io.ReadAll(io.LimitReader(file, maxFormBytes))
	if err != nil {
		http.Error(w, `{"error":"read audio failed"}`, http.StatusBadRequest)
		return
	}

	filename := "recording.webm"
	if header != nil && header.Filename != "" {
		filename = header.Filename
	}

	text, err := h.whisper.Transcribe(r.Context(), filename, audio)
	if err != nil {
		http.Error(w, `{"error":"transcription failed"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(TranscribeResponse{Text: text})
}
