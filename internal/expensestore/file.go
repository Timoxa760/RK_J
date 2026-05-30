package expensestore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Record — сохранённая трата пользователя.
type Record struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}

type filePayload struct {
	Records []Record `json:"records"`
}

// FileStore — JSON-файл для локальной разработки без PostgreSQL.
type FileStore struct {
	path string
	mu   sync.Mutex
}

// NewFileStore создаёт хранилище по пути из EXPENSE_STORE_PATH или data/expenses.json.
func NewFileStore(path string) (*FileStore, error) {
	if path == "" {
		path = DefaultPath()
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir expense store: %w", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(`{"records":[]}`), 0o644); err != nil {
			return nil, fmt.Errorf("init expense store: %w", err)
		}
	}
	return &FileStore{path: path}, nil
}

// DefaultPath возвращает путь к файлу расходов.
func DefaultPath() string {
	if p := os.Getenv("EXPENSE_STORE_PATH"); p != "" {
		return p
	}
	return "data/expenses.json"
}

// Append добавляет траты пользователя.
func (f *FileStore) Append(userID, source string, date time.Time, items []Item) ([]Record, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	payload, err := f.readLocked()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	out := make([]Record, 0, len(items))
	for i, it := range items {
		rec := Record{
			ID:          fmt.Sprintf("file-%d-%d", now.UnixNano(), i),
			UserID:      userID,
			Amount:      it.Amount,
			Category:    it.Category,
			Description: it.Description,
			Date:        date,
			Source:      source,
			CreatedAt:   now,
		}
		payload.Records = append(payload.Records, rec)
		out = append(out, rec)
	}

	if err := f.writeLocked(payload); err != nil {
		return nil, err
	}
	return out, nil
}

// ListSince возвращает траты пользователя начиная с даты since включительно.
func (f *FileStore) ListSince(userID string, since time.Time) ([]Record, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	payload, err := f.readLocked()
	if err != nil {
		return nil, err
	}

	since = startOfDay(since)
	out := make([]Record, 0)
	for _, rec := range payload.Records {
		if rec.UserID != userID {
			continue
		}
		if !rec.Date.Before(since) {
			out = append(out, rec)
		}
	}
	return out, nil
}

// Item — входная позиция для Append.
type Item struct {
	Amount      float64
	Category    string
	Description string
}

func (f *FileStore) readLocked() (*filePayload, error) {
	raw, err := os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("read expense store: %w", err)
	}
	var payload filePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("decode expense store: %w", err)
	}
	return &payload, nil
}

func (f *FileStore) writeLocked(payload *filePayload) error {
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("encode expense store: %w", err)
	}
	tmp := f.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return fmt.Errorf("write expense store tmp: %w", err)
	}
	return os.Rename(tmp, f.path)
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
