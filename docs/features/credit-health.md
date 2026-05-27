# 3.2 Credit Health Dashboard

**Цель:** Пользователь видит реальную кредитную нагрузку и рекомендации.

## Расчёты

- **DTI:** `sum(monthly_payments) / avg(income)`
- **Stress-test:** `savings / avg(spend)`

## Технически

credit-service (Go) хранит кредиты в PostgreSQL.
AI-скан договоров через YandexGPT в ai-enrichment.
