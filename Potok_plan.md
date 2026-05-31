# Поток — Global Design Document

**Финансовый навигатор**

*Хакатон Финтех — 29–31 мая 2026 — Кейсодатель: Клерк.Ру*

> Оглавление GDD. Полная карта «где что искать» — **[NAVI.md](./NAVI.md)**.

---

## Содержание

**Часть 1. ПРОДУКТ**

1. [Обзор](./docs/product/overview.md) — что это, аудитория, инсайт, решение
2. [UX-сценарии](./docs/product/ux-scenarios.md) — 5 основных сценариев
3. [Онбординг](./docs/product/onboarding.md) — `/onboarding`, ~1 минута
4. [Способы ввода](./docs/product/input-methods.md) — голос, чек, ФНС
5. [Финансовая модель](./docs/product/financial-model.md) — единая модель и цикл
6. [Финансовое здоровье](./docs/product/financial-health.md) — метрики и подача
7. [Философия UX](./docs/product/philosophy.md) — тон, формулировки
8. [Монетизация](./docs/product/monetization.md) — ипотечный разбор
9. [Кейс Клерк.Ру](./docs/product/case-alignment.md) — соответствие ТЗ
10. [Roadmap](./docs/product/roadmap.md) — post-MVP, трудозатраты

**Часть 2. РЕАЛИЗАЦИЯ**

11. [Архитектура](./docs/architecture/overview.md)
12. [Kafka events](./docs/architecture/kafka-events.md)
13. [Защита архитектуры](./docs/architecture/defense.md) — pitch backend
14. [Фичи → сервисы](./docs/features/)
15. [API Contract](./docs/api/API_Contract.md) + [TypeScript](./docs/api/typescript-types.md) + [OpenAPI 3.1](./docs/contracts/openapi.yaml)
16. [Модель данных](./docs/database/)
17. [Фазы разработки](./docs/phases/phases.md)
18. [Deployment](./docs/deployment/) — back/front, env, demo
19. [MVP статус](./docs/mvp/README.md) + [Pitch materials](./docs/pitch/README.md)

**Гипотезы (низкий приоритет)**

- [Social](./docs/features/social.md) — челленджи
- [Аукцион отказов](./docs/features/auction.md)

---

## Elevator Pitch

**Поток** анализирует финансовое положение и превращает сложные данные в простые рекомендации.

Пользователь получает: финансовое здоровье, оценку устойчивости, прогноз достижения целей, анализ привычек и **одно** персональное действие.

Главный вопрос продукта:

> **Что мне сделать сегодня, чтобы быстрее получить то, что я хочу?**

---

## Инсайт

Людям не нужны графики. Людям нужны **ответы**.

| Не так | А так |
|--------|-------|
| Вы потратили 12 000 ₽ на маркетплейсы | Из-за текущих привычек вы накопите на отпуск через 16 месяцев |
| Таблица категорий | Если сократить импульсивные покупки на 15%, цель ближе на 7 месяцев |
| DTI 0.42 | Маркетплейсы отодвигают цель на 3 месяца |

---

## Связь docs ↔ код

| Документ | Ветка | Примечание |
|----------|-------|------------|
| Архитектура, API, БД | `back` | Порты — `back/docker-compose.yml` |
| Экраны, маршруты | `front` | Nuxt 4, см. [NAVI.md](./NAVI.md) |
| Продукт и UX | `docs` | Источник правды по смыслу |

**Не ломать:** имена сервисов, порты, топики Kafka, структура каталогов `back`/`front`.
