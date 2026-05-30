import { defineStore } from 'pinia'
import type { CategoriesResponse, TimeMachineResponse } from '~/types/api'
import { mockCategories, mockTimeMachine } from '~/store/mocks'
import { normalizeCategories, normalizeTimeMachine } from '~/utils/apiNormalize'

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    categories: null as CategoriesResponse | null,
    timemachine: null as TimeMachineResponse | null,
    loading: false,
    error: null as string | null
  }),

  actions: {
    async loadAll(options?: { silent?: boolean }) {
      const { apiFetchWithDemo, demoMode } = useApi()
      if (!options?.silent) {
        this.loading = true
      }
      this.error = null

      try {
        const [categoriesRaw, timemachineRaw] = await Promise.all([
          apiFetchWithDemo('/dashboard/categories', mockCategories),
          apiFetchWithDemo('/dashboard/timemachine', mockTimeMachine)
        ])

        this.categories = normalizeCategories(categoriesRaw)
        this.timemachine = normalizeTimeMachine(timemachineRaw)
      } catch (e) {
        this.error = e instanceof Error ? e.message : 'Ошибка загрузки дашборда'
        if (demoMode.value) {
          this.categories = normalizeCategories(mockCategories)
          this.timemachine = normalizeTimeMachine(mockTimeMachine)
        }
      } finally {
        if (!options?.silent) {
          this.loading = false
        }
      }
    }
  }
})
