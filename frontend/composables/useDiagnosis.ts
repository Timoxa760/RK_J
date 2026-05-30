import type { AiDiagnosisResponse } from '~/types/api'

export function useDiagnosis() {
  const { apiFetch } = useApi()

  const diagnosis = ref<AiDiagnosisResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchDiagnosis() {
    loading.value = true
    error.value = null

    try {
      diagnosis.value = await apiFetch<AiDiagnosisResponse>('/ai/diagnosis')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить картину'
      diagnosis.value = null
    } finally {
      loading.value = false
    }
  }

  return {
    diagnosis,
    loading,
    error,
    fetchDiagnosis
  }
}
