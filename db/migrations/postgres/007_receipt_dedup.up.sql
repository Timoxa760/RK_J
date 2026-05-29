CREATE TABLE IF NOT EXISTS receipt_dedup (
    hash VARCHAR(64) PRIMARY KEY,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_receipt_dedup_expires_at ON receipt_dedup (expires_at);
