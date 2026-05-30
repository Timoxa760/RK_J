package dashboard

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"backend_project/internal/auth"
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
		if items == nil {
			items = []receiptListItem{}
		}
		writeJSON(w, receiptsListResponse{Receipts: items})
		return
	}

	if h.expenseFile != nil {
		rows, err := h.loadUserExpensesFile(userID, since)
		if err != nil {
			http.Error(w, `{"error":"failed to load receipts"}`, http.StatusInternalServerError)
			return
		}
		items := fileRowsToReceiptList(rows, limit)
		writeJSON(w, receiptsListResponse{Receipts: items})
		return
	}

	if h.demoMode {
		writeJSON(w, receiptsListResponse{Receipts: demoReceiptList()})
		return
	}

	writeJSON(w, receiptsListResponse{Receipts: []receiptListItem{}})
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

func fileRowsToReceiptList(rows []pgExpenseRow, limit int) []receiptListItem {
	if limit <= 0 {
		limit = 50
	}
	items := make([]receiptListItem, 0, len(rows))
	for i, row := range rows {
		if i >= limit {
			break
		}
		store := row.Description
		if store == "" {
			store = "Не указан"
		}
		items = append(items, receiptListItem{
			ID:       store + "-" + row.ExpenseDate.Format("20060102"),
			Store:    store,
			Amount:   row.Amount,
			Date:     row.ExpenseDate.Format("2006-01-02"),
			Category: row.Category,
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
