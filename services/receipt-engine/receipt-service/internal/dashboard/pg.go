package dashboard

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend_project/internal/auth"
	"backend_project/internal/expensestore"
)

type pgExpenseRow struct {
	Category    string
	Description string
	Amount      float64
	ExpenseDate time.Time
}

func (h *Handler) loadUserExpenses(ctx context.Context, userID string, since time.Time) ([]pgExpenseRow, error) {
	userID = auth.NormalizePhone(userID)
	if h.pgPool != nil {
		rows, err := h.loadUserExpensesPG(ctx, userID, since)
		if err != nil {
			if h.expenseFile != nil {
				return h.loadUserExpensesFile(userID, since)
			}
			return nil, err
		}
		return rows, nil
	}
	if h.expenseFile != nil {
		return h.loadUserExpensesFile(userID, since)
	}
	return nil, fmt.Errorf("expense storage unavailable")
}

func (h *Handler) loadUserExpensesPG(ctx context.Context, userID string, since time.Time) ([]pgExpenseRow, error) {
	if h.pgPool == nil {
		return nil, fmt.Errorf("postgres unavailable")
	}
	rows, err := h.pgPool.Query(ctx, `
		SELECT category, description, amount, expense_date
		FROM manual_expenses
		WHERE user_id = $1 AND expense_date >= $2
		ORDER BY expense_date DESC, created_at DESC
	`, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []pgExpenseRow
	for rows.Next() {
		var row pgExpenseRow
		if err := rows.Scan(&row.Category, &row.Description, &row.Amount, &row.ExpenseDate); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (h *Handler) loadUserExpensesFile(userID string, since time.Time) ([]pgExpenseRow, error) {
	recs, err := h.expenseFile.ListSince(userID, since)
	if err != nil {
		return nil, err
	}
	out := make([]pgExpenseRow, 0, len(recs))
	for _, rec := range recs {
		out = append(out, pgExpenseRow{
			Category:    rec.Category,
			Description: rec.Description,
			Amount:      rec.Amount,
			ExpenseDate: rec.Date,
		})
	}
	return out, nil
}

func (h *Handler) queryCategoriesPG(ctx context.Context, userID string) *categoriesResponse {
	since := startOfMonth(time.Now())
	rows, err := h.loadUserExpenses(ctx, userID, since)
	if err != nil || len(rows) == 0 {
		return nil
	}

	type itemKey struct {
		category string
		name     string
	}
	byCategory := map[string]float64{}
	itemsByCategory := map[string][]catItem{}

	for _, row := range rows {
		category := row.Category
		if category == "" {
			category = "Прочее"
		}
		name := row.Description
		if name == "" {
			name = category
		}
		byCategory[category] += row.Amount
		itemsByCategory[category] = append(itemsByCategory[category], catItem{
			Name:     name,
			Price:    row.Amount,
			Quantity: 1,
			Total:    row.Amount,
		})
	}

	cats := make([]catCategory, 0, len(byCategory))
	for category, total := range byCategory {
		items := itemsByCategory[category]
		cats = append(cats, catCategory{
			Name:  category,
			Total: total,
			Subcategories: []catSubcategory{{
				Name:  category,
				Total: total,
				Items: items,
			}},
		})
	}
	if len(cats) == 0 {
		return nil
	}
	return &categoriesResponse{Categories: cats}
}

func (h *Handler) queryStoresPG(ctx context.Context, userID string) *storesResponse {
	since := startOfMonth(time.Now().AddDate(0, -2, 0))
	rows, err := h.loadUserExpenses(ctx, userID, since)
	if err != nil || len(rows) == 0 {
		return nil
	}

	type storeAgg struct {
		name     string
		purchases int
		total    float64
	}
	agg := map[string]*storeAgg{}
	for _, row := range rows {
		store := row.Description
		if store == "" {
			store = "Не указан"
		}
		s, ok := agg[store]
		if !ok {
			s = &storeAgg{name: store}
			agg[store] = s
		}
		s.purchases++
		s.total += row.Amount
	}

	stores := make([]storeItem, 0, len(agg))
	for _, s := range agg {
		avg := s.total
		if s.purchases > 0 {
			avg = s.total / float64(s.purchases)
		}
		stores = append(stores, storeItem{
			Name:         s.name,
			AvgCheck:     avg,
			Purchases:    s.purchases,
			Total:        s.total,
			ImpulseRatio: 0,
		})
	}
	if len(stores) == 0 {
		return nil
	}
	return &storesResponse{Stores: stores}
}

func buildSankeyFromCategories(cats []catCategory) sankeyResponse {
	var total float64
	for _, c := range cats {
		total += c.Total
	}
	if total <= 0 {
		return sankeyResponse{}
	}

	nodes := []sankeyNode{{Name: "Расходы", Value: total}}
	links := make([]sankeyLink, 0, len(cats))
	for _, c := range cats {
		nodes = append(nodes, sankeyNode{Name: c.Name, Value: c.Total})
		links = append(links, sankeyLink{Source: "Расходы", Target: c.Name, Value: c.Total})
	}
	return sankeyResponse{Nodes: nodes, Links: links}
}

func (h *Handler) queryComparePG(ctx context.Context, userID string, monthCount int) *compareResponse {
	if h.pgPool == nil {
		return nil
	}
	userID = auth.NormalizePhone(userID)
	if monthCount < 2 {
		monthCount = 2
	}

	rows, err := h.pgPool.Query(ctx, `
		SELECT DISTINCT date_trunc('month', expense_date)::date AS m
		FROM manual_expenses
		WHERE user_id = $1
		ORDER BY m DESC
		LIMIT $2
	`, userID, monthCount)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var monthDates []time.Time
	for rows.Next() {
		var m time.Time
		if err := rows.Scan(&m); err != nil {
			return nil
		}
		monthDates = append(monthDates, startOfMonth(m))
	}
	if len(monthDates) == 0 {
		return nil
	}

	if len(monthDates) == 1 {
		prev := startOfMonth(monthDates[0].AddDate(0, -1, 0))
		monthDates = append([]time.Time{prev}, monthDates[0])
	}

	sort.Slice(monthDates, func(i, j int) bool {
		return monthDates[i].Before(monthDates[j])
	})

	months := make([]compareMonth, 0, len(monthDates))
	for _, m := range monthDates {
		cats, err := h.queryMonthCategories(ctx, userID, m)
		if err != nil {
			return nil
		}
		months = append(months, compareMonth{
			Label:      formatMonthLabel(m),
			Categories: cats,
		})
	}

	resp := &compareResponse{Months: months}
	if len(months) >= 2 {
		resp.Insights = buildCompareInsights(months[len(months)-2], months[len(months)-1])
	}
	return resp
}

func (h *Handler) queryMonthCategories(ctx context.Context, userID string, monthStart time.Time) ([]compareCategory, error) {
	monthEnd := monthStart.AddDate(0, 1, 0)
	rows, err := h.pgPool.Query(ctx, `
		SELECT category, sum(amount) AS total
		FROM manual_expenses
		WHERE user_id = $1
		  AND expense_date >= $2
		  AND expense_date < $3
		GROUP BY category
		ORDER BY total DESC
	`, userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []compareCategory
	for rows.Next() {
		var name string
		var total float64
		if err := rows.Scan(&name, &total); err != nil {
			return nil, err
		}
		cats = append(cats, compareCategory{Name: name, Total: total})
	}
	if cats == nil {
		cats = []compareCategory{}
	}
	return cats, rows.Err()
}

func buildCompareInsights(prev, curr compareMonth) compareInsights {
	if len(curr.Categories) == 0 {
		return compareInsights{}
	}
	best := curr.Categories[0]
	prevTotal := findCategoryTotal(prev.Categories, best.Name)
	delta := int(best.Total - prevTotal)
	pct := 0
	if prevTotal > 0 {
		pct = int((best.Total - prevTotal) / prevTotal * 100)
	}
	return compareInsights{
		BiggestChange: compareChange{
			Category:     best.Name,
			Delta:        delta,
			DeltaPercent: pct,
		},
	}
}

var ruMonths = []string{
	"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
	"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
}

func formatMonthLabel(t time.Time) string {
	m := int(t.Month())
	if m >= 1 && m <= 12 {
		return fmt.Sprintf("%s %d", ruMonths[m-1], t.Year())
	}
	return t.Format("2006-01")
}

func findCategoryTotal(cats []compareCategory, name string) float64 {
	for _, c := range cats {
		if c.Name == name {
			return c.Total
		}
	}
	return 0
}

func (h *Handler) queryTimemachinePG(ctx context.Context, userID string) *timemachineResponse {
	if h.pgPool == nil {
		return nil
	}
	userID = auth.NormalizePhone(userID)
	rows, err := h.pgPool.Query(ctx, `
		SELECT to_char(date_trunc('month', expense_date), 'YYYY-MM') AS month,
		       sum(amount) AS spent
		FROM manual_expenses
		WHERE user_id = $1
		GROUP BY 1
		ORDER BY 1 ASC
		LIMIT 60
	`, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	type monthRow struct {
		month string
		spent float64
	}
	var data []monthRow
	for rows.Next() {
		var row monthRow
		if err := rows.Scan(&row.month, &row.spent); err != nil {
			return nil
		}
		data = append(data, row)
	}
	if len(data) == 0 {
		return nil
	}

	months := make([]string, len(data))
	real := make([]float64, len(data))
	opt := make([]float64, len(data))
	var cumulative float64
	for i, row := range data {
		months[i] = row.month
		cumulative += row.spent
		real[i] = cumulative
		opt[i] = cumulative * 0.85
	}

	return &timemachineResponse{
		Months:           months,
		RealSavings:      real,
		OptimizedSavings: opt,
		DifferenceFinal:  int(opt[len(opt)-1] - real[len(real)-1]),
	}
}

func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// SetPostgres подключает пул PostgreSQL для чтения manual_expenses.
func (h *Handler) SetPostgres(pool *pgxpool.Pool) {
	h.pgPool = pool
}

// SetExpenseFile подключает JSON-хранилище расходов для локальной разработки.
func (h *Handler) SetExpenseFile(store *expensestore.FileStore) {
	h.expenseFile = store
}
