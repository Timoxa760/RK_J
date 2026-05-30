import type { AiDiagnosisResponse } from '~/types/api'
import { useAuthStore } from '~/store/authStore'

export function useFnsMco() {
  const { apiFetch } = useApi()
  const { user } = useAuthStore()
  const toast = useToast()

  const loading = ref(false)
  const authStep = ref<'idle' | 'code_sent' | 'linked'>('idle')
  const error = ref<string | null>(null)

  const phone = computed(() => user?.phone ?? '')

  async function startAuth() {
    if (!phone.value) {
      error.value = 'Сначала войдите в аккаунт'
      return
    }
    loading.value = true
    error.value = null
    try {
      await apiFetch<{ success: boolean; message?: string }>('/fns/mco/auth', {
        method: 'POST',
        body: { phone: phone.value }
      })
      authStep.value = 'code_sent'
      toast.show('Код отправлен — в demo используйте 0000', 'success')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось начать подключение'
      toast.show(error.value, 'error')
    } finally {
      loading.value = false
    }
  }

  async function verifyAuth(code: string) {
    if (!phone.value) return
    loading.value = true
    error.value = null
    try {
      await apiFetch<{ success: boolean }>('/fns/mco/auth/verify', {
        method: 'POST',
        body: { phone: phone.value, code: code.trim() }
      })
      authStep.value = 'linked'
      toast.show('«Мои чеки» подключены', 'success')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Неверный код'
      toast.show(error.value, 'error')
    } finally {
      loading.value = false
    }
  }

  async function syncReceipts() {
    if (!phone.value) return 0
    loading.value = true
    error.value = null
    try {
      const res = await apiFetch<{ success?: boolean; count?: number }>('/fns/mco/sync', {
        method: 'POST',
        body: { phone: phone.value }
      })
      const count = res.count ?? 0
      toast.show(
        count > 0 ? `Загружено чеков: ${count}` : 'Новых чеков пока нет',
        count > 0 ? 'success' : 'default'
      )
      return count
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось синхронизировать'
      toast.show(error.value, 'error')
      return 0
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    authStep,
    error,
    phone,
    startAuth,
    verifyAuth,
    syncReceipts
  }
}

export function useOnboardingDiagnosis() {
  const { apiFetch } = useApi()
  const diagnosis = ref<AiDiagnosisResponse | null>(null)
  const loading = ref(false)

  async function loadDiagnosis() {
    loading.value = true
    try {
      diagnosis.value = await apiFetch<AiDiagnosisResponse>('/ai/diagnosis')
    } catch {
      diagnosis.value = null
    } finally {
      loading.value = false
    }
  }

  return { diagnosis, loading, loadDiagnosis }
}
