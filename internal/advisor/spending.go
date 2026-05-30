package advisor

import (
	"context"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend_project/internal/auth"
	"backend_project/internal/expensestore"
)

type CategorySpend struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Share  float64 `json:"share"`
}

type SpendingSummary struct {
	MonthTotal  float64         `json:"month_total"`
	Categories  []CategorySpend `json:"categories"`
	RecentCount int             `json:"recent_count"`
}

type SpendingProvider interface {
	MonthSummary(ctx context.Context, userID string) SpendingSummary
}

type PGSpendingProvider struct {
	pool *pgxpool.Pool
	file *expensestore.FileStore
}

func NewPGSpendingProvider(pool *pgxpool.Pool, file *expensestore.FileStore) *PGSpendingProvider {
	return &PGSpendingProvider{pool: pool, file: file}
}

func (p *PGSpendingProvider) MonthSummary(ctx context.Context, userID string) SpendingSummary {
	userID = auth.NormalizePhone(userID)
	since := startOfMonth(time.Now())

	var rows []expenseRow
	if p.pool != nil {
		if pgRows, err := p.loadPG(ctx, userID, since); err == nil {
			rows = pgRows
		} else if p.file != nil {
			rows = p.loadFile(userID, since)
		}
	} else if p.file != nil {
		rows = p.loadFile(userID, since)
	}
	return aggregateSpending(rows)
}

type expenseRow struct {
	Category string
	Amount   float64
}

func (p *PGSpendingProvider) loadPG(ctx context.Context, userID string, since time.Time) ([]expenseRow, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT category, amount
		FROM manual_expenses
		WHERE user_id = $1 AND expense_date >= $2
	`, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []expenseRow
	for rows.Next() {
		var row expenseRow
		if err := rows.Scan(&row.Category, &row.Amount); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (p *PGSpendingProvider) loadFile(userID string, since time.Time) []expenseRow {
	recs, err := p.file.ListSince(userID, since)
	if err != nil {
		return nil
	}
	out := make([]expenseRow, 0, len(recs))
	for _, rec := range recs {
		out = append(out, expenseRow{Category: rec.Category, Amount: rec.Amount})
	}
	return out
}

func aggregateSpending(rows []expenseRow) SpendingSummary {
	if len(rows) == 0 {
		return SpendingSummary{}
	}

	totals := map[string]float64{}
	for _, row := range rows {
		if row.Amount <= 0 {
			continue
		}
		name := NormalizeCategory(row.Category)
		totals[name] += row.Amount
	}

	grandTotal := 0.0
	for _, amount := range totals {
		grandTotal += amount
	}
	if grandTotal <= 0 {
		return SpendingSummary{RecentCount: len(rows)}
	}

	items := make([]CategorySpend, 0, len(totals))
	for name, amount := range totals {
		items = append(items, CategorySpend{
			Name:   name,
			Amount: amount,
			Share:  amount / grandTotal,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Amount > items[j].Amount
	})

	const maxItems = 8
	if len(items) > maxItems {
		head := items[:maxItems-1]
		tail := 0.0
		for _, item := range items[maxItems-1:] {
			tail += item.Amount
		}
		merged := false
		for i := range head {
			if head[i].Name == "Прочие расходы" {
				head[i].Amount += tail
				head[i].Share = head[i].Amount / grandTotal
				merged = true
				break
			}
		}
		if !merged {
			head = append(head, CategorySpend{
				Name:   "Прочие расходы",
				Amount: tail,
				Share:  tail / grandTotal,
			})
		}
		items = head
	}

	return SpendingSummary{
		MonthTotal:  grandTotal,
		Categories:  items,
		RecentCount: len(rows),
	}
}

func startOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func TopCategory(s SpendingSummary) *CategorySpend {
	if len(s.Categories) == 0 {
		return nil
	}
	top := s.Categories[0]
	return &top
}
