import type { DashboardSummary } from '~/utils/dashboardSummary'
import type { InsightItem, TimeMachineResponse } from '~/types/api'
import { storeToRefs } from 'pinia'
import { useAiPlanStore } from '~/store/aiPlanStore'

export function useAiPlan() {
  const store = useAiPlanStore()
  const { plan, diagnosisFromPlan, loading, error } = storeToRefs(store)

  async function fetchPlan(
    input: {
      summary: DashboardSummary
      timemachine: TimeMachineResponse | null
      topInsight: InsightItem | null
    },
    options?: { force?: boolean }
  ) {
    await store.fetchPlan(input, options)
  }

  function invalidatePlan() {
    store.invalidate()
  }

  return {
    plan,
    diagnosisFromPlan,
    loading,
    error,
    hasCache: computed(() => store.hasCache),
    fetchPlan,
    invalidatePlan
  }
}
