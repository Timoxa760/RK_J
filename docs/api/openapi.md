# API Архитектура

Базовый URL: `/api/v1`

## Аутентификация

| Метод | Endpoint | Сервис | Описание |
|-------|----------|--------|----------|
| POST | `/api/v1/auth/register` | user-service | Регистрация |
| POST | `/api/v1/auth/login` | user-service | Вход (JWT) |

## Чеки

| Метод | Endpoint | Сервис | Описание |
|-------|----------|--------|----------|
| GET | `/api/v1/receipts` | receipt-service | Лента чеков (itemized) |
| POST | `/api/v1/providers/connect` | user-service | Привязать LK |
| POST | `/api/v1/providers/{name}/sync` | scraper-service | Синхронизация |

## Бюджеты и цели

| Метод | Endpoint | Сервис | Описание |
|-------|----------|--------|----------|
| POST | `/api/v1/budgets` | budget-service | Создать бюджет |
| POST | `/api/v1/goals` | goal-service | Создать цель |

## Кредиты

| Метод | Endpoint | Сервис | Описание |
|-------|----------|--------|----------|
| GET | `/api/v1/credits/dashboard` | credit-service | DTI, подушка |
| POST | `/api/v1/credits/scan` | credit-service | Скан договора |

## Аналитика

| Метод | Endpoint | Сервис | Описание |
|-------|----------|--------|----------|
| POST | `/api/v1/scenarios/simulate` | analytics-service | Time Machine |
| GET | `/api/v1/insights` | analytics-service | Инсайты |
| GET | `/api/v1/digest/latest` | reporting-service | Дайджест |

## Социальное

| Метод | Endpoint | Сервис | Описание |
|-------|----------|--------|----------|
| POST | `/api/v1/challenges` | social-service | Челлендж |
