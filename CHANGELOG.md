# Changelog

Все заметные изменения документации ветки **`docs`** проекта RK_J («Поток»).

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.1.0/),
версионирование — [Semantic Versioning](https://semver.org/lang/ru/).

**Базовая версия для сравнения:** `b0a74ac` — *docs: add full project documentation* (27.05.2026, бренд **MoneyMind**).

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

[2.0.0]: https://github.com/Timoxa760/RK_J/compare/b0a74ac...docs
[1.0.0]: https://github.com/Timoxa760/RK_J/commit/b0a74ac
