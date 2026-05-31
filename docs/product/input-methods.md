# Способы ввода расходов

## Ключевое UX-решение

Пользователь **не думает** о способе ввода. Видит только кнопку **«Добавить»**:

1. **Голосом** — свободная фраза
2. **Вручную** — сумма, магазин, категория

Модель строится на этих данных и постепенно уточняется через советника и dashboard.

---

## 1. Голосовой ввод (приоритет)

Свободный текст → структурирование транзакций.

**Пример:** «Вчера после работы зашла в пятёрочку, взяла продукты на 5000, платье за 3400 и кроссовки за 16000»

LLM (ai-processor):

- несколько транзакций;
- категории: groceries, clothing, sport/shoes;
- относительная дата («вчера»).

**Бэкенд:** `POST /api/v1/expenses/manual` с `source: voice` — [API_Contract.md §9](../api/API_Contract.md)

---

## 2. Ручной ввод

Поля: сумма, магазин, категория, дата.

**Бэкенд:** `POST /api/v1/expenses/manual` — [API_Contract.md §9](../api/API_Contract.md)

---

## 3. LK ритейлеров — **Out of MVP**

X5 Club, Magnit — **legacy-код** в `scraper-service`, не продуктовый ingest.

---

## Связи

- **Продукт**: [ux-scenarios.md](./ux-scenarios.md) №3
- **Фича**: [receipt-magic.md](../features/receipt-magic.md) (pipeline расходов)
- **Не менять в коде**: топики Kafka `receipt.raw` → `receipt.parsed` → enriched (legacy ingest)
