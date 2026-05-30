import { defineStore } from 'pinia'
import type {
  CategoriesResponse,
  CompareResponse,
  SankeyResponse,
  StoresResponse,
  TimeMachineResponse
} from '~/types/api'
import {
  mockCategories,
  mockCompare,
  mockSankey,
  mockStores,
  mockTimeMachine
} from '~/store/mocks'
import {
  normalizeCategories,
  normalizeCompare,
  normalizeSankey,
  normalizeStores,
  normalizeTimeMachine
} from '~/utils/apiNormalize'

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
    async loadAll() {
      const { apiFetchWithDemo, demoMode } = useApi()
      this.loading = true
      this.error = null

      try {
        const [sankeyRaw, storesRaw, categoriesRaw, compareRaw, timemachineRaw] =
          await Promise.all([
            apiFetchWithDemo('/dashboard/sankey', mockSankey),
            apiFetchWithDemo('/dashboard/stores', mockStores),
            apiFetchWithDemo('/dashboard/categories', mockCategories),
            apiFetchWithDemo('/dashboard/compare?months=2', mockCompare),
            apiFetchWithDemo('/dashboard/timemachine', mockTimeMachine)
          ])

        this.sankey = normalizeSankey(sankeyRaw)
        this.stores = normalizeStores(storesRaw)
        this.categories = normalizeCategories(categoriesRaw)
        this.compare = normalizeCompare(compareRaw)
        this.timemachine = normalizeTimeMachine(timemachineRaw)
      } catch (e) {
        this.error = e instanceof Error ? e.message : 'Ошибка загрузки дашборда'
        if (demoMode.value) {
          this.sankey = mockSankey
          this.stores = normalizeStores(mockStores)
          this.categories = normalizeCategories(mockCategories)
          this.compare = normalizeCompare(mockCompare)
          this.timemachine = normalizeTimeMachine(mockTimeMachine)
        }
      } finally {
        this.loading = false
      }
    }
  }
})
