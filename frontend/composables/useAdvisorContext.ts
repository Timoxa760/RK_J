import type { AdvisorContext } from '~/utils/advisorChat'
import { buildDashboardSummary } from '~/utils/dashboardSummary'

/** Данные для советника в сайдбаре (на всех страницах приложения). */
export function useAdvisorContext() {
  const { profile, loadProfile } = useFinancialProfile()
  const { primaryGoal } = useGoals()
  const { diagnosis, fetchDiagnosis } = useDiagnosis()
  const { topInsight, fetchInsights } = useInsights()
  const { timemachine, loadAll: loadDashboardSlice } = useDashboard()
  const { dashboard: credits, fetchDashboard: fetchCredits } = useCredits()

  const loading = ref(false)

  const summary = computed(() =>
    buildDashboardSummary({
      profile: profile.value,
      timemachine: timemachine.value,
      credits: credits.value,
      topInsight: topInsight.value
    })
  )

  const advisorContext = computed<AdvisorContext>(() => ({
    diagnosis: diagnosis.value,
    topInsight: topInsight.value,
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
        fetchDiagnosis()
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
