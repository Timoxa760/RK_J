# PostgreSQL — модель данных

> Миграции: `back/db/migrations/postgres/001–008_*.up.sql`  
> БД: `moneymind`, пользователь `postgres` (compose).

## Миграции

| # | Файл | Таблицы |
|---|------|---------|
| 001 | `001_users` | `users`, `user_providers` |
| 002 | `002_receipts` | `receipts` |
| 003 | `003_credits` | `credits` |
| 004 | `004_categories_budgets_goals` | `categories`, `budgets`, `goals` |
| 005 | `005_challenges_achievements` | `challenges`, `challenge_participants`, `achievements`, `daily_activity` |
| 006 | `006_bank_tokens` | `bank_tokens`, `scraper_sessions` |
| 007 | `007_receipt_dedup` | `receipt_dedup` |
| 008 | `008_manual_expenses` | `manual_expenses` |

```bash
# back
make migrate      # up
make migrate-down # down 1
```

## Таблицы (детально)

### users

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | SERIAL PK | |
| phone | VARCHAR(20) UNIQUE | Логин |
| email | VARCHAR(255) UNIQUE | Опционально |
| role | VARCHAR(50) | default `user` |
| created_at | TIMESTAMPTZ | |

### user_providers

Привязки LK / OAuth. `credentials` — AES-256-GCM (см. scraper).

| Колонка | Тип |
|---------|-----|
| user_id | FK → users |
| provider_name | VARCHAR(50) |
| credentials | TEXT |
| status | default `active` |

UNIQUE `(user_id, provider_name)`.

### receipts

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | VARCHAR(64) PK | ID чека |
| user_id | VARCHAR(64) | |
| provider | VARCHAR(50) | x5club, mock, email, … |
| store_name | VARCHAR(255) | |
| total_amount | NUMERIC(10,2) | |
| purchased_at | TIMESTAMPTZ | |
| items | JSONB | Позиции itemized |

Дедуп: отдельная таблица `receipt_dedup` (hash sha256).

### manual_expenses — ключевой ввод «Поток»

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | UUID PK | |
| user_id | VARCHAR(64) | |
| raw_text | TEXT | Голос / свободный текст |
| amount | DECIMAL(12,2) | |
| category | VARCHAR(64) | |
| description | TEXT | |
| expense_date | DATE | |
| source | VARCHAR(16) | `manual` \| `voice` |
| created_at | TIMESTAMPTZ | |

Индексы: `user_id`, `expense_date`.

### goals

| Колонка | Тип | Описание |
|---------|-----|----------|
| name | VARCHAR(200) | Цель |
| target_amount | NUMERIC(12,2) | |
| current_amount | NUMERIC(12,2) | default 0 |
| deadline | DATE | |
| status | VARCHAR(20) | default `active` |

### credits

| Колонка | Тип |
|---------|-----|
| bank_name | VARCHAR(100) |
| amount, remaining_amount | NUMERIC(12,2) |
| interest_rate | NUMERIC(5,2) |
| term_months | INT |
| monthly_payment | NUMERIC(10,2) |
| next_payment_date | DATE |

### receipt_dedup

| Колонка | Тип |
|---------|-----|
| hash | VARCHAR(64) PK |
| expires_at | TIMESTAMPTZ |

TTL-кэш дедупликации (см. receipt-service).

### Прочие (демо / legacy)

| Таблица | Назначение |
|---------|------------|
| `categories`, `budgets` | Legacy UX |
| `challenges`, `achievements`, `daily_activity` | Social / gamification |
| `bank_tokens` | OAuth банков (optional) |
| `scraper_sessions` | Сессии sync провайдеров |

## Связи

- [../../product/financial-model.md](../../product/financial-model.md)
- [../clickhouse/schema.md](../clickhouse/schema.md)
- [../../architecture/kafka-events.md](../../architecture/kafka-events.md)
