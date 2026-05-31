# LLM-интеграция (Gemini + Antigravity Tools)

> Источник правды по коду: `backend/internal/llm/`. Env: [environment.md](../deployment/environment.md), настройка прокси: [antigravity-setup.md](../deployment/antigravity-setup.md).

## Зачем два режима

| Режим | Base URL | Протокол | Когда |
|-------|----------|----------|-------|
| **Google direct** | `https://generativelanguage.googleapis.com/v1beta` | Native Gemini `generateContent` / `streamGenerateContent`, ключ в query `?key=` | Прод с Google AI Studio, регион без блокировок |
| **Antigravity Tools** | `http://127.0.0.1:8045/v1` (локально) | OpenAI-compatible `POST /v1/chat/completions` (+ SSE stream) | Dev, локальный прокси к Code Assist / Claude / GPT через один API key |

Клиент (`llm.Client`) выбирает режим по base URL:

- содержит `generativelanguage.googleapis.com` → native Gemini;
- иначе → OpenAI-маршрут с заголовками `Authorization: Bearer` и `x-api-key`.

## Antigravity Tools (рекомендуется для dev)

**Antigravity Tools** — локальный GUI-прокси (порт **8045**), который мультиплексирует запросы на Google Code Assist и другие модели.

### Почему не native Gemini через Antigravity

На части Google-аккаунтов Code Assist возвращает:

```text
User location is not supported for the API use. (FAILED_PRECONDITION)
```

для всех `gemini-*` моделей через native `/v1beta/models/gemini-*:generateContent`.

**Рабочий обход:** OpenAI-маршрут `/v1/chat/completions` с моделями вроде:

- `claude-sonnet-4-6` ✅
- `gpt-oss-120b-medium` ✅
- `gemini-*` через OpenAI route — часто ❌ (тот же location block)

### Переменные окружения

| Переменная | Локально | Docker (`ai-processor`, `credit-service`) |
|------------|----------|-------------------------------------------|
| `GEMINI_PROVIDER` | `antigravity` | `antigravity` |
| `GEMINI_API_KEY` | API Key из Antigravity → API Proxy | тот же ключ |
| `GEMINI_MODEL` | `claude-sonnet-4-6` | **зафиксировано** в `docker-compose.yml` |
| `GEMINI_BASE_URL` | `http://127.0.0.1:8045/v1` | `http://host.docker.internal:8045/v1` |

> **Важно:** shell-переменная `GEMINI_MODEL=gemini-*` перебивает `.env` при `docker compose up`. В compose модель зафиксирована как `claude-sonnet-4-6`, чтобы не попасть в location block.

### Нормализация URL

Если в `.env` указан legacy `.../v1beta` для Antigravity, клиент автоматически заменяет на `.../v1`.

## Где используется LLM

| Функция | Пакет | Fallback |
|---------|-------|----------|
| Чат советника | `internal/advisor/chat.go` | `heuristicChat` |
| Stream чат | `client.StreamComplete` + SSE handler | heuristic |
| Финансовый план / диагноз | `internal/advisor/plan.go` | шаблоны |
| Парсинг голоса / ручного расхода | `ai-processor/internal/expense` | regex / dict |
| Категоризация чеков | dict (MVP), LLM — фаза 2 | словарь |

Успешный ответ LLM в API помечается `source: "gemini"` (историческое имя; фактически может быть Claude через Antigravity).

## Проверка

```bash
# Antigravity жив
curl -s http://127.0.0.1:8045/healthz

# OpenAI route
curl -s -X POST http://127.0.0.1:8045/v1/chat/completions \
  -H "Authorization: Bearer $GEMINI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-sonnet-4-6","messages":[{"role":"user","content":"OK"}],"max_tokens":10}'

# E2E backend
cd backend && ./scripts/smoke_auth_chat.sh
# ожидание: source: gemini
```

## Связи

- [advisor-system.md](./advisor-system.md) — snapshot → chat → stream
- [overview.md](./overview.md) — место ai-processor в архитектуре
- [../product/advisor.md](../product/advisor.md) — продуктовая роль советника
