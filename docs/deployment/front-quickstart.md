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
Вход: телефон и пароль на `/login`.

## Docker

```bash
# из корня front
docker compose up --build
# http://localhost:80
```

## Маршруты (MVP)

| Путь | Описание |
|------|----------|
| `/` | Welcome / landing |
| `/login` | Вход |
| `/onboarding` | Wizard (~1 мин), см. [onboarding.md](../product/onboarding.md) |
| `/dashboard` | Главный экран: narrative, план, метрики, «Что если» |
| `/advisor` | Полноэкранный чат советника |
| `/receipts` | Расходы, кнопка «Добавить» |
| `/credits` | Кредитный светофор (PDF scan) |
| `/profile` | Профиль, цели, ФНС (mock) |
| `/analytics` | Redirect → `/dashboard` |

**Не MVP:** `/social`, `/digest` — удалены из навигации.

## App shell

Flex-layout без shadcn offcanvas:

- `components/app/AppShellLayout.client.vue`
- `components/app/AppSidebar.vue` — `.mm-app-sidebar`
- `SharedBackgroundFlow` — только на `/` и `/login` (`layouts/default.vue`)

## Структура `frontend/`

```
frontend/
├── pages/           # маршруты
├── components/      # app/, advisor/, dashboard/, charts/
├── composables/     # useDashboard, useAdvisorChat, useAiPlan, …
├── store/           # Pinia + mocks при ошибке API
├── types/api.ts     # типы ответов API
└── nuxt.config.ts
```

## API

Контракт: [API_Contract.md](../api/API_Contract.md) (ветка `docs`).

`NUXT_PUBLIC_API_BASE=http://localhost:8000`

## Стек

Nuxt 4 · Vue 3 · TypeScript · Tailwind · Pinia · ECharts · PWA

## Связи

- [environment.md](./environment.md)
- [../product/ux-scenarios.md](../product/ux-scenarios.md)
- [../pitch/frontend.md](../pitch/frontend.md)
