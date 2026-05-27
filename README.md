# MONEYMIND

**Первый в России поштучный финансовый ассистент**

Хакатон Финтех — 29–31 мая 2026 | Кейсодатель: Клерк.Ру

---

## Ветки

| Ветка | Назначение |
|-------|------------|
| `main` | Базовая (README, точка входа) |
| `back` | Бэкенд: Go + Python микросервисы, инфраструктура |
| `front` | Фронтенд: React 19 + Vite + Tailwind CSS 4 |
| `docs` | Документация: архитектура, API, БД, фазы |

## Быстрый старт

```bash
# Бэкенд
git checkout back
docker compose up -d

# Фронтенд
git checkout front
npm install
npm run dev
```

## Ресурсы

- [Global Design Document](./MoneyMind_plan.md)
- [Архитектура](./docs/architecture/)
- [API спецификация](./docs/api/)
- [Схемы БД](./docs/database/)
- [Фичи](./docs/features/)
- [Фазы разработки](./docs/phases/)
