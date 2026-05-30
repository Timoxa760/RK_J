import type { CreditScanResponse, CreditsDashboardResponse } from '~/types/api'
import { normalizeCredits } from '~/utils/apiNormalize'

const emptyDashboard = (): CreditsDashboardResponse => ({
  dti: 0,
  savings: 0,
  total_debt: 0,
  monthly_payments: 0,
  monthly_income: 0,
  credits: []
})

export function useCredits() {
  const { apiFetch } = useApi()

  const dashboard = ref<CreditsDashboardResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const scanResult = ref<CreditScanResponse | null>(null)
  const scanLoading = ref(false)

  const hasCredits = computed(() => (dashboard.value?.credits?.length ?? 0) > 0)

  async function fetchDashboard() {
    loading.value = true
    error.value = null
    try {
      const raw = await apiFetch<CreditsDashboardResponse>('/credits/dashboard')
      dashboard.value = normalizeCredits(raw)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка загрузки'
      dashboard.value = null
    } finally {
      loading.value = false
    }
  }

  async function scanContract(file: File) {
    scanLoading.value = true
    error.value = null
    scanResult.value = null
    try {
      const form = new FormData()
      form.append('file', file)
      scanResult.value = await apiFetch<CreditScanResponse>('/credits/scan', {
        method: 'POST',
        body: form
      })
      await fetchDashboard()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось распознать договор'
    } finally {
      scanLoading.value = false
    }
  }

  return {
    dashboard,
    loading,
    error,
    scanResult,
    scanLoading,
    hasCredits,
    fetchDashboard,
    scanContract,
    emptyDashboard
  }
}
