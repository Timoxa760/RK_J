# NAVI — где что искать

Карта репозитория **RK_J** (продукт **«Поток»**). Начни отсюда, если не знаешь, в какой файл идти.

---

## С чего начать

| Задача | Куда |
|--------|------|
| Что изменилось в docs | [CHANGELOG.md](./CHANGELOG.md) |
| Понять продукт за 5 минут | [Potok_plan.md](./Potok_plan.md) → [docs/product/overview.md](./docs/product/overview.md) |
| UX, сценарии, онбординг | [docs/product/](./docs/product/) |
| Запустить бэкенд | [docs/deployment/back-quickstart.md](./docs/deployment/back-quickstart.md) |
| Запустить фронт | [docs/deployment/front-quickstart.md](./docs/deployment/front-quickstart.md) |
| Защита архитектуры (pitch) | [docs/architecture/defense.md](./docs/architecture/defense.md) |
| Соответствие ТЗ Клерка | [docs/product/case-alignment.md](./docs/product/case-alignment.md) |
| Roadmap | [docs/product/roadmap.md](./docs/product/roadmap.md) |
| API эндпоинты | [docs/api/API_Contract.md](./docs/api/API_Contract.md) |
| OpenAPI YAML | [docs/contracts/openapi.yaml](./docs/contracts/openapi.yaml) |
| Порты сервисов | [docs/architecture/overview.md](./docs/architecture/overview.md) |
| Статус хакатона | [docs/phases/phases.md](./docs/phases/phases.md) |

---

## Ветки git

| Ветка | Что внутри | Не трогать без нужды |
|-------|------------|----------------------|
| **`docs`** | Вся продуктовая и техдокументация (ты здесь) | — |
| **`back`** | Go-сервисы, `docker-compose.yml`, миграции, Makefile | порты, имена сервисов |
| **`front`** | Nuxt 4, `frontend/pages/`, composables | структура `frontend/` |
| **`main`** | Минимальный README | — |

```bash
git checkout docs    # документация
git checkout back    # бэкенд
git checkout front   # фронтенд
```

---

## Документация (`docs/`)

### Продукт — «зачем и как для пользователя»

| Файл | Что внутри |
|------|------------|
| [product/overview.md](./docs/product/overview.md) | ЦА, проблема, инсайт, метрики |
| [product/ux-scenarios.md](./docs/product/ux-scenarios.md) | 5 сценариев, главный экран, цикл |
| [product/onboarding.md](./docs/product/onboarding.md) | Онбординг, маршрут `/onboarding` |
| [product/input-methods.md](./docs/product/input-methods.md) | Голос, ручной ввод, кнопка «Добавить» |
| [product/financial-model.md](./docs/product/financial-model.md) | Единая модель, микро-действия |
| [product/financial-health.md](./docs/product/financial-health.md) | Здоровье, runway, подача |
| [product/philosophy.md](./docs/product/philosophy.md) | Тон, без давления |
| [product/monetization.md](./docs/product/monetization.md) | Ипотечный разбор (monetization) |

### Архитектура и инфра

| Файл | Что внутри |
|------|------------|
| [architecture/overview.md](./docs/architecture/overview.md) | Стек, 15 сервисов, порты |
| [architecture/stack-audit.md](./docs/architecture/stack-audit.md) | Context7, расхождения docs ↔ код |
| [architecture/kafka-events.md](./docs/architecture/kafka-events.md) | Kafka топики, pipeline |
| [architecture/defense.md](./docs/architecture/defense.md) | Защита архитектуры для pitch |
| [guides/context7.md](./docs/guides/context7.md) | MCP Context7 |
| [guides/mcp-setup.md](./docs/guides/mcp-setup.md) | MCP GitHub + Context7 setup |

### Deployment

| Файл | Что внутри |
|------|------------|
| [deployment/back-quickstart.md](./docs/deployment/back-quickstart.md) | Make, migrate, тесты |
| [deployment/front-quickstart.md](./docs/deployment/front-quickstart.md) | Nuxt, маршруты, структура |
| [deployment/environment.md](./docs/deployment/environment.md) | `.env` back + front |
| [deployment/demo.md](./docs/deployment/demo.md) | DEMO_MODE, seed, **demo_flow.sh** |
| [deployment/scripts/](./docs/deployment/scripts/) | health_check, demo_flow |
| [deployment/docker-compose.md](./docs/deployment/docker-compose.md) | Infra из `back` |
| [deployment/vps-deploy.md](./docs/deployment/vps-deploy.md) | Prod VPS: Docker, Caddy, Antigravity |
| [deployment/vercel-deploy.md](./docs/deployment/vercel-deploy.md) | Prod Vercel: домен, env |

### Продукт (доп.)

| Файл | Что внутри |
|------|------------|
| [product/case-alignment.md](./docs/product/case-alignment.md) | Маппинг ТЗ Клерк.Ру |
| [product/roadmap.md](./docs/product/roadmap.md) | Post-MVP, оценки в dev-days |

### API и данные

| Файл | Что внутри |
|------|------------|
| [api/API_Contract.md](./docs/api/API_Contract.md) | **Полный контракт** — JSON, приоритеты |
| [api/typescript-types.md](./docs/api/typescript-types.md) | TS-типы front (`types/api.ts`) |
| [api/openapi.md](./docs/api/openapi.md) | Индекс API |
| [contracts/openapi.yaml](./docs/contracts/openapi.yaml) | OpenAPI 3.1 (Swagger) |
| [database/postgresql/schema.md](./docs/database/postgresql/schema.md) | Таблицы PG, `manual_expenses` |
| [database/clickhouse/schema.md](./docs/database/clickhouse/schema.md) | OLAP, dashboard |

### Фичи → код

| Файл | Продукт | Сервис `back` |
|------|---------|---------------|
| [features/receipt-magic.md](./docs/features/receipt-magic.md) | Ввод расходов, pipeline | scraper, receipt, ai-processor |
| [features/credit-health.md](./docs/features/credit-health.md) | Фин. здоровье | credit-service |
| [features/credit-scan.md](./docs/features/credit-scan.md) | PDF scan кредитов | credit-service |
| [features/credit-constructor.md](./docs/features/credit-constructor.md) | Ипотека | credit, bank, analytics |
| [features/predictive-ai.md](./docs/features/predictive-ai.md) | Прогноз цели | analytics-service |
| [features/time-machine.md](./docs/features/time-machine.md) | «Если не менять» | dashboard/timemachine |
| [features/detective.md](./docs/features/detective.md) | Инсайты-причины | analytics insights |
| [features/digest.md](./docs/features/digest.md) | Действие недели | reporting-service |
| [features/fotochecking.md](./docs/features/fotochecking.md) | _(удалено)_ | — |
| [features/social.md](./docs/features/social.md) | **Гипотеза**, низкий приоритет | social-service |
| [features/auction.md](./docs/features/auction.md) | **Гипотеза**, не в demo | — |

### План работ

| Файл | Что внутри |
|------|------------|
| [phases/phases.md](./docs/phases/phases.md) | MVP фазы 0–8; фаза 6 — [advisor.md](./docs/product/advisor.md) |

---

## Фронтенд (`front` → `frontend/`)

| Маршрут | Экран | Документ |
|---------|-------|----------|
| `/` | Welcome, лендинг | [ux-scenarios.md](./docs/product/ux-scenarios.md) |
| `/login` | Вход (демо `0000`) | — |
| **`/onboarding`** | **Онбординг ~1 мин** (⏳ сделать) | [onboarding.md](./docs/product/onboarding.md) |
| `/dashboard` | Главный экран, графики | [ux-scenarios.md](./docs/product/ux-scenarios.md) |
| `/receipts` | Расходы, «Добавить» | [input-methods.md](./docs/product/input-methods.md) |
| `/credits` | Кредиты, ипотека | [monetization.md](./docs/product/monetization.md) |
| `/analytics` | Прогноз | [predictive-ai.md](./docs/features/predictive-ai.md) |
| `/digest` | Дайджест | [digest.md](./docs/features/digest.md) |
| `/social` | Social (гипотеза) | [social.md](./docs/features/social.md) |
| `/profile` | Профиль | — |

**Composables:** `frontend/composables/useDashboard.ts`, `useCredits.ts`, …  
**API base:** `NUXT_PUBLIC_API_BASE=http://localhost:8000`

---

## Бэкенд (`back`)

| Что | Где |
|-----|-----|
| Все сервисы | `back/services/` (core-api, receipt-engine, finance-core, …) |
| Миграции PG | `back/db/migrations/postgres/` |
| Миграции CH | `back/db/migrations/clickhouse/` |
| Gateway routes | `back/services/core-api/api-gateway/` |
| Compose | `back/docker-compose.yml` |
| Seed / demo | `back/scripts/seed_data.go`; эталон: [docs/deployment/scripts/](./docs/deployment/scripts/) |

**Ключевые API для продукта:**

| Продукт | Endpoint |
|---------|----------|
| Голос/ручной | `POST /api/v1/expenses/manual` |
| Dashboard | `GET /api/v1/dashboard/*` |
| Цель / профиль | `GET/PATCH /api/v1/users/me/profile` |
| Кредиты (PDF) | `POST /api/v1/credits/scan`, `GET /api/v1/credits/dashboard` |
| Советник | `GET /api/v1/ai/plan`, `POST /api/v1/ai/chat` |
| Ипотека | `/api/v1/banks/*`, `/api/v1/scenarios/*` |

---

## Cursor / агент

| Файл | Назначение |
|------|------------|
| [.cursor/skills/potok-docs/SKILL.md](./.cursor/skills/potok-docs/SKILL.md) | Skill для правки docs |
| [.cursor/context7-libraries.json](./.cursor/context7-libraries.json) | ID библиотек Context7 |
| [.cursor/rules/potok-docs-workflow.mdc](./.cursor/rules/potok-docs-workflow.mdc) | Правила workflow |

---

## Устаревшее (не использовать)

| Путь | Почему |
|------|--------|
| `services/docker-compose.yml` в ветке `docs` | Legacy; актуальный compose в **`back`** |
| `front/API_Contract.md` | Перенесён в **`docs/api/API_Contract.md`** |
| `MoneyMind_plan.md` | Переименован → **Potok_plan.md** |

---

## Быстрые ответы

**«Где описан онбординг?»** → [docs/product/onboarding.md](./docs/product/onboarding.md), маршрут `/onboarding` на `front`.

**«Где ипотека для жюри?»** → [docs/product/monetization.md](./docs/product/monetization.md) + `/credits` на `front`.

**«Social в demo?»** → Нет, только [гипотеза](./docs/features/social.md).

**«Какой порт у credit-service?»** → **8009** ([overview.md](./docs/architecture/overview.md)).
