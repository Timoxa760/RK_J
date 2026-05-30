CREATE TABLE IF NOT EXISTS manual_expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(64) NOT NULL,
    raw_text TEXT NOT NULL DEFAULT '',
    amount DECIMAL(12,2) NOT NULL,
    category VARCHAR(64) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    expense_date DATE NOT NULL DEFAULT CURRENT_DATE,
    source VARCHAR(16) NOT NULL DEFAULT 'manual',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_manual_expenses_user_id ON manual_expenses (user_id);
CREATE INDEX IF NOT EXISTS idx_manual_expenses_date ON manual_expenses (expense_date);
