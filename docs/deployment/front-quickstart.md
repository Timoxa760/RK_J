# Frontend — быстрый старт

> Ветка **`front`**. Nuxt 4 «Поток».

## Команды

```bash
git checkout front
cd frontend
npm install
cp .env.example .env
npm run dev
```

Приложение: http://localhost:3000  
Демо-вход: код **`0000`** на `/login`.

## Docker

```bash
# из корня front
docker compose up --build
# http://localhost:80
```

## Маршруты

| Путь | Описание |
|------|----------|
| `/` | Welcome |
| `/login` | Вход |
| `/onboarding` | ⏳ wizard (спека в [onboarding.md](../product/onboarding.md)) |
| `/dashboard` | Главный экран |
| `/receipts` | Расходы |
| `/credits` | Кредитный светофор |
| `/analytics` | Прогноз |
| `/digest` | Дайджест |
| `/social` | Гипотеза |
| `/profile` | Профиль |

## Структура `frontend/`

```
frontend/
├── pages/           # маршруты
├── components/      # charts/, shared/
├── composables/     # useDashboard, useCredits, …
├── store/           # Pinia
├── types/api.ts     # типы ответов API
└── nuxt.config.ts
```

## API

Контракт: [API_Contract.md](../api/API_Contract.md) (ветка `docs`).

## Стек

Nuxt 4 · Vue 3 · TypeScript · Tailwind · Pinia · ECharts · PWA

## Связи

- [environment.md](./environment.md)
- [../product/ux-scenarios.md](../product/ux-scenarios.md)
