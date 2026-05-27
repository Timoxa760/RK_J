# Фазы разработки (48 часов)

| # | Фаза | Оценка | Критичность |
|---|------|--------|-------------|
| **0** | Инфраструктура и DevOps | 2–3ч | Блокер |
| **1** | X5 Club provider | 2ч | Блокер |
| **2** | Магнит provider | 2ч | Блокер |
| **3** | Scraper ядро (Scheduler + Kafka) | 2ч | Блокер |
| **4** | Receipt-service (Kafka consumer) | 2ч | Блокер |
| **5** | Auth + привязка провайдеров | 2ч | Блокер |
| **6** | AI-категоризация (YandexGPT) | 2–3ч | Ключевая |
| **7** | Email + ФНС fallback | 2ч | Важная |
| **8** | Dashboard + аналитика | 3ч | Важная |
| **9** | Лента, ВкусВилл, Ozon, WB | 3–4ч | Дополнительно |
| **10** | Credit Health Dashboard | 2ч | WOW-фича |
| **11** | Time Machine + Predictive AI | 2–3ч | WOW-фича |
| **12** | Геймификация + Social | 2–3ч | WOW-фича |
| **13** | Детектор + Дайджест | 2ч | Для демо |
| **14** | Demo Polish + хардкод | 2–3ч | Финал |

**Итого:** ~30–36 часов.
**Критический путь:** 0 → 1 → 3 → 4 → 5 → 8 → 14 ≈ 17 часов.

## Критический путь (17 часов)

### Фаза 0: Инфраструктура и DevOps
Docker Compose: PostgreSQL 16, ClickHouse 24, Redis 7, Kafka + ZK, MinIO.

### Фаза 1: X5 Club provider
Auth POST /api/v2/auth/login, history GET /api/v2/history.

### Фаза 3: Scraper-service ядро
Provider interface, Scheduler (горутины), Kafka producer, AES-256-GCM.

### Фаза 4: Receipt-service
Kafka consumer → валидация → дедупликация → PostgreSQL.

### Фаза 5: Auth + привязка провайдеров
Регистрация по телефону, JWT (access + refresh), привязка LK.

### Фаза 8: Dashboard + аналитика
GET /api/v1/dashboard, /analytics/categories, /analytics/trends.

### Фаза 14: Demo Polish
Прекэширование. Фейковый fast-path. 6 актов < 30 секунд.
