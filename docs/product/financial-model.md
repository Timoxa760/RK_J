# Единая финансовая модель

Формируется из **profile**, PDF-сканов кредитов и транзакций.

## Компоненты

| Компонент | Источник |
|-----------|----------|
| Доход | profile (`active_income`, `passive_income`; skip → unknown) |
| Подушка | profile (`emergency_fund`, breakdown) |
| Цель | profile (`goal_kind`, `goal_title`, `goal_amount`) — **не goal-service** |
| Обязательные расходы | profile `fixed_expenses` или inferred из трат |
| Кредиты | **только** `user_credits` из PDF scan |
| Фактические траты | manual_expenses, receipt_items |

## UserFinanceSnapshot (для ИИ)

```go
Profile          // incl. goal_*, skipped_*
Credits          // из PDF scans
RateBenchmarks   // scan rate vs aggregator
Metrics          // income, expenses, freeCashflow, runwayMonths, dti
RecentExpenses   // категории 30/90 дней
DataCompleteness // known | skipped | inferred
```

## data_completeness

| Поле | skipped | inferred | known |
|------|---------|----------|-------|
| income | не считать 0 | credits / median spend | profile |
| cushion | runway=null | — | emergency_fund |
| goal | generic plan | — | profile.goal_* |
| expenses | fixed=0 | receipts sum | fixed_expenses |
| credits | empty dashboard | — | после scan |

## Расчёты

- **freeCashflow** = income − expenses
- **runwayMonths** = emergency_fund / expenses (null если skipped_cushion или expenses=0)
- **dti** = monthly_payments / monthly_income × 100 (из scans + profile income)

## Цикл обновления

1. Действие (расход, scan PDF, PATCH profile)
2. Обновление snapshot
3. Invalidate `/ai/plan` cache
4. Dashboard + advisor

См. [advisor.md](./advisor.md), [credit-scan.md](../features/credit-scan.md).
