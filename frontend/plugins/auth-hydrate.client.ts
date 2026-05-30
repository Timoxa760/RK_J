import { useAuthStore } from '~/store/authStore'

/** Сессия из localStorage — после инициализации Pinia, до mount. */
export default defineNuxtPlugin((nuxtApp) => {
  useAuthStore().hydrate()

  nuxtApp.hook('app:mounted', () => {
    useAuthStore().hydrate()
  })
})
