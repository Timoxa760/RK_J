import { decodeAskQuery } from '~/utils/advisorChat'

const advisorShellKey = Symbol('advisor-shell')

export type AdvisorShellApi = ReturnType<typeof createAdvisorShell>

function createAdvisorShell() {
  const route = useRoute()
  const router = useRouter()
  const { requestOpen: requestMobileSidebar } = usePendingMobileSidebar()
  const { advisorContext, refreshAdvisorContext } = useAdvisorContext()
  const {
    messages,
    typing,
    error: chatError,
    initChat,
    sendMessage,
    handleAskQuery
  } = useAdvisorChat(() => advisorContext.value)

  async function processAskFromQuery() {
    const ask = decodeAskQuery(route.query.ask)
    if (!ask) return
    requestMobileSidebar()
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
    await refreshAdvisorContext({ silent: true })
    initChat()
    await processAskFromQuery()
  }

  watch(
    () => route.query.ask,
    async (ask) => {
      if (!ask) return
      await processAskFromQuery()
    }
  )

  return {
    messages,
    typing,
    chatError,
    sendMessage,
    refreshAdvisorContext,
    bootstrap
  }
}

/** Только внутри SidebarProvider (AppShellAdvisorHost). */
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
    messages: ref([]),
    typing: ref(false),
    chatError: ref<string | null>(null),
    sendMessage: noop,
    refreshAdvisorContext: noop,
    bootstrap: noop
  }
}
