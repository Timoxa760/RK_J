# Docker Compose — инфраструктура

> **Источник правды:** `docker-compose.yml` в ветке **`back`**.  
> В ветке `docs` лежит упрощённый `services/docker-compose.yml` (legacy, только infra-заготовка) — для разработки используйте `back`.

## Сервисы инфра (`back`)

| Сервис | Образ | Назначение |
|--------|-------|------------|
| PostgreSQL | `postgres:18.0` | OLTP, БД `moneymind` |
| ClickHouse | `clickhouse/clickhouse-server:25.12` | OLAP |
| Redis | `redis:8.8.0` | Кэш |
| ZooKeeper | `confluentinc/cp-zookeeper:latest` | Kafka |
| Kafka | `confluentinc/cp-kafka:4.0.2` | Event bus |
| Garage | `dxflrs/garage:v2.3.0` | S3-compatible storage |

Сеть: `moneymind_network`.

## Быстрый старт (`back`)

```bash
git checkout back
cp .env.example .env
docker compose up -d postgres clickhouse redis kafka garage
make migrate
make migrate-ch
make build
# запуск сервисов — см. back/README.md
```

## Health-checks

```bash
# ClickHouse
curl 'http://localhost:8123/?query=SELECT%201'

# PostgreSQL
docker exec backend_postgres pg_isready -U postgres

# api-gateway (после старта сервисов)
curl http://localhost:8000/health
```

## Frontend (`front`)

```bash
git checkout front
docker compose up --build   # nginx :80
# или локально: cd frontend && npm run dev  # :3000
```

## Связи

- **Архитектура**: [../architecture/overview.md](../architecture/overview.md)
- **Миграции**: [../database/postgresql/schema.md](../database/postgresql/schema.md)
