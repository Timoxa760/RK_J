# Ввод расходов и pipeline

> Продукт: [input-methods.md](../product/input-methods.md)  
> Legacy Kafka pipeline — технический ingest, не продуктовый UX.

## Цель в «Поток»

Пользователь добавляет трату **голосом** или **вручную**; система интерпретирует и обновляет **финансовую модель**, показывая влияние на цель — не «15 000 ₽ на маркетплейсы», а «цель сдвигается на N месяцев».

## User Flow (UX)

1. Кнопка **«Добавить»** → голос / вручную
2. Система разбивает и категоризирует
3. Обновляется прогноз цели и один мягкий инсайт

## Технический pipeline (`back`)

**Основной путь (MVP):**

```
голос / ручной → ai-processor → manual_expenses → PostgreSQL + ClickHouse → dashboard
```

**Legacy (demo/mock, не в UX):**

```
scraper / mock → Kafka receipt.raw → receipt-service → ai-processor → ClickHouse
```

## Сервисы

| Этап | Сервис |
|------|--------|
| Голос, ручной | ai-processor → `manual_expenses` |
| Legacy ingest (mock, email) | scraper-service |
| Persist, dashboard | receipt-service |

## Связи

- **Инсайты**: [detective.md](./detective.md)
