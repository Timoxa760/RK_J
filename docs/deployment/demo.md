# Demo и seed

## DEMO_MODE

`DEMO_MODE=true` в `.env` (`back`):

- user-service — in-memory
- dashboard — hardcoded JSON (без ClickHouse)
- быстрый onboard без Docker

Подробнее: [defense.md §10](../architecture/defense.md).

## Seed (`back/scripts/seed_data.go`)

```bash
git checkout back
go run scripts/seed_data.go
```

POST demo-расходы → `/api/v1/expenses/manual`.

## Скрипты demo (готовые в docs)

Эталонные скрипты лежат в **`docs/deployment/scripts/`**. Скопировать в `back/scripts/` при необходимости:

```bash
git checkout back
cp ../docs/deployment/scripts/demo_flow.sh scripts/
cp ../docs/deployment/scripts/health_check.sh scripts/
chmod +x scripts/*.sh
```

| Скрипт | Назначение |
|--------|------------|
| [scripts/demo_flow.sh](./scripts/demo_flow.sh) | 6 шагов API для жюри |
| [scripts/health_check.sh](./scripts/health_check.sh) | gateway + infra ping |

### demo_flow.sh

1. Login (`0000`)
2. POST `/expenses/manual` (голос)
3. GET dashboard sankey / timemachine
4. GET credits dashboard
5. GET insights

Сценарий UI: [case-alignment.md](../product/case-alignment.md).

### health_check.sh

Проверяет `api-gateway/health`, dashboard (200/401), ClickHouse HTTP.

```bash
API_BASE=http://localhost:8000 ./docs/deployment/scripts/health_check.sh
```

## Front demo

`NUXT_PUBLIC_DEMO_MODE=true` — моки без бэкенда.

## Связи

- [phases.md](../phases/phases.md) — фаза 8 (demo polish)
- [back-quickstart.md](./back-quickstart.md)
