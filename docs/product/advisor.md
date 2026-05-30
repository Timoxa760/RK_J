# Финансовый советник (ИИ)

> Чат + план на `/dashboard`. API: [API_Contract.md](../api/API_Contract.md) § Advisor.

## Роль

Единая ИИ-модель «Поток» на базе **Google Gemini**:

- строит **финансовый план** (3 шага + цель + runway);
- даёт **диагноз** (score, indicators, main_action);
- отвечает в **чате** на вопросы по плану, цели, кредитам.

Контекст собирается **на сервере** (`UserFinanceSnapshot`), клиент шлёт только `message` + `history`.

## Источники snapshot

| Данные | Источник |
|--------|----------|
| Доход, подушка, цель, расходы | `user_financial_profiles` + skip-флаги |
| Кредиты, DTI | `user_credits` (только PDF-сканы) |
| Ставки vs рынок | rates-aggregator |
| Траты | `manual_expenses`, ClickHouse |

## Экран `/advisor`

Полноэкранный чат — вкладка в mobile tab bar и пункт в sidebar. CTA «Спросить советника» ведут на `/advisor?ask=`.

## Skip-aware

Если `skipped_income: true` — модель **не** трактует `active_income=0` как факт. Предлагает заполнить профиль или опирается на inferred данные.

## Типовые интенты чата

- «составь план» / «что делать»
- «где урезать» / «сократить траты»
- «когда дойду до цели»
- «ставка высокая» — при наличии scan + benchmark

## Эндпоинты

| Метод | Путь |
|-------|------|
| GET | `/api/v1/ai/plan` |
| GET | `/api/v1/ai/diagnosis` |
| GET | `/api/v1/ai/chat/history` |
| DELETE | `/api/v1/ai/chat/history` |
| POST | `/api/v1/ai/chat` |
| POST | `/api/v1/ai/chat/stream` |

Фронт: `useAdvisorChat`, `FinancialPlanCard`, `buildFinancialPlan` → данные с API.
