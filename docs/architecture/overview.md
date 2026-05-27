# Системная архитектура

## Технологический стек

| Слой | Технология |
|------|------------|
| **Backend Core** | Go 1.23 (chi, pgx, redis/go-redis, kafka-go, colly, goquery) |
| **AI / Data** | Python 3.12 (FastAPI, Prophet, scikit-learn, beautifulsoup4) |
| **Frontend** | React 19 + Vite + Tailwind CSS 4 + Recharts + PWA |
| **AI: NLP** | YandexGPT API, GigaChat, RuBERT |
| **AI: Forecasting** | Prophet (Facebook) |
| **Scraping** | Colly, goquery, go-imap, mitmproxy |
| **Infra** | Docker Compose, Kafka, 6 БД |
| **Auth** | JWT + OAuth2 (Яндекс, Mail.ru) |

## Базы данных

| БД | Назначение |
|----|------------|
| **PostgreSQL 16** | OLTP-ядро: пользователи, чеки, товары (JSONB), бюджеты, цели, кредиты |
| **ClickHouse 24** | OLAP-аналитика: агрегированные траты, дайджесты, прогнозы |
| **Redis 7** | Кэш (ФНС, ЦБ) + очереди + leaderboards |
| **Kafka** | Event Bus: receipt.raw, receipt.parsed, receipt.enriched, insight.found |
| **MinIO** | Object Storage: изображения чеков, экспорт CSV/XLSX/PDF |

## Микросервисы

| Сервис | Порт | Язык | Назначение |
|--------|------|------|------------|
| api-gateway | 8000 | Go | Единый вход, JWT, rate limiting |
| user-service | 8001 | Go | Регистрация, профиль |
| receipt-service | 8002 | Go | Приём чеков, валидация, дедупликация |
| scraper-service | 8003 | Go | Сбор чеков (8 провайдеров) |
| category-service | 8004 | Go | CRUD категорий |
| budget-service | 8005 | Go | Бюджеты, лимиты |
| goal-service | 8006 | Go | Цели, прогресс |
| credit-service | 8009 | Go | DTI, подушка, таймеры |
| notification-service | 8008 | Go | Push (FCM), Telegram Bot |
| reporting-service | 8010 | Go | Дайджесты, экспорт |
| ai-enrichment | 8100 | Python | Категоризация YandexGPT |
| analytics-service | 8101 | Python | Детектор, прогнозы |
| social-service | 8104 | Python | Челленджи, лидерборды |
| gamification | 8007 | Go | Ачивки, XP, уровни |
| fnps-service | 8084 | Go | ФНС/QR fallback |

## Схема взаимодействия

```
[Frontend] → [api-gateway:8000] → [service:80XX]
                                      ↓
                              [Kafka Event Bus]
                              ↓        ↓        ↓
                       [receipt]  [ai-enrich]  [analytics]
                              ↓        ↓        ↓
                       [PostgreSQL / ClickHouse / Redis]
```
