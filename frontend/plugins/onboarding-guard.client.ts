import { needsOnboarding } from '~/composables/useOnboarding'
import { useAuthStore } from '~/store/authStore'

const PUBLIC_PATHS = ['/', '/login']
const ONBOARDING_PATH = '/onboarding'

/** Дублирует проверку middleware после гидрации — SSR не видит localStorage. */
export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook('app:mounted', () => {
    const authStore = useAuthStore()
    authStore.hydrate()
    if (!authStore.isAuthenticated) return

    const route = useRoute()
    if (PUBLIC_PATHS.includes(route.path) || route.path === ONBOARDING_PATH) return

    if (needsOnboarding(authStore.user?.phone, authStore.user?.id)) {
      navigateTo(ONBOARDING_PATH)
    }
  })
})
