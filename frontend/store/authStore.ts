import { defineStore } from 'pinia'

export interface AuthUser {
  id: string
  phone: string
  name?: string
}

const TOKEN_KEY = 'money_mind_token'
const USER_KEY = 'money_mind_user'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as AuthUser | null,
    token: null as string | null
  }),

  getters: {
    isAuthenticated: (state) => Boolean(state.token)
  },

  actions: {
    hydrate() {
      if (!import.meta.client) return
      const token = localStorage.getItem(TOKEN_KEY)
      const userRaw = localStorage.getItem(USER_KEY)
      if (token) this.token = token
      if (userRaw) {
        try {
          this.user = JSON.parse(userRaw) as AuthUser
        } catch {
          this.user = null
        }
      }
    },

    setSession(token: string, user: AuthUser) {
      this.token = token
      this.user = user
      if (import.meta.client) {
        localStorage.setItem(TOKEN_KEY, token)
        localStorage.setItem(USER_KEY, JSON.stringify(user))
      }
    },

    logout() {
      this.token = null
      this.user = null
      if (import.meta.client) {
        localStorage.removeItem(TOKEN_KEY)
        localStorage.removeItem(USER_KEY)
      }
    }
  }
})
