package dashboard

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"backend_project/internal/auth"
	"backend_project/internal/expensestore"
)

type receiptListItem struct {
	ID           string  `json:"id"`
	Store        string  `json:"store"`
	Amount       float64 `json:"amount"`
	Date         string  `json:"date"`
	Category     string  `json:"category,omitempty"`
	ImpulseCount int     `json:"impulse_count,omitempty"`
}

type receiptsListResponse struct {
	Receipts []receiptListItem `json:"receipts"`
}

// listReceipts возвращает ленту расходов пользователя из manual_expenses.
func (h *Handler) listReceipts(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromRequest(r, h.jwtSecret)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	limit := parseReceiptLimit(r)
	since := time.Now().AddDate(-2, 0, 0)

	if h.pgPool != nil {
		items, err := h.queryReceiptsPG(r.Context(), userID, since, limit)
		if err != nil {
			http.Error(w, `{"error":"failed to load receipts"}`, http.StatusInternalServerError)
			return
		}
		if len(items) > 0 {
			writeJSON(w, receiptsListResponse{Receipts: items})
			return
		}
	}

	if h.expenseFile != nil {
		recs, err := h.expenseFile.ListSince(userID, since)
		if err != nil {
			http.Error(w, `{"error":"failed to load receipts"}`, http.StatusInternalServerError)
			return
		}
		items := recordsToReceiptList(recs, limit)
		writeJSON(w, receiptsListResponse{Receipts: items})
		return
	}

	if h.demoMode {
		writeJSON(w, receiptsListResponse{Receipts: demoReceiptList()})
		return
	}

	writeJSON(w, receiptsListResponse{Receipts: []receiptListItem{}})
}

// deleteReceipt удаляет трату пользователя из manual_expenses или file store.
func (h *Handler) deleteReceipt(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromRequest(r, h.jwtSecret)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	if h.demoMode {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	userID = auth.NormalizePhone(userID)

	if h.pgPool != nil {
		tag, err := h.pgPool.Exec(r.Context(), `
			DELETE FROM manual_expenses
			WHERE id = $1::uuid AND user_id = $2
		`, id, userID)
		if err != nil {
			http.Error(w, `{"error":"failed to delete receipt"}`, http.StatusInternalServerError)
			return
		}
		if tag.RowsAffected() > 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	if h.expenseFile != nil {
		ok, err := h.expenseFile.Delete(userID, id)
		if err != nil {
			http.Error(w, `{"error":"failed to delete receipt"}`, http.StatusInternalServerError)
			return
		}
		if ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
}

func (h *Handler) queryReceiptsPG(ctx context.Context, userID string, since time.Time, limit int) ([]receiptListItem, error) {
	userID = auth.NormalizePhone(userID)
	rows, err := h.pgPool.Query(ctx, `
		SELECT id::text, description, amount, category, expense_date, source
		FROM manual_expenses
		WHERE user_id = $1 AND expense_date >= $2
		ORDER BY expense_date DESC, created_at DESC
		LIMIT $3
	`, userID, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []receiptListItem
	for rows.Next() {
		var id, description, category, source string
		var amount float64
		var expenseDate time.Time
		if err := rows.Scan(&id, &description, &amount, &category, &expenseDate, &source); err != nil {
			return nil, err
		}
		store := description
		if store == "" {
			store = "Не указан"
		}
		items = append(items, receiptListItem{
			ID:       id,
			Store:    store,
			Amount:   amount,
			Date:     expenseDate.Format("2006-01-02"),
			Category: category,
		})
	}
	return items, rows.Err()
}

func recordsToReceiptList(recs []expensestore.Record, limit int) []receiptListItem {
	if limit <= 0 {
		limit = 50
	}
	sort.Slice(recs, func(i, j int) bool {
		if recs[i].Date.Equal(recs[j].Date) {
			return recs[i].CreatedAt.After(recs[j].CreatedAt)
		}
		return recs[i].Date.After(recs[j].Date)
	})
	items := make([]receiptListItem, 0, len(recs))
	for i, rec := range recs {
		if i >= limit {
			break
		}
		store := rec.Description
		if store == "" {
			store = "Не указан"
		}
		items = append(items, receiptListItem{
			ID:       rec.ID,
			Store:    store,
			Amount:   rec.Amount,
			Date:     rec.Date.Format("2006-01-02"),
			Category: rec.Category,
		})
	}
	return items
}

func parseReceiptLimit(r *http.Request) int {
	const defaultLimit = 50
	v := r.URL.Query().Get("limit")
	if v == "" {
		return defaultLimit
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 || n > 200 {
		return defaultLimit
	}
	return n
}

func demoReceiptList() []receiptListItem {
	return []receiptListItem{
		{ID: "demo-1", Store: "Пятёрочка", Amount: 1240, Date: "2026-05-28", Category: "Продукты"},
		{ID: "demo-2", Store: "Яндекс.Такси", Amount: 450, Date: "2026-05-27", Category: "Транспорт"},
	}
}
