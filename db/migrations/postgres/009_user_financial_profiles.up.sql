CREATE TABLE IF NOT EXISTS user_financial_profiles (
    user_id TEXT PRIMARY KEY,
    active_income NUMERIC NOT NULL DEFAULT 0,
    passive_income NUMERIC NOT NULL DEFAULT 0,
    emergency_fund NUMERIC NOT NULL DEFAULT 0,
    emergency_breakdown JSONB NOT NULL DEFAULT '{"cash":0,"deposit":0,"investments":0}',
    fixed_expenses JSONB NOT NULL DEFAULT '[]',
    goal_kind TEXT NOT NULL DEFAULT 'save',
    goal_title TEXT NOT NULL DEFAULT '',
    goal_amount NUMERIC NOT NULL DEFAULT 0,
    skipped_income BOOLEAN NOT NULL DEFAULT FALSE,
    skipped_cushion BOOLEAN NOT NULL DEFAULT FALSE,
    skipped_goal BOOLEAN NOT NULL DEFAULT FALSE,
    skipped_expenses BOOLEAN NOT NULL DEFAULT FALSE,
    survey_input_mode TEXT,
    onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
