import { useAuthStore } from '~/store/authStore'
import { readOnboardingCompleted } from '~/composables/useOnboarding'

const PUBLIC_PATHS = ['/', '/login']
const ONBOARDING_PATH = '/onboarding'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()
  if (import.meta.client && !authStore.token) {
    authStore.hydrate()
  }

  if (to.path === '/welcome') {
    return navigateTo('/')
  }

  const onboardingDone = () =>
    readOnboardingCompleted(authStore.user?.phone, authStore.user?.id)

  if (PUBLIC_PATHS.includes(to.path)) {
    if (authStore.isAuthenticated && to.path === '/login') {
      return navigateTo(onboardingDone() ? '/dashboard' : ONBOARDING_PATH)
    }
    return
  }

  if (!authStore.isAuthenticated) {
    return navigateTo('/login')
  }

  if (!import.meta.client) return

  const completed = onboardingDone()
  const allowedDuringOnboarding = [...PUBLIC_PATHS, ONBOARDING_PATH]

  if (!completed && !allowedDuringOnboarding.includes(to.path)) {
    return navigateTo(ONBOARDING_PATH)
  }

  if (completed && to.path === ONBOARDING_PATH) {
    return navigateTo('/dashboard')
  }
})
