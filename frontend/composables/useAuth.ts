import type { AuthUser, LoginResponse } from '~/types/api'
import { useAuthStore } from '~/store/authStore'

const MOCK_CODE = '0000'

export function useAuth() {
  const { apiFetch, demoMode } = useApi()
  const authStore = useAuthStore()

  async function register(phone: string) {
    if (demoMode.value) return
    await apiFetch<{ message: string; expires_in?: number }>('/auth/register', {
      method: 'POST',
      body: { phone }
    })
  }

  async function login(phone: string, code: string) {
    if (demoMode.value && code === MOCK_CODE) {
      const phoneKey = phone.replace(/\D/g, '') || String(Date.now())
      const user: AuthUser = {
        id: `demo-${phoneKey}`,
        phone,
        role: 'user',
        name: 'Пользователь'
      }
      authStore.setSession(`mock-jwt-demo-${phoneKey}`, user)
      return { token: `mock-jwt-demo-${phoneKey}`, user }
    }

    const res = await apiFetch<LoginResponse>('/auth/login', {
      method: 'POST',
      body: { phone, code }
    })

    const user: AuthUser = {
      id: res.user.id,
      phone: res.user.phone,
      role: res.user.role,
      name: res.user.name ?? 'Пользователь'
    }
    authStore.setSession(res.access_token, user, res.refresh_token)
    return { token: res.access_token, user }
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
