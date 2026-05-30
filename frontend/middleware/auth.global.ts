import { useAuthStore } from '~/store/authStore'
import { isOnboardingComplete, needsOnboarding } from '~/composables/useOnboarding'

const PUBLIC_PATHS = ['/', '/login']
const ONBOARDING_PATH = '/onboarding'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()
  const onServer = Boolean(useRequestEvent())

  if (!onServer) {
    authStore.hydrate()
  }

  if (to.path === '/welcome') {
    return navigateTo('/')
  }

  const onboardingDone = () =>
    isOnboardingComplete(authStore.user?.phone, authStore.user?.id)

  if (PUBLIC_PATHS.includes(to.path)) {
    if (authStore.isAuthenticated && to.path === '/login') {
      return navigateTo(onboardingDone() ? '/dashboard' : ONBOARDING_PATH)
    }
    return
  }

  // Токен в localStorage — на SSR сессии нет, проверяем только на клиенте после hydrate
  if (!authStore.isAuthenticated) {
    if (onServer) return
    return navigateTo('/login')
  }

  if (onServer) return

  const completed = onboardingDone()
  const allowedDuringOnboarding = [...PUBLIC_PATHS, ONBOARDING_PATH, '/profile', '/receipts']

  if (needsOnboarding(authStore.user?.phone, authStore.user?.id) && !allowedDuringOnboarding.includes(to.path)) {
    return navigateTo(ONBOARDING_PATH)
  }

  if (completed && to.path === ONBOARDING_PATH) {
    return navigateTo('/dashboard')
  }
})
