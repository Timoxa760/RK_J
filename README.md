# MoneyMind — Backend

**Поштучный финансовый ассистент**

Хакатон Финтех 2026 · Кейсодатель: Клерк.Ру

Ветка `back` — 15 Go-микросервисов + инфраструктура.

---

## Стек

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.25 |
| HTTP-роутер | chi |
| PostgreSQL | 18 |
| ClickHouse | 25.12 |
| Redis | 8.8 |
| Kafka | 4.0.2 |
| S3 | Garage 2.3.0 |
| Аутентификация | JWT (HS256) |
| Шифрование | AES-256-GCM |

## Сервисы и порты

| Сервис | Порт | Категория |
|--------|------|-----------|
| api-gateway | 8000 | core-api |
| user-service | 8001 | core-api |
| receipt-service | 8002 | receipt-engine |
| scraper-service | 8003 | receipt-engine |
| category-service | 8004 | finance-core |
| budget-service | 8005 | finance-core |
| ~~goal-service~~ | ~~8006~~ | **removed** — цель в profile |
| credit-service | 8007 | finance-core |
| bank-service | 8008 | finance-core |
| ai-processor | 8009 | money-intelligence |
| analytics-service | 8010 | money-intelligence |
| notification-service | 8011 | reporting |
| reporting-service | 8100 | reporting |
| social-service | 8101 | social-game |
| gamification | 8102 | social-game |

## Реализованные фазы

- **Phase 0** — Инфраструктура (Docker Compose, миграции PG/CH, Makefile, 15 Dockerfile)
- **Phase 1** — Email + ФНС (OAuth Яндекс/Mail.ru, IMAP-клиент, парсинг HTML-чеков, FNS API)
- **Phase 2** — X5 Club провайдер (HTTP-клиент, Worker Pool, маппинг)
- **Phase 3** — Магнит провайдер (HTTP-клиент, INN detector, маппинг)
- **Phase 4** — Scraper ядро (Provider Interface, Scheduler, Kafka Producer, Credentials, Rate Limiter, Mock)
- **Phase 5** — Receipt-service (Kafka Consumer, Validator, дедупликация, PostgreSQL + JSONB)
- **Phase 6** — Auth + JWT (регистрация, логин, JWT middleware, привязка провайдеров)

## Архитектура

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client   │────▶│ api-gw   │────▶│  user    │
└──────────┘     │  (:8000) │     │ (:8001)  │
                 └────┬─────┘     └──────────┘
                      │
              ┌───────┴────────┐
              │  scraper-svc   │
              │    (:8003)     │
              │ ┌───────────┐  │
              │ │ x5club    │  │
              │ │ magnit    │──│──▶ Kafka ──▶ receipt-svc
              │ │ email/FNS │  │    receipt.raw  (:8002)
              │ │ scheduler │  │                │
              │ └───────────┘  │                ▼
              └────────────────┘           PostgreSQL
```

## Быстрый старт

```bash
# 1. Скопировать и заполнить .env
cp .env.example .env

# 2. Собрать все сервисы
make build

# 3. Запустить инфраструктуру
docker compose up -d postgres clickhouse redis kafka

# 4. Накатить миграции
make migrate
make migrate-ch

# 5. Запустить сервисы
make up
```

## Тесты

```bash
go test -count=1 ./services/...

# Результат: 81 тест, все PASS
# go vet — чисто
```

## Ветки проекта

| Ветка | Назначение |
|-------|------------|
| `main` | Базовая (README, точка входа) |
| `back` | Бэкенд ← |
| `front` | Фронтенд |
| `docs` | Документация |
