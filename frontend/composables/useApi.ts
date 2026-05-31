import { useAuthStore } from '~/store/authStore'
import { useAiPlanStore } from '~/store/aiPlanStore'

type FetchOptions = Parameters<typeof $fetch>[1]

export function useApi() {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const apiV1Base = computed(() => {
    // В dev браузер ходит через nitro devProxy (/api → gateway), без CORS.
    if (import.meta.dev && import.meta.client) {
      return '/api/v1'
    }
    const base = String(config.public.apiBase).replace(/\/$/, '')
    return `${base}/api/v1`
  })

  const demoMode = computed(() => config.public.demoMode)

  function authHeaders(): Record<string, string> {
    return authStore.token ? { Authorization: `Bearer ${authStore.token}` } : {}
  }

  async function apiFetch<T>(path: string, options: FetchOptions = {}): Promise<T> {
    try {
      return await $fetch<T>(path, {
        ...options,
        baseURL: apiV1Base.value,
        headers: {
          ...authHeaders(),
          ...(options.headers as Record<string, string> | undefined)
        }
      })
    } catch (error: unknown) {
      const status = (error as { statusCode?: number })?.statusCode
      if (import.meta.client && status === 401) {
        const route = useRoute()
        // На онбординге и login не редиректим — ошибку покажет форма.
        if (route.path !== '/onboarding' && route.path !== '/login') {
          authStore.logout()
          useAiPlanStore().clearCache()
          await navigateTo('/login')
        }
      }
      throw error
    }
  }

  /** Demo: mock. Production: только API, ошибки пробрасываются вызывающему коду. */
  async function apiFetchWithDemo<T>(
    path: string,
    mock: T,
    options: FetchOptions = {}
  ): Promise<T> {
    if (demoMode.value) return mock
    return apiFetch<T>(path, options)
  }

  return {
    apiV1Base,
    demoMode,
    apiFetch,
    apiFetchWithDemo
  }
}
