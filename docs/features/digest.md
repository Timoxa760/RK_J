# Дайджест и главное действие недели

> **Статус (2026-05-31):** отдельная страница `/digest` **удалена**. Контент перенесён на `/dashboard` в блок `PageNarrative` (совет недели, mindfulness score, доход/траты).

## Состав (на dashboard)

- Текущая картина (narrative из diagnosis + insights)
- **Mindfulness score** — из AI diagnosis (`score/100`)
- **Главное действие недели** — `main_action` из `/ai/diagnosis`
- Доход / траты — отдельные карточки в narrative

## API (legacy)

`GET /api/v1/digest/*` — reporting-service (optional, не используется front MVP).

## UI

| Было | Стало |
|------|-------|
| `/digest` | `/dashboard` → `PageNarrative`, `DashboardMindfulnessScore` |

## Связи

- [ux-scenarios.md](../product/ux-scenarios.md) №1, №4
- [advisor.md](../product/advisor.md)
