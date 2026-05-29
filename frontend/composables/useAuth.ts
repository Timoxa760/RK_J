import { useAuthStore, type AuthUser } from '~/store/authStore'

export function useAuth() {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const MOCK_CODE = '0000'
  const MOCK_TOKEN = 'mock-jwt-demo-token'

  async function register(phone: string) {
    try {
      await $fetch(`${config.public.apiBase}/auth/register`, {
        method: 'POST',
        body: { phone }
      })
    } catch {
      // mock: registration always succeeds offline
    }
  }

  async function login(phone: string, code: string) {
    if (code === MOCK_CODE) {
      const user: AuthUser = {
        id: 'mock-user-1',
        phone,
        name: 'Пользователь'
      }
      authStore.setSession(MOCK_TOKEN, user)
      return { token: MOCK_TOKEN, user }
    }

    try {
      const res = await $fetch<{ token: string; user: AuthUser }>(
        `${config.public.apiBase}/auth/login`,
        {
          method: 'POST',
          body: { phone, code }
        }
      )
      authStore.setSession(res.token, res.user)
      return res
    } catch {
      throw new Error('Неверный код или сервер недоступен.')
    }
  }

  function logout() {
    authStore.logout()
  }

  return {
    authStore,
    register,
    login,
    logout
  }
}
