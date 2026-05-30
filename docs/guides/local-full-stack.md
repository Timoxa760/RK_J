# Локальный полный стек «Поток»

Запуск бэкенда с PostgreSQL, Kafka, ClickHouse и всеми 15 Go-сервисами. Поток данных для UI: **login → manual expense → PostgreSQL → dashboard**.

Whisper и OnlySQ подключаются отдельно (следующий шаг); без них manual expenses работают через **regex** и сохраняются в БД.

## Связи

- **Зависит от:** Docker Desktop, Go 1.22+, `make build`.
- **Используется:** Nuxt front (`NUXT_PUBLIC_API_BASE=http://127.0.0.1:8000`, `NUXT_PUBLIC_DEMO_MODE=false`).
- **Связанные документы:** [docs/backend/STATUS.md](../backend/STATUS.md), [docs/backend/DEPLOY.md](../backend/DEPLOY.md), [docs/api/API_Contract.md](../api/API_Contract.md).

## Быстрый старт (Windows)

```powershell
cd C:\backend_project
powershell -File scripts\start_stack.ps1
powershell -File scripts\verify_e2e.ps1
```

## Быстрый старт (macOS / Linux, Docker)

```bash
cd backend
cp .env.example .env   # JWT_SECRET один для gateway + user-service + ai-processor
make up-full           # docker compose up + migrate-pg + migrate-ch
./scripts/smoke_api.sh
```

Фронт:

```bash
cd frontend
NUXT_PUBLIC_API_BASE=http://127.0.0.1:8000 NUXT_PUBLIC_DEMO_MODE=false npm run dev -- --port 3000
```

Логин: телефон + пароль (регистрация через `/auth/register`). В `DEMO_MODE=true` — код **0000**.

## Что поднимает `start_stack.ps1`

| Шаг | Действие |
|-----|----------|
| 1 | `docker compose up` — postgres, redis, zookeeper, kafka, clickhouse |
| 2 | Миграции PG + CH (`docker compose --profile tools run migrate-*`) |
| 3 | `make build` — 15 бинарников в `bin/` |
| 4 | `start_services.ps1` — все сервисы, `DEMO_MODE=false` |

## Порты

| Компонент | Порт |
|-----------|------|
| api-gateway | 8000 |
| user-service | 8001 |
| receipt-service | 8002 |
| ai-processor | 8100 |
| analytics-service | 8101 |
| credit-service | 8009 |
| postgres | 5432 |
| kafka | 9092 |
| clickhouse HTTP | 8123 |
| whisper (profile `ai`) | 9001 |

## Проверка потока данных

```powershell
powershell -File scripts\verify_e2e.ps1
```

Сценарий: register → login → `POST /receipt/manual` → `GET /dashboard/categories` с суммой из БД.

## Подключение AI (фаза 2)

```powershell
# Whisper
docker compose --profile ai up -d whisper

# Gemini — ключ в .env (не коммитить)
$env:GEMINI_API_KEY = "your-gemini-key"
powershell -File scripts\stop_stack.ps1
powershell -File scripts\start_services.ps1
```

Voice (`POST /receipt/voice`) требует запущенный Whisper. Manual (`POST /receipt/manual`) работает без AI.

## Остановка

```powershell
powershell -File scripts\stop_stack.ps1          # только Go-сервисы
powershell -File scripts\stop_stack.ps1 -WithInfra # + docker compose down
```

## Устранение неполадок

| Симптом | Решение |
|---------|---------|
| `pgxpool` timeout | Docker Desktop запущен? `docker compose up -d postgres` |
| **502** на `/dashboard/*`, `/insights` | Поднять **receipt-service** (:8002) и **analytics-service** (:8101): `make up-full` или `docker compose up -d` |
| **500** на `PATCH /users/me/profile` | Миграции не применены: `make up-full` или `docker compose --profile tools run --rm migrate-pg`. Нужна таблица `user_financial_profiles` (009) |
| Dashboard пустой | `DEMO_MODE=false`, миграции применены, JWT одинаковый во всех сервисах |
| 502 на `/banks/*` | Пересобрать: `docker compose up -d --build bank-service` |
| Voice 503 | `docker compose --profile ai up -d whisper` |

Логи сервисов: `logs\*.log`, `logs\*.err`.
