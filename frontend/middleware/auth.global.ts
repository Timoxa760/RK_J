import { useAuthStore } from '~/store/authStore'

const PUBLIC_PATHS = ['/', '/login']

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()
  if (import.meta.client && !authStore.token) {
    authStore.hydrate()
  }

  if (to.path === '/welcome') {
    return navigateTo('/')
  }

  if (PUBLIC_PATHS.includes(to.path)) {
    if (authStore.isAuthenticated && to.path === '/login') {
      return navigateTo('/dashboard')
    }
    return
  }

  if (!authStore.isAuthenticated) {
    return navigateTo('/login')
  }
})
