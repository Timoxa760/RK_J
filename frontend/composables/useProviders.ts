import type { ProviderConnectResponse, ProviderId } from '~/types/api'
import { PROVIDERS } from '~/types/api'

export function useProviders() {
  const { apiFetch, demoMode } = useApi()

  const connected = ref<Partial<Record<ProviderId, string>>>({})
  const loading = ref<string | null>(null)
  const error = ref<string | null>(null)
  const success = ref<string | null>(null)

  const providers = PROVIDERS

  async function connect(
    provider: ProviderId,
    credentials: { phone: string; password: string }
  ) {
    loading.value = provider
    error.value = null
    success.value = null
    try {
      if (demoMode.value) {
        await new Promise((r) => setTimeout(r, 600))
        connected.value[provider] = 'active'
        success.value = `${providers.find((p) => p.id === provider)?.label} подключён`
        return
      }
      const res = await apiFetch<ProviderConnectResponse>(
        `/providers/connect?provider=${provider}`,
        {
          method: 'POST',
          body: { credentials }
        }
      )
      connected.value[provider] = res.status
      success.value = res.message
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось подключить магазин'
    } finally {
      loading.value = null
    }
  }

  async function sync(provider: ProviderId) {
    loading.value = provider
    error.value = null
    try {
      if (demoMode.value) {
        await new Promise((r) => setTimeout(r, 400))
        success.value = 'Синхронизация запущена'
        return
      }
      await apiFetch(`/providers/${provider}/sync`, { method: 'POST' })
      success.value = 'Синхронизация запущена'
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка синхронизации'
    } finally {
      loading.value = null
    }
  }

  return {
    providers,
    connected,
    loading,
    error,
    success,
    connect,
    sync
  }
}
