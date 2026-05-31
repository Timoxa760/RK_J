# Питч для тимлида / жюри (5–10 мин)

## 1. Elevator pitch (30 сек)

**Поток** — финансовый навигатор для людей 25–45 лет, которые не хотят вести Excel, но хотят понять: *куда уходят деньги и что сделать сегодня, чтобы быстрее достичь цели*.

Мы не показываем десять графиков. Мы даём **один план**, **один диагноз** и **один чат**, который знает контекст пользователя.

> Инсайт: людям нужны **ответы**, не таблицы.

---

## 2. Проблема и аудитория (1 мин)

- Регулярный доход, но нет ощущения контроля.
- Классические приложения требуют ручного учёта — пользователи бросают через несколько дней.
- Целевая аудитория: найм, фриланс, малый бизнес; цели — отпуск, машина, ипотека, подушка.

---

## 3. Решение (1 мин)

| Блок | Что даёт пользователю |
|------|------------------------|
| Онбординг ~1 мин | Доход, подушка, цель — можно пропустить блоки |
| Dashboard | Narrative, финансовый план (3 шага), метрики, «Что если» |
| Советник `/advisor` | Персональный чат со streaming |
| Кредиты | PDF договора → DTI и сравнение со ставкой рынка |
| Расходы | Голос или ручной ввод за секунды |

---

## 4. Demo flow (2–3 мин)

Рекомендуемый порядок — см. [mvp/README.md](../mvp/README.md):

1. Лендинг → вход.
2. Онбординг (skip ok) → dashboard с планом.
3. «Добавить расход» голосом.
4. `/advisor` — вопрос «где урезать».
5. Загрузка PDF кредита → светофор на `/credits`.
6. *(Опционально)* ФНС в профиле — mock, для wow-эффекта импорта чеков.

---

## 5. Архитектура (1–2 мин)

```
[Nuxt :3000] → [api-gateway :8000] → микросервисы Go
                      │
        user-service · receipt-service · credit-service
        ai-processor · analytics-service · scraper-service
                      │
              PostgreSQL + ClickHouse + Kafka
```

- **Stateless API** — масштабирование репликами сервисов.
- **OLTP / OLAP split** — транзакции в PG, аналитика dashboard в ClickHouse.
- **Advisor** — snapshot собирается на сервере; клиент шлёт только текст сообщения.
- **LLM** — Google Gemini; в dev опционально Antigravity → Claude через `:8045`.

Живые сервисы для demo: gateway, user, receipt, credit, ai-processor, analytics (+ infra).

---

## 6. Что готово vs mock (30 сек)

| Готово | Mock / demo-only |
|--------|------------------|
| Auth, profile, skip-flags | ФНС — UI mock на front |
| Голос/ручной, dashboard API | DEMO_MODE seed для быстрого старта |
| PDF credits, rates benchmark | File-store профиля (PG — следующий шаг) |
| `/ai/plan`, `/ai/chat`, SSE | Heuristic fallback если LLM недоступен |

---

## 7. Риски и mitigations (1 мин)

| Риск | Mitigation |
|------|------------|
| LLM недоступен | Regex/heuristic fallback, badge `source` в UI |
| Неполный профиль (skip) | Skip-aware snapshot — не считаем доход = 0 |
| Расхождение docs ↔ код | Ветка `docs` — источник правды; NAVI + API Contract |
| Масштаб >10k MAU | PG/CH на отдельных машинах, реплики stateless сервисов |

---

## 8. Post-MVP (30 сек)

- PG-backed snapshot без file-store.
- Реальная интеграция ФНС (не mock).
- Push-уведомления, банковская агрегация (Open Banking).
- Social/challenges — только как гипотеза, не в core.

---

## Q&A — частые вопросы

**«Почему микросервисы на хакатоне?»** — разделение ingest, analytics и advisor; gateway уже маршрутизирует; demo не требует все 15 сервисов одновременно.

**«Откуда данные для советника?»** — `BuildSnapshot`: profile + credits (PDF) + expenses; см. [advisor-system.md](../architecture/advisor-system.md).

**«Social в продукте?»** — нет, out of scope MVP.
