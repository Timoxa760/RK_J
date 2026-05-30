package manual

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend_project/internal/expensestore"
	"backend_project/services/money-intelligence/ai-processor/internal/expense"
)
type fallbackStorage struct {
	primary Storage
	file    *expensestore.FileStore
}

func newFallbackStorage(primary Storage, file *expensestore.FileStore) Storage {
	return &fallbackStorage{primary: primary, file: file}
}

// NewFallbackStorage экспортирует fallback Storage для main.
func NewFallbackStorage(primary Storage, file *expensestore.FileStore) Storage {
	return newFallbackStorage(primary, file)
}

func (s *fallbackStorage) SaveExpenses(ctx context.Context, userID, source string, date time.Time, items []expense.Item) ([]ExpenseDTO, error) {
	if s.primary != nil {
		out, err := s.primary.SaveExpenses(ctx, userID, source, date, items)
		if err == nil {
			s.mirrorToFile(userID, source, date, items)
			return out, nil
		}
		log.Printf("expense pg save failed, fallback to file: %v", err)
	}
	return s.saveToFile(userID, source, date, items)
}

func (s *fallbackStorage) mirrorToFile(userID, source string, date time.Time, items []expense.Item) {
	if s.file == nil {
		return
	}
	fileItems := make([]expensestore.Item, len(items))
	for i, it := range items {
		fileItems[i] = expensestore.Item{
			Amount:      it.Amount,
			Category:    it.Category,
			Description: it.Description,
		}
	}
	if _, err := s.file.Append(userID, source, date, fileItems); err != nil {
		log.Printf("expense file mirror failed: %v", err)
	}
}

func (s *fallbackStorage) saveToFile(userID, source string, date time.Time, items []expense.Item) ([]ExpenseDTO, error) {
	if s.file == nil {
		return nil, fmt.Errorf("save failed")
	}
	fileItems := make([]expensestore.Item, len(items))
	for i, it := range items {
		fileItems[i] = expensestore.Item{
			Amount:      it.Amount,
			Category:    it.Category,
			Description: it.Description,
		}
	}
	saved, err := s.file.Append(userID, source, date, fileItems)
	if err != nil {
		return nil, err
	}
	out := make([]ExpenseDTO, len(saved))
	for i, rec := range saved {
		out[i] = ExpenseDTO{
			ID:          rec.ID,
			Amount:      rec.Amount,
			Category:    rec.Category,
			Description: rec.Description,
		}
	}
	return out, nil
}
