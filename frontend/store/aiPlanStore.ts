import { defineStore } from 'pinia'
import type { AiDiagnosisResponse, AiPlanApiResponse, InsightItem, TimeMachineResponse } from '~/types/api'
import type { DashboardSummary } from '~/utils/dashboardSummary'
import type { FinancialPlan } from '~/utils/financialPlan'
import { applyProfileGoalToPlan, buildFinancialPlan } from '~/utils/financialPlan'
import { goalFromProfile } from '~/composables/useGoals'
import {
  clearStoredAiPlanForCurrentUser,
  readStoredAiPlan,
  writeStoredAiPlan
} from '~/utils/aiPlanStorage'

export const useAiPlanStore = defineStore('aiPlan', {
  state: () => ({
    plan: null as FinancialPlan | null,
    diagnosisFromPlan: null as AiDiagnosisResponse | null,
    loading: false,
    error: null as string | null,
    loadedAt: 0,
    hydratedFromStorage: false
  }),

  getters: {
    hasCache(state): boolean {
      return Boolean(state.plan && state.loadedAt > 0)
    }
  },

  actions: {
    hydrateFromStorage() {
      if (this.hydratedFromStorage || this.loadedAt > 0) return
      const stored = readStoredAiPlan()
      if (!stored) {
        this.hydratedFromStorage = true
        return
      }

      const { profile, loadProfile } = useFinancialProfile()
      loadProfile()
      const primaryGoal = goalFromProfile(profile.value)

      this.plan = stored.plan
        ? applyProfileGoalToPlan(stored.plan, primaryGoal)
        : null
      this.diagnosisFromPlan = stored.diagnosis
      this.loadedAt = stored.loadedAt
      this.hydratedFromStorage = true
    },

    persistToStorage() {
      if (!this.plan || !this.loadedAt) return
      writeStoredAiPlan({
        plan: this.plan,
        diagnosis: this.diagnosisFromPlan,
        loadedAt: this.loadedAt
      })
    },

    invalidate() {
      this.loadedAt = 0
      clearStoredAiPlanForCurrentUser()
    },

    clearCache() {
      this.plan = null
      this.diagnosisFromPlan = null
      this.loadedAt = 0
      this.error = null
      this.hydratedFromStorage = false
    },

    async fetchPlan(
      input: {
        summary: DashboardSummary
        timemachine: TimeMachineResponse | null
        topInsight: InsightItem | null
      },
      options?: { force?: boolean }
    ) {
      this.hydrateFromStorage()
      if (!options?.force && this.hasCache) return

      const { apiFetch } = useApi()
      const { profile, loadProfile, fetchProfileFromApi } = useFinancialProfile()

      this.loading = true
      this.error = null
      loadProfile()
      await fetchProfileFromApi()

      const primaryGoal = goalFromProfile(profile.value)

      try {
        const res = await apiFetch<AiPlanApiResponse>('/ai/plan')
        this.plan = applyProfileGoalToPlan(res.plan, primaryGoal)
        this.diagnosisFromPlan = res.diagnosis
        this.loadedAt = Date.now()
        this.persistToStorage()
      } catch (e) {
        this.error = e instanceof Error ? e.message : 'Не удалось загрузить план'
        this.plan = buildFinancialPlan({
          primaryGoal,
          summary: input.summary,
          timemachine: input.timemachine,
          diagnosis: null,
          topInsight: input.topInsight
        })
        this.diagnosisFromPlan = null
        this.loadedAt = Date.now()
        this.persistToStorage()
      } finally {
        this.loading = false
      }
    }
  }
})
