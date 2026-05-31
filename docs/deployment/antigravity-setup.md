# Настройка Antigravity Tools для «Поток»

> LLM-архитектура: [../architecture/llm-integration.md](../architecture/llm-integration.md).

## 1. Установка

Antigravity Tools — локальное приложение с API-прокси на порту **8045**.

1. Установить Antigravity Tools (GUI).
2. Войти Google-аккаунт с Code Assist (если нужен доступ к моделям Google).
3. **API Proxy** → скопировать API Key → `GEMINI_API_KEY` в `backend/.env`.

## 2. Конфиг backend

```bash
cd backend
cp .env.example .env
```

```env
GEMINI_PROVIDER=antigravity
GEMINI_API_KEY=sk-...          # из Antigravity API Proxy
GEMINI_MODEL=claude-sonnet-4-6
GEMINI_BASE_URL=http://127.0.0.1:8045/v1
```

Docker подхватывает ключ из `.env`, но **модель и base URL для контейнеров** заданы в `docker-compose.yml`:

- `GEMINI_MODEL=claude-sonnet-4-6`
- `GEMINI_BASE_URL=http://host.docker.internal:8045/v1`

## 3. Upstream proxy (опционально)

Если native Gemini блокируется по IP, в Antigravity можно включить upstream SOCKS/HTTP proxy:

```json
{
  "upstream_proxy": {
    "enabled": true,
    "url": "socks5://127.0.0.1:10808"
  },
  "proxy": {
    "enabled": true,
    "auto_start": true
  }
}
```

Файл конфига: `~/.antigravity_tools/gui_config.json` (или через GUI).

> Даже с VPN часть аккаунтов Google Code Assist остаётся заблокированной для `gemini-*`. Используйте OpenAI-маршрут с `claude-sonnet-4-6`.

## 4. Проверка

```bash
curl -s http://127.0.0.1:8045/healthz
# {"status":"ok","version":"..."}

curl -s -X POST http://127.0.0.1:8045/v1/chat/completions \
  -H "Authorization: Bearer $GEMINI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-sonnet-4-6","messages":[{"role":"user","content":"say OK"}],"max_tokens":5}'
```

Из Docker-контейнера:

```bash
docker exec backend_ai_processor wget -qO- \
  --header="Authorization: Bearer $GEMINI_API_KEY" \
  --header="Content-Type: application/json" \
  --post-data='{"model":"claude-sonnet-4-6","messages":[{"role":"user","content":"OK"}],"max_tokens":5}' \
  http://host.docker.internal:8045/v1/chat/completions
```

## 5. Запуск стека

```bash
cd backend
docker compose up -d postgres redis user-service api-gateway ai-processor
./scripts/smoke_auth_chat.sh
```

Ожидание в шаге 5: `source: gemini`.

## Troubleshooting

| Симптом | Причина | Решение |
|---------|---------|---------|
| `source: heuristic` | LLM недоступен или ошибка | Проверить healthz, ключ, модель |
| `location is not supported` | gemini-* на аккаунте | Переключить на `claude-sonnet-4-6` |
| Контейнер не достучится до :8045 | `127.0.0.1` внутри контейнера | `host.docker.internal:8045/v1` |
| Shell `GEMINI_MODEL=gemini-*` | перебивает `.env` | `unset GEMINI_MODEL` или полагаться на compose |

## Связи

- [environment.md](./environment.md)
- [back-quickstart.md](./back-quickstart.md)
