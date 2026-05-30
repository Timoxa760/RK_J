import type { InsightsResponse } from '~/types/api'
import { normalizeInsights } from '~/utils/apiNormalize'

export function useInsights() {
  const { apiFetch } = useApi()

  const insights = ref<InsightsResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const topInsight = computed(() => insights.value?.insights?.[0] ?? null)

  async function fetchInsights() {
    loading.value = true
    error.value = null
    try {
      const raw = await apiFetch<InsightsResponse>('/insights')
      insights.value = normalizeInsights(raw)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить советы'
      insights.value = null
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
