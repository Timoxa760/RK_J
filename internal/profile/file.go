package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileStore хранит профили в data/profiles/{user_id}.json (demo и fallback).
type FileStore struct {
	dir string
	mu  sync.Mutex
}

func NewFileStore(dir string) (*FileStore, error) {
	if dir == "" {
		dir = "data/profiles"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &FileStore{dir: dir}, nil
}

func (s *FileStore) path(userID string) string {
	safe := filepath.Base(userID)
	return filepath.Join(s.dir, safe+".json")
}

func (s *FileStore) Get(userID string) (FinancialProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path(userID))
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return FinancialProfile{}, err
	}
	var p FinancialProfile
	if err := json.Unmarshal(data, &p); err != nil {
		return FinancialProfile{}, err
	}
	return p, nil
}

func (s *FileStore) Save(userID string, p FinancialProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path(userID) + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path(userID))
}

func (s *FileStore) Patch(userID string, patch FinancialProfile, merge func(cur, patch FinancialProfile) FinancialProfile) (FinancialProfile, error) {
	cur, err := s.Get(userID)
	if err != nil {
		return FinancialProfile{}, err
	}
	next := merge(cur, patch)
	if err := s.Save(userID, next); err != nil {
		return FinancialProfile{}, fmt.Errorf("save profile: %w", err)
	}
	return next, nil
}

// MergePatch применяет partial update (zero = keep unless explicitly set via pointers — simplified: non-zero fields win for numbers when patch sent).
func MergePatch(cur, patch FinancialProfile) FinancialProfile {
	out := cur
	if patch.ActiveIncome > 0 || patch.SkippedIncome {
		out.ActiveIncome = patch.ActiveIncome
		out.SkippedIncome = patch.SkippedIncome
	}
	if patch.PassiveIncome > 0 || patch.SkippedIncome {
		out.PassiveIncome = patch.PassiveIncome
	}
	if patch.EmergencyFund > 0 || patch.SkippedCushion {
		out.EmergencyFund = patch.EmergencyFund
		out.SkippedCushion = patch.SkippedCushion
	}
	if patch.EmergencyBreakdown.Cash > 0 || patch.EmergencyBreakdown.Deposit > 0 || patch.EmergencyBreakdown.Investments > 0 {
		out.EmergencyBreakdown = patch.EmergencyBreakdown
	}
	if len(patch.FixedExpenses) > 0 || patch.SkippedExpenses {
		out.FixedExpenses = patch.FixedExpenses
		out.SkippedExpenses = patch.SkippedExpenses
	}
	if patch.GoalTitle != "" || patch.GoalAmount > 0 || patch.SkippedGoal {
		out.GoalKind = patch.GoalKind
		if patch.GoalKind != "" {
			out.GoalKind = patch.GoalKind
		}
		if patch.GoalTitle != "" {
			out.GoalTitle = patch.GoalTitle
		}
		out.GoalAmount = patch.GoalAmount
		out.SkippedGoal = patch.SkippedGoal
	}
	if patch.SurveyInputMode != "" {
		out.SurveyInputMode = patch.SurveyInputMode
	}
	if patch.OnboardingCompleted {
		out.OnboardingCompleted = true
	}
	return out
}
