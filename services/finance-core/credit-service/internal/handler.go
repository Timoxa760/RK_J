package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	iroot "backend_project/internal/auth"
	"backend_project/internal/creditstore"
	"backend_project/internal/onlysq"
	"backend_project/internal/pdfextract"
	"backend_project/internal/profile"
	"backend_project/internal/rates"

	"github.com/go-chi/chi/v5"
)

const creditScanPrompt = `Из текста кредитного договора извлеки условия. Ответ — только JSON:
{"amount": 0, "rate": 0, "term_months": 0, "monthly_payment": 0, "bank": ""}
amount и monthly_payment в рублях, rate — годовая ставка в процентах.`

// Handler — credits из PDF scan only.
type Handler struct {
	credits  *creditstore.FileStore
	profiles *profile.FileStore
	rates    *rates.Client
	llm      *onlysq.Client
}

func NewHandler(credits *creditstore.FileStore, profiles *profile.FileStore, llm *onlysq.Client) *Handler {
	return &Handler{
		credits:  credits,
		profiles: profiles,
		rates:    rates.NewClient(),
		llm:      llm,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/credits/dashboard", h.dashboard)
	r.Post("/api/v1/credits/scan", h.scan)
	r.Delete("/api/v1/credits/{id}", h.delete)
}

func (h *Handler) userID(r *http.Request) (string, error) {
	return iroot.UserIDFromRequest(r, "")
}

func (h *Handler) monthlyIncome(uid string) float64 {
	if h.profiles == nil {
		return 0
	}
	p, err := h.profiles.Get(uid)
	if err != nil || p.SkippedIncome {
		return 0
	}
	return p.ActiveIncome + p.PassiveIncome
}

func (h *Handler) emergencyFund(uid string) float64 {
	if h.profiles == nil {
		return 0
	}
	p, err := h.profiles.Get(uid)
	if err != nil || p.SkippedCushion {
		return 0
	}
	return p.EmergencyFund
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	dash := h.credits.Dashboard(uid, h.monthlyIncome(uid), h.emergencyFund(uid))
	writeJSON(w, http.StatusOK, dash)
}

func (h *Handler) scan(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file required"})
		return
	}
	defer file.Close()
	if header != nil && !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "pdf required"})
		return
	}
	body, err := io.ReadAll(io.LimitReader(file, 4<<20))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "read file"})
		return
	}
	hash := sha256.Sum256(body)
	hashStr := hex.EncodeToString(hash[:])

	text, extractErr := pdfextract.TextFromPDF(body)
	parsed, confidence, err := h.parseContract(r.Context(), text, extractErr)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	bench, _ := h.rates.Fetch(r.Context(), "consumer", parsed.Amount, parsed.TermMonths)
	vs := creditstore.CompareRate(parsed.Rate, bench.BenchmarkRate)

	credit, err := h.credits.Add(uid, creditstore.Credit{
		Bank:           parsed.Bank,
		Amount:         parsed.Amount,
		Rate:           parsed.Rate,
		TermMonths:     parsed.TermMonths,
		Remaining:      parsed.Amount,
		MonthlyPayment: parsed.MonthlyPayment,
		BenchmarkRate:  bench.BenchmarkRate,
		RateVsMarket:   vs,
		NextPayment:    creditstore.FormatNextPayment(),
		SourceFileHash: hashStr,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"parsed": map[string]any{
			"amount":          parsed.Amount,
			"rate":            parsed.Rate,
			"term_months":     parsed.TermMonths,
			"monthly_payment": parsed.MonthlyPayment,
			"bank":            parsed.Bank,
		},
		"benchmark_rate":  bench.BenchmarkRate,
		"rate_vs_market":  vs,
		"confidence":      confidence,
		"credit_id":       credit.ID,
	})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	id := chi.URLParam(r, "id")
	if err := h.credits.Delete(uid, id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

type scanParsed struct {
	Amount         float64 `json:"amount"`
	Rate           float64 `json:"rate"`
	TermMonths     int     `json:"term_months"`
	MonthlyPayment float64 `json:"monthly_payment"`
	Bank           string  `json:"bank"`
}

func (h *Handler) parseContract(ctx context.Context, text string, extractErr error) (scanParsed, float64, error) {
	confidence := 0.75
	trimmed := strings.TrimSpace(text)
	if h.llm != nil && h.llm.Enabled() && len(trimmed) > 20 {
		raw, err := h.llm.Complete(ctx, creditScanPrompt, trimmed)
		if err == nil {
			var p scanParsed
			if json.Unmarshal([]byte(extractJSON(raw)), &p) == nil {
				if err := creditstore.ValidateScan(p.Amount, p.Rate, p.TermMonths); err == nil {
					return p, 0.87, nil
				}
			}
		}
	}
	if demoMode() {
		p := scanParsed{
			Amount:         envFloat("CREDIT_SCAN_DEMO_AMOUNT", 1200000),
			Rate:           envFloat("CREDIT_SCAN_DEMO_RATE", 14.5),
			TermMonths:     envInt("CREDIT_SCAN_DEMO_TERM", 36),
			MonthlyPayment: envFloat("CREDIT_SCAN_DEMO_PAYMENT", 42000),
			Bank:           envStr("CREDIT_SCAN_DEMO_BANK", "Т-Банк"),
		}
		if err := creditstore.ValidateScan(p.Amount, p.Rate, p.TermMonths); err == nil {
			return p, confidence, nil
		}
	}
	if extractErr != nil {
		return scanParsed{}, 0, fmt.Errorf("could not extract text from pdf: %v", extractErr)
	}
	return scanParsed{}, 0, fmt.Errorf("could not parse contract fields from pdf")
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func demoMode() bool {
	return os.Getenv("DEMO_MODE") == "true"
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
