# Поток

**Финансовый навигатор — ответы вместо графиков**

Хакатон Финтех — 29–31 мая 2026 | Кейсодатель: Клерк.Ру

> *«Что мне сделать сегодня, чтобы быстрее получить то, что я хочу?»*

**Не знаешь, куда идти?** → **[NAVI.md](./NAVI.md)** — карта всего репозитория.

---

## Ветки

| Ветка | Назначение | Стек |
|-------|------------|------|
| `main` | Базовая точка входа | README |
| `back` | Бэкенд: 15 Go-микросервисов, Kafka, миграции | Go 1.25, chi, pgx, kafka-go |
| `front` | Фронтенд приложения «Поток» | **Nuxt 4.3** (Vue 3, Pinia, Tailwind, ECharts, PWA) |
| `docs` | Документация: продукт, архитектура, API, БД, фазы | Markdown |

> Репозиторий: [Timoxa760/RK_J](https://github.com/Timoxa760/RK_J). Код живёт в `back` и `front`; ветка `docs` — источник правды по продукту.

## Быстрый старт

```bash
# Бэкенд (ветка back)
git checkout back
cp .env.example .env
make build && docker compose up -d
make migrate && make migrate-ch

# Фронтенд (ветка front)
git checkout front
cd frontend && npm install && cp .env.example .env
npm run dev   # http://localhost:3000, демо-код 0000
```

## Документация

| Раздел | Описание |
|--------|----------|
| **[NAVI.md](./NAVI.md)** | **Где что искать** — начни здесь |
| [CHANGELOG.md](./CHANGELOG.md) | История изменений docs (1.0.0 → 2.0.0) |
| [Potok_plan.md](./Potok_plan.md) | GDD, оглавление |
| [Продукт и UX](./docs/product/) | Онбординг, сценарии, философия |
| [Архитектура](./docs/architecture/) | Стек, сервисы, порты |
| [API Contract](./docs/api/API_Contract.md) | Полный контракт + [OpenAPI 3.1](./docs/contracts/openapi.yaml) |
| [БД](./docs/database/) | PostgreSQL + ClickHouse |
| [Фичи → сервисы](./docs/features/) | Маппинг на `back` |
| [Фазы](./docs/phases/) | План хакатона и статус |
| [Context7](./docs/guides/context7.md) | Реестр для агента |
| [Deployment](./docs/deployment/) | back/front, env, demo |

## Суть продукта (кратко)

**Поток** — финансовый навигатор для людей 25–45 лет, которые не хотят вести сложный учёт, но хотят понимать своё положение и быстрее достигать целей.

Система показывает **траекторию жизни в деньгах**: финансовое здоровье, прогноз цели, мягкие рекомендации и одно действие на неделю.

**Ввод расходов** (на выбор): голос, скан чека, ФНС (опционально).
