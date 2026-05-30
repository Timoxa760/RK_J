# Переменные окружения

> Источник: `backend/.env.example`, `frontend/.env.example`.  
> Скопировать: `cp .env.example .env` (ветка `back` или `front`).

## Backend (`back`)

| Переменная | Пример | Назначение |
|------------|--------|------------|
| `DEMO_MODE` | `false` | `true` = legacy код `0000` без SMS; `false` = OTP + SMS.ru |
| `JWT_SECRET` | random string | HS256 — **одинаковый** на gateway, user-service, ai-processor |
| `DATABASE_URL` | `postgres://...` | PostgreSQL (users, profile, expenses) |
| `REDIS_URL` | `redis://localhost:6379/0` | OTP codes (TTL 5 min) |
| `SMS_PROVIDER` | `smsru` | `smsru` \| `console` (только CI) |
| `SMSRU_API_ID` | из [sms.ru](https://sms.ru) | API-ключ для отправки SMS |
| `SMSRU_TEST` | `0` | `1` = SMS.ru test mode (без доставки, для CI) |
| `OTP_TTL` | `300` | Срок жизни кода (сек) |
| `OTP_LENGTH` | `6` | Длина OTP |
| `OTP_RATE_LIMIT` | `5` | Max resend за 15 мин на номер |
| `GEMINI_API_KEY` | — | Advisor, parse, credit PDF |
| `PROFILE_DATA_DIR` | `/app/data/profiles` | File mirror для advisor (Docker volume) |
| `CREDIT_DATA_DIR` | `/app/data/credits` | Credit scans для advisor |
| `CLICKHOUSE_*` | см. compose | OLAP dashboard |
| `ENCRYPTION_KEY` | 32 bytes | AES для creds провайдеров |
| `RUCAPTCHA_KEY` | — | FNS MCO (не demo) |

### Настройка SMS.ru

1. Регистрация на [sms.ru](https://sms.ru) под **тем же номером**, что для входа в приложение.
2. Скопировать `api_id` → `SMSRU_API_ID`.
3. **Бесплатно:** до 5 SMS/день на свой номер ([sms.ru/free](https://sms.ru/free)).
4. Для других номеров — пополнить баланс (~200 ₽ на тесты).
5. Текст OTP: `Поток: код XXXXXX. Действует 5 мин.` (≤70 символов).

Проверка:

```bash
cd backend
docker compose up -d postgres redis user-service api-gateway ai-processor
docker compose --profile tools run --rm migrate-pg
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+7XXXXXXXXXX"}'
# код придёт в SMS
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
