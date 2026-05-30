package manual

import (
	"context"
	"time"

	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)

type repoStorage struct {
	repo *Repo
}

func newRepoStorage(repo *Repo) *repoStorage {
	return &repoStorage{repo: repo}
}

// NewRepoStorage экспортирует адаптер Storage для тестов.
func NewRepoStorage(repo *Repo) Storage {
	return newRepoStorage(repo)
}

// SaveExpenses сохраняет траты в PostgreSQL и ClickHouse.
func (s *repoStorage) SaveExpenses(ctx context.Context, userID, source string, date time.Time, items []expense.Item) ([]ExpenseDTO, error) {
	out := make([]ExpenseDTO, 0, len(items))
	for _, it := range items {
		e := &Expense{
			UserID:      userID,
			Amount:      it.Amount,
			Category:    it.Category,
			Description: it.Description,
			Source:      source,
			Date:        date,
		}
		if err := s.repo.Insert(ctx, e); err != nil {
			return nil, err
		}
		out = append(out, ExpenseDTO{
			ID:          e.ID,
			Amount:      e.Amount,
			Category:    e.Category,
			Description: e.Description,
		})
	}
	return out, nil
}
