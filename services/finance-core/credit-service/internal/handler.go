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
	"backend_project/internal/creditparse"
	"backend_project/internal/creditstore"
	"backend_project/internal/mortgage"
	"backend_project/internal/llm"
	"backend_project/internal/pdfextract"
	"backend_project/internal/profile"
	"backend_project/internal/rates"

	"github.com/go-chi/chi/v5"
)

const creditScanPrompt = `Из текста договора кредита, займа, займа наличными или микрозайма извлеки условия. Ответ — только JSON:
{"amount": 0, "rate": 0, "term_months": 0, "monthly_payment": 0, "bank": ""}
amount и monthly_payment в рублях, rate — годовая ставка или ПСК в процентах, term_months — срок в месяцах (если в договоре дни — переведи в месяцы, минимум 1).
bank — название банка/МФО или тип продукта («Займ», «Кредит наличными»).`

// Handler — credits из PDF scan only.
type Handler struct {
	credits  *creditstore.FileStore
	profiles profile.Store
	rates    *rates.Client
	llm      *llm.Client
}

func NewHandler(credits *creditstore.FileStore, profiles profile.Store, llmClient *llm.Client) *Handler {
	return &Handler{
		credits:  credits,
		profiles: profiles,
		rates:    rates.NewClient(),
		llm:      llmClient,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/credits/dashboard", h.dashboard)
	r.Post("/api/v1/credits/scan", h.scan)
	r.Post("/api/v1/credits/mortgage/analyze", h.analyzeMortgage)
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
	dash = enrichDashboardPayments(dash)
	writeJSON(w, http.StatusOK, dash)
}

func enrichDashboardPayments(dash creditstore.Dashboard) creditstore.Dashboard {
	var monthlyPayments float64
	for i := range dash.Credits {
		payment := dash.Credits[i].MonthlyPayment
		if payment <= 0 {
			payment = creditparse.EstimateMonthlyPayment(
				dash.Credits[i].Amount,
				dash.Credits[i].Rate,
				dash.Credits[i].TermMonths,
			)
			dash.Credits[i].MonthlyPayment = payment
		}
		monthlyPayments += payment
	}
	if monthlyPayments > 0 {
		dash.MonthlyPayments = monthlyPayments
		if dash.MonthlyIncome > 0 {
			dash.DTI = monthlyPayments / dash.MonthlyIncome * 100
		}
		if dash.Savings > 0 {
			dash.StressTestMonths = dash.Savings / monthlyPayments
		}
	}
	return dash
}

type mortgageAnalyzeRequest struct {
	MortgageAmount     float64 `json:"mortgage_amount"`
	MonthlyIncome      float64 `json:"monthly_income"`
	Savings            float64 `json:"savings"`
	ExistingDTI        float64 `json:"existing_dti"`
	StressTestMonths   float64 `json:"stress_test_months"`
}

func (h *Handler) analyzeMortgage(w http.ResponseWriter, r *http.Request) {
	uid, err := h.userID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	var req mortgageAnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MortgageAmount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "mortgage_amount required"})
		return
	}

	income := req.MonthlyIncome
	savings := req.Savings
	dti := req.ExistingDTI
	stress := req.StressTestMonths
	existingPayments := float64(0)

	if income <= 0 {
		income = h.monthlyIncome(uid)
	}
	if savings <= 0 {
		savings = h.emergencyFund(uid)
	}
	if dti <= 0 && income > 0 {
		dash := h.credits.Dashboard(uid, income, savings)
		dti = dash.DTI
		stress = dash.StressTestMonths
		existingPayments = dash.MonthlyPayments
	}

	bench, _ := h.rates.Fetch(r.Context(), "mortgage", req.MortgageAmount, 240)
	resp := mortgage.BuildBreakdown(mortgage.AnalyzeInput{
		MortgageAmount:   req.MortgageAmount,
		MonthlyIncome:    income,
		Savings:          savings,
		ExistingDTI:      dti,
		StressTestMonths: stress,
		BenchmarkRate:    bench.BenchmarkRate,
		ExistingPayments: existingPayments,
	})
	writeJSON(w, http.StatusOK, resp)
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
	filename := ""
	if header != nil {
		filename = header.Filename
	}
	parsed, confidence, err := h.parseContract(r.Context(), text, extractErr, filename)
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
			"amount":             parsed.Amount,
			"rate":               parsed.Rate,
			"term_months":        parsed.TermMonths,
			"monthly_payment":    parsed.MonthlyPayment,
			"bank":               parsed.Bank,
			"payment_estimated":  parsed.PaymentEstimated,
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
	Amount           float64 `json:"amount"`
	Rate             float64 `json:"rate"`
	TermMonths       int     `json:"term_months"`
	MonthlyPayment   float64 `json:"monthly_payment"`
	Bank             string  `json:"bank"`
	PaymentEstimated bool    `json:"payment_estimated,omitempty"`
}

func (h *Handler) parseContract(ctx context.Context, text string, extractErr error, filename string) (scanParsed, float64, error) {
	confidence := 0.75
	trimmed := strings.TrimSpace(text)
	if h.llm != nil && h.llm.Enabled() && len(trimmed) > 20 {
		raw, err := h.llm.Complete(ctx, creditScanPrompt, trimmed)
		if err == nil {
			var p scanParsed
			if json.Unmarshal([]byte(extractJSON(raw)), &p) == nil {
				if err := creditstore.ValidateScan(p.Amount, p.Rate, p.TermMonths); err == nil {
					return normalizeParsed(p), 0.87, nil
				}
			}
		}
	}
	if len(trimmed) > 20 {
		if fields, ok := creditparse.ParseFromText(trimmed); ok {
			if err := creditstore.ValidateScan(fields.Amount, fields.Rate, fields.TermMonths); err == nil {
				return scanParsed{
					Amount:           fields.Amount,
					Rate:             fields.Rate,
					TermMonths:       fields.TermMonths,
					MonthlyPayment:   fields.MonthlyPayment,
					Bank:             fields.Bank,
					PaymentEstimated: fields.PaymentEstimated,
				}, 0.72, nil
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
		return scanParsed{}, 0, fmt.Errorf("scan_no_text")
	}
	if creditparse.IsGeneralConditionsTemplate(trimmed, filename) {
		return scanParsed{}, 0, fmt.Errorf("scan_general_conditions")
	}
	return scanParsed{}, 0, fmt.Errorf("scan_parse_failed")
}

func normalizeParsed(p scanParsed) scanParsed {
	if strings.TrimSpace(p.Bank) == "" {
		p.Bank = "Кредит"
	}
	if p.MonthlyPayment <= 0 && p.Amount > 0 && p.TermMonths > 0 {
		p.MonthlyPayment = creditparse.EstimateMonthlyPayment(p.Amount, p.Rate, p.TermMonths)
		p.PaymentEstimated = true
	}
	return p
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
