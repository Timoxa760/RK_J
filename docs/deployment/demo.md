# Demo и seed

## DEMO_MODE

| `DEMO_MODE` | Auth | Dashboard | Advisor |
|-------------|------|-----------|---------|
| `true` | код `0000`, без SMS | hardcoded JSON | file-store + Gemini |
| `false` | SMS.ru OTP, Redis, Postgres | PG / ClickHouse | profile PG + shared volume |

Production stack: `DEMO_MODE=false`, `SMSRU_API_ID`, `JWT_SECRET`, `make migrate`.

Smoke:

```bash
# Production path (SMS на телефон)
SMOKE_PHONE=+79991234567 ./scripts/smoke_auth_chat.sh

# Legacy demo (код 0000)
DEMO_MODE=true SMOKE_OTP=0000 ./scripts/smoke_auth_chat.sh
```

## SMS.ru (production auth)

См. [environment.md](./environment.md) — регистрация, `api_id`, 5 free SMS/день на свой номер.

## Seed (`back/scripts/seed_data.go`)

```bash
git checkout back
go run scripts/seed_data.go
```

## Скрипты

| Скрипт | Назначение |
|--------|------------|
| [scripts/smoke_auth_chat.sh](../../backend/scripts/smoke_auth_chat.sh) | register → SMS → login → profile → `/ai/chat` |
| [scripts/demo_flow.sh](./scripts/demo_flow.sh) | API tour для жюри |
| [scripts/health_check.sh](./scripts/health_check.sh) | gateway + infra ping |

## Front

`NUXT_PUBLIC_DEMO_MODE=false` — реальный login через API (6-значный SMS-код).  
`true` — mock JWT без бэка (только offline UI).

## Связи

- [environment.md](./environment.md)
- [back-quickstart.md](./back-quickstart.md)
