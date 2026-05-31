# Питч для frontend-разработчика (5–10 мин)

## 1. Стек и структура (1 мин)

| Слой | Технология |
|------|------------|
| Framework | **Nuxt 4.3**, Vue 3, TypeScript |
| State | Pinia (`store/`) |
| UI | Tailwind, shadcn-vue (точечно), PrimeVue |
| Charts | ECharts + обёртки `components/charts/` |
| PWA | `manifest.json`, offline shell |

```
frontend/
├── pages/           # маршруты (file-based routing)
├── components/      # app/, advisor/, dashboard/, charts/
├── composables/     # useDashboard, useAdvisorChat, …
├── utils/           # advisorMarkdown, financialPlan, dashboardSummary
├── store/           # dashboardStore + mocks/
└── layouts/         # default, onboarding, welcome
```

Запуск: `cd frontend && npm install && npm run dev` → http://localhost:3000  
Док: [front-quickstart.md](../deployment/front-quickstart.md).

---

## 2. App shell — без shadcn offcanvas (1 мин)

**Проблема:** встроенный `SidebarProvider` уводил иконки за viewport на mobile.

**Решение:** простой flex-layout в `AppShellLayout.client.vue`:

```
.mm-app-shell (flex row)
├── AppSidebar.vue     — <aside class="mm-app-sidebar">
└── .mm-app-shell-content — slot страницы + advisor host
```

- Навигация: plain `<NuxtLink>` с классом `.mm-app-sidebar-link`.
- Бренд «Поток» — по центру header sidebar, симметричные отступы.
- Mobile: tab bar (`MobileTabBar.vue`) + hamburger → sidebar overlay.

Стили: `assets/css/main.css` (`.mm-app-sidebar`, `.mm-app-shell`).

---

## 3. Маршруты MVP (1 мин)

| Путь | Страница | MVP |
|------|----------|-----|
| `/` | Landing | ✅ welcome layout |
| `/login` | Auth | ✅ |
| `/onboarding` | Wizard | ✅ |
| `/dashboard` | Главный экран | ✅ narrative + plan + charts |
| `/advisor` | Чат советника | ✅ tab + sidebar link |
| `/receipts` | Список расходов | ✅ |
| `/credits` | Кредитный светофор | ✅ |
| `/profile` | Профиль, ФНС, цели | ✅ |
| `/analytics` | — | redirect → `/dashboard` |
| `/social` | — | redirect / stub, **не MVP** |

**BackgroundFlow** (декоративные линии) — только на bare pages (`/`, `/login`) в `layouts/default.vue`, не в app shell.

---

## 4. Ключевые composables (2 мин)

| Composable | Назначение |
|------------|------------|
| `useDashboard` | sankey, categories, timemachine — API first |
| `useAiPlan` | `GET /ai/plan` — plan + diagnosisFromPlan |
| `useAdvisorChat` | `POST /ai/chat`, SSE stream, history |
| `useAdvisorContext` | контекст для CTA «Спросить» |
| `useOpenAdvisorChat` | navigate `/advisor?ask=…` |
| `useFinancialProfile` | profile + skip-flags |
| `useCredits` | scan PDF, dashboard DTI |
| `useFns` | **mock** ФНС: SMS flow, import 13 receipts |
| `useAddExpenseSheet` | sheet «Добавить» голос/ручной |

**Принцип:** advisor-поток без клиентских моков плана — данные с `/ai/plan`. Charts могут fallback на `store/mocks/` при ошибке API.

---

## 5. Advisor UI (2 мин)

**Компоненты:**

- `AdvisorChat.vue` — сообщения, streaming, actions
- `FinancialPlanCard.vue` — цель, 3 шага, runway
- `AdvisorAskButton.vue` — CTA на dashboard / receipts
- `AppShellAdvisorHost.vue` — host для overlay (desktop)

**Structured blocks:** ответ LLM парсится в `utils/advisorStructured.ts` → markdown + action chips.

**Text repair:** `utils/advisorMarkdown.ts` — allowlist-исправление «склеенных» русских слов (без агрессивного split по «на/за/до» внутри морфем). Тесты: `npm run test:advisor-markdown`.

**Badge:** UI показывает `source: gemini | heuristic` из ответа API.

---

## 6. Credits + ФНС UI (1 мин)

- **`/credits`** — upload PDF, DTI gauge, список договоров.
- **`ProfileFnsSection.vue`** — двухшаговый dialog: телефон → SMS (любой код) → import mock receipts ~50k ₽.
- Copy без слова «demo» — см. `constants/productCopy.ts`.

---

## 7. Demo checklist (1 мин)

1. `NUXT_PUBLIC_API_BASE=http://localhost:8000`
2. Login → onboarding → dashboard loads plan
3. Add expense → categories refresh
4. `/advisor` → streaming works
5. Upload credit PDF on `/credits`
6. (Opt) FNS connect in profile

---

## 8. Known limitations

| Limitation | Workaround |
|------------|------------|
| PWA cache stale | hard refresh / disable SW in dev |
| Chart mocks on API fail | intentional fallback in composables |
| ФНС — front mock only | не вызывает back FNS endpoints в demo |
| `/digest`, `/social` removed from nav | redirects or stubs |

---

## Связи

- Product advisor: [product/advisor.md](../product/advisor.md)
- API types: [api/typescript-types.md](../api/typescript-types.md)
- UX: [product/ux-scenarios.md](../product/ux-scenarios.md)
