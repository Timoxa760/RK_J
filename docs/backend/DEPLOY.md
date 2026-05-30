# Деплой бэкенда на домашний сервер (Docker Compose)

> Для hunamuna123 и команды. Gateway — единственная точка входа для Vercel-фронта.

## Связи

- **Зависит от:** Docker, `.env`, OnlySQ API key, Cloudflare Tunnel (или HTTPS reverse proxy)
- **Используется:** Nuxt front (Vercel), `docs/api/API_Contract.md`
- **Связанные документы:** [API_Contract.md](../api/API_Contract.md), [STATUS.md](./STATUS.md)

---

## 1. Подготовка

```bash
git clone <repo-url> && cd backend_project
cp .env.example .env
```

Заполнить в `.env`:

| Переменная | Пример | Обязательно |
|------------|--------|-------------|
| `JWT_SECRET` | длинная случайная строка | да |
| `ONLYSQ_API_KEY` | ключ с api.onlysq.ru | да (для LLM и advice) |
| `DEMO_MODE` | `true` (demo) / `false` (PG) | да |
| `ONLYSQ_BASE_URL` | `https://api.onlysq.ru/v1` | по умолчанию ок |
| `WHISPER_URL` | `http://whisper:9000/v1/audio/transcriptions` | в compose по умолчанию |

---

## 2. Запуск

```bash
# инфраструктура + AI + gateway
docker compose up -d --build postgres redis clickhouse kafka zookeeper \
  whisper api-gateway user-service ai-processor \
  receipt-service analytics-service credit-service goal-service reporting-service

# миграции (при DEMO_MODE=false)
make migrate
```

Проверка:

```bash
curl -s http://localhost:8000/health
curl -s -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"+79991234567","code":"0000"}'
```

---

## 3. HTTPS для Vercel

Фронт на Vercel — HTTPS. Браузер блокирует HTTP API с HTTPS-страницы.

**Рекомендация:** Cloudflare Tunnel на `:8000`:

```bash
cloudflared tunnel --url http://localhost:8000
```

Полученный `https://....trycloudflare.com` передать фронтендеру:

```env
NUXT_PUBLIC_API_BASE=https://....trycloudflare.com/api/v1
```

---

## 4. Обновление после git pull

```bash
git pull origin back
docker compose up -d --build api-gateway user-service ai-processor whisper
```

---

## 5. Минимальный набор для «Поток» (голос)

| Сервис | Порт | Наружу |
|--------|------|--------|
| api-gateway | 8000 | да (через tunnel) |
| user-service | 8001 | нет |
| ai-processor | 8100 | нет |
| whisper | 9000 | нет |
| postgres | 5432 | нет |

Whisper и OnlySQ доступны ai-processor только внутри Docker-сети `moneymind_network`.

---

## 6. Smoke

```bash
export DEMO_MODE=true JWT_SECRET=<из .env>
./scripts/smoke_api.sh
bash scripts/smoke_critical.sh
```

---

## 7. Troubleshooting

| Симптом | Причина |
|---------|---------|
| 401 на API | Нет JWT или другой `JWT_SECRET` |
| 503 на `/expenses/voice` | Whisper не поднят или `WHISPER_URL` неверный |
| Нет `advice` в ответе | Нет `ONLYSQ_API_KEY` — работает regex fallback |
| Vercel не достучался | Нет HTTPS tunnel / неверный URL в env фронта |
