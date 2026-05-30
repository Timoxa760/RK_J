package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PGStore хранит профили в user_financial_profiles.
type PGStore struct {
	pool *pgxpool.Pool
}

func NewPGStore(pool *pgxpool.Pool) *PGStore {
	return &PGStore{pool: pool}
}

func (s *PGStore) Get(userID string) (FinancialProfile, error) {
	ctx := context.Background()
	var p FinancialProfile
	var breakdownJSON, expensesJSON []byte
	var updatedAt time.Time
	err := s.pool.QueryRow(ctx, `
		SELECT active_income, passive_income, emergency_fund, emergency_breakdown,
		       fixed_expenses, goal_kind, goal_title, goal_amount,
		       skipped_income, skipped_cushion, skipped_goal, skipped_expenses,
		       survey_input_mode, onboarding_completed, updated_at
		FROM user_financial_profiles WHERE user_id = $1
	`, userID).Scan(
		&p.ActiveIncome, &p.PassiveIncome, &p.EmergencyFund, &breakdownJSON,
		&expensesJSON, &p.GoalKind, &p.GoalTitle, &p.GoalAmount,
		&p.SkippedIncome, &p.SkippedCushion, &p.SkippedGoal, &p.SkippedExpenses,
		&p.SurveyInputMode, &p.OnboardingCompleted, &updatedAt,
	)
	if err != nil {
		return Default(), nil
	}
	_ = json.Unmarshal(breakdownJSON, &p.EmergencyBreakdown)
	_ = json.Unmarshal(expensesJSON, &p.FixedExpenses)
	if p.FixedExpenses == nil {
		p.FixedExpenses = []FixedExpense{}
	}
	p.UpdatedAt = updatedAt
	return p, nil
}

func (s *PGStore) Save(userID string, p FinancialProfile) error {
	ctx := context.Background()
	breakdownJSON, err := json.Marshal(p.EmergencyBreakdown)
	if err != nil {
		return err
	}
	expensesJSON, err := json.Marshal(p.FixedExpenses)
	if err != nil {
		return err
	}
	if p.FixedExpenses == nil {
		expensesJSON = []byte("[]")
	}
	if p.GoalKind == "" {
		p.GoalKind = "save"
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO user_financial_profiles (
			user_id, active_income, passive_income, emergency_fund, emergency_breakdown,
			fixed_expenses, goal_kind, goal_title, goal_amount,
			skipped_income, skipped_cushion, skipped_goal, skipped_expenses,
			survey_input_mode, onboarding_completed, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			active_income = EXCLUDED.active_income,
			passive_income = EXCLUDED.passive_income,
			emergency_fund = EXCLUDED.emergency_fund,
			emergency_breakdown = EXCLUDED.emergency_breakdown,
			fixed_expenses = EXCLUDED.fixed_expenses,
			goal_kind = EXCLUDED.goal_kind,
			goal_title = EXCLUDED.goal_title,
			goal_amount = EXCLUDED.goal_amount,
			skipped_income = EXCLUDED.skipped_income,
			skipped_cushion = EXCLUDED.skipped_cushion,
			skipped_goal = EXCLUDED.skipped_goal,
			skipped_expenses = EXCLUDED.skipped_expenses,
			survey_input_mode = EXCLUDED.survey_input_mode,
			onboarding_completed = EXCLUDED.onboarding_completed,
			updated_at = NOW()
	`, userID, p.ActiveIncome, p.PassiveIncome, p.EmergencyFund, breakdownJSON,
		expensesJSON, p.GoalKind, p.GoalTitle, p.GoalAmount,
		p.SkippedIncome, p.SkippedCushion, p.SkippedGoal, p.SkippedExpenses,
		p.SurveyInputMode, p.OnboardingCompleted)
	if err != nil {
		return fmt.Errorf("save profile: %w", err)
	}
	return nil
}

// DualStore пишет в PG и file-store (advisor читает file при shared volume).
type DualStore struct {
	primary Store
	mirror  *FileStore
}

func NewDualStore(primary Store, mirror *FileStore) *DualStore {
	return &DualStore{primary: primary, mirror: mirror}
}

func (s *DualStore) Get(userID string) (FinancialProfile, error) {
	return s.primary.Get(userID)
}

func (s *DualStore) Save(userID string, p FinancialProfile) error {
	if err := s.primary.Save(userID, p); err != nil {
		return err
	}
	if s.mirror != nil {
		_ = s.mirror.Save(userID, p)
	}
	return nil
}
