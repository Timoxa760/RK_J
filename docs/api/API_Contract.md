# API Contract — Поток

> **Единый источник правды (ветка `front`, папка `docs/`).** Scope: [scope.md](../product/scope.md).  
> **Machine-readable:** [openapi.yaml](../contracts/openapi.yaml) (OpenAPI 3.1)  
> **Карта проекта:** [NAVI.md](../../NAVI.md)

| | |
|--|--|
| **Base URL** | `http://localhost:8000/api/v1` |
| **Gateway** | `api-gateway:8000` → reverse proxy |
| **Auth** | `Authorization: Bearer <jwt>` (кроме `/auth/register`, `/auth/login`) |
| **Форматы** | JSON; даты ISO 8601; суммы в **рублях** (`float64`) |

## Маршрутизация gateway → сервис (обновлённая)

| Префикс | Сервис | Порт |
|---------|--------|------|
| `/auth/` | user-service | 8001 |
| `/dashboard/`, `/receipts/` | receipt-service | 8002 |
| `/receipt/manual`, `/receipt/voice`, `/receipt/fns/` | scraper-service | 8003 |
| `/ai/` | ai-processor | 8100 |
| `/credits/` | credit-service | 8009 |
| `/users/me/` | user-service | 8001 |
| `/banks/` | bank-service | 8011 |
| `/budgets/`, `/categories/` | finance-core | 8005–8004 |
| `/insights/`, `/forecast/`, `/scenarios/` | analytics-service | 8101 |
| `/digest/` | reporting-service | 8010 |

---

## СУММАРНАЯ ТАБЛИЦА ЭНДПОИНТОВ

| Приоритет | Метод | Путь | Описание | Группа |
|-----------|-------|------|----------|--------|
| 🔴 critical | POST | `/auth/register` | Регистрация (телефон + пароль) | Core API |
| 🔴 critical | POST | `/auth/login` | Вход, получение JWT | Core API |
| 🟡 important | POST | `/auth/password/forgot` | Запрос кода сброса пароля | Core API |
| 🟡 important | POST | `/auth/password/reset` | Сброс пароля по коду | Core API |
| 🔴 critical | POST | `/receipt/manual` | Ручной ввод расхода | Receipt Engine |
| 🔴 critical | POST | `/receipt/fns/scan` | Чек по QR ФНС | Receipt Engine |
| 🟢 optional | POST | `/receipt/voice` | Голосовой ввод расхода | Receipt Engine |
| 🔴 critical | GET | `/dashboard/sankey` | Санки-диаграмма | Receipt Engine |
| 🔴 critical | GET | `/dashboard/categories` | Круговая с детализацией | Receipt Engine |
| 🟡 important | GET | `/dashboard/timemachine` | Накопления за 60 мес | Receipt Engine |
| 🟡 important | GET | `/dashboard/stores` | Пузырьковая по магазинам | Receipt Engine |
| 🟡 important | GET | `/dashboard/compare` | Сравнение месяцев | Receipt Engine |
| 🔴 critical | GET | `/ai/plan` | Финансовый план + diagnosis | ai-processor |
| 🔴 critical | GET | `/ai/diagnosis` | Финансовый диагноз | ai-processor |
| 🔴 critical | POST | `/ai/chat` | Чат с AI-ассистентом | ai-processor |
| 🟡 important | GET | `/users/me/profile` | Финансовый профиль | user-service |
| 🟡 important | PATCH | `/users/me/profile` | Обновить профиль | user-service |
| 🟡 important | POST | `/users/me/onboarding/complete` | Завершить онбординг | user-service |
| 🟡 important | POST | `/onboarding/parse` | Парсинг ответа опроса | ai-processor |
| 🟡 important | GET | `/credits/dashboard` | DTI из PDF-сканов | credit-service |
| 🟡 important | POST | `/credits/scan` | AI-скан PDF договора | credit-service |
| ~~removed~~ | ~~POST~~ | ~~`/goals`~~ | ~~CRUD целей~~ | ~~goal-service~~ |
| ~~removed~~ | ~~GET~~ | ~~`/ai/goal/{id}/forecast`~~ | — | — |
| ~~removed~~ | ~~*~~ | ~~`/challenges/*`~~ | — | — |
| 🟡 important | GET | `/insights` | Инсайты (подписки, дубли, переплаты) | Money Intelligence |
| 🟡 important | POST | `/scenarios/simulate` | Time Machine симуляция | Money Intelligence |
| 🟢 optional | GET | `/forecast` | Прогноз трат на 7 дней | Money Intelligence |
| 🟢 optional | GET | `/digest/latest` | Ежемесячный дайджест | Reporting |

---

## 1. Auth API (Core API)

### POST /api/v1/auth/register — 🔴 critical

Регистрация по номеру телефона и паролю (мин. 8 символов).

**Body:**
```json
{"phone": "+79991234567", "password": "secret12345"}
```
**200 OK:**
```json
{"message": "registered"}
```
**400** — неверный формат телефона или короткий пароль | **409** — пользователь уже существует

### POST /api/v1/auth/login — 🔴 critical

Вход по телефону и паролю, получение JWT.

**Body:**
```json
{"phone": "+79991234567", "password": "secret12345"}
```
**200 OK:**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "dGhpcy...",
  "expires_in": 900,
  "user": {"id": "uuid", "phone": "+79991234567", "role": "user"}
}
```
**401** — неверный телефон или пароль

### POST /api/v1/auth/password/forgot — 🟡 important

Запрос кода для сброса пароля. Ответ одинаковый, существует аккаунт или нет.

**Body:**
```json
{"phone": "+79991234567"}
```
**200 OK:**
```json
{"message": "If the account exists, a reset code has been sent", "expires_in": 300}
```

### POST /api/v1/auth/password/reset — 🟡 important

Установка нового пароля по коду из SMS.

**Body:**
```json
{"phone": "+79991234567", "code": "482913", "new_password": "newsecret1"}
```
**200 OK:**
```json
{"message": "password updated"}
```
**401** — неверный или просроченный код | **400** — короткий пароль

---

## 2. Dashboard API (Receipt Engine)

### GET /api/v1/dashboard/sankey — 🔴 critical

Санки-диаграмма потока денег. **Только уровень Доход → Категории** (двухуровневая). Детализация до товаров — в `/dashboard/categories` (круговая диаграмма).

Доходы берутся из mock-данных (демо-режим) или из `bank_transactions`. Расходы — из ClickHouse (MV `spending_by_category`).

**200 OK:**
```json
{
  "nodes": [
    {"name": "Зарплата", "value": 180000},
    {"name": "Накопления", "value": 35000},
    {"name": "Продукты", "value": 52000},
    {"name": "Кафе и рестораны", "value": 28000},
    {"name": "Транспорт", "value": 15000},
    {"name": "Развлечения", "value": 12000}
  ],
  "links": [
    {"source": "Зарплата", "target": "Накопления", "value": 35000},
    {"source": "Зарплата", "target": "Продукты", "value": 52000}
  ]
}
```
**401**

### GET /api/v1/dashboard/timemachine — 🟡 important

Накопления за 60 месяцев: реальные vs оптимистичные.

**200 OK:**
```json
{
  "months": ["2026-05", "2026-06", "2026-07"],
  "real_savings": [500000, 512000, 524500],
  "optimized_savings": [500000, 516000, 532500],
  "difference_final": 467000
}
```
**401**

### GET /api/v1/dashboard/stores — 🟡 important

Пузырьковая диаграмма по магазинам.

**200 OK:**
```json
{
  "stores": [
    {"name": "Пятёрочка", "avg_check": 650, "purchases": 14, "total": 9100, "impulse_ratio": 0.25},
    {"name": "Магнит", "avg_check": 720, "purchases": 10, "total": 7200, "impulse_ratio": 0.20},
    {"name": "ВкусВилл", "avg_check": 980, "purchases": 7, "total": 6860, "impulse_ratio": 0.10},
    {"name": "Ozon", "avg_check": 2100, "purchases": 4, "total": 8400, "impulse_ratio": 0.65},
    {"name": "Wildberries", "avg_check": 1850, "purchases": 5, "total": 9250, "impulse_ratio": 0.70},
    {"name": "Лента", "avg_check": 820, "purchases": 8, "total": 6560, "impulse_ratio": 0.15}
  ]
}
```
**401**

### GET /api/v1/dashboard/compare — 🟡 important

Сравнение расходов за N последних месяцев.

**Query:** `?months=2` (default 2)

**200 OK:**
```json
{
  "months": [
    {
      "label": "Февраль 2026",
      "categories": [
        {"name": "Продукты", "total": 48000},
        {"name": "Кафе и рестораны", "total": 25000}
      ]
    },
    {
      "label": "Март 2026",
      "categories": [
        {"name": "Продукты", "total": 52000},
        {"name": "Кафе и рестораны", "total": 28000}
      ]
    }
  ],
  "insights": {
    "biggest_change": {"category": "Кафе и рестораны", "delta": 3000, "delta_percent": 12}
  }
}
```
**401**

### GET /api/v1/dashboard/categories — 🔴 critical

Расходы по категориям с детализацией по товарам (для круговой диаграммы с кликом).

**200 OK:**
```json
{
  "categories": [
    {
      "name": "Продукты",
      "total": 52000,
      "subcategories": [
        {
          "name": "Молочные",
          "total": 8500,
          "items": [
            {"name": "Молоко 3.2%", "price": 78, "quantity": 12, "total": 936},
            {"name": "Творог 5%", "price": 120, "quantity": 6, "total": 720}
          ]
        }
      ]
    }
  ]
}
```
**401**

---

## 3. Receipt API (Receipt Engine)

### POST /api/v1/receipt/manual — 🔴 critical

Добавить расход вручную.

**Body:**
```json
{
  "store": "Пятёрочка",
  "amount": 1032.50,
  "category": "Продукты",
  "date": "2026-05-30T14:32:00Z"
}
```
200 OK:
```json
{
  "receipt_id": "uuid",
  "store": "Пятёрочка",
  "amount": 1032.50,
  "category": "Продукты",
  "date": "2026-05-30T14:32:00Z",
  "status": "saved"
}
```
400 — неверный формат суммы или категории | 401

### POST /api/v1/receipt/voice — 🟢 optional
Добавить расход голосом. Аудио распознаётся через Whisper API, AI извлекает магазин, товары и сумму.

Body: multipart/form-data — поле audio (mp3/wav)

200 OK:
```json
{
  "receipt_id": "uuid",
  "store": "Пятёрочка",
  "items": [
    {"name": "Молоко", "price": 89.90, "quantity": 1},
    {"name": "Хлеб", "price": 45.50, "quantity": 1}
  ],
  "total": 135.40,
  "category": "Продукты",
  "confidence": 0.92
}
```
400 — аудио не распознано | 401

### POST /api/v1/receipt/fns/scan — 🔴 critical
Отсканировать QR-код чека ФНС.

Body:
```json
{
  "fn": "9285000100351475",
  "fd": "1234567890",
  "fp": "1234567890"
}
```
200 OK:
```json
{
  "receipt_id": "uuid",
  "store": "Пятёрочка",
  "inn": "7725007364",
  "date": "2026-05-30T14:32:00Z",
  "total": 1032.50,
  "items": [
    {"name": "Молоко", "price": 89.90, "quantity": 1},
    {"name": "Хлеб", "price": 45.50, "quantity": 2}
  ],
  "category": "Продукты"
}
```
400 — неверные fn/fd/fp | 401 | 404 — чек не найден в ФНС

---

## 4. Credits API (Finance Core)

> **Единственный источник кредитных данных** — PDF-скан. См. [credit-scan.md](../features/credit-scan.md).

### GET /api/v1/credits/dashboard — 🟡 important

Агрегат по таблице `user_credits` пользователя. **Пустой**, если сканов не было.

**200 OK (есть сканы):**
```json
{
  "dti": 28,
  "stress_test_months": 4.2,
  "savings": 0,
  "total_debt": 980000,
  "monthly_payments": 42000,
  "monthly_income": 0,
  "credits": [
    {
      "id": "uuid",
      "bank": "Т-Банк",
      "amount": 1200000,
      "rate": 14.5,
      "benchmark_rate": 12.1,
      "rate_vs_market": "above",
      "term_months": 36,
      "remaining": 980000,
      "monthly_payment": 42000,
      "next_payment": "2026-06-15"
    }
  ]
}
```

`dti` — **проcentы** 0–100. `monthly_income` для DTI берётся из profile при наличии.

**401**

### POST /api/v1/credits/scan — 🟡 important

PDF договора → OnlySQ → сохранение в PG → rates-aggregator.

**Body:** `multipart/form-data` — поле `file` (PDF)

**200 OK:**
```json
{
  "parsed": {
    "amount": 1200000, "rate": 14.5, "term_months": 36,
    "monthly_payment": 42000, "bank": "Т-Банк"
  },
  "benchmark_rate": 12.1,
  "rate_vs_market": "above",
  "confidence": 0.87,
  "credit_id": "uuid"
}
```
**400** — не PDF | **422** — не распознано

### DELETE /api/v1/credits/{id} — 🟢 optional

Удалить сохранённый скан.

---

## 5. Insights API (Money Intelligence)

### GET /api/v1/insights — 🟡 important

Инсайты: подписки, дубли, переплаты.

**200 OK:**
```json
{
  "insights": [
    {
      "type": "subscription",
      "severity": "warning",
      "title": "Найдена скрытая подписка",
      "description": "Списывается 299 ₽ каждый месяц",
      "amount": 299,
      "merchant": "Яндекс.Плюс"
    },
    {
      "type": "duplicate",
      "severity": "info",
      "title": "Дублирование в чеке",
      "description": "Товар 'Молоко 3.2%' пробит дважды",
      "amount": 156
    },
    {
      "type": "overprice",
      "severity": "warning",
      "title": "Переплата за товар",
      "description": "Молоко 3.2% куплено за 95 ₽, средняя — 78 ₽",
      "amount": 17,
      "store": "Пятёрочка"
    }
  ]
}
```
**401**

### GET /api/v1/ai/diagnosis — 🔴 critical

Финансовая картина пользователя: оценка, показатели и главное действие.

**200 OK:**
```json
{
  "score": 72,
  "grade": "B",
  "indicators": [
    {"name": "Долговая нагрузка", "value": 28, "norm": "<30", "status": "good"},
    {"name": "Подушка безопасности", "value": 4.2, "norm": ">3", "status": "good"},
    {"name": "Накопления от дохода", "value": 15, "norm": ">20", "status": "warning"},
    {"name": "Импульсивные траты", "value": 32, "norm": "<25", "status": "critical"},
    {"name": "Стабильность доходов", "value": 85, "norm": ">70", "status": "good"}
  ],
  "main_action": {
    "title": "Сократите доставку еды",
    "description": "Вы тратите 9 000 ₽ в месяц на доставку. Готовьте дома 3 раза в неделю — это сэкономит 4 500 ₽ в месяц.",
    "potential_savings": 4500,
    "difficulty": "easy"
  },
  "next_check_days": 30
}
```

| Поле `status` | Значение |
|---------------|----------|
| `good` | в норме |
| `warning` | стоит улучшить |
| `critical` | требует внимания |

**401**

### POST /api/v1/scenarios/simulate — 🟡 important

Симулятор «Машина времени». Что если сократить категорию на N%.

**Body:**
```json
{
  "scenario": "reduce_delivery",
  "reduction_percent": 50,
  "months": 60
}
```
**Допустимые `scenario`:** `reduce_delivery`, `reduce_cafe`, `reduce_entertainment`, `custom`

**200 OK:**
```json
{
  "months": ["2026-05", "2026-06", "2026-07"],
  "real_savings": [500000, 512000, 524500],
  "optimized_savings": [500000, 525000, 551250],
  "difference_final": 467000,
  "scenario": {
    "name": "reduce_delivery",
    "monthly_saving": 4500,
    "annual_saving": 54000
  }
}
```
**400** — неверный сценарий | **401**

### GET /api/v1/forecast — 🟢 optional

Прогноз трат на N дней.

**Query:** `?days=7`

**200 OK:**
```json
{
  "dates": ["2026-05-28", "2026-05-29", "2026-05-30"],
  "forecast": [5200, 5100, 5300],
  "upper_bound": [6240, 6120, 6360],
  "lower_bound": [4160, 4080, 4240]
}
```
**401**

---

## 6. Social API (Social & Game)

### POST /api/v1/challenges — 🟢 optional

Создать челлендж.

**Body:**
```json
{
  "type": "least_spend",
  "title": "Кто меньше потратит на доставку",
  "duration_days": 7,
  "max_participants": 10
}
```
**Types:** `least_spend | most_saved | streak`

**200 OK:**
```json
{
  "id": "challenge-uuid",
  "type": "least_spend",
  "status": "active",
  "invite_token": "invite-xxxx",
  "participants": 1,
  "created_at": "2026-05-27T12:00:00Z"
}
```
**400** | **401**

### GET /api/v1/challenges/{id}/leaderboard — 🟢 optional

Лидерборд челленджа (анонимизированный).

**200 OK:**
```json
{
  "challenge_id": "challenge-uuid",
  "type": "least_spend",
  "leaderboard": [
    {"position": 1, "username": "Анна", "avatar": null, "relative_score": 0.0},
    {"position": 2, "username": "Иван", "avatar": null, "relative_score": 0.35}
  ],
  "my_position": {"position": 2, "total_participants": 3}
}
```
**401** | **403** — не участник | **404**

---

## 7. Digest API (Reporting)

### GET /api/v1/digest/latest — 🟢 optional

Ежемесячный дайджест.

**200 OK:**
```json
{
  "period": {"from": "2026-04-01", "to": "2026-04-30"},
  "total_spent": 145000,
  "total_income": 180000,
  "saved": 35000,
  "by_category": [
    {"name": "Продукты", "total": 52000, "percent": 35.9, "trend": "+8.3%"}
  ],
  "word_cloud": ["молоко", "латте", "хлеб", "сыр", "такси"],
  "top_stores": [{"name": "Пятёрочка", "total": 9100, "visits": 14}],
  "mindfulness_rating": 72,
  "ai_advice": "Попробуйте сократить доставку — это 9 000 ₽ в месяц",
  "insights_summary": "Найдено 2 скрытые подписки и 3 переплаты"
}
```
**401**

---

## 8. Banks API (Finance Core)

### GET /api/v1/banks/accounts — 🟢 optional

Список счетов пользователя из подключённых банков.

**200 OK:**
```json
{
  "accounts": [
    {"id": "uuid", "bank": "Т-Банк", "name": "Дебетовая Tinkoff Black", "balance": 340000, "currency": "RUB"}
  ]
}
```
**401**

### GET /api/v1/banks/transactions — 🟢 optional

Транзакции за период.

**Query:** `?from=2026-04-01&to=2026-04-30`

**200 OK:**
```json
{
  "transactions": [
    {"id": "uuid", "date": "2026-04-15", "amount": 180000, "description": "Зарплата", "category": "income"},
    {"id": "uuid", "date": "2026-04-16", "amount": -1200, "description": "Пятёрочка", "category": "food"}
  ]
}
```
**401**

---

## 9. Expenses API (AI Processor) — 🔴 critical для «Поток»

### POST /api/v1/expenses/manual

Голосовой или ручной ввод. Текст в `raw_text` парсится на бэкенде (`parser.Parse`).

**Body:**
```json
{
  "user_id": "uuid-or-phone-id",
  "raw_text": "купил продукты на 5000 и кроссовки за 16000",
  "amount": 0,
  "category": "",
  "description": "",
  "date": "2026-05-30",
  "source": "voice"
}
```

| Поле | Обязательно | Описание |
|------|-------------|----------|
| `user_id` | да | ID пользователя |
| `raw_text` | нет* | Текст для LLM/парсера |
| `amount` | нет* | Если задан — используется напрямую |
| `source` | нет | `manual` (default) \| `voice` |

\* Нужен `raw_text` с распознанной суммой **или** `amount` > 0.

**200 OK:**
```json
{
  "success": true,
  "id": "uuid",
  "amount": 5000,
  "category": "Продукты",
  "parsed": true
}
```

**400** — `user_id required`, `amount required`, invalid JSON | **500** — save failed

> **Roadmap:** один запрос → несколько транзакций из одной фразы (сейчас — одна запись).

---

## 10. FNS & Ingest API (Scraper Service)

Опциональный автослой. ФНС не обязательна для MVP.

### POST /api/v1/fns/ticket — 🔴 critical

Проверка чека по QR / ticket data.

**Body:**
```json
{
  "qr": "t=20260530T1200&s=5000.00&fn=...&i=...&fp=...&n=1"
}
```

**200 OK:** нормализованный чек → Kafka `receipt.raw` → `receipt-service`.

**400** | **502** — ФНС недоступна

### POST /api/v1/fns/qr

Fallback по QR-строке (аналог ticket).

### POST /api/v1/fns/mco/auth

Начало OAuth-потока MCO (мобильный кабинет налогоплательщика).

### POST /api/v1/fns/mco/auth/verify

Подтверждение кода MCO.

---

## 11. Goals API — REMOVED

> Цель хранится в **profile** (`goal_kind`, `goal_title`, `goal_amount`, `skipped_goal`).  
> `goal-service` и `/goals/*` **не используются** в MVP.

---

## 12. Profile & Onboarding

### GET /api/v1/users/me/profile — 🟡 important

**200 OK:**
```json
{
  "active_income": 150000,
  "passive_income": 20000,
  "emergency_fund": 340000,
  "emergency_breakdown": {"cash": 50000, "deposit": 200000, "investments": 90000},
  "fixed_expenses": [{"title": "Аренда", "amount": 45000}],
  "goal_kind": "save",
  "goal_title": "Отпуск",
  "goal_amount": 150000,
  "skipped_income": false,
  "skipped_cushion": false,
  "skipped_goal": false,
  "skipped_expenses": false,
  "survey_input_mode": "voice",
  "onboarding_completed": true
}
```

### PATCH /api/v1/users/me/profile

Partial update. Те же поля, все optional.

### POST /api/v1/users/me/onboarding/complete

**Body:** `{"onboarding_completed": true}`

### POST /api/v1/onboarding/parse

**Body:** `{"step": "income|cushion|goal|expenses", "raw_text": "...", "locale": "ru"}`

**Response:** `{"parsed": true, "step": "...", "patch": {...}}`

---

## 13. Advisor API (ai-processor)

См. [advisor.md](../product/advisor.md).

### GET /api/v1/ai/plan — 🔴 critical

**200 OK:**
```json
{
  "plan": {
    "goalTitle": "Отпуск",
    "goalProgress": "Накоплено 12% — осталось 132 000 ₽",
    "steps": [
      {"title": "...", "description": "..."},
      {"title": "...", "description": "..."},
      {"title": "...", "description": "..."}
    ],
    "runwayText": "Запас примерно на 5 мес.",
    "freeCashflowText": "После расходов остаётся 85 000 ₽/мес.",
    "updatedAt": 1717000000000
  },
  "diagnosis": {
    "score": 72,
    "grade": "B",
    "indicators": [],
    "main_action": {
      "title": "Сократите доставку",
      "description": "...",
      "potential_savings": 4500,
      "difficulty": "easy"
    },
    "next_check_days": 30
  }
}
```

### GET /api/v1/ai/diagnosis

Только `diagnosis` объект (или тот же shape без plan).

### POST /api/v1/ai/chat — 🔴 critical

**Body:**
```json
{
  "message": "Составь план",
  "history": [{"role": "user", "content": "..."}, {"role": "assistant", "content": "..."}]
}
```

**200 OK:** `{"reply": "..."}`

Контекст (profile, credits, expenses) — server-side по JWT.

---

## 14. Challenges API — REMOVED

> `/social`, `/challenges/*`, `social-service` — out of scope. См. [social.md](../features/social.md).

---

## 15. Согласование с `front`

Типы в `frontend/types/api.ts` — эталон для dashboard. Отличия от примеров ниже:

| Endpoint | Контракт (legacy) | Front (`types/api.ts`) |
|----------|-------------------|------------------------|
| `/dashboard/timemachine` | `months`, `real_savings[]` | `points[]` с `actual`, `optimistic`, `delta` |
| `/dashboard/compare` | `label`, nested categories | `month`, `categories[]` с `share` |
| `/credits/dashboard` | `dti: 0.28` (доля) | `dti: 38` (**проценты** 0–100) |

При интеграции без demo-mode — **привести back к front-типам** или обновить composables.

---

## 16. Типы данных (TypeScript для Nuxt)

```typescript
// Все типы строго соответствуют JSON-ответам выше.
// Ниже — только перечисление ключевых типов для кодогенерации.

interface User { id: string; phone: string; email?: string; role: 'user' | 'admin'; created_at: string }

interface Receipt {
  id: string; user_id: string; provider: ProviderType; store_name: string;
  total: number; purchased_at: string; items: ReceiptItem[]; checksum: string
}

interface ReceiptItem {
  name: string; price: number; quantity: number; category: string; is_impulsive: boolean
}

interface Credit {
  id: string; bank: string; amount: number; rate: number; term_months: number;
  remaining: number; monthly_payment: number; next_payment: string
}

interface Challenge {
  id: string; type: 'least_spend' | 'most_saved' | 'streak'; title: string;
  status: 'active' | 'completed' | 'cancelled'; duration_days: number;
  max_participants: number; participants_count: number; invite_token: string;
  created_by: string; created_at: string
}

type ProviderType = 'x5club' | 'magnit' | 'lenta' | 'vkusvill' | 'ozon' | 'wb' | 'fns' | 'email'

type InsightType = 'subscription' | 'duplicate' | 'overprice'
type Severity = 'info' | 'warning' | 'critical'
```

> **Важно:** Все денежные значения — `Float64` (копейки). Даты — ISO 8601 (`string`). Enum-поля валидируются на бэкенде.
