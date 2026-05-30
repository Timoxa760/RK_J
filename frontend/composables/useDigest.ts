import type { DigestResponse } from '~/types/api'

export function useDigest() {
  const { apiFetch } = useApi()

  const digest = ref<DigestResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function loadDigest() {
    loading.value = true
    error.value = null
    try {
      digest.value = await apiFetch<DigestResponse>('/digest/latest')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить сводку'
      digest.value = null
    } finally {
      loading.value = false
    }
  }

  return { digest, loading, error, loadDigest }
}
