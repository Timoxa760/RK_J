# 3.4 Финансовый детектив

AI сканирует всю историю трат и находит утечки.

## Детектируемые паттерны

| Паттерн | Логика |
|---------|--------|
| Подписки | `GROUP BY amount HAVING count > 3 AND deviation < 5%` |
| Дубли | Одинаковые item.name в одном чеке |
| Переплаты | `item.price > avg(price over 30d)` per store |

## Pipeline

```
Результаты → Kafka insight.found → notification-service → push/Telegram
```
