# Ввод расходов и pipeline чеков

> Продукт: [input-methods.md](../product/input-methods.md)  
> Старый заголовок «Receipt Magic» — технический pipeline ingest, не позиционирование.

## Цель в «Поток»

Пользователь добавляет трату одним из способов; система интерпретирует и обновляет **финансовую модель**, показывая влияние на цель — не «15 000 ₽ на маркетплейсы», а «цель сдвигается на N месяцев».

## User Flow (UX)

1. Кнопка **«Добавить»** → голос / чек / ФНС (на выбор)
2. Система разбивает и категоризирует
3. Обновляется прогноз цели и один мягкий инсайт

## Технический pipeline (`back`)

```
scraper / manual / voice
        ↓
  Kafka receipt.raw
        ↓
receipt-service: validate → dedup (sha256) → receipt.parsed
        ↓
ai-processor: categorize → enriched
        ↓
PostgreSQL + ClickHouse → dashboard / analytics
```

## Сервисы

| Этап | Сервис |
|------|--------|
| FNS, X5, Magnit, email | scraper-service |
| Голос, ручной | ai-processor → `manual_expenses` |
| Persist, dashboard | receipt-service |

## Связи

- **ФНС**: [fotochecking.md](./fotochecking.md)
- **Инсайты**: [detective.md](./detective.md)
