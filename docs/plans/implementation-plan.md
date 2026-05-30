# Implementation Plan — Whisper + OnlySQ для «Поток»

## Цель и контекст

Подключить источник данных для голосового ассистента: Whisper (STT на сервере hunamuna123 в Docker) и OnlySQ (LLM для парсинга трат и совета). Единая точка входа — `ai-processor` через gateway `:8000`.

## Объём работ

**Входит:** API Contract, OnlySQ client, LLM-парсер с fallback на regex, `POST /expenses/voice`, docker-compose env, DEPLOY.md.

**Не входит (следующие итерации):** OnlySQ categorizer для чеков, live dashboard из PG, credits scan LLM.

## Архитектурное решение

```
Фронт → api-gateway → ai-processor
                          ├→ whisper (Docker, внутр. сеть)
                          ├→ OnlySQ API (HTTPS)
                          └→ postgres (prod) / in-memory (demo)
```

Парсинг: OnlySQ → regex fallback. Сохранение: все `expenses[]` из ответа LLM.

## Стек

- Go 1.25, `net/http` для OnlySQ/Whisper (OpenAI-совместимые API)
- Docker Compose: существующие сервисы + `whisper` (openai-whisper-asr-webservice)
- OnlySQ: `https://api.onlysq.ru/v1`, модель `gpt-4o-mini` по умолчанию

## Контракты

См. `docs/api/API_Contract.md` — секция 9 (manual + voice).

## Риски

- Latency Whisper+LLM: увеличен `WriteTimeout` ai-processor до 120s
- Без `ONLYSQ_API_KEY` — regex fallback, без `advice` от LLM
- Whisper-образ требует CPU/RAM на домашнем сервере

## Связи

- **Зависит от:** OnlySQ API, Whisper container, api-gateway
- **Используется:** Nuxt front (Vercel)
- **Связанные документы:** [API_Contract.md](../api/API_Contract.md), [DEPLOY.md](../backend/DEPLOY.md)
