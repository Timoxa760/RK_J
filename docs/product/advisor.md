# Финансовый советник (ИИ)

> Чат на `/advisor`, план и narrative на `/dashboard`. API: [API_Contract.md](../api/API_Contract.md) § Advisor. Архитектура: [advisor-system.md](../architecture/advisor-system.md).

## Роль

Единая ИИ-модель «Поток» (LLM через Google Gemini **или** Antigravity → Claude):

- строит **финансовый план** (3 шага + цель + runway);
- даёт **диагноз** (score, indicators, main_action);
- отвечает в **чате** с streaming, actions и историей в PostgreSQL.

Контекст собирается **на сервере** (`UserFinanceSnapshot`), клиент шлёт только `message` + `history`.

## Экраны

| Маршрут | Содержание |
|---------|------------|
| `/dashboard` | PageNarrative (совет недели, mindfulness score, доход/траты), план, метрики, симулятор «Что если» |
| `/advisor` | Полноэкранный чат — вкладка в mobile tab bar и пункт sidebar |

CTA «Спросить советника» → `/advisor?ask=…`. Встроенный чат в sidebar **удалён** (2026-05-31).

## Источники snapshot

| Данные | Источник |
|--------|----------|
| Доход, подушка, цель, расходы | `user_financial_profiles` + skip-флаги |
| Кредиты, DTI | `user_credits` (только PDF-сканы) |
| Ставки vs рынок | rates-aggregator |
| Траты | `manual_expenses`, ClickHouse |

## Skip-aware

Если `skipped_income: true` — модель **не** трактует `active_income=0` как факт.

## Chat features

| Возможность | Реализация |
|-------------|------------|
| Streaming | `POST /ai/chat/stream` (SSE) |
| История | `GET/DELETE /ai/chat/history` |
| Actions | navigate, add_expense, followup |
| Fallback | `source: heuristic` при недоступности LLM |
| Badge | UI показывает source (gemini / heuristic) |

## Эндпоинты

| Метод | Путь |
|-------|------|
| GET | `/api/v1/ai/plan` |
| GET | `/api/v1/ai/diagnosis` |
| GET | `/api/v1/ai/chat/history` |
| DELETE | `/api/v1/ai/chat/history` |
| POST | `/api/v1/ai/chat` |
| POST | `/api/v1/ai/chat/stream` |

Фронт: `useAdvisorChat`, `useAiPlan`, `FinancialPlanCard`, `AdvisorChatActions`.

## LLM в dev

Antigravity Tools на `:8045`, модель `claude-sonnet-4-6`. См. [antigravity-setup.md](../deployment/antigravity-setup.md).
