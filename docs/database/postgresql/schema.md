# PostgreSQL — Модель данных

## Таблицы

| Таблица | Назначение |
|---------|------------|
| `users` | Пользователи (телефон, JWT, роль) |
| `user_providers` | Привязки LK (креды зашифрованы AES-256-GCM) |
| `receipts` | Чеки (itemized, items как JSONB, чексума UNIQUE) |
| `categories` | Категории (дерево: parent_id) |
| `budgets` | Лимиты по категориям |
| `goals` | Цели с target_amount и auto_save_percent |
| `credits` | Кредиты со ставкой, графиком, DTI |
| `challenges` | Челленджи |
| `achievements` | Ачивки и XP |
| `daily_activity` | Streak трекинг |
| `scraper_sessions` | Кэш сессий скрепера |
