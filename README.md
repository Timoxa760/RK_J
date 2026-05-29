# Поток

Голосовой помощник по личным тратам: рассказываете, что купили — сервис разбирает расходы и предлагает, что улучшить. Без таблиц, бухгалтерии и подключения банка.

Фронтенд на **Nuxt 4** (Vue 3, TypeScript, Tailwind, Pinia, ECharts, PWA).

## Быстрый старт

```bash
cd frontend
npm install
cp .env.example .env   # при необходимости
npm run dev
```

Приложение: [http://localhost:3000](http://localhost:3000)

Демо-вход: код **`0000`** на `/login`.

## Маршруты

| Путь | Доступ | Описание |
|------|--------|----------|
| `/` | публичный | Welcome, лендинг голосового помощника |
| `/login` | публичный | Вход |
| `/dashboard` | авторизация | Диагноз, графики расходов |
| `/receipts` | авторизация | Расходы / чеки |
| `/credits` | авторизация | Кредиты, DTI |
| `/analytics` | авторизация | Прогноз и аналитика |
| `/social` | авторизация | Социальные сценарии |
| `/digest` | авторизация | Дайджест |
| `/profile` | авторизация | Профиль |
| `/welcome` | — | редирект на `/` |

## Переменные окружения

См. [`frontend/.env.example`](frontend/.env.example):

| Переменная | По умолчанию | Назначение |
|------------|--------------|------------|
| `NUXT_PUBLIC_API_BASE` | `http://localhost:8000` | Базовый URL API |
| `NUXT_PUBLIC_DEMO_MODE` | `true` | Моки и demo-режим без бэкенда |

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

Фронтенд будет на [http://localhost:80](http://localhost:80) (nginx + статика из Nuxt).

## Структура

```
frontend/
├── app.vue
├── nuxt.config.ts
├── assets/css/           # design tokens, фон, UI
├── components/
│   ├── charts/           # ECharts
│   ├── shared/           # BackgroundFlow, HeroFlowWord, shell
│   └── …
├── composables/          # useDashboard, useCredits, …
├── constants/            # line1 paths, background flow
├── layouts/              # default (app), welcome (лендинг)
├── middleware/           # auth.global.ts
├── pages/
├── public/               # line1.svg, manifest, favicon
├── store/                # Pinia
└── types/                # API-типы
```

## API

Контракт бэкенда: [`API_Contract.md`](API_Contract.md).

## Стек

Nuxt 4 · Vue 3 · TypeScript · Tailwind CSS · Pinia · ECharts · Lucide · PWA (`@vite-pwa/nuxt`)
