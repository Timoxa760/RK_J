package manual

import (
	"context"
	"fmt"
	"sync"
	"time"

	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

type demoStorage struct {
	mu     sync.Mutex
	byUser map[string][]*Expense
}

func newDemoStorage() *demoStorage {
	return &demoStorage{byUser: make(map[string][]*Expense)}
}

// SaveExpenses сохраняет траты в памяти (DEMO_MODE).
func (d *demoStorage) SaveExpenses(ctx context.Context, userID, source string, date time.Time, items []expense.Item) ([]ExpenseDTO, error) {
	out := make([]ExpenseDTO, 0, len(items))
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, it := range items {
		e := &Expense{
			ID:          fmt.Sprintf("demo-%d-%d", time.Now().UnixNano(), i),
			UserID:      userID,
			Amount:      it.Amount,
			Category:    it.Category,
			Description: it.Description,
			Source:      source,
			Date:        date,
			CreatedAt:   time.Now(),
		}
		d.byUser[userID] = append(d.byUser[userID], e)
		out = append(out, ExpenseDTO{
			ID:          e.ID,
			Amount:      e.Amount,
			Category:    e.Category,
			Description: e.Description,
		})
	}
	return out, nil
}
