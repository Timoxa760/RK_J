import type { AdvisorContext } from '~/utils/advisorChat'
import { buildDashboardSummary } from '~/utils/dashboardSummary'
import { resolveSavingsTimemachine } from '~/utils/dashboardProjections'

/** Данные для советника в сайдбаре (на всех страницах приложения). */
export function useAdvisorContext() {
  const { profile, loadProfile } = useFinancialProfile()
  const { primaryGoal } = useGoals()
  const { diagnosis, fetchDiagnosis } = useDiagnosis()
  const { topInsight, fetchInsights } = useInsights()
  const { categories, timemachine, loadAll: loadDashboardSlice } = useDashboard()
  const { dashboard: credits, fetchDashboard: fetchCredits } = useCredits()

  const loading = ref(false)

  const projectedTimemachine = computed(() =>
    resolveSavingsTimemachine(timemachine.value, profile.value, categories.value)
  )

  const summary = computed(() =>
    buildDashboardSummary({
      profile: profile.value,
      timemachine: projectedTimemachine.value,
      credits: credits.value,
      topInsight: topInsight.value,
      categories: categories.value
    })
  )

  const advisorContext = computed<AdvisorContext>(() => ({
    diagnosis: diagnosis.value,
    topInsight: topInsight.value,
    timemachine: projectedTimemachine.value,
    primaryGoal: primaryGoal.value,
    goalForecast: summary.value.goalForecast,
    categories: categories.value
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
