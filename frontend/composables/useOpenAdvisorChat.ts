const CHAT_HASH = '#advisor-chat'

export function useOpenAdvisorChat() {
  const route = useRoute()
  const { requestOpen: requestMobileSidebar } = usePendingMobileSidebar()

  async function openAdvisorChat(prompt: string) {
    const trimmed = prompt.trim()
    const path =
      route.path === '/login' || route.path === '/' || route.path === '/onboarding'
        ? '/dashboard'
        : route.path

    if (!trimmed) {
      await navigateTo({ path, hash: CHAT_HASH })
      requestMobileSidebar()
      scrollToAdvisorChat()
      return
    }

    await navigateTo({
      path,
      query: { ask: encodeURIComponent(trimmed) },
      hash: CHAT_HASH
    })
    requestMobileSidebar()
    scrollToAdvisorChat()
  }

  return { openAdvisorChat }
}

export function scrollToAdvisorChat() {
  if (!import.meta.client) return
  nextTick(() => {
    document.getElementById('advisor-chat')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}
