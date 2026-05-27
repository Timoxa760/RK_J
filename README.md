# MoneyMind — Frontend

Первый в России поштучный финансовый ассистент.

## Стек

- React 19 + Vite
- Tailwind CSS 4
- Recharts (графики)
- PWA

## Структура

```
src/
├── app/                          # Страницы (роутинг)
│   ├── auth/                     # Вход / регистрация
│   ├── dashboard/                # Главная
│   ├── receipts/                 # Чеки
│   ├── analytics/                # Аналитика
│   ├── budget/                   # Бюджеты
│   ├── goals/                    # Цели
│   ├── credits/                  # Кредиты
│   ├── social/                   # Социальное
│   ├── digest/                   # Дайджест
│   ├── profile/                  # Профиль
│   └── settings/                 # Настройки
├── components/                   # Компоненты
│   ├── ui/                       # Базовые UI-kit
│   ├── layout/                   # Шапка, сайдбар
│   ├── charts/                   # Графики Recharts
│   ├── receipts/                 # Блоки чеков
│   ├── budget/                   # Бюджетные виджеты
│   ├── goals/                    # Цели
│   ├── credits/                  # Кредитный хелс
│   ├── social/                   # Челленджи
│   └── digest/                   # Дайджест
├── composables/                  # Хуки
├── stores/                       # Состояние (Pinia)
├── api/                          # API клиент
├── types/                        # TypeScript типы
├── utils/                        # Хелперы
├── assets/                       # Ресурсы
│   ├── images/
│   └── fonts/
├── public/                       # Статика
└── tests/                        # Тесты
```

## Установка и запуск

```bash
npm install
npm run dev
```

## Экраны

| Путь | Раздел |
|------|--------|
| `/dashboard` | Главная |
| `/receipts` | Чеки |
| `/analytics` | Аналитика |
| `/budget` | Бюджеты |
| `/goals` | Цели |
| `/credits` | Кредиты |
| `/social` | Социальное |
| `/digest` | Дайджест |
| `/profile` | Профиль |
| `/settings` | Настройки |
