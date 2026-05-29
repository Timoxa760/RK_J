import type { CreditsDashboardResponse } from '~/types/api'
import { useAuthStore } from '~/store/authStore'

const mockCredits: CreditsDashboardResponse = {
  dti: 38,
  stress_test_dti: 52,
  monthly_income: 120000,
  credits: [
    { id: '1', name: 'Ипотека', balance: 3200000, payment: 28000, rate: 12.5 },
    { id: '2', name: 'Автокредит', balance: 450000, payment: 12000, rate: 15.9 },
    { id: '3', name: 'Кредитная карта', balance: 85000, payment: 8500, rate: 24.9 }
  ]
}

export function useCredits() {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const dashboard = ref<CreditsDashboardResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchDashboard() {
    loading.value = true
    error.value = null
    try {
      if (config.public.demoMode) {
        dashboard.value = mockCredits
        return
      }
      dashboard.value = await $fetch<CreditsDashboardResponse>('/api/v1/credits/dashboard', {
        baseURL: config.public.apiBase,
        headers: authStore.token
          ? { Authorization: `Bearer ${authStore.token}` }
          : undefined
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка загрузки'
      dashboard.value = mockCredits
    } finally {
      loading.value = false
    }
  }

  return { dashboard, loading, error, fetchDashboard }
}
