const CHAT_HASH = '#advisor-chat'
const ADVISOR_PATH = '/advisor'

export function useOpenAdvisorChat() {
  const route = useRoute()
  const { requestOpen: requestMobileSidebar } = usePendingMobileSidebar()

  async function openAdvisorChat(prompt: string) {
    const encoded = prompt ? encodeURIComponent(prompt) : undefined
    const query = encoded ? { ask: encoded } : {}

    if (route.path === ADVISOR_PATH) {
      if (encoded) {
        await navigateTo({ path: ADVISOR_PATH, query, hash: CHAT_HASH })
      }
      scrollToAdvisorChat()
      return
    }

    await navigateTo({ path: ADVISOR_PATH, query, hash: CHAT_HASH })
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
