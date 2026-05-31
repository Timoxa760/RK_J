/** Снимает старый SW после удаления модулей (например useFnsMco) — один раз за версию сборки. */
export default defineNuxtPlugin(async () => {
  if (!('serviceWorker' in navigator)) return

  const marker = 'potok-sw-cleared-v2'
  if (localStorage.getItem(marker)) return

  try {
    const registrations = await navigator.serviceWorker.getRegistrations()
    await Promise.all(registrations.map((registration) => registration.unregister()))
    if ('caches' in window) {
      const keys = await caches.keys()
      await Promise.all(keys.map((key) => caches.delete(key)))
    }
    localStorage.setItem(marker, '1')
  } catch {
    // ignore — offline or restricted storage
  }
})
