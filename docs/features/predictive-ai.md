# Прогноз цели (Predictive AI)

> Продукт: «через сколько достигну цель при текущем поведении».

## Механика

- Временной ряд расходов (ClickHouse + PG)
- Прогноз свободного потока и даты достижения `target_amount`
- Anomaly: трата сильно выше нормы → мягкий алерт

## API

- `analytics-service` — `/api/v1/forecast/*`
- Dashboard trajectory — `/api/v1/dashboard/timemachine`

## UI (`front`)

- `/analytics` — ForecastChart
- `/dashboard` — TimeMachineChart («если ничего не менять»)

## Связи

- [financial-model.md](../product/financial-model.md)
- [time-machine.md](./time-machine.md)
