import type { ForecastResponse, InsightsResponse, TimeMachineResponse } from '~/types/api'
import { useAuthStore } from '~/store/authStore'

const mockTimeMachine: TimeMachineResponse = {
  points: Array.from({ length: 12 }, (_, i) => ({
    month: `2025-${String(i + 1).padStart(2, '0')}`,
    actual: 50000 + i * 4200,
    optimistic: 52000 + i * 4800
  })),
  delta: 18400
}

const mockForecast: ForecastResponse = {
  points: [
    { month: '2026-06', amount: 72000 },
    { month: '2026-07', amount: 68500 },
    { month: '2026-08', amount: 71000 }
  ]
}

const mockInsights: InsightsResponse = {
  insights: [
    {
      id: '1',
      title: 'Рост трат в кафе',
      body: 'Расходы на кофе выросли на 18%.',
      severity: 'warning'
    },
    {
      id: '2',
      title: 'Накопления в норме',
      body: 'Вы откладываете 22% дохода.',
      severity: 'success'
    }
  ]
}

export function useAnalytics() {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const timeMachine = ref<TimeMachineResponse | null>(null)
  const forecast = ref<ForecastResponse | null>(null)
  const insights = ref<InsightsResponse | null>(null)
  const loading = ref(false)
  const scenarioResult = ref<string | null>(null)

  async function fetchJson<T>(path: string, mock: T): Promise<T> {
    if (config.public.demoMode) return mock
    try {
      return await $fetch<T>(path, {
        baseURL: config.public.apiBase,
        headers: authStore.token
          ? { Authorization: `Bearer ${authStore.token}` }
          : undefined
      })
    } catch {
      return mock
    }
  }

  async function loadAll() {
    loading.value = true
    try {
      timeMachine.value = await fetchJson('/api/v1/dashboard/timemachine', mockTimeMachine)
      forecast.value = await fetchJson('/api/v1/analytics/forecast', mockForecast)
      insights.value = await fetchJson('/api/v1/analytics/insights', mockInsights)
    } finally {
      loading.value = false
    }
  }

  async function simulateScenario(params: { cutCategory: string; percent: number }) {
    const reduction = Math.round(params.percent * 720)
    scenarioResult.value = `Экономия ~${reduction.toLocaleString('ru-RU')} ₽/мес`
    await fetchJson('/api/v1/analytics/simulate', { savings: reduction })
  }

  return {
    timeMachine,
    forecast,
    insights,
    loading,
    scenarioResult,
    loadAll,
    simulateScenario
  }
}
