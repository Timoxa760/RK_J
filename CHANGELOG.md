# Changelog

Все заметные изменения документации ветки **`docs`** проекта RK_J («Поток»).

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.1.0/),
версионирование — [Semantic Versioning](https://semver.org/lang/ru/).

**Базовая версия для сравнения:** `b0a74ac` — *docs: add full project documentation* (27.05.2026, бренд **MoneyMind**).

---

## [2.2.0] — 2026-05-30

### Кратко

Синхронизация MVP-документации с ветками `front`/`back`, якорь [mvp/README.md](./docs/mvp/README.md), три устных питча для защиты.

### Added

- **`docs/mvp/README.md`** — in/out scope, mock vs real, demo script 3 мин, команды запуска.
- **`docs/pitch/README.md`**, **`teamlead.md`**, **`frontend.md`**, **`backend.md`** — устные питчи 5–10 мин.

### Changed

- **`docs/product/scope.md`** — ФНС = front mock; `/digest`, social out.
- **`docs/product/overview.md`**, **`ux-scenarios.md`** — `/dashboard` + `/advisor`, demo path.
- **`docs/product/advisor.md`** — streaming, text repair, structured blocks.
- **`docs/phases/phases.md`** — фаза 9, FNS mock, блок «Защита MVP».
- **`docs/deployment/front-quickstart.md`** — актуальные маршруты, app shell.
- **`docs/features/credit-scan.md`** — кириллические PDF.
- **`NAVI.md`**, **`Potok_plan.md`** — секции MVP и Pitch.
- Слияние `front/docs` → `docs` worktree (сохранены `advisor-system.md`, `llm-integration.md`, `antigravity-setup.md`).

---

## [2.1.0] — 2026-05-31

### Кратко

Документированы **Antigravity Tools + LLM dual-mode**, полная архитектура **финансового советника** (SSE, actions, PG history), изменения dashboard/front (удаление `/digest`, narrative, mindfulness). Синхронизация с ветками `back` и `front`.

### Added

- **`docs/architecture/llm-integration.md`** — Google direct vs Antigravity OpenAI route, env, troubleshooting.
- **`docs/architecture/advisor-system.md`** — snapshot, chat pipeline, SSE, actions, frontend map.
- **`docs/deployment/antigravity-setup.md`** — пошаговая настройка прокси :8045 для dev.

### Changed

- **`docs/architecture/overview.md`** — LLM-слой, Antigravity в схеме, SSE endpoints.
- **`docs/product/advisor.md`** — `/advisor`, streaming, history, actions, Antigravity.
- **`docs/product/ux-scenarios.md`** — без `/digest`, narrative на dashboard.
- **`docs/deployment/environment.md`** — `GEMINI_MODEL=claude-sonnet-4-6`, `/v1` base URL.
- **`docs/deployment/front-quickstart.md`** — маршруты `/advisor`, удалён `/digest`.
- **`docs/features/digest.md`** — пометка deprecated, контент на dashboard.
- **`docs/phases/phases.md`** — фаза 9 (LLM + UX polish).
- **`NAVI.md`**, **`docs/guides/git-and-branches.md`** — новые ссылки и worktree layout.

### Backend (`back`) — задокументировано

| Изменение | Файлы |
|-----------|-------|
| LLM dual-mode client | `internal/llm/client.go`, `types.go` |
| Antigravity OpenAI route | `GEMINI_BASE_URL=/v1`, model `claude-sonnet-4-6` в compose |
| Smoke stream fix | `scripts/smoke_auth_chat.sh` |

### Frontend (`front`) — задокументировано

| Изменение | Файлы |
|-----------|-------|
| Удалён `/digest` | `pages/digest.vue`, nav, mocks |
| Narrative + mindfulness | `PageNarrative.vue`, `DashboardMindfulnessScore.vue` |
| Sidebar advisor убран | `AppSidebarAdvisor.vue` удалён |
| Симулятор «Что если» | `ScenarioSimulator.vue`, dashboard layout |

---

## [2.0.0] — 2026-05-30

### Кратко

Ребрендинг **MoneyMind → «Поток»** (финансовый навигатор), полная переработка GDD под кейс Клерк.Ру, синхронизация docs с фактическими ветками `back` / `front`, единый API-контракт и инфраструктура для demo/защиты.

### Added

#### Точки входа и GDD

- **`NAVI.md`** — карта репозитория («где что искать»): ветки, маршруты front, порты back, быстрые ответы.
- **`Potok_plan.md`** — новое оглавление GDD (продукт + реализация), заменяет монолитный `MoneyMind_plan.md`.
- **`docs/product/`** (10 документов):
  - `overview.md`, `ux-scenarios.md`, `onboarding.md`, `input-methods.md`
  - `financial-model.md`, `financial-health.md`, `philosophy.md`, `monetization.md`
  - `case-alignment.md` — маппинг на ТЗ Клерк.Ру (светофор, диагноз, цели)
  - `roadmap.md` — post-MVP и оценки в dev-days

#### API и контракты

- **`docs/api/API_Contract.md`** (~691 строк) — единый контракт: auth, dashboard, expenses/manual, FNS, goals, credits, gateway-таблица, приоритеты и JSON-примеры.
- **`docs/api/typescript-types.md`** — выжимка из `front/frontend/types/api.ts` для интеграции UI.
- **`docs/contracts/openapi.yaml`** — OpenAPI 3.1 (критические paths).
- **`docs/contracts/README.md`** — описание пакета контрактов.

#### Архитектура и данные

- **`docs/architecture/defense.md`** — pitch backend (перенос из `back/docs/ARCHITECTURE_DEFENSE.txt`).
- **`docs/architecture/kafka-events.md`** — топики Kafka, pipeline `receipt.raw` → `receipt.parsed`.
- **`docs/architecture/stack-audit.md`** — сверка docs ↔ код через Context7.

#### Deployment и demo

- **`docs/deployment/back-quickstart.md`**, **`front-quickstart.md`**, **`environment.md`**, **`demo.md`**.
- **`docs/deployment/scripts/demo_flow.sh`** — 6 шагов API для жюри (login → manual expense → dashboard → credits → insights).
- **`docs/deployment/scripts/health_check.sh`** — health gateway и infra.

#### Guides и Cursor

- **`docs/guides/context7.md`**, **`docs/guides/mcp-setup.md`** — MCP Context7 + GitHub для RK_J.
- **`.cursor/skills/potok-docs/`** — skill агента (замена moneymind-docs).
- **`.cursor/rules/`**, **`.cursor/context7-libraries.json`** — workflow и реестр библиотек.

#### Прочее

- **`services/README.md`** — пометка, что legacy compose в `docs/services/` устарел; актуальный — `back/docker-compose.yml`.

### Changed

#### Позиционирование продукта

| Было (1.0.0) | Стало (2.0.0) |
|--------------|---------------|
| MoneyMind — «поштучный финансовый ассистент» | **Поток** — «финансовый навигатор, ответы вместо графиков» |
| Фокус на receipt-level tracking | Фокус на здоровье, прогнозе цели, одном действии на неделю |
| Social и auction — в core GDD | **Гипотезы**, не в demo (явно в features + phases) |
| Онбординг не выделен | Отдельная страница **`/onboarding`** (⏳ на front) |

#### README

- Актуальный стек: **Nuxt 4.3** (не React 19 + Vite), **Go 1.25**, 15 микросервисов.
- Быстрый старт через `make migrate`, demo-код `0000`.
- Таблица документации со ссылкой на NAVI и API Contract.

#### Архитектура (`docs/architecture/overview.md`)

- Порты и сервисы сверены с `back/docker-compose.yml` (credit-service **8009**, bank **8011**, ai-processor **8100**, analytics **8101**).
- Стек инфра: PG 18, ClickHouse 25.12, Redis 8.8, Kafka 4.0.2, Garage S3.
- Убраны устаревшие ссылки на Python-сервисы как основной стек.

#### База данных

- **`docs/database/postgresql/schema.md`** — расширено по миграциям `001`–`008`: users, receipts, manual_expenses, goals, credits, social, scraper_sessions и др.
- **`docs/database/clickhouse/schema.md`** — таблицы OLAP, materialized views из миграций dashboard.

#### API (`docs/api/openapi.md`)

- Индекс перенаправляет на **API_Contract.md** и **typescript-types.md** как источники правды.

#### Фичи (`docs/features/*.md`)

- Каждая фича привязана к сервису `back` и экрану `front`.
- Переформулировки под «навигатор»: narrative, микро-действия, фин. здоровье.

#### Фазы (`docs/phases/phases.md`)

- Статусы ✅/🟡/⏳ по сверке с `back` / `front`.
- Перестроен критический путь: **онбординг → ввод → здоровье → прогноз → ипотека**.
- Фаза 14: ссылка на эталонный `demo_flow.sh` в `docs/deployment/scripts/`.

#### Deployment (`docs/deployment/docker-compose.md`)

- Актуальный compose и env — в ветке **`back`**, не в `docs/services/`.

### Removed

- **`MoneyMind_plan.md`** — содержимое декомпозировано в `Potok_plan.md` + `docs/product/` + `docs/architecture/`.

### Deprecated

| Путь | Замена |
|------|--------|
| `MoneyMind_plan.md` | `Potok_plan.md` |
| `front/API_Contract.md` | `docs/api/API_Contract.md` |
| `back/docs/API_Contract.md` (пустой) | `docs/api/API_Contract.md` |
| `docs/services/docker-compose.yml` | `back/docker-compose.yml` |
| Skill `moneymind-docs` | `.cursor/skills/potok-docs/` |

### Fixed

- Расхождения docs ↔ код: порты credit/bank/ai/analytics, Nuxt вместо React, версии БД и Kafka.
- Битые relative links в `docs/database/postgresql/schema.md`.
- DTI и timemachine: зафиксированы delta front vs back в API_Contract (§13).

### Migration guide (1.0.0 → 2.0.0)

1. Вместо `MoneyMind_plan.md` открывай **`Potok_plan.md`** или **`NAVI.md`**.
2. API для интеграции — **`docs/api/API_Contract.md`**, не `front/API_Contract.md`.
3. Запуск infra — **`git checkout back`**, не `docs/services/`.
4. Demo-скрипты — **`docs/deployment/scripts/`** (при необходимости копировать в `back/scripts/`).
5. Social / auction — только как гипотеза; для защиты кейса см. **`docs/product/case-alignment.md`**.

---

## [1.0.0] — 2026-05-27

### Added

- Первый набор документации в ветке `docs` (коммит `b0a74ac`).
- **`MoneyMind_plan.md`** — оглавление GDD (монолит, 48 строк).
- **`README.md`** — ветки, быстрый старт, ссылки на docs.
- **`docs/architecture/overview.md`** — краткий обзор стека и сервисов.
- **`docs/api/openapi.md`** — краткая API-спецификация.
- **`docs/database/`** — черновики схем PostgreSQL и ClickHouse.
- **`docs/features/`** — 10 фич (receipt-magic, credit-health, time-machine, detective, social, predictive-ai, auction, digest, fotochecking, credit-constructor).
- **`docs/phases/phases.md`** — 15 фаз хакатона (48 ч), критический путь scraper → receipt → dashboard.
- **`docs/deployment/docker-compose.md`** — описание compose.
- **`services/docker-compose.yml`**, **`services/.env.example`** — legacy infra в ветке docs.

### Known limitations (1.0.0)

- Бренд и продукт описаны как **MoneyMind** (трекер по чекам), не «Поток».
- Front указан как React 19 + Vite — не совпадает с фактическим Nuxt 4 в ветке `front`.
- Нет продуктовой папки `docs/product/`, case-alignment, roadmap, NAVI, OpenAPI YAML, полного API Contract.
- Social и auction — в core GDD без пометки «гипотеза».
- Схемы БД — краткие черновики без миграций `back`.

---

[2.2.0]: https://github.com/Timoxa760/RK_J/compare/c512fa6...docs
[2.1.0]: https://github.com/Timoxa760/RK_J/compare/90930e9...docs
[2.0.0]: https://github.com/Timoxa760/RK_J/compare/b0a74ac...docs
[1.0.0]: https://github.com/Timoxa760/RK_J/commit/b0a74ac
