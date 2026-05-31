# Demo и seed

## DEMO_MODE

| `DEMO_MODE` | Auth | Dashboard | Advisor |
|-------------|------|-----------|---------|
| `true` | телефон + пароль | hardcoded JSON | file-store + Gemini |
| `false` | телефон + пароль, Postgres | PG / ClickHouse | profile PG + shared volume |

Production stack: `DEMO_MODE=false`, `JWT_SECRET`, `DATABASE_URL`, `make migrate`.

Smoke:

```bash
SMOKE_PHONE=+79991234567 SMOKE_PASSWORD=secret12345 ./scripts/smoke_auth_chat.sh
```

## Seed (`back/scripts/seed_data.go`)

```bash
git checkout back
go run scripts/seed_data.go
```

## Скрипты

| Скрипт | Назначение |
|--------|------------|
| [scripts/smoke_auth_chat.sh](../../backend/scripts/smoke_auth_chat.sh) | register → login → profile → `/ai/chat` |
| [scripts/demo_flow.sh](./scripts/demo_flow.sh) | API tour для жюри |
| [scripts/health_check.sh](../../backend/scripts/health_check.sh) | gateway + infra ping |

## Front

`NUXT_PUBLIC_DEMO_MODE=false` — реальный login через API (телефон + пароль).  
`true` — mock JWT без бэка (только offline UI).

## Связи

- [environment.md](./environment.md)
- [back-quickstart.md](./back-quickstart.md)
