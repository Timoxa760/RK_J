/** true после mount — не показываем UI, зависящий от localStorage, до гидрации. */
export function useShellReady() {
  const ready = ref(false)
  onMounted(() => {
    ready.value = true
  })
  return ready
}
