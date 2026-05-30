# Changelog

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.1.0/).

## [Unreleased]

### Добавлено

- `POST /api/v1/expenses/voice` — голосовой ввод через Whisper + OnlySQ.
- Поля `advice`, `expenses[]`, `parsed_by`, `transcript` в ответах expenses API.
- Клиент OnlySQ (`internal/onlysq`), клиент Whisper, LLM-парсер с fallback на regex.
- Сервис Whisper в docker-compose, `docs/backend/DEPLOY.md`.
- Скрипты полного локального стека: `start_infra.ps1`, `migrate.ps1`, `start_stack.ps1`, `verify_e2e.ps1`.
- Demo API в bank-, category-, budget-, social-service (контракт optional).
- Руководство `docs/guides/local-full-stack.md`.

### Изменено

- `POST /api/v1/expenses/manual` — парсинг через OnlySQ при наличии API-ключа.
- `start_services.ps1` — все 15 сервисов, `DEMO_MODE=false`, Kafka/ClickHouse env.
- docker-compose: migrate profile `tools`, Whisper на host-порту 9001 (profile `ai`).
