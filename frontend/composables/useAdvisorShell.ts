import type { AdvisorChatAction } from '~/types/api'
import { decodeAskQuery, type AdvisorContext } from '~/utils/advisorChat'
import { scrollToAdvisorChat } from '~/composables/useOpenAdvisorChat'

const advisorShellKey = Symbol('advisor-shell')

export type AdvisorShellApi = ReturnType<typeof createAdvisorShell>

function createAdvisorShell() {
  const route = useRoute()
  const router = useRouter()
  const { advisorContext, refreshAdvisorContext } = useAdvisorContext()
  const {
    messages,
    typing,
    error: chatError,
    initChat,
    resetChat,
    sendMessage,
    handleAskQuery
  } = useAdvisorChat(() => advisorContext.value)

  async function runAction(action: AdvisorChatAction) {
    switch (action.type) {
      case 'navigate':
        await navigateTo({ path: action.path ?? '/dashboard', hash: action.hash })
        break
      case 'open_add_expense':
        useAddExpenseSheet().show()
        break
      case 'open_profile': {
        const query = action.profileField ? { section: action.profileField } : undefined
        await navigateTo({ path: '/profile', query })
        break
      }
      case 'ask_followup':
        if (action.ask) await sendMessage(action.ask)
        break
    }
  }

  async function processAskFromQuery() {
    const ask = decodeAskQuery(route.query.ask)
    if (!ask) return
    if (route.path !== '/advisor') {
      await navigateTo({ path: '/advisor', query: route.query, hash: '#advisor-chat' })
      return
    }
    await handleAskQuery(ask)
    scrollToAdvisorChat()
    const { ask: _removed, ...rest } = route.query
    await router.replace({ path: route.path, query: rest, hash: '#advisor-chat' })
  }

  const bootstrapped = useState('advisor-shell-boot', () => false)

  async function bootstrap() {
    if (bootstrapped.value) {
      await processAskFromQuery()
      return
    }
    bootstrapped.value = true
    typing.value = false
    await refreshAdvisorContext({ silent: true })
    await initChat()
    await processAskFromQuery()
  }

  watch(
    () => route.query.ask,
    async (ask) => {
      if (!ask) return
      await processAskFromQuery()
    }
  )

  const { addedVersion } = useAddExpenseSheet()
  watch(addedVersion, () => {
    refreshAdvisorContext({ silent: true })
  })

  return {
    advisorContext,
    messages,
    typing,
    chatError,
    sendMessage,
    resetChat,
    runAction,
    refreshAdvisorContext,
    bootstrap
  }
}

/** Только внутри AppShellAdvisorScope (предок сайдбара и страниц). */
export function useAdvisorShellProvider() {
  const api = createAdvisorShell()
  provide(advisorShellKey, api)
  return api
}

/** API советника; без провайдера — заглушка (не вызывать useSidebar). */
export function useAdvisorShell() {
  const injected = inject<AdvisorShellApi | null>(advisorShellKey, null)
  if (injected) return injected

  const noop = async () => {}
  return {
    advisorContext: computed(() => null),
    messages: ref([]),
    typing: ref(false),
    chatError: ref<string | null>(null),
    sendMessage: noop,
    resetChat: noop,
    runAction: noop,
    refreshAdvisorContext: noop,
    bootstrap: noop
  }
}
