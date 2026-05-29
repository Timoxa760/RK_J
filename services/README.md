# services/ (ветка docs)

> **Не использовать для разработки.** Актуальный Docker Compose и все сервисы — в ветке **`back`**.

| Файл здесь | Назначение |
|------------|------------|
| `docker-compose.yml` | Legacy-заготовка (PG 16, CH 24, MinIO) — **устарела** |
| `.env.example` | Пример переменных |

```bash
git checkout back
docker compose up -d
```

См. [../docs/deployment/docker-compose.md](../docs/deployment/docker-compose.md).
