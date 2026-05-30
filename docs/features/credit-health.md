# Финансовое здоровье (Credit Health → Health Dashboard)

> Продукт: [financial-health.md](../product/financial-health.md)

## Цель

Пользователь видит **устойчивость и запас**, а не бухгалтерский отчёт:

- финансовое здоровье 🟢/🟡/🔴;
- runway (подушка в месяцах);
- влияние текущих трат на цель.

## Расчёты

| Метрика | Формула / логика |
|---------|------------------|
| DTI | `sum(monthly_payments) / avg(income)` |
| Runway | `savings / avg(monthly_spend)` |
| Устойчивость | composite 0–100 из дохода, обязательств, подушки |

## Реализация

| Слой | Компонент |
|------|-----------|
| API | `credit-service` — `/api/v1/credits/` |
| Dashboard | `receipt-service` — sankey, compare |
| UI | `front` — `/credits`, блоки на `/dashboard` |

## Подача (философия)

Без давления: «у тебя запас на X месяцев», «траты замедляют путь к цели».

## Связи

- [philosophy.md](../product/philosophy.md)
- [monetization.md](../product/monetization.md) — платные сценарии на базе здоровья
