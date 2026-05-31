# Кредиты: PDF-скан и rates-aggregator

> Единственный источник кредитных данных в MVP.

## Поток

```mermaid
sequenceDiagram
  participant UI as Front /credits
  participant GW as api-gateway
  participant CS as credit-service
  participant LLM as OnlySQ
  participant RA as rates-aggregator
  participant PG as PostgreSQL

  UI->>GW: POST /credits/scan PDF
  GW->>CS: multipart file
  CS->>LLM: extract contract fields
  CS->>RA: benchmark rate
  CS->>PG: INSERT user_credits
  CS-->>UI: parsed + rate_comparison
  UI->>GW: GET /credits/dashboard
  GW->>CS: aggregate user_credits
  CS-->>UI: dti, credits[]
```

## POST /credits/scan

- `multipart/form-data`, поле `file` (PDF)
- Text extraction поддерживает **кириллические** договоры (pdf parser + OnlySQ)
- OnlySQ извлекает: bank, amount, rate, term_months, monthly_payment
- rates-aggregator возвращает `benchmark_rate` для сравнения
- Запись в `user_credits`

## GET /credits/dashboard

- Агрегат по `user_credits` пользователя
- Пустой ответ (dti=0, credits=[]) до первого скана
- `dti` в **процентах** 0–100 (как на front)

## rates-aggregator

Env: `RATES_AGGREGATOR_URL`, `RATES_AGGREGATOR_API_KEY`.

MVP: mock + ключевая ставка ЦБ (fallback). Провайдер (Banki.ru и т.д.) — позже.

Internal: `GET /internal/rates?product=consumer&amount=&term=`

## Не используется

- bank-service для кредитов
- demo hardcode в handler
- ручной ввод кредита без PDF
