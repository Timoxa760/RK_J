import type { AiDiagnosisResponse } from '~/types/api'
import { mockDiagnosis } from '~/store/mocks/diagnosis'

export function useDiagnosis() {
  const { apiFetchWithDemo, demoMode } = useApi()

  const diagnosis = ref<AiDiagnosisResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchDiagnosis() {
    loading.value = true
    error.value = null

    try {
      diagnosis.value = await apiFetchWithDemo<AiDiagnosisResponse>(
        '/ai/diagnosis',
        mockDiagnosis
      )
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить картину'
      if (demoMode.value) {
        diagnosis.value = mockDiagnosis
      }
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
