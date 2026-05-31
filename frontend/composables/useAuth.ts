import type { AuthUser, LoginResponse } from '~/types/api'
import { useAuthStore } from '~/store/authStore'
import { useAiPlanStore } from '~/store/aiPlanStore'

export function useAuth() {
  const { apiFetch } = useApi()
  const authStore = useAuthStore()

  async function register(phone: string, password: string) {
    await apiFetch<{ message: string }>('/auth/register', {
      method: 'POST',
      body: { phone, password }
    })
  }

  async function login(phone: string, password: string) {
    const res = await apiFetch<LoginResponse>('/auth/login', {
      method: 'POST',
      body: { phone, password }
    })

    const user: AuthUser = {
      id: res.user.id,
      phone: res.user.phone,
      role: res.user.role
    }
    authStore.setSession(res.access_token, user, res.refresh_token)
    return { token: res.access_token, user }
  }

  async function requestPasswordReset(phone: string) {
    await apiFetch<{ message: string; expires_in?: number }>('/auth/password/forgot', {
      method: 'POST',
      body: { phone }
    })
  }

  async function resetPassword(phone: string, code: string, newPassword: string) {
    await apiFetch<{ message: string }>('/auth/password/reset', {
      method: 'POST',
      body: { phone, code, new_password: newPassword }
    })
  }

  function logout() {
    authStore.logout()
    useAiPlanStore().clearCache()
  }

  return {
    authStore,
    register,
    login,
    requestPasswordReset,
    resetPassword,
    logout
  }
}
