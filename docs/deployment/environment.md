# Переменные окружения

> Источник: `backend/.env.example`, `frontend/.env.example`.  
> Скопировать: `cp .env.example .env` (ветка `back` или `front`).

## Backend (`back`)

| Переменная | Пример | Назначение |
|------------|--------|------------|
| `DEMO_MODE` | `false` | `true` = legacy shortcuts в части сервисов; auth всегда телефон + пароль |
| `JWT_SECRET` | random string | HS256 — **одинаковый** на gateway, user-service, ai-processor |
| `DATABASE_URL` | `postgres://...` | PostgreSQL (users, profile, expenses) |
| `REDIS_URL` | `redis://localhost:6379/0` | Коды сброса пароля (TTL, опционально — in-memory fallback) |
| `OTP_TTL` | `300` | Срок жизни кода сброса пароля (сек) |
| `GEMINI_API_KEY` | — | Ключ Google AI Studio **или** API Key из Antigravity Tools → API Proxy |
| `GEMINI_PROVIDER` | `antigravity` | `antigravity` (локальный прокси :8045) \| `google` (прямой API) |
| `GEMINI_MODEL` | `gpt-oss-120b-medium` | Модель в Antigravity (`gpt-oss-120b-medium`, `claude-sonnet-4-6`, …). Native Gemini может быть заблокирован по региону аккаунта |
| `GEMINI_BASE_URL` | `http://127.0.0.1:8045/v1` | OpenAI-маршрут Antigravity; в Docker: `http://host.docker.internal:8045/v1` |
| `PROFILE_DATA_DIR` | `/app/data/profiles` | File mirror для advisor (Docker volume) |
| `CREDIT_DATA_DIR` | `/app/data/credits` | Credit scans для advisor |
| `CLICKHOUSE_*` | см. compose | OLAP dashboard |
| `ENCRYPTION_KEY` | 32 bytes | AES для creds провайдеров |

### Auth (телефон + пароль)

Регистрация и вход — `POST /auth/register` и `POST /auth/login` с телефоном и паролем (мин. 8 символов). SMS не используется.

Сброс пароля: `POST /auth/password/forgot` → `POST /auth/password/reset` с кодом из Redis (локально допускается stub `0000`).

Проверка:

```bash
cd backend
docker compose up -d postgres redis user-service api-gateway ai-processor
docker compose --profile tools run --rm migrate-pg
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+7XXXXXXXXXX","password":"secret12345"}'
```

## Frontend (`front`)

| Переменная | По умолчанию | Назначение |
|------------|--------------|------------|
| `NUXT_PUBLIC_API_BASE` | `http://localhost:8000` | api-gateway |
| `NUXT_PUBLIC_DEMO_MODE` | `false` | `true` = mock JWT без бэка |

## Связи

- [back-quickstart.md](./back-quickstart.md)
- [front-quickstart.md](./front-quickstart.md)
- [demo.md](./demo.md)
