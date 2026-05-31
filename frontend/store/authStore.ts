import { defineStore } from 'pinia'
import type { AuthUser } from '~/types/api'
import { clearStoredAiPlanForCurrentUser } from '~/utils/aiPlanStorage'

const TOKEN_KEY = 'money_mind_token'
const REFRESH_KEY = 'money_mind_refresh'
const USER_KEY = 'money_mind_user'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as AuthUser | null,
    token: null as string | null,
    refreshToken: null as string | null
  }),

  getters: {
    isAuthenticated: (state) => Boolean(state.token)
  },

  actions: {
    hydrate() {
      if (!import.meta.client) return
      const token = localStorage.getItem(TOKEN_KEY)
      const refresh = localStorage.getItem(REFRESH_KEY)
      const userRaw = localStorage.getItem(USER_KEY)
      if (token) this.token = token
      if (refresh) this.refreshToken = refresh
      if (userRaw) {
        try {
          this.user = JSON.parse(userRaw) as AuthUser
        } catch {
          this.user = null
        }
      }
    },

    setSession(token: string, user: AuthUser, refreshToken?: string) {
      this.token = token
      this.user = user
      this.refreshToken = refreshToken ?? null
      if (import.meta.client) {
        localStorage.setItem(TOKEN_KEY, token)
        localStorage.setItem(USER_KEY, JSON.stringify(user))
        if (refreshToken) {
          localStorage.setItem(REFRESH_KEY, refreshToken)
        } else {
          localStorage.removeItem(REFRESH_KEY)
        }
      }
    },

    logout() {
      if (import.meta.client) {
        clearStoredAiPlanForCurrentUser()
      }
      this.token = null
      this.user = null
      this.refreshToken = null
      if (import.meta.client) {
        localStorage.removeItem(TOKEN_KEY)
        localStorage.removeItem(REFRESH_KEY)
        localStorage.removeItem(USER_KEY)
      }
    }
  }
})
