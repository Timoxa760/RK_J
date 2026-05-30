# Бэкенд «Поток» — статус, changelog, roadmap

> Ветка: `back` · Gateway: `http://localhost:8000/api/v1`  
> Контракт API: [docs/api/API_Contract.md](../api/API_Contract.md)  
> Последняя верификация: **2026-05-30** (локально, `DEMO_MODE=true`)

## Связи

- **Зависит от:** PostgreSQL, Kafka, ClickHouse (prod-режим); для demo-smoke — только Go-бинарники.
- **Используется:** Nuxt front (`front`), api-gateway.
- **Связанные документы:** [API_Contract.md](../api/API_Contract.md), корневой [README.md](../../README.md).

---

## Уверенность в тестах

Проверки выполнялись **в этой среде** (Linux, без Docker Compose и без Kafka/PostgreSQL на хосте):

| Проверка | Команда | Результат |
|----------|---------|-----------|
| Основной smoke | `./scripts/smoke_api.sh` | **ALL API SMOKE PASSED** |
| Критические + scraper | `bash scripts/smoke_critical.sh` | **0 FAIL**, 1 **PARTIAL** (`/fns/mco/sync` → 500 без Kafka) |
| Unit (ключевые модули) | см. раздел «Как перепроверить» | **OK** (при `env -u JWT_SECRET`) |

**Ограничения (честно):**

- Dashboard, credits, goals, digest в smoke — **demo/in-memory** или заглушки, не live PG.
- `POST /fns/mco/sync` с валидным `phone` отвечает **500**, если Kafka недоступен — это ожидаемо в локальном прогоне без брокера.
- `POST /credits/scan` требует **multipart** с полем `file`; JSON-тело даёт 400 (не баг API).
- Полный **docker compose** (PG + Kafka + все 15 сервисов) в CI этого прогона **не гонялся**.
- Тест `user-service/middleware` падает, если в shell выставлен `JWT_SECRET`, отличный от дефолта в тесте — при прогоне использовать `env -u JWT_SECRET`.

Перед релизом рекомендуется повторить smoke на машине с поднятым `docker compose` и зафиксировать дату в этом файле.

---

## Changelog (ветка `back`)

Формат по мотивам [Keep a Changelog](https://keepachangelog.com/ru/1.1.0/). Версии не тегировались; ниже — по коммитам git.

### [Unreleased]

#### Добавлено

- `POST /api/v1/expenses/voice` — Whisper STT + OnlySQ parse + advice.
- Поля `advice`, `expenses[]`, `parsed_by`, `transcript` в expenses API.
- Клиент OnlySQ, Whisper, `docs/backend/DEPLOY.md`, сервис `whisper` в docker-compose.

#### Изменено

- `POST /api/v1/expenses/manual` — OnlySQ с fallback на regex.
- ai-processor: общий Processor для manual и voice.

### 82b09d5 — docs(backend): статус API, changelog и smoke_critical

#### Добавлено

- Demo API целей: `POST /goals`, `GET /goals`, `GET /goals/{id}` (goal-service, in-memory).
- Алиасы analytics для front: `GET /api/v1/analytics/insights`, `/forecast`, `POST /analytics/simulate`.
- Прокси gateway: префикс `/api/v1/analytics/`.
- Поле `token` в ответе login (дубликат `access_token` для Nuxt).
- Скрипт `scripts/smoke_critical.sh` — прогон critical/important из контракта.

#### Исправлено

- Маршрутизация gateway (полный path, `registerProxy` base + `base/*`).
- Auth/register/login/providers по API_Contract.
- ai-processor: manual expenses в DEMO без PG.

### 06913e6 — feat(api): goals, digest demo and analytics front aliases

См. коммит выше (Unreleased по сути совпадает с HEAD).

### 3ea102a — feat(api): credits dashboard, insights and scenarios demo API

- `GET /credits/dashboard`, `POST /credits/scan` (demo).
- `GET /insights`, `GET /forecast`, `POST /scenarios/simulate` (demo).
- Расширен `smoke_api.sh`.

### 506d984 — fix(api): sync auth and gateway with API_Contract

- Gateway: корректный proxy path, `*_SERVICE_URL` для локального запуска.
- Auth: `phone` + demo `code` `0000`, providers `?provider=`.
- Скрипты: `demo_flow.sh`, `smoke_api.sh`, `smoke_local.sh`.

### 12b7f12 — feat: ai-processor, dashboard, MCO, dedup

- ai-processor (manual, parser, categorizer, ClickHouse).
- Dashboard handlers, FNS MCO provider, seed, dedup.

### dafda67 — feat: backend Phases 0–6

- 15 сервисов, scraper, receipt pipeline, auth JWT.

---

## Что работает (проверено smoke)

Статус **E2E через gateway** (`smoke_api` + `smoke_critical`, 2026-05-30).

| Приоритет | Эндпоинт | Статус | Режим |
|-----------|----------|--------|-------|
| critical | `POST /auth/register` | OK | demo |
| critical | `POST /auth/login` | OK | demo code `0000`, JWT + `token` |
| critical | `POST /providers/connect` | OK | demo |
| critical | `GET /dashboard/sankey` | OK | demo |
| critical | `GET /dashboard/categories` | OK | demo |
| critical | `POST /expenses/manual` | OK | ai-processor, OnlySQ или regex |
| critical | `POST /expenses/voice` | NEW | Whisper + OnlySQ (нужен compose) |
| critical | `POST /fns/ticket` | OK | scraper demo (нужен `SCRAPER_SERVICE_URL` в smoke_critical) |
| important | `GET /dashboard/timemachine` | OK | 60 месяцев в JSON |
| important | `GET /dashboard/stores` | OK | demo |
| important | `GET /dashboard/compare` | OK | demo |
| important | `GET /credits/dashboard` | OK | demo, DTI 0..1 |
| important | `POST /credits/scan` | OK | multipart `file` |
| important | `GET /insights` | OK | demo |
| important | `GET /forecast` | OK | demo |
| important | `GET /analytics/insights` | OK | алиас front |
| important | `POST /scenarios/simulate` | OK | demo |
| important | `POST /goals` | OK | in-memory demo |
| important | `GET /goals/{id}` | OK | in-memory demo |
| optional | `GET /digest/latest` | OK | demo |
| — | JWT без токена / невалидный | OK | 401 |

### Частично

| Эндпоинт | Статус | Причина |
|----------|--------|---------|
| `POST /fns/mco/sync` | PARTIAL | HTTP 500 без Kafka; валидация `phone` работает |
| Kafka → receipt pipeline | PARTIAL | consumer стартует; без брокера — ошибки в логах |
| Dashboard на реальных чеках PG | PARTIAL | нужны PG + seed / ingest |

### Не реализовано / 502 через gateway

| Эндпоинт | Статус |
|----------|--------|
| `PATCH /users/me/profile` | Roadmap (onboarding) |
| `POST /users/me/onboarding/complete` | Roadmap |
| Live x5club / magnit sync | Код есть, E2E smoke нет |
| PG-backed credits / insights | Demo JSON, не из `manual_expenses` |

С 2026-05-30 optional API (`/banks/*`, `/categories`, `/budgets`, `/challenges`) отвечают demo JSON при запущенных сервисах (см. `start_services.ps1`).

---

## Что ещё нужно реализовать (приоритет)

### P0 — для демо с живыми данными

1. **Полный локальный стек:** `powershell -File scripts\start_stack.ps1` — см. [local-full-stack.md](../guides/local-full-stack.md).
2. Проверка потока: `powershell -File scripts\verify_e2e.ps1` (register → expense → dashboard PG).
3. **MCO sync** — стабильный ответ при наличии Kafka (или graceful demo без 500).
4. Связать **receipt-service** с PG для dashboard не в DEMO_MODE — реализовано для `manual_expenses`.
5. Синхронизация с **front**: credits DTI (% vs доля), timemachine, `useAuth().token`.

### P1 — контракт important

1. **bank-service** — accounts, transactions.
2. **category-service** / **budget-service** — CRUD по контракту.
3. **POST /providers/{name}/sync** — форс-синк чеков.
4. Onboarding: **PATCH /users/me/profile**, complete flag.

### P2 — optional / хакатон+

1. **social-service** — challenges, leaderboard.
2. **notification-service**, **gamification**.
3. Live scraper smoke (x5, magnit) с тестовыми кредами.
4. OpenAPI в `docs/contracts/` синхронизировать с фактическими путями (`/analytics/*`).

---

## Как перепроверить

```bash
cd RK_J   # корень ветки back
export DEMO_MODE=true JWT_SECRET=test-secret
export USER_SERVICE_URL=http://127.0.0.1:8001
export RECEIPT_SERVICE_URL=http://127.0.0.1:8002
export SCRAPER_SERVICE_URL=http://127.0.0.1:8003
export AI_PROCESSOR_URL=http://127.0.0.1:8100
export CREDIT_SERVICE_URL=http://127.0.0.1:8009
export ANALYTICS_SERVICE_URL=http://127.0.0.1:8101
export GOAL_SERVICE_URL=http://127.0.0.1:8006
export REPORTING_SERVICE_URL=http://127.0.0.1:8010

./scripts/smoke_api.sh
bash scripts/smoke_critical.sh

# unit (без конфликта JWT_SECRET)
env -u JWT_SECRET go test \
  ./services/core-api/user-service/... \
  ./services/receipt-engine/receipt-service/internal/... \
  ./services/receipt-engine/scraper-service/internal/providers/fns/... \
  ./services/money-intelligence/ai-processor/internal/manual/... \
  ./services/finance-core/credit-service/internal/... \
  ./services/finance-core/goal-service/internal/... \
  ./services/money-intelligence/analytics-service/internal/... \
  ./services/reporting/reporting-service/internal/... \
  -count=1
```

Кеши Go: `~/.cache/potok-go-build`, `~/.cache/potok-go-mod` (см. `.gitignore`).

---

## Карта сервисов (фактические порты в коде)

| Сервис | Порт (локально) | Gateway prefix |
|--------|-----------------|----------------|
| api-gateway | 8000 | — |
| user-service | 8001 | `/auth/`, `/providers/` |
| receipt-service | 8002 | `/dashboard/`, `/receipts/` |
| scraper-service | 8003 | `/fns/`, `/x5club/`, `/magnit/`, `/email/` |
| category-service | 8004 | `/categories/` |
| budget-service | 8005 | `/budgets/` |
| goal-service | 8006 | `/goals/` |
| credit-service | 8009 | `/credits/` |
| reporting-service | 8010 | `/digest/` |
| bank-service | 8011 | `/banks/` |
| ai-processor | 8100 | `/expenses/` |
| analytics-service | 8101 | `/insights/`, `/forecast/`, `/scenarios/`, `/analytics/` |
| social-service | 8102 | `/challenges/` |

Примечание: порты в корневом README.md могут расходиться с `main.go` — источник правды для smoke: `.env.example` и `scripts/smoke_*.sh`.
