CREATE TABLE IF NOT EXISTS user_credits (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    bank TEXT NOT NULL DEFAULT '',
    amount NUMERIC NOT NULL DEFAULT 0,
    rate NUMERIC NOT NULL DEFAULT 0,
    term_months INT NOT NULL DEFAULT 0,
    remaining NUMERIC NOT NULL DEFAULT 0,
    monthly_payment NUMERIC NOT NULL DEFAULT 0,
    benchmark_rate NUMERIC,
    rate_vs_market TEXT,
    next_payment DATE,
    scanned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    source_file_hash TEXT
);

CREATE INDEX IF NOT EXISTS idx_user_credits_user_id ON user_credits (user_id);
