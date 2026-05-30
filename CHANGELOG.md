# Changelog

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.1.0/).

## [Unreleased]

### Добавлено

- `POST /api/v1/expenses/voice` — голосовой ввод через Whisper + OnlySQ.
- Поля `advice`, `expenses[]`, `parsed_by`, `transcript` в ответах expenses API.
- Клиент OnlySQ (`internal/onlysq`), клиент Whisper, LLM-парсер с fallback на regex.
- Сервис Whisper в docker-compose, `docs/backend/DEPLOY.md`.

### Изменено

- `POST /api/v1/expenses/manual` — парсинг через OnlySQ при наличии API-ключа.
