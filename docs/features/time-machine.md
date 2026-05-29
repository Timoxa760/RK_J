# Траектория цели (Time Machine)

> Продукт: «если ничего не менять, путь будет таким» vs «если сократить X — цель ближе на N».

## Смысл для пользователя

Не симулятор ради графика — **сравнение двух будущих** с понятным выводом.

## API

- `GET /api/v1/dashboard/timemachine` — receipt-service
- `POST /api/v1/scenarios/simulate` — analytics-service (ипотека, крупные решения)

## UI

`front` — `ChartsTimeMachineChart` на `/dashboard`.

## Связи

- [predictive-ai.md](./predictive-ai.md)
- [monetization.md](../product/monetization.md)
