# MVP «Поток» — статус и demo

> Актуально на **2026-05-30**. Ветки: `front` (UI), `back` (API), `docs` (эта папка — источник правды по смыслу).

## Что это

**Поток** — финансовый навигатор: ответы и один план вместо перегруженной аналитики.  
Главный экран — **`/dashboard`**. Полноэкранный чат — **`/advisor`**.

## In scope (показываем на защите)

| Область | Статус | Примечание |
|---------|--------|------------|
| Онбординг 4 блока (skippable) | ✅ | `/onboarding` → profile |
| Голос + ручной расход | ✅ | `POST /expenses/manual`, `/receipt/voice` |
| PDF-скан кредита | ✅ | `POST /credits/scan` → DTI на `/credits` |
| ИИ-план + диагноз | ✅ | `GET /ai/plan`, `/ai/diagnosis` |
| Чат советника | ✅ | `/advisor`, streaming SSE |
| Dashboard narrative + «Что если» | ✅ | без отдельного `/digest` |
| ФНС «Мои чеки» | 🟡 mock | UI в `/profile`, 13 чеков ~50k ₽, любой SMS-код |

## Out of scope

- `goal-service`, CRUD `/goals/*` — цель только в profile
- `/social`, challenges, gamification
- bank-service sync для кредитов
- X5 Club / Magnit LK как ingest
- Отдельный маршрут `/digest` (контент на dashboard)

Подробнее: [scope.md](../product/scope.md).

## Mock vs production-ready

| Компонент | Demo / mock | Production-ready |
|-----------|-------------|------------------|
| Auth, JWT | ✅ | ✅ |
| Profile + skip-flags | file-store / PG | PG (миграции есть) |
| Расходы голос/ручной | ✅ | ✅ |
| Dashboard charts | API first, fallback mock | ✅ API |
| Credits PDF scan | OnlySQ + regex fallback | ✅ |
| Advisor plan/chat | Gemini + heuristic fallback | ✅ |
| ФНС connect | **front mock only** | API scraper — не в demo path |
| LLM dev | Antigravity :8045 опционально | Google Gemini direct |

## Запуск (3 команды)

**Backend** — [back-quickstart.md](../deployment/back-quickstart.md):

```bash
cd backend   # worktree `back`
make infra && make migrate && make up
```

**Frontend** — [front-quickstart.md](../deployment/front-quickstart.md):

```bash
cd frontend
npm install && cp .env.example .env && npm run dev
```

Открыть http://localhost:3000 → `/login` (телефон + пароль).

**Smoke:** `docs/deployment/scripts/demo_flow.sh` или `backend/scripts/smoke_critical.sh`.

## Demo script (~3 мин)

1. **Welcome** `/` — инсайт «ответы, не графики» → «Войти».
2. **Онбординг** `/onboarding` — доход / подушка / цель (можно skip) → dashboard.
3. **Dashboard** `/dashboard` — narrative, план (3 шага), метрики, симулятор «Что если».
4. **Расход** — «Добавить» → голос или ручной → обновление категорий.
5. **Советник** `/advisor` — «составь план» / «где урезать» → streaming ответ.
6. **Кредиты** `/profile` или `/credits` — загрузить PDF договора → DTI и сравнение ставки.
7. **(Опц.) ФНС** `/profile` — подключить «Мои чеки», любой код → импорт чеков.

## Питчи для защиты

| Аудитория | Файл | Время |
|-----------|------|-------|
| Тимлид / жюри | [pitch/teamlead.md](../pitch/teamlead.md) | 5–10 мин |
| Frontend | [pitch/frontend.md](../pitch/frontend.md) | 5–10 мин |
| Backend | [pitch/backend.md](../pitch/backend.md) | 5–10 мин |

Оглавление: [pitch/README.md](../pitch/README.md).

## Связи

- Фазы: [phases/phases.md](../phases/phases.md)
- UX: [product/ux-scenarios.md](../product/ux-scenarios.md)
- API: [api/API_Contract.md](../api/API_Contract.md)
- Архитектура: [architecture/overview.md](../architecture/overview.md)
