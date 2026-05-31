# Kafka — Event Bus

> Источник: `back` — scraper-service, receipt-service, ai-processor.  
> Брокер: Kafka 4.0.2 (`cp-kafka`), клиент: `segmentio/kafka-go`.

## Топики (факт в коде)

| Топик | Producer | Consumer | Group ID | Payload |
|-------|----------|----------|----------|---------|
| `receipt.raw` | scraper-service | receipt-service | `receipt-service` | `RawReceipt` (JSON) |
| `receipt.parsed` | receipt-service | ai-processor | `ai-processor` | `RawReceipt` + validation |

### Планируемые (docs / legacy GDD)

| Топик | Назначение | Статус |
|-------|------------|--------|
| `receipt.enriched` | После категоризации → CH | 🟡 частично (CH write напрямую) |
| `insight.found` | Инсайты детектива | ⏳ roadmap |

## Pipeline

```
scraper / mock / email (legacy)
        │
        ▼
  receipt.raw ──► receipt-service
        │            validate → dedup (sha256)
        │            INSERT receipts (PG)
        ▼
 receipt.parsed ──► ai-processor
        │            categorize (dict / LLM)
        │            INSERT receipt_items (CH)
        │            manual_expenses (PG, HTTP path)
        ▼
   dashboard API / analytics
```

## Формат сообщения (ориентир)

```json
{
  "id": "uuid",
  "user_id": "string",
  "provider": "x5club|mock|email",
  "store_name": "Пятёрочка",
  "total_amount": 1520.50,
  "purchased_at": "2026-05-30T12:00:00Z",
  "items": [
    {"name": "Молоко 3.2%", "price": 78, "quantity": 2}
  ]
}
```

## Consumer groups

| Group | Сервис | Назначение |
|-------|--------|------------|
| `receipt-service` | receipt-service | Persist + dedup |
| `ai-processor` | ai-processor | Categorize + ClickHouse |

При падении consumer Kafka хранит offset — после рестарта дочитывает backlog.

## HTTP vs Kafka

| Путь | Транспорт | Почему |
|------|-----------|--------|
| Чеки из scraper | Kafka | Batch, async, replay |
| `POST /expenses/manual` | HTTP | Нужен ответ пользователю сразу |

## Связи

- [defense.md §4](./defense.md) — почему Kafka
- [../features/receipt-magic.md](../features/receipt-magic.md)
- [../deployment/environment.md](../deployment/environment.md) — `KAFKA_URL`
