# ClickHouse — аналитика

> Миграция: `back/db/migrations/clickhouse/001_receipt_items.sql`

```bash
make migrate-ch   # back
```

## Таблица receipt_items

| Колонка | Тип | Описание |
|---------|-----|----------|
| user_id | String | |
| store | String | Магазин |
| category | String | Категория |
| item_name | String | Позиция |
| price | Float64 | |
| quantity | UInt32 | |
| purchased_at | DateTime | |
| is_impulsive | UInt8 | 0/1 |

Engine: `MergeTree()`, partition `toYYYYMM(purchased_at)`, order `(user_id, purchased_at)`.

## Materialized Views

### spending_by_category

Агрегат трат по категориям и месяцам. Используется dashboard / analytics.

```sql
-- sum(price * quantity) GROUP BY user_id, category, month
```

Engine: `SummingMergeTree()`, partition by month.

### store_aggregates

Агрегат по магазинам: purchases, avg_check, total, impulse_ratio.

```sql
-- для GET /dashboard/stores
```

Engine: `SummingMergeTree()`, order `(user_id, store)`.

## Поток данных

```
ai-processor (categorizer)
        │
        ▼
 receipt_items (CH)  ◄── dual-write с manual_expenses (PG)
        │
        ├── MV spending_by_category
        └── MV store_aggregates
        │
        ▼
 receipt-service dashboard handlers
```

## Связи

- [postgresql/schema.md](../postgresql/schema.md)
- [../../api/typescript-types.md](../../api/typescript-types.md)
- [../../architecture/kafka-events.md](../../architecture/kafka-events.md)
