import type {
  MortgageAnalyzeRequest,
  MortgageBreakdownResponse
} from '~/types/api'

export function useMortgage() {
  const { apiFetch } = useApi()

  const breakdown = ref<MortgageBreakdownResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function analyze(input: MortgageAnalyzeRequest) {
    loading.value = true
    error.value = null
    try {
      breakdown.value = await apiFetch<MortgageBreakdownResponse>(
        '/credits/mortgage/analyze',
        {
          method: 'POST',
          body: input
        }
      )
      return breakdown.value
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось рассчитать разбор'
      breakdown.value = null
      return null
    } finally {
      loading.value = false
    }
  }

  function reset() {
    breakdown.value = null
    error.value = null
  }

  return {
    breakdown,
    loading,
    error,
    analyze,
    reset
  }
}
