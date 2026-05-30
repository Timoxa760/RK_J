package creditstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type FileStore struct {
	dir string
	mu  sync.Mutex
}

func NewFileStore(dir string) (*FileStore, error) {
	if dir == "" {
		dir = "data/credits"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &FileStore{dir: dir}, nil
}

func (s *FileStore) path(userID string) string {
	return filepath.Join(s.dir, filepath.Base(userID)+".json")
}

func (s *FileStore) List(userID string) ([]Credit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path(userID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var items []Credit
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *FileStore) Add(userID string, c Credit) (Credit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := s.readLocked(userID)
	if err != nil {
		return Credit{}, err
	}
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	if c.ScannedAt.IsZero() {
		c.ScannedAt = time.Now().UTC()
	}
	c.UserID = userID
	items = append(items, c)
	if err := s.writeLocked(userID, items); err != nil {
		return Credit{}, err
	}
	return c, nil
}

func (s *FileStore) Delete(userID, creditID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, err := s.readLocked(userID)
	if err != nil {
		return err
	}
	next := make([]Credit, 0, len(items))
	for _, c := range items {
		if c.ID != creditID {
			next = append(next, c)
		}
	}
	return s.writeLocked(userID, next)
}

func (s *FileStore) Dashboard(userID string, monthlyIncome, emergencyFund float64) Dashboard {
	items, _ := s.List(userID)
	var totalDebt, monthlyPayments float64
	for _, c := range items {
		rem := c.Remaining
		if rem <= 0 {
			rem = c.Amount
		}
		totalDebt += rem
		monthlyPayments += c.MonthlyPayment
	}
	dti := float64(0)
	if monthlyIncome > 0 && monthlyPayments > 0 {
		dti = monthlyPayments / monthlyIncome * 100
	}
	stress := float64(0)
	if monthlyPayments > 0 && emergencyFund > 0 {
		stress = emergencyFund / monthlyPayments
	}
	return Dashboard{
		DTI:              dti,
		StressTestMonths: stress,
		Savings:          emergencyFund,
		TotalDebt:        totalDebt,
		MonthlyPayments:  monthlyPayments,
		MonthlyIncome:    monthlyIncome,
		Credits:          items,
	}
}

func (s *FileStore) readLocked(userID string) ([]Credit, error) {
	data, err := os.ReadFile(s.path(userID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var items []Credit
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *FileStore) writeLocked(userID string, items []Credit) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path(userID) + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path(userID))
}

func CompareRate(actual, benchmark float64) string {
	if benchmark <= 0 {
		return "unknown"
	}
	if actual <= benchmark+0.5 {
		return "at_or_below"
	}
	return "above"
}

func FormatNextPayment() string {
	return time.Now().AddDate(0, 1, 15).Format("2006-01-02")
}

func ValidateScan(amount, rate float64, term int) error {
	if amount <= 0 || rate <= 0 || term <= 0 {
		return fmt.Errorf("incomplete contract fields")
	}
	return nil
}
