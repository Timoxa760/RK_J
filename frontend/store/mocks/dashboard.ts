/**
 * Mock-данные dashboard/credits/insights — по docs/api/API_Contract.md (v3).
 */
import type {
  CategoriesResponse,
  CompareResponse,
  CreditsDashboardResponse,
  InsightsResponse,
  SankeyResponse,
  StoresResponse,
  TimeMachineApiResponse
} from '~/types/api'

export const mockSankey: SankeyResponse = {
  nodes: [
    { name: 'Зарплата', value: 180000 },
    { name: 'Накопления', value: 35000 },
    { name: 'Продукты', value: 52000 },
    { name: 'Кафе и рестораны', value: 28000 },
    { name: 'Транспорт', value: 15000 },
    { name: 'Развлечения', value: 12000 }
  ],
  links: [
    { source: 'Зарплата', target: 'Накопления', value: 35000 },
    { source: 'Зарплата', target: 'Продукты', value: 52000 },
    { source: 'Зарплата', target: 'Кафе и рестораны', value: 28000 },
    { source: 'Зарплата', target: 'Транспорт', value: 15000 },
    { source: 'Зарплата', target: 'Развлечения', value: 12000 }
  ]
}

export const mockStores: StoresResponse = {
  stores: [
    { name: 'Пятёрочка', avg_check: 650, purchases: 14, total: 9100, impulse_ratio: 0.25 },
    { name: 'Магнит', avg_check: 720, purchases: 10, total: 7200, impulse_ratio: 0.2 },
    { name: 'ВкусВилл', avg_check: 980, purchases: 7, total: 6860, impulse_ratio: 0.1 },
    { name: 'Ozon', avg_check: 2100, purchases: 4, total: 8400, impulse_ratio: 0.65 },
    { name: 'Wildberries', avg_check: 1850, purchases: 5, total: 9250, impulse_ratio: 0.7 },
    { name: 'Лента', avg_check: 820, purchases: 8, total: 6560, impulse_ratio: 0.15 }
  ]
}

export const mockCategories: CategoriesResponse = {
  categories: [
    {
      name: 'Продукты',
      total: 52000,
      subcategories: [
        {
          name: 'Молочные',
          total: 8500,
          items: [
            { name: 'Молоко 3.2%', price: 78, quantity: 12, total: 936 },
            { name: 'Творог 5%', price: 120, quantity: 6, total: 720 }
          ]
        }
      ]
    }
  ]
}

export const mockCompare: CompareResponse = {
  months: [
    {
      label: 'Февраль 2026',
      categories: [
        { name: 'Продукты', total: 48000 },
        { name: 'Кафе и рестораны', total: 25000 }
      ]
    },
    {
      label: 'Март 2026',
      categories: [
        { name: 'Продукты', total: 52000 },
        { name: 'Кафе и рестораны', total: 28000 }
      ]
    }
  ],
  insights: {
    biggest_change: { category: 'Кафе и рестораны', delta: 3000, delta_percent: 12 }
  }
}

export const mockTimeMachine: TimeMachineApiResponse = {
  months: ['2026-05', '2026-06', '2026-07'],
  real_savings: [500000, 512000, 524500],
  optimized_savings: [500000, 516000, 532500],
  difference_final: 467000
}

export const mockInsights: InsightsResponse = {
  insights: [
    {
      type: 'subscription',
      severity: 'warning',
      title: 'Найдена скрытая подписка',
      description: 'Списывается 299 ₽ каждый месяц',
      amount: 299,
      merchant: 'Яндекс.Плюс'
    },
    {
      type: 'duplicate',
      severity: 'info',
      title: 'Дублирование в чеке',
      description: "Товар 'Молоко 3.2%' пробит дважды",
      amount: 156
    },
    {
      type: 'overprice',
      severity: 'warning',
      title: 'Переплата за товар',
      description: 'Молоко 3.2% куплено за 95 ₽, средняя — 78 ₽',
      amount: 17,
      store: 'Пятёрочка'
    }
  ]
}

export const mockCredits: CreditsDashboardResponse = {
  dti: 0.28,
  stress_test_months: 4.2,
  savings: 340000,
  total_debt: 1200000,
  monthly_payments: 42000,
  monthly_income: 180000,
  credits: [
    {
      id: 'uuid',
      bank: 'Т-Банк',
      amount: 1200000,
      rate: 14.5,
      term_months: 36,
      remaining: 980000,
      monthly_payment: 42000,
      next_payment: '2026-06-15'
    }
  ]
}
