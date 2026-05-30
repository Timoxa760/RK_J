import type { AiDiagnosisResponse, AiPlanApiResponse, InsightItem, TimeMachineResponse } from '~/types/api'
import type { DashboardSummary } from '~/utils/dashboardSummary'
import type { FinancialPlan } from '~/utils/financialPlan'
import { buildFinancialPlan } from '~/utils/financialPlan'
import { goalFromProfile } from '~/composables/useGoals'

export function useAiPlan() {
  const { apiFetch } = useApi()
  const { profile, loadProfile } = useFinancialProfile()

  const plan = ref<FinancialPlan | null>(null)
  const diagnosisFromPlan = ref<AiDiagnosisResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchPlan(input: {
    summary: DashboardSummary
    timemachine: TimeMachineResponse | null
    topInsight: InsightItem | null
  }) {
    loading.value = true
    error.value = null
    loadProfile()

    try {
      const res = await apiFetch<AiPlanApiResponse>('/ai/plan')
      plan.value = res.plan
      diagnosisFromPlan.value = res.diagnosis
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить план'
      plan.value = buildFinancialPlan({
        primaryGoal: goalFromProfile(profile.value),
        summary: input.summary,
        timemachine: input.timemachine,
        diagnosis: null,
        topInsight: input.topInsight
      })
      diagnosisFromPlan.value = null
    } finally {
      loading.value = false
    }
  }

  return {
    plan,
    diagnosisFromPlan,
    loading,
    error,
    fetchPlan
  }
}
