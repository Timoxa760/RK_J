# Аудит стека (Context7 + ветки `back`/`front`, 2026-05-30)

## Связи

- **Используется**: `overview.md`, skill `potok-docs`
- **Источник факта**: `git show back:go.mod`, `back/docker-compose.yml`, `front:frontend/package.json`
- **Context7**: `/golang/go`, `/websites/nuxt_4_x`, `/go-chi/docs`, `/clickhouse/clickhouse-docs`

## Результат сверки

| Компонент | Старые docs | Факт (`back`/`front`) | Context7 / статус |
|-----------|-------------|------------------------|-------------------|
| **Go** | 1.26 (целевой) | **1.25.0** в go.mod | 1.26.0 exists; back на 1.25 — OK, не апгрейдить в хакатоне |
| **Frontend** | React 19 + Vite 8 | **Nuxt 4.3.0** | Nuxt 4 stable (`/websites/nuxt_4_x`) |
| **Vue** | — | 3 (via Nuxt) | OK |
| **Tailwind** | v4 | `@nuxtjs/tailwindcss` 6.x | v4 через Nuxt module |
| **PostgreSQL** | 16 | **18.0** image | PG 18 supported |
| **ClickHouse** | 24 → 25 | **25.12** | Актуально |
| **Redis** | 7 | **8.8.0** | Новее docs |
| **Kafka** | 3.x / latest | **4.0.2** (cp-kafka) | Зафиксировано в compose |
| **S3** | MinIO | **Garage 2.3.0** | back использует Garage, не MinIO |
| **chi, pgx, kafka-go** | — | как в go.mod | Context7 OK |
| **Python FastAPI** | ai-enrichment | **нет в back** | AI на Go ai-processor |
| **OpenAPI** | markdown | **API_Contract.md + openapi.yaml** | ✅ в `docs/` |

## Расхождения docs ↔ код (зафиксировано, код не трогаем)

| Тема | docs (было) | код |
|------|-------------|-----|
| Бренд | MoneyMind (старое) | **Поток**, навигатор |
| Frontend stack | React | **Nuxt 4** |
| credit-service port | 8009 | **8009** в back (compose + gateway) |
| bank-service port | 8008 | **8011** в back |
| dashboard | analytics-service | **receipt-service** `/api/v1/dashboard/*` |
| Object storage | MinIO | **Garage** |

## Действия (только docs, выполнено в этой ветке)

1. ✅ `overview.md` — порты и стек из `back`
2. ✅ `README.md` — Nuxt 4, Поток
3. ✅ `deployment/docker-compose.md` — описание infra из `back`
4. ✅ `docs/api/API_Contract.md` + `docs/contracts/openapi.yaml`
5. ⏳ `front/API_Contract.md` — deprecated, ссылка на docs

## Не менять без согласования с `back`/`front`

- Имена и порты микросервисов
- Топики Kafka
- Структура каталогов `services/*`
- package name `money-mind-frontend` (front) — косметика, не в scope хакатона

## Источники Context7 (2026-05-30)

- Nuxt 4: `/websites/nuxt_4_x` — active v4, semver releases
- Go: `/golang/go` — go1.25.0, go1.26.0
- chi: `/go-chi/docs`
- ClickHouse: `/clickhouse/clickhouse-docs`
