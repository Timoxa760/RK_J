import type { AdvisorContext } from '~/utils/advisorChat'
import { buildDashboardSummary } from '~/utils/dashboardSummary'
import { mockInsights } from '~/store/mocks'
import { normalizeInsights } from '~/utils/apiNormalize'
import type { InsightsResponse } from '~/types/api'

/** Данные для советника в сайдбаре (на всех страницах приложения). */
export function useAdvisorContext() {
  const { profile, loadProfile } = useFinancialProfile()
  const { primaryGoal, fetchGoals } = useGoals()
  const { diagnosis, fetchDiagnosis } = useDiagnosis()
  const { topInsight, fetchInsights } = useInsights()
  const { timemachine, loadAll: loadDashboardSlice } = useDashboard()
  const { dashboard: credits, fetchDashboard: fetchCredits } = useCredits()
  const { apiFetchWithDemo, demoMode } = useApi()

  const insightsData = ref<InsightsResponse | null>(null)
  const loading = ref(false)

  async function loadInsightCard() {
    try {
      const raw = await apiFetchWithDemo('/insights', mockInsights)
      insightsData.value = normalizeInsights(raw)
    } catch {
      if (demoMode.value) {
        insightsData.value = normalizeInsights(mockInsights)
      }
    }
  }

  const summary = computed(() =>
    buildDashboardSummary({
      profile: profile.value,
      timemachine: timemachine.value,
      credits: credits.value,
      topInsight: topInsight.value ?? insightsData.value?.insights?.[0] ?? null
    })
  )

  const advisorContext = computed<AdvisorContext>(() => ({
    diagnosis: diagnosis.value,
    topInsight: topInsight.value ?? insightsData.value?.insights?.[0] ?? null,
    timemachine: timemachine.value,
    primaryGoal: primaryGoal.value,
    goalForecast: summary.value.goalForecast
  }))

  async function refreshAdvisorContext(options?: { silent?: boolean }) {
    if (!options?.silent) loading.value = true
    loadProfile()
    try {
      await Promise.all([
        loadDashboardSlice({ silent: options?.silent }),
        fetchCredits(),
        fetchInsights(),
        loadInsightCard(),
        fetchDiagnosis(),
        fetchGoals()
      ])
    } finally {
      if (!options?.silent) loading.value = false
    }
  }

  return {
    advisorContext,
    refreshAdvisorContext,
    loading
  }
}
