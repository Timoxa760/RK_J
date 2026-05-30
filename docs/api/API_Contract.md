# API Contract — Поток

> **Единый источник правды (ветка `docs`).** Любые изменения API — сначала здесь, затем код в `back` / `front`.  
> **Machine-readable:** [openapi.yaml](../contracts/openapi.yaml) (OpenAPI 3.1)  
> **Карта проекта:** [NAVI.md](../../NAVI.md)

| | |
|--|--|
| **Base URL** | `http://localhost:8000/api/v1` |
| **Gateway** | `api-gateway:8000` → reverse proxy |
| **Auth** | `Authorization: Bearer <jwt>` (кроме `/auth/register`, `/auth/login`) |
| **Форматы** | JSON; даты ISO 8601; суммы в **рублях** (`float64`) |

### Маршрутизация gateway → сервис

| Префикс | Сервис | Порт |
|---------|--------|------|
| `/auth/`, `/providers/` | user-service | 8001 |
| `/dashboard/`, `/receipts/` | receipt-service | 8002 |
| `/fns/`, `/x5club/`, `/magnit/`, `/email/` | scraper-service | 8003 |
| `/expenses/` | ai-processor | 8100 |
| `/credits/` | credit-service | 8009 |
| `/banks/` | bank-service | 8011 |
| `/goals/`, `/budgets/`, `/categories/` | finance-core | 8006–8005–8004 |
| `/insights/`, `/forecast/`, `/scenarios/` | analytics-service | 8101 |
| `/digest/` | reporting-service | 8010 |
| `/challenges/` | social-service | 8102 |

---

## СУММАРНАЯ ТАБЛИЦА ЭНДПОИНТОВ

| Приоритет | Метод | Путь | Описание | Группа |
|-----------|-------|------|----------|--------|
| 🔴 critical | POST | `/auth/register` | Регистрация по телефону | Core API |
| 🔴 critical | POST | `/auth/login` | Вход, получение JWT | Core API |
| 🔴 critical | POST | `/providers/connect` | Привязать магазин | Core API |
| 🔴 critical | POST | `/expenses/manual` | Текстовый ввод трат | AI Processor |
| 🔴 critical | POST | `/expenses/voice` | Голосовой ввод (audio → Whisper → LLM) | AI Processor |
| 🔴 critical | POST | `/fns/ticket` | Чек по QR ФНС | Scraper |
| 🟡 important | POST | `/fns/mco/sync` | Синк истории MCO | Scraper |
| 🟡 important | POST | `/goals` | Создать цель | Finance Core |
| 🔴 critical | GET | `/dashboard/sankey` | Санки-диаграмма | Receipt Engine |
| 🔴 critical | GET | `/dashboard/categories` | Круговая с детализацией | Receipt Engine |
| 🟡 important | GET | `/dashboard/timemachine` | Накопления за 60 мес | Receipt Engine |
| 🟡 important | GET | `/dashboard/stores` | Пузырьковая по магазинам | Receipt Engine |
| 🟡 important | GET | `/dashboard/compare` | Сравнение месяцев | Receipt Engine |
| 🟡 important | GET | `/credits/dashboard` | DTI, подушка, stress-test | Finance Core |
| 🟡 important | POST | `/credits/scan` | AI-скан договора | Finance Core |
| 🟡 important | GET | `/insights` | Инсайты (подписки, дубли) | Money Intelligence |
| 🟡 important | POST | `/scenarios/simulate` | Time Machine симуляция | Money Intelligence |
| 🟢 optional | GET | `/forecast` | Прогноз трат на 7 дней | Money Intelligence |
| 🟢 optional | POST | `/challenges` | Создать челлендж | Social & Game |
| 🟢 optional | GET | `/challenges/{id}/leaderboard` | Лидерборд челленджа | Social & Game |
| 🟢 optional | GET | `/digest/latest` | Ежемесячный дайджест | Reporting |
| 🟢 optional | POST | `/providers/{name}/sync` | Форсировать синхронизацию | Receipt Engine |
| 🟢 optional | GET | `/banks/accounts` | Счета из банка | Finance Core |
| 🟢 optional | GET | `/banks/transactions` | Транзакции из банка | Finance Core |

---

## 1. Auth API (Core API)

### POST /api/v1/auth/register — 🔴 critical

Регистрация по номеру телефона. В демо-режиме код всегда `0000`.

**Body:**
```json
{"phone": "+79991234567"}
```
**200 OK:**
```json
{"message": "SMS sent", "expires_in": 300}
```
**400** — неверный формат телефона

Повторный запрос для уже зарегистрированного номера возвращает **200** (`SMS sent`) — повторная отправка кода, как в demo-режиме.

### POST /api/v1/auth/login — 🔴 critical

Подтверждение SMS-кода, получение JWT.

**Body:**
```json
{"phone": "+79991234567", "code": "0000"}
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
**401** — неверный код

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

## 3. Providers API (Core API)

### POST /api/v1/providers/connect — 🔴 critical

Привязать магазин/провайдера.

**Query:** `?provider=x5club`

**Body:**
```json
{
  "credentials": {"phone": "+79991234567", "password": "***"}
}
```
**200 OK:**
```json
{"message": "Provider connected", "provider": "x5club", "status": "active"}
```
**400** — неверные credentials | **409** — уже привязан

### POST /api/v1/providers/{name}/sync — 🟢 optional

Форсировать синхронизацию провайдера.

**Path:** `x5club | magnit | lenta | vkusvill | ozon | wb | email | fns`

**200 OK:**
```json
{"message": "Sync started", "provider": "x5club"}
```
**202** — уже синхронизируется | **404** — провайдер не привязан

---

## 4. Credits API (Finance Core)

### GET /api/v1/credits/dashboard — 🟡 important

Кредитный дашборд: DTI, подушка, список кредитов.

**200 OK:**
```json
{
  "dti": 0.28,
  "stress_test_months": 4.2,
  "savings": 340000,
  "total_debt": 1200000,
  "monthly_payments": 42000,
  "monthly_income": 180000,
  "credits": [
    {
      "id": "uuid",
      "bank": "Т-Банк",
      "amount": 1200000,
      "rate": 14.5,
      "term_months": 36,
      "remaining": 980000,
      "monthly_payment": 42000,
      "next_payment": "2026-06-15"
    }
  ]
}
```
**401**

### POST /api/v1/credits/scan — 🟡 important

Загрузить PDF договора, OnlySQ извлекает условия.

**Body:** `multipart/form-data` — поле `file` (PDF)

**200 OK:**
```json
{
  "parsed": {
    "amount": 1200000, "rate": 14.5, "term_months": 36,
    "monthly_payment": 42000, "bank": "Т-Банк"
  },
  "confidence": 0.87
}
```
**400** — не PDF | **422** — не распознано

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

Пайплайн: **Whisper** (только `/voice`) → **OnlySQ** (парсинг + совет) → fallback **regex** (`parser.Parse`) → сохранение.

Auth: `Authorization: Bearer <jwt>` (через gateway).

### POST /api/v1/expenses/manual — 🔴 critical

Текстовый ввод. Парсинг: OnlySQ при наличии `ONLYSQ_API_KEY`, иначе regex.

**Body:**
```json
{
  "user_id": "uuid-or-phone-id",
  "raw_text": "купил продукты на 5000 и кроссовки за 16000",
  "amount": 0,
  "category": "",
  "description": "",
  "date": "2026-05-30",
  "source": "manual"
}
```

| Поле | Обязательно | Описание |
|------|-------------|----------|
| `user_id` | да | ID пользователя |
| `raw_text` | нет* | Текст для OnlySQ / regex |
| `amount` | нет* | Если задан — подставляется в первую трату |
| `category` | нет | Переопределяет категорию первой траты |
| `description` | нет | Описание первой траты |
| `date` | нет | ISO date `YYYY-MM-DD` |
| `source` | нет | `manual` (default) \| `voice` (если текст уже с фронта) |

\* Нужен `raw_text` **или** `amount` > 0.

**200 OK:**
```json
{
  "success": true,
  "id": "uuid",
  "amount": 5000,
  "category": "Продукты",
  "parsed": true,
  "parsed_by": "onlysq",
  "advice": "Продуктовая трата в рамках недели; кроссовки — разовая покупка вне обычного ритма.",
  "expenses": [
    {"id": "uuid-1", "amount": 5000, "category": "Продукты", "description": "продукты"},
    {"id": "uuid-2", "amount": 16000, "category": "Одежда", "description": "кроссовки"}
  ]
}
```

| Поле ответа | Описание |
|-------------|----------|
| `id`, `amount`, `category` | Первая трата (обратная совместимость с Nuxt) |
| `parsed_by` | `onlysq` \| `regex` |
| `advice` | Одна рекомендация от LLM (пусто при regex-only) |
| `expenses` | Все траты из фразы (1..N) |

**400** — `user_id required`, `amount required`, invalid JSON  
**500** — save failed  
**503** — OnlySQ недоступен и regex не смог извлечь сумму

---

### POST /api/v1/expenses/voice — 🔴 critical

Голосовой ввод. Multipart: аудио → Whisper → OnlySQ → сохранение.

**Content-Type:** `multipart/form-data`

| Поле | Обязательно | Описание |
|------|-------------|----------|
| `file` | да | Аудио: `webm`, `wav`, `mp3`, `ogg`, `m4a` (max 10 MB) |
| `user_id` | да | ID пользователя |
| `date` | нет | ISO date `YYYY-MM-DD` |

**200 OK:** тот же формат, что `/expenses/manual`, плюс:

```json
{
  "transcript": "Только что вышел из продуктового, потратил 1332 рубля",
  "source": "voice"
}
```

**400** — `user_id required`, `file required`, файл слишком большой  
**503** — Whisper недоступен (`WHISPER_URL` не задан или сервис не отвечает)  
**500** — save failed

---

### Legacy: POST /api/v1/receipt/manual и POST /api/v1/receipt/voice

Адаптер для текущего Nuxt-front (поля `store`, `receipt_id`, `audio`).  
`user_id` берётся из JWT (`sub`), тело front не меняется.

| Front | Адаптер → expenses |
|-------|-------------------|
| `POST /receipt/manual` `{store, amount, category, date}` | → `/expenses/manual` |
| `POST /receipt/voice` multipart `audio` | → Whisper + `/expenses/manual` pipeline |

**200 OK /receipt/manual:** `{receipt_id, store, amount, category, date, status}`  
**200 OK /receipt/voice:** `{receipt_id, store, items[], total, category, confidence}`

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

### POST /api/v1/fns/mco/sync — 🟡 important

Загрузка истории чеков после авторизации MCO.

### POST /api/v1/email/receipts — 🟢 optional

Pull чеков из IMAP (OAuth Яндекс / Mail.ru): `GET /api/v1/auth/oauth/{provider}`.

### POST /api/v1/x5club/send-code | `/x5club/sync`

OTP и синхронизация X5 Club.

### POST /api/v1/magnit/send-code | `/magnit/sync`

OTP и синхронизация Магнит.

---

## 11. Goals API (Finance Core)

### POST /api/v1/goals — 🟡 important

Создание цели (онбординг, сценарий №2).

**Body (ориентир):**
```json
{
  "title": "Отпуск",
  "target_amount": 150000,
  "target_date": "2026-12-01",
  "auto_save_percent": 10
}
```

**200 OK:**
```json
{
  "id": "uuid",
  "title": "Отпуск",
  "target_amount": 150000,
  "current_amount": 0,
  "progress_percent": 0
}
```

Проксируется через gateway → `goal-service:8006`. CRUD: `/goals/{id}`.

---

## 12. Onboarding Profile — ⏳ roadmap

Эндпоинты для wizard `/onboarding` (ещё не в `back`):

| Метод | Путь | Описание |
|-------|------|----------|
| PATCH | `/users/me/profile` | active_income, passive_income, emergency_fund |
| POST | `/users/me/onboarding/complete` | `onboarding_completed: true` |

До реализации — данные можно собирать локально на front + `POST /goals` + `POST /expenses/manual`.

---

## 13. Согласование с `front`

Типы в `frontend/types/api.ts` — эталон для dashboard. Отличия от примеров ниже:

| Endpoint | Контракт (legacy) | Front (`types/api.ts`) |
|----------|-------------------|------------------------|
| `/dashboard/timemachine` | `months`, `real_savings[]` | `points[]` с `actual`, `optimistic`, `delta` |
| `/dashboard/compare` | `label`, nested categories | `month`, `categories[]` с `share` |
| `/credits/dashboard` | `dti: 0.28` (доля) | `dti: 38` (**проценты** 0–100) |

При интеграции без demo-mode — **привести back к front-типам** или обновить composables.

---

## 14. Типы данных (TypeScript для Nuxt)

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
