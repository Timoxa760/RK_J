# Docker Compose — Инфраструктура

## Сервисы

| Сервис | Версия | Назначение |
|--------|--------|------------|
| PostgreSQL | 16 | OLTP-ядро |
| ClickHouse | 24 | OLAP-аналитика |
| Redis | 7 | Кэш + очереди |
| Kafka | latest | Event Bus |
| ZooKeeper | latest | Координация Kafka |
| MinIO | latest | Object Storage |

## Быстрый старт

```bash
docker compose up -d
docker compose ps
```

## Health-checks

```bash
curl http://localhost:8123/?query=SELECT%201
curl http://localhost:8082/health
```
