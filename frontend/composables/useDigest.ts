import type { DigestResponse } from '~/types/api'
import { mockDigest } from '~/store/mocks'

export function useDigest() {
  const { apiFetchWithDemo, demoMode } = useApi()

  const digest = ref<DigestResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function loadDigest() {
    loading.value = true
    error.value = null
    try {
      digest.value = await apiFetchWithDemo<DigestResponse>('/digest/latest', mockDigest)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить сводку'
      if (demoMode.value) {
        digest.value = mockDigest
      }
    } finally {
      loading.value = false
    }
  }

  return { digest, loading, error, loadDigest }
}
