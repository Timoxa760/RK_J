# Backend — быстрый старт

> Ветка **`back`**. Порты — [overview.md](../architecture/overview.md) (источник правды: `docker-compose.yml`).

## Стек

Go 1.25 · chi · PostgreSQL 18 · ClickHouse 25.12 · Redis 8.8 · Kafka 4.0.2 · Garage S3

## Команды

```bash
git checkout back
cp .env.example .env
make build
docker compose up -d postgres clickhouse redis kafka garage
make migrate
make migrate-ch
make up   # или docker compose up для всех сервисов
```

## Тесты

```bash
go test -count=1 ./services/...
# ~81 тест (на момент backend-v0.1)
```

## Seed / demo-данные

```bash
go run scripts/seed_data.go
# POST demo expenses → ai-processor /expenses/manual
```

См. [demo.md](./demo.md).

## Реализованные фазы (back)

| Фаза | Содержание |
|------|------------|
| 0 | Infra, миграции, Dockerfile |
| 1 | Auth + gateway |
| 2 | Ingest: голос/ручной, email (legacy) |
| 3 | Scraper + Kafka + receipt-service + dashboard API |
| 4–5 | Profile, credits PDF, rates |
| 6 | ai-processor advisor (`/ai/plan`, `/ai/chat`) |

Полная карта: [phases.md](../phases/phases.md).

## Связи

- [environment.md](./environment.md)
- [../api/API_Contract.md](../api/API_Contract.md)
- [../architecture/defense.md](../architecture/defense.md)
