# TypeScript-типы API (front)

> **Эталон для интеграции:** `front/frontend/types/api.ts`  
> Контракт: [API_Contract.md](./API_Contract.md)  
> При расхождении — **приоритет у front-типов** для UI; back подгоняется или docs фиксирует delta (см. API_Contract §13).

## Dashboard

```typescript
interface SankeyResponse {
  nodes: { name: string; category?: 'income' | 'category' | 'savings'; value?: number }[]
  links: { source: string; target: string; value: number }[]
}

interface StoresResponse {
  stores: {
    id: string
    name: string
    avg_check: number
    visits: number
    total: number
    impulse_ratio: number
  }[]
}

interface CategoriesResponse {
  categories: {
    name: string
    amount: number
    share: number
    subcategories: {
      name: string
      items: { name: string; amount: number }[]
    }[]
  }[]
}

interface CompareResponse {
  months: { month: string; categories: { name: string; amount: number; share: number }[] }[]
}

interface TimeMachineResponse {
  points: { month: string; actual: number; optimistic: number }[]
  delta: number
}
```

## Credits (кредитный светофор)

```typescript
interface CreditsDashboardResponse {
  dti: number              // проценты 0–100, не доля 0–1
  stress_test_dti: number
  monthly_income: number
  credits: {
    id: string
    name: string
    balance: number
    payment: number
    rate: number
  }[]
}
```

## Analytics

```typescript
interface InsightsResponse {
  insights: {
    id: string
    title: string
    body: string
    severity: 'info' | 'warning' | 'success'
  }[]
}

interface ForecastResponse {
  points: { month: string; amount: number }[]
}
```

## Expenses (ручной / голос)

```typescript
interface ManualExpenseRequest {
  user_id: string
  raw_text?: string
  amount?: number
  category?: string
  description?: string
  date?: string       // YYYY-MM-DD
  source?: 'manual' | 'voice'
}

interface ManualExpenseResponse {
  success: boolean
  id?: string
  amount: number
  category: string
  parsed: boolean
}
```

## Auth (из API_Contract, типов в front пока нет)

```typescript
interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: { id: string; phone: string; role: string }
}
```

## Composables → endpoints

| Composable | Endpoint |
|------------|----------|
| `useDashboard` | `/dashboard/sankey`, `/categories`, `/stores`, `/compare`, `/timemachine` |
| `useCredits` | `/credits/dashboard` |
| `useAnalytics` | `/insights`, `/forecast` |
| `useSocial` | `/challenges/*` |

## Связи

- [API_Contract.md](./API_Contract.md)
- [../deployment/front-quickstart.md](../deployment/front-quickstart.md)
