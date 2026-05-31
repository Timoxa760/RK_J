import type { CreditScanResponse, CreditsDashboardResponse } from '~/types/api'
import { normalizeCredits } from '~/utils/apiNormalize'
import { enrichCreditsDashboard } from '~/utils/creditsDashboard'
import { hasCreditsData } from '~/utils/dashboardSummary'
import { formatApiError } from '~/utils/apiError'

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
  const { profile, loadProfile } = useFinancialProfile()

  if (import.meta.client) {
    loadProfile()
  }

  const dashboard = useState<CreditsDashboardResponse | null>('credits-dashboard', () => null)
  const loading = useState('credits-dashboard-loading', () => false)
  const fetchError = useState<string | null>('credits-fetch-error', () => null)
  const scanResult = useState<CreditScanResponse | null>('credits-scan-result', () => null)
  const scanLoading = useState('credits-scan-loading', () => false)
  const scanError = useState<string | null>('credits-scan-error', () => null)
  const deleting = useState('credits-delete-loading', () => false)

  const hasCredits = computed(() => hasCreditsData(dashboard.value))

  const enrichedDashboard = computed(() =>
    dashboard.value ? enrichCreditsDashboard(dashboard.value, profile.value) : null
  )

  async function fetchDashboard() {
    loading.value = true
    fetchError.value = null
    try {
      const raw = await apiFetch<CreditsDashboardResponse>('/credits/dashboard')
      dashboard.value = normalizeCredits(raw)
    } catch (e) {
      fetchError.value = formatApiError(e, 'Не удалось загрузить данные по кредитам')
      dashboard.value = null
    } finally {
      loading.value = false
    }
  }

  async function scanContract(file: File) {
    scanLoading.value = true
    scanError.value = null
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
      scanError.value = formatApiError(e, 'Не удалось распознать договор')
    } finally {
      scanLoading.value = false
    }
  }

  async function deleteCredit(id: string) {
    deleting.value = true
    scanError.value = null
    try {
      await apiFetch(`/credits/${encodeURIComponent(id)}`, { method: 'DELETE' })
      if (scanResult.value?.credit_id === id) {
        scanResult.value = null
      }
      await fetchDashboard()
    } catch (e) {
      scanError.value = formatApiError(e, 'Не удалось удалить кредит')
    } finally {
      deleting.value = false
    }
  }

  return {
    dashboard,
    enrichedDashboard,
    loading,
    fetchError,
    error: fetchError,
    scanResult,
    scanLoading,
    scanError,
    deleting,
    hasCredits,
    fetchDashboard,
    scanContract,
    deleteCredit,
    emptyDashboard
  }
}
