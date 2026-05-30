import type { AiDiagnosisResponse, AiPlanApiResponse, InsightItem, TimeMachineResponse } from '~/types/api'
import type { DashboardSummary } from '~/utils/dashboardSummary'
import type { FinancialPlan } from '~/utils/financialPlan'
import { mockDiagnosis } from '~/store/mocks/diagnosis'
import { buildFinancialPlan } from '~/utils/financialPlan'
import { goalFromProfile } from '~/composables/useGoals'

export function useAiPlan() {
  const { apiFetchWithDemo, demoMode } = useApi()
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

    const fallbackPlan = buildFinancialPlan({
      primaryGoal: goalFromProfile(profile.value),
      summary: input.summary,
      timemachine: input.timemachine,
      diagnosis: mockDiagnosis,
      topInsight: input.topInsight
    })

    try {
      const res = await apiFetchWithDemo<AiPlanApiResponse>('/ai/plan', {
        plan: fallbackPlan,
        diagnosis: mockDiagnosis
      })
      plan.value = res.plan
      diagnosisFromPlan.value = res.diagnosis
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить план'
      if (demoMode.value) {
        plan.value = fallbackPlan
        diagnosisFromPlan.value = mockDiagnosis
      }
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
