import { defineStore } from 'pinia'
import type { CategoriesResponse, TimeMachineResponse } from '~/types/api'
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
      const { apiFetch } = useApi()
      if (!options?.silent) {
        this.loading = true
      }
      this.error = null

      try {
        const [categoriesRaw, timemachineRaw] = await Promise.all([
          apiFetch<CategoriesResponse>('/dashboard/categories'),
          apiFetch<TimeMachineResponse>('/dashboard/timemachine')
        ])

        this.categories = normalizeCategories(categoriesRaw)
        this.timemachine = normalizeTimeMachine(timemachineRaw)
      } catch (e) {
        this.error = e instanceof Error ? e.message : 'Ошибка загрузки дашборда'
        this.categories = null
        this.timemachine = null
      } finally {
        if (!options?.silent) {
          this.loading = false
        }
      }
    }
  }
})
