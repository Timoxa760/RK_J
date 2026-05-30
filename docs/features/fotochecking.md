# Скан чека

> Часть [input-methods.md](../product/input-methods.md) — «Сфотографировать чек».

## UX

Пользователь не думает о pipeline — только «Добавить → сфоткать чек».

## Реализация

- `ai-processor` — парсинг позиций
- `receipt-service` — persist + dedup
- QR/ФНС fallback — `scraper-service` `/api/v1/fns/*`

## Связи

- [receipt-magic.md](./receipt-magic.md)
