# Implementation Plan — полный локальный стек

## Цель и контекст

Обеспечить end-to-end поток данных для «Поток»: фронт + бэкенд + PostgreSQL + Kafka + ClickHouse, без обязательного Whisper/OnlySQ. AI подключается отдельно.

## Объём работ

**Входит:** infra scripts, миграции, 15 сервисов локально, E2E verify, demo handlers для optional API.

**Не входит:** live bank sync, PG-backed credits/insights, production deploy.

## Архитектурное решение

- **Infra:** Docker Compose (postgres, kafka, clickhouse, redis).
- **App:** Go-бинарники на хосте (быстрая итерация), gateway :8000.
- **Данные manual expenses:** ai-processor → PostgreSQL `manual_expenses` → receipt-service dashboard.
- **Чеки Kafka:** scraper → `receipt.raw` → receipt-service → `receipt.parsed` → ai-processor → ClickHouse (при наличии брокера).

## Стек

| Компонент | Версия / образ | Назначение |
|-----------|----------------|------------|
| PostgreSQL | 18.0 | manual_expenses, receipts |
| Kafka | cp-kafka 4.0.2 | receipt pipeline |
| ClickHouse | 25.12 | аналитика чеков |
| Go services | 15 бинарников | API |

## Риски

- Docker Desktop обязателен на Windows.
- Credits/insights пока demo JSON — dashboard categories из PG.

## Альтернативы

- Полный `docker compose up` всех сервисов — медленнее сборка; выбран hybrid (infra в Docker, app на хосте).

## Связи

- **Зависит от:** `docker-compose.yml`, `db/migrations/`.
- **Связанные документы:** [local-full-stack.md](../guides/local-full-stack.md), [STATUS.md](../backend/STATUS.md).
