# Scope продукта «Поток» (MVP)

> Актуально для веток `front` + `back`. Полный контракт: [API_Contract.md](../api/API_Contract.md).

## In scope

| Область | Описание |
|---------|----------|
| Онбординг | 4 блока (доход, подушка, цель, расходы), skippable → **profile** |
| Финансовый профиль | `GET/PATCH /users/me/profile`, skip-флаги, цель в профиле |
| Расходы | Голос, ручной ввод |
| ФНС «Мои чеки» | **UI mock на front** — SMS flow, импорт тестовых чеков; back API не в demo path |
| Кредиты | **Только** PDF-скан договора → `POST /credits/scan` → dashboard |
| Ставки | **rates-aggregator** — этalon vs ставка из скана |
| Советник | `/ai/plan`, `/ai/diagnosis`, `/ai/chat`, SSE stream |
| Dashboard | Narrative, план, метрики, симулятор «Что если»; чат на `/advisor` |

## Out of scope (removed)

| Было | Статус |
|------|--------|
| `goal-service`, CRUD `/goals/*` | **Removed** — цель только в profile |
| `/social`, challenges, gamification | **Removed** |
| `GET /ai/goal/{id}/forecast` | **Removed** — прогноз в `/ai/plan` |
| Demo/stub credits dashboard | **Removed** — данные только из сканов |
| bank-service sync для кредитов | **Not in MVP** |
| X5 Club / Magnit LK как ingest | **Not in MVP** |
| Отдельный `/digest` | **Removed** — narrative на dashboard |

## Dual-branch

| Папка | Ветка |
|-------|-------|
| `frontend/`, локальные `docs/` | `front` |
| `backend/` worktree | `back` |
| Источник правды docs | `docs` (worktree `RK_J-docs`) |

См. [git-and-branches.md](../guides/git-and-branches.md), [mvp/README.md](../mvp/README.md).
