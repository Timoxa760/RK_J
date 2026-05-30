import type { InsightsResponse } from '~/types/api'
import { mockInsights } from '~/store/mocks'
import { normalizeInsights } from '~/utils/apiNormalize'

export function useInsights() {
  const { apiFetchWithDemo, demoMode } = useApi()

  const insights = ref<InsightsResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const topInsight = computed(() => insights.value?.insights?.[0] ?? null)

  async function fetchInsights() {
    loading.value = true
    error.value = null
    try {
      const raw = await apiFetchWithDemo('/insights', mockInsights)
      insights.value = normalizeInsights(raw)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить советы'
      if (demoMode.value) {
        insights.value = normalizeInsights(mockInsights)
      }
    } finally {
      loading.value = false
    }
  }

  return {
    insights,
    topInsight,
    loading,
    error,
    fetchInsights
  }
}
