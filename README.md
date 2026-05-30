# Поток

Голосовой помощник по личным тратам: рассказываете, что купили — сервис разбирает расходы и предлагает, что улучшить. Без таблиц, бухгалтерии и подключения банка.

Фронтенд на **Nuxt 4** (Vue 3, TypeScript, Tailwind, Pinia, ECharts, PWA).

## Ветки

| Ветка | Назначение |
|-------|------------|
| `main` | Базовая |
| `front` | Фронтенд Nuxt 4 |
| `back` | Бэкенд (Go + Python) |
| `docs` | Документация продукта и API |

## Документация

- [NAVI.md](NAVI.md) — карта репозитория
- [docs/api/API_Contract.md](docs/api/API_Contract.md) — контракт API
- [docs/product/](docs/product/) — UX-сценарии, онбординг
- [docs/features/](docs/features/) — спеки фич
- [docs/deployment/front-quickstart.md](docs/deployment/front-quickstart.md) — быстрый старт фронта

Локальная копия контракта (ветка `front`): [API_Contract.md](API_Contract.md).

## Быстрый старт

```bash
cd frontend
npm install
cp .env.example .env   # при необходимости
npm run dev
```

Приложение: [http://localhost:3000](http://localhost:3000)

Вход на `/login`. В demo-режиме код **`0000`**.

## Маршруты

| Путь | Доступ | Описание |
|------|--------|----------|
| `/` | публичный | Лендинг |
| `/login` | публичный | Вход |
| `/dashboard` | авторизация | Диагноз, графики расходов |
| `/receipts` | авторизация | Расходы / чеки |
| `/credits` | авторизация | Кредиты, DTI |
| `/analytics` | авторизация | Прогноз и аналитика |
| `/digest` | авторизация | Дайджест |
| `/profile` | авторизация | Профиль, магазины |

## Переменные окружения

См. [`frontend/.env.example`](frontend/.env.example):

| Переменная | По умолчанию | Назначение |
|------------|--------------|------------|
| `NUXT_PUBLIC_API_BASE` | `http://localhost:8000` | Базовый URL API (без `/api/v1`) |
| `NUXT_PUBLIC_DEMO_MODE` | `true` | Моки при недоступном API |

Для staging: задайте `NUXT_PUBLIC_API_BASE` на URL бэкенда и `NUXT_PUBLIC_DEMO_MODE=false`. Демо-тур: откройте любой экран приложения с `?tour=1` при включённом demo mode.

## Скрипты

```bash
npm run dev       # разработка
npm run build     # production-сборка
npm run preview   # превью сборки
npm run generate  # статическая генерация
```

## Docker

Из корня репозитория:

```bash
docker compose up --build
```

Фронтенд будет на [http://localhost:80](http://localhost:80).

## Стек

Nuxt 4 · Vue 3 · TypeScript · Tailwind CSS · Pinia · ECharts · Lucide · PWA (`@vite-pwa/nuxt`)
