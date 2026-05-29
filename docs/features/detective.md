# Поведенческие инсайты (Финансовый детектив)

> Продукт: причины, не цифры — «маркетплейсы отодвигают цель на 3 месяца».

## Примеры инсайтов

| Паттерн | Пример формулировки |
|---------|---------------------|
| После зарплаты | «Чаще покупаете в первые 5 дней после дохода» |
| Импульс | «40% свободных денег — спонтанные покупки» |
| Подписки | «3 похожих списания — проверьте подписки» |
| Маркетплейсы | «MP отодвигают цель на 3 месяца» |

## Pipeline

```
analytics-service → insights
        ↓
notification-service (optional push)
        ↓
front: InsightsPanel, /analytics
```

## API

`GET /api/v1/insights/*` — analytics-service

## Связи

- [ux-scenarios.md](../product/ux-scenarios.md) №3, №4
- [philosophy.md](../product/philosophy.md)
