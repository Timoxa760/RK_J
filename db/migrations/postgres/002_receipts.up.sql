CREATE TABLE IF NOT EXISTS receipts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    provider VARCHAR(50) NOT NULL,
    store_name VARCHAR(255) NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL, -- Копейки хранм в NUMERIC/Float64
    purchased_at TIMESTAMP WITH TIME ZONE NOT NULL,
    checksum VARCHAR(64) UNIQUE NOT NULL, -- SHA-256 для дедупликации
    items JSONB NOT NULL,                 -- Состав чека по ТЗ в JSONB
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);