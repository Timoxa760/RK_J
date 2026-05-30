package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend_project/services/money-intelligence/ai-processor/internal/manual"
	"backend_project/services/money-intelligence/ai-processor/internal/whisper"
)

const maxFormBytes = 11 << 20

// Handler адаптирует legacy-маршруты /receipt/* к pipeline expenses.
type Handler struct {
	whisper   *whisper.Client
	processor *manual.Processor
}

// NewHandler создаёт HTTP-адаптер для front-контракта receipt.
func NewHandler(w *whisper.Client, p *manual.Processor) *Handler {
	return &Handler{whisper: w, processor: p}
}

// ManualCreate обрабатывает POST /receipt/manual.
func (h *Handler) ManualCreate(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromRequest(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req ManualRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	store := strings.TrimSpace(req.Store)
	if store == "" {
		store = "Не указан"
	}
	category := strings.TrimSpace(req.Category)
	if category == "" {
		category = "Прочее"
	}
	date := normalizeDate(req.Date)

	rawText := store
	if req.Amount > 0 {
		rawText = strings.TrimSpace(fmt.Sprintf("%s %.0f руб", store, req.Amount))
	}

	resp, code, err := h.processor.Create(r.Context(), manual.CreateRequest{
		UserID:      userID,
		RawText:     rawText,
		Amount:      req.Amount,
		Category:    category,
		Description: store,
		Date:        date,
		Source:      "manual",
	})
	if err != nil {
		writeProcessorError(w, code, err)
		return
	}

	out := ManualResponse{
		ReceiptID: firstNonEmpty(resp.ID, "receipt-"+time.Now().Format("20060102150405")),
		Store:     store,
		Amount:    resp.Amount,
		Category:  resp.Category,
		Date:      toISO8601(date),
		Status:    "saved",
	}
	if req.Date != "" {
		out.Date = req.Date
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// VoiceCreate обрабатывает POST /receipt/voice (поле audio вместо file).
func (h *Handler) VoiceCreate(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromRequest(r)
	if err != nil {
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

	transcript, err := h.whisper.Transcribe(r.Context(), filename, audio)
	if err != nil {
		http.Error(w, `{"error":"transcription failed"}`, http.StatusServiceUnavailable)
		return
	}

	resp, code, err := h.processor.CreateFromTranscript(r.Context(), userID, transcript, r.FormValue("date"))
	if err != nil {
		writeProcessorError(w, code, err)
		return
	}

	out := mapVoiceResponse(resp, transcript)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func mapVoiceResponse(resp manual.CreateResponse, transcript string) VoiceResponse {
	items := make([]VoiceLineItem, 0, len(resp.Expenses))
	var total float64
	store := ""

	for _, e := range resp.Expenses {
		name := strings.TrimSpace(e.Description)
		if name == "" || len(name) > 80 {
			name = e.Category
		}
		if name == "" {
			name = "Покупка"
		}
		items = append(items, VoiceLineItem{
			Name:     name,
			Price:    e.Amount,
			Quantity: 1,
			Category: e.Category,
		})
		total += e.Amount
		if store == "" && name != "" && len(name) <= 40 {
			store = name
		}
	}
	if store == "" {
		store = "Покупка"
	}
	if total == 0 {
		total = resp.Amount
	}
	if len(items) == 0 && resp.Amount > 0 {
		items = []VoiceLineItem{{
			Name:     resp.Category,
			Price:    resp.Amount,
			Quantity: 1,
			Category: resp.Category,
		}}
		total = resp.Amount
	}

	confidence := 0.75
	if resp.ParsedBy == "gemini" {
		confidence = 0.92
	} else if resp.Parsed {
		confidence = 0.85
	}

	return VoiceResponse{
		ReceiptID:  firstNonEmpty(resp.ID, "receipt-voice-"+time.Now().Format("20060102150405")),
		Store:      store,
		Items:      items,
		Total:      total,
		Category:   resp.Category,
		Confidence: confidence,
	}
}

func writeProcessorError(w http.ResponseWriter, code int, err error) {
	status := http.StatusBadRequest
	if code == 500 {
		status = http.StatusInternalServerError
	}
	http.Error(w, `{"error":"`+err.Error()+`"}`, status)
}

func normalizeDate(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Now().Format("2006-01-02")
	}
	if len(raw) >= 10 {
		if t, err := time.Parse("2006-01-02", raw[:10]); err == nil {
			return t.Format("2006-01-02")
		}
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.Format("2006-01-02")
	}
	return time.Now().Format("2006-01-02")
}

func toISO8601(dateYYYYMMDD string) string {
	if t, err := time.Parse("2006-01-02", dateYYYYMMDD); err == nil {
		return t.UTC().Format(time.RFC3339)
	}
	return time.Now().UTC().Format(time.RFC3339)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
