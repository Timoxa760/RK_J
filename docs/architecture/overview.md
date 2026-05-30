# Системная архитектура

> Версии сверены с Context7 (2026-05-30) и **фактическим кодом ветки `back`**. Подробности: [stack-audit.md](./stack-audit.md).

## Продукт и клиент

**Поток** — финансовый навигатор (ветка `front`: Nuxt 4, порт 3000 dev / 80 docker).  
API: `api-gateway:8000` → микросервисы.

## Технологический стек

| Слой | Технология | Ветка |
|------|------------|-------|
| **Backend** | Go **1.25** (chi, pgx, segmentio/kafka-go, go-jwt, oauth2) | `back` |
| **AI / enrichment** | Go `ai-processor` + **Google Gemini** (категоризация, голос/ручной, advisor) | `back` |
| **Analytics** | Go `analytics-service` (insights, forecast, scenarios) | `back` |
| **Frontend** | **Nuxt 4.3** (Vue 3, Pinia, Tailwind, ECharts, PWA) | `front` |
| **Scraping / ingest** | scraper-service: **FNS**, MCO, email/IMAP (X5/Magnit — legacy, не MVP) | `back` |
| **Infra** | Docker Compose: PG 18, CH 25.12, Redis 8.8, Kafka 4.0.2, Garage S3 | `back` |
| **Auth** | JWT HS256, OAuth Яндекс/Mail.ru (email receipts) | `back` |

> В старых docs указаны React/Vite и Python-сервисы — **в текущем `back`/`front` их нет**; analytics и AI реализованы на Go.

## Базы данных

| БД | Версия (`back`) | Назначение |
|----|-----------------|------------|
| **PostgreSQL** | 18.0 | OLTP: users, receipts, user_financial_profiles, user_credits, manual_expenses |
| **ClickHouse** | 25.12 | OLAP: receipt_items, агрегаты dashboard |
| **Redis** | 8.8.0 | Кэш, сессии, leaderboards |
| **Kafka** | 4.0.2 (cp-kafka) | `receipt.raw` → parsed → enriched |
| **Garage** | 2.3.0 | Object storage (S3-compatible) |

## Микросервисы (как в `back/docker-compose.yml`)

| Сервис | Порт | Категория | Назначение |
|--------|------|-----------|------------|
| api-gateway | 8000 | core-api | Единый вход, JWT, reverse proxy |
| user-service | 8001 | core-api | Регистрация, login, providers |
| receipt-service | 8002 | receipt-engine | Чеки, дедуп, **dashboard API** |
| scraper-service | 8003 | receipt-engine | **FNS**, email (X5/Magnit — legacy) |
| category-service | 8004 | finance-core | CRUD категорий |
| budget-service | 8005 | finance-core | Бюджеты (legacy, не core UX) |
| credit-service | 8009 | finance-core | PDF scan, DTI, rates client |
| reporting-service | 8010 | reporting | Дайджест |
| bank-service | 8011 | finance-core | Банки, ипотечный разбор |
| ai-processor | 8100 | money-intelligence | Голос/ручной, advisor, onboarding parse |
| analytics-service | 8101 | money-intelligence | Insights, forecast, scenarios |

**Removed from MVP:** `goal-service` (цель в profile), `social-service` / gamification (см. [scope.md](../product/scope.md)).

> Порты и имена **не менять** — зафиксированы в `back/docker-compose.yml` и api-gateway routes.  
> `back/README.md` может расходиться с compose — приоритет у compose + gateway.

## Схема взаимодействия

```
[Nuxt front :3000] → [api-gateway :8000]
                           │
     ┌─────────────────────┼─────────────────────┐
     ▼                     ▼                     ▼
 user-service      scraper-service        ai-processor
 receipt-service         │                      │
     │                   ▼                      │
     │              Kafka receipt.*             │
     ▼                   ▼                      ▼
 PostgreSQL         ClickHouse            analytics-service
```

## Маппинг продукта → сервисы

| Продукт | Сервисы |
|---------|---------|
| Онбординг, цель в профиле | user-service (`/users/me/profile`) |
| Кредиты (PDF-only) | credit-service + internal/rates |
| Советник, план | ai-processor (`/ai/plan`, `/ai/chat`) |
| Голос / ручной ввод | ai-processor |
| ФНС / чеки | scraper-service, receipt-service |
| Финансовое здоровье | credit-service, receipt-service (dashboard) |
| Инсайты, прогноз | analytics-service |
| Ипотека (monetization) | credit-service, bank-service, analytics-service |

## Связи

- **Зависит от**: ветка `back` (источник правды по инфра)
- **Используется**: [../product/](../product/), [../features/](../features/), [../api/API_Contract.md](../api/API_Contract.md), [../../NAVI.md](../../NAVI.md)
- **Связанные документы**: [stack-audit.md](./stack-audit.md), [../deployment/docker-compose.md](../deployment/docker-compose.md)
