# Переменные окружения

> Источник: `back/.env.example`, `front/frontend/.env.example`.  
> Скопировать: `cp .env.example .env` (ветка `back` или `front`).

## Backend (`back`)

| Переменная | Пример | Назначение |
|------------|--------|------------|
| `DEMO_MODE` | `true` | In-memory / hardcoded dashboard без полной infra |
| `JWT_SECRET` | random string | HS256 для access token |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/moneymind?sslmode=disable` | PostgreSQL |
| `CLICKHOUSE_HOST` | `localhost` | ClickHouse native |
| `CLICKHOUSE_USER` / `PASSWORD` / `DB` | см. compose | OLAP |
| `REDIS_URL` | `localhost:6379` | Кэш |
| `KAFKA_URL` | `localhost:9092` | Event bus |
| `ONLYSQ_URL` | `http://localhost:3000` | LLM (roadmap) |
| `ENCRYPTION_KEY` | 32 bytes | AES-256-GCM для creds провайдеров |
| `RUCAPTCHA_KEY` | — | FNS MCO captcha |
| `MCO_TOKEN_DIR` | `mco_tokens` | Токены FNS MCO |
| `2CAPTCHA_API_KEY` | — | Lenta provider (legacy) |
| `OAUTH_YANDEX_*` | — | Email receipts OAuth |
| `OAUTH_MAIL_*` | — | Mail.ru OAuth |

## Frontend (`front`)

| Переменная | По умолчанию | Назначение |
|------------|--------------|------------|
| `NUXT_PUBLIC_API_BASE` | `http://localhost:8000` | api-gateway |
| `NUXT_PUBLIC_DEMO_MODE` | `true` | Моки без бэкенда |

## Связи

- [back-quickstart.md](./back-quickstart.md)
- [front-quickstart.md](./front-quickstart.md)
- [../architecture/defense.md](../architecture/defense.md) — DEMO_MODE
