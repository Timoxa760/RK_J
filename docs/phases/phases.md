# Фазы разработки (48 часов)

> Статус сверен с веткой **`back`** (коммиты `dafda67`, `12b7f12`) и **`front`** (`9ff425c`).

| # | Фаза | Оценка | Критичность | Статус |
|---|------|--------|-------------|--------|
| **0** | Инфраструктура и DevOps | 2–3ч | Блокер | ✅ `back` |
| **1** | Email + ФНС | 2ч | Блокер | ✅ scraper FNS/MCO |
| **2** | X5 Club provider | 2ч | Блокер | ✅ |
| **3** | Магнит provider | 2ч | Блокер | ✅ |
| **4** | Scraper ядро (Scheduler + Kafka) | 2ч | Блокер | ✅ |
| **5** | Receipt-service | 2ч | Блокер | ✅ + dashboard API |
| **6** | Auth + JWT + providers | 2ч | Блокер | ✅ |
| **7** | AI: голос/ручной + категоризация | 2–3ч | Ключевая | ✅ ai-processor |
| **8** | Dashboard + аналитика | 3ч | Важная | 🟡 back API + front charts |
| **9** | LK / доп. провайдеры | 3–4ч | Доп. | 🟡 X5/Magnit есть |
| **10** | Финансовое здоровье (UI) | 2ч | WOW | 🟡 `/credits`, product docs |
| **11** | Прогноз цели + сценарии | 2–3ч | WOW | 🟡 timemachine, forecast API |
| **12** | Social / gamification | 2–3ч | Демо | 🟡 страница `/social` |
| **13** | Инсайты + дайджест | 2ч | Демо | 🟡 `/digest`, insights |
| **14** | Demo Polish + онбординг UX | 2–3ч | Финал | ⏳ онбординг, narrative |

## Критический путь (продукт «Поток»)

```
Онбординг → первое действие (голос/чек/ФНС) → фин. здоровье → прогноз цели → ипотечный разбор (демо)
```

Технический путь (уже пройден):

```
0 → 1–6 (back) → 8 (dashboard) → 14 (polish front + seed)
```

## Фаза 14 — Demo Polish (актуальные задачи)

1. **`/onboarding`** на `front` — wizard по [onboarding.md](../product/onboarding.md)
2. Кнопка «Добавить» → голос / чек / ФНС
3. Narrative на dashboard: ответы вместо голых графиков
4. Seed + [demo_flow.sh](../deployment/scripts/demo_flow.sh) (эталон в `docs`; копия в `back/scripts/`)
5. Ипотечный сценарий ([monetization.md](../product/monetization.md))

**Не в demo:** social, auction — [гипотезы](../features/social.md).

## Связи

- **Продукт**: [../product/ux-scenarios.md](../product/ux-scenarios.md)
- **Архитектура**: [../architecture/overview.md](../architecture/overview.md)
