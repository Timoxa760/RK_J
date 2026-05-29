CREATE TABLE IF NOT EXISTS receipt_items (
    user_id String,
    store String,
    category String,
    item_name String,
    price Float64,
    quantity UInt32,
    purchased_at DateTime,
    is_impulsive UInt8
) ENGINE = MergeTree()
ORDER BY (user_id, purchased_at);

CREATE MATERIALIZED VIEW IF NOT EXISTS spending_by_category
ENGINE = SummingMergeTree()
ORDER BY (user_id, category, toStartOfMonth(purchased_at))
AS SELECT
    user_id,
    category,
    toStartOfMonth(purchased_at) as month,
    sum(price * quantity) as total
FROM receipt_items
GROUP BY user_id, category, month;

CREATE MATERIALIZED VIEW IF NOT EXISTS store_aggregates
ENGINE = SummingMergeTree()
ORDER BY (user_id, store)
AS SELECT
    user_id,
    store,
    count(*) as purchases,
    avg(price * quantity) as avg_check,
    sum(price * quantity) as total,
    avg(is_impulsive) as impulse_ratio
FROM receipt_items
GROUP BY user_id, store;

CREATE MATERIALIZED VIEW IF NOT EXISTS monthly_category_spending
ENGINE = SummingMergeTree()
ORDER BY (user_id, category, month)
AS SELECT
    user_id,
    category,
    toStartOfMonth(purchased_at) as month,
    sum(price * quantity) as total
FROM receipt_items
GROUP BY user_id, category, month;