import { defineStore } from 'pinia'
import { useAuthStore } from '~/store/authStore'
import type {
  CategoriesResponse,
  CompareResponse,
  SankeyResponse,
  StoresResponse,
  TimeMachineResponse
} from '~/types/api'

const mockSankey: SankeyResponse = {
  nodes: [
    { name: 'Зарплата', category: 'income' },
    { name: 'Фриланс', category: 'income' },
    { name: 'Накопления', category: 'savings' },
    { name: 'Продукты', category: 'category' },
    { name: 'Кафе', category: 'category' },
    { name: 'Транспорт', category: 'category' },
    { name: 'Развлечения', category: 'category' },
    { name: 'ЖКХ', category: 'category' }
  ],
  links: [
    { source: 'Зарплата', target: 'Продукты', value: 28000 },
    { source: 'Зарплата', target: 'Кафе', value: 8000 },
    { source: 'Зарплата', target: 'Транспорт', value: 6000 },
    { source: 'Зарплата', target: 'ЖКХ', value: 12000 },
    { source: 'Зарплата', target: 'Накопления', value: 15000 },
    { source: 'Фриланс', target: 'Развлечения', value: 5000 },
    { source: 'Фриланс', target: 'Накопления', value: 7000 }
  ]
}

const mockStores: StoresResponse = {
  stores: [
    { id: '1', name: 'Пятёрочка', avg_check: 850, visits: 24, total: 20400, impulse_ratio: 0.15 },
    { id: '2', name: 'Starbucks', avg_check: 420, visits: 18, total: 7560, impulse_ratio: 0.72 },
    { id: '3', name: 'Ozon', avg_check: 3200, visits: 8, total: 25600, impulse_ratio: 0.45 },
    { id: '4', name: 'Метро', avg_check: 1200, visits: 12, total: 14400, impulse_ratio: 0.08 },
    { id: '5', name: 'DNS', avg_check: 8900, visits: 2, total: 17800, impulse_ratio: 0.55 },
    { id: '6', name: 'Аптека 36.6', avg_check: 650, visits: 6, total: 3900, impulse_ratio: 0.22 }
  ]
}

const mockCategories: CategoriesResponse = {
  categories: [
    {
      name: 'Продукты',
      amount: 28000,
      share: 0.32,
      subcategories: [
        {
          name: 'Молочные',
          items: [
            { name: 'Молоко', amount: 2400 },
            { name: 'Творог', amount: 1800 },
            { name: 'Сыр', amount: 3200 }
          ]
        }
      ]
    },
    { name: 'Кафе', amount: 8000, share: 0.09, subcategories: [] },
    { name: 'Транспорт', amount: 6000, share: 0.07, subcategories: [] },
    { name: 'ЖКХ', amount: 12000, share: 0.14, subcategories: [] },
    { name: 'Развлечения', amount: 5000, share: 0.06, subcategories: [] }
  ]
}

const mockCompare: CompareResponse = {
  months: [
    {
      month: '2026-04',
      categories: [
        { name: 'Продукты', amount: 26000, share: 0.3 },
        { name: 'Кафе', amount: 9000, share: 0.1 }
      ]
    },
    {
      month: '2026-05',
      categories: [
        { name: 'Продукты', amount: 28000, share: 0.32 },
        { name: 'Кафе', amount: 8000, share: 0.09 }
      ]
    }
  ]
}

const mockTimeMachine: TimeMachineResponse = {
  points: Array.from({ length: 12 }, (_, i) => ({
    month: `2025-${String(i + 1).padStart(2, '0')}`,
    actual: 50000 + i * 4200,
    optimistic: 52000 + i * 4800
  })),
  delta: 18400
}

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    sankey: null as SankeyResponse | null,
    stores: null as StoresResponse | null,
    categories: null as CategoriesResponse | null,
    compare: null as CompareResponse | null,
    timemachine: null as TimeMachineResponse | null,
    loading: false,
    error: null as string | null
  }),

  actions: {
    async fetchJson<T>(path: string, mock: T): Promise<T> {
      const config = useRuntimeConfig()
      if (config.public.demoMode) return mock

      const authStore = useAuthStore()
      try {
        return await $fetch<T>(path, {
          baseURL: config.public.apiBase,
          headers: authStore.token
            ? { Authorization: `Bearer ${authStore.token}` }
            : undefined
        })
      } catch {
        return mock
      }
    },

    async loadAll() {
      this.loading = true
      this.error = null
      try {
        const [sankey, stores, categories, compare, timemachine] = await Promise.all([
          this.fetchJson('/api/v1/dashboard/sankey', mockSankey),
          this.fetchJson('/api/v1/dashboard/stores', mockStores),
          this.fetchJson('/api/v1/dashboard/categories', mockCategories),
          this.fetchJson('/api/v1/dashboard/compare?months=2', mockCompare),
          this.fetchJson('/api/v1/dashboard/timemachine', mockTimeMachine)
        ])
        this.sankey = sankey
        this.stores = stores
        this.categories = categories
        this.compare = compare
        this.timemachine = timemachine
      } catch (e) {
        this.error = e instanceof Error ? e.message : 'Ошибка загрузки дашборда'
        this.sankey = mockSankey
        this.stores = mockStores
        this.categories = mockCategories
        this.compare = mockCompare
        this.timemachine = mockTimeMachine
      } finally {
        this.loading = false
      }
    }
  }
})
