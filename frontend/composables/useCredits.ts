import type { CreditScanResponse, CreditsDashboardResponse } from '~/types/api'
import { mockCredits } from '~/store/mocks'
import { normalizeCredits } from '~/utils/apiNormalize'

export function useCredits() {
  const { apiFetch, apiFetchWithDemo, demoMode } = useApi()

  const dashboard = ref<CreditsDashboardResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const scanResult = ref<CreditScanResponse | null>(null)
  const scanLoading = ref(false)

  async function fetchDashboard() {
    loading.value = true
    error.value = null
    try {
      const raw = await apiFetchWithDemo<CreditsDashboardResponse>(
        '/credits/dashboard',
        mockCredits
      )
      dashboard.value = normalizeCredits(raw)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка загрузки'
      if (demoMode.value) {
        dashboard.value = normalizeCredits(mockCredits)
      }
    } finally {
      loading.value = false
    }
  }

  async function scanContract(file: File) {
    scanLoading.value = true
    error.value = null
    try {
      const form = new FormData()
      form.append('file', file)
      scanResult.value = await apiFetch<CreditScanResponse>('/credits/scan', {
        method: 'POST',
        body: form
      })
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
    fetchDashboard,
    scanContract
  }
}
