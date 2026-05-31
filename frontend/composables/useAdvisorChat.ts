import type { AdvisorChatAction, AiChatHistoryResponse, AiChatResponse, AdvisorReplyBlock } from '~/types/api'
import { useAuthStore } from '~/store/authStore'
import { currentUserStorageKey } from '~/utils/userStorage'
import { historyToTurns, toApiHistory, type AdvisorContext } from '~/utils/advisorChat'
import { streamAdvisorChat } from '~/utils/advisorStream'
import { formatApiError } from '~/utils/apiError'

const CHAT_PREFIX = 'potok:advisor-chat'

export interface ChatTurn {
  id: string
  role: 'user' | 'assistant'
  content: string
  title?: string
  blocks?: AdvisorReplyBlock[]
  createdAt: number
  actions?: AdvisorChatAction[]
  source?: 'gemini' | 'heuristic' | 'local'
  streaming?: boolean
}

function readStoredChat(): ChatTurn[] {
  if (!import.meta.client) return []
  try {
    const raw = sessionStorage.getItem(currentUserStorageKey(CHAT_PREFIX))
    if (!raw) return []
    const parsed = JSON.parse(raw) as ChatTurn[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function writeStoredChat(messages: ChatTurn[]) {
  if (!import.meta.client) return
  sessionStorage.setItem(currentUserStorageKey(CHAT_PREFIX), JSON.stringify(messages))
}

function newId() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function withoutLocalAssistant(turns: ChatTurn[]): ChatTurn[] {
  return turns.filter((turn) => !(turn.role === 'assistant' && turn.source === 'local'))
}

function hasAssistantPayload(res: AiChatResponse): boolean {
  return Boolean(res.reply?.trim() || res.blocks?.length || res.title?.trim())
}

export function useAdvisorChat(_getContext: () => AdvisorContext) {
  const { apiFetch, apiV1Base } = useApi()
  const authStore = useAuthStore()

  const messages = useState<ChatTurn[]>('advisor-chat-messages', () => [])
  const typing = useState('advisor-chat-typing', () => false)
  const error = useState<string | null>('advisor-chat-error', () => null)
  const historyLoaded = useState('advisor-chat-history-loaded', () => false)

  function persist() {
    writeStoredChat(messages.value)
  }

  async function loadServerHistory() {
    try {
      const res = await apiFetch<AiChatHistoryResponse>('/ai/chat/history?limit=50')
      if (res.messages?.length) {
        messages.value = withoutLocalAssistant(historyToTurns(res.messages))
        persist()
        historyLoaded.value = true
        return true
      }
    } catch {
      /* fallback to local */
    }
    const local = withoutLocalAssistant(readStoredChat())
    if (local.length) {
      messages.value = local
    }
    historyLoaded.value = true
    return false
  }

  async function initChat(opts?: { reload?: boolean }) {
    if (historyLoaded.value && !opts?.reload) {
      return
    }
    await loadServerHistory()
  }

  async function resetChat() {
    try {
      await apiFetch('/ai/chat/history', { method: 'DELETE' })
    } catch {
      /* local reset still works */
    }
    messages.value = []
    historyLoaded.value = true
    error.value = null
    persist()
  }

  async function requestReply(trimmed: string, assistantId: string) {
    const history = toApiHistory(messages.value)

    const applyResponse = (
      reply: string,
      meta?: Partial<ChatTurn> & { blocks?: AdvisorReplyBlock[]; title?: string }
    ) => {
      const idx = messages.value.findIndex((m) => m.id === assistantId)
      const turn: ChatTurn = {
        id: assistantId,
        role: 'assistant',
        content: reply,
        createdAt: Date.now(),
        streaming: false,
        ...meta
      }
      if (idx >= 0) messages.value[idx] = turn
      else messages.value.push(turn)
      persist()
    }

    const failReply = (message: string) => {
      messages.value = messages.value.filter((m) => m.id !== assistantId)
      error.value = message
      persist()
    }

    messages.value.push({
      id: assistantId,
      role: 'assistant',
      content: '',
      createdAt: Date.now(),
      streaming: true
    })

    const token = authStore.token

    if (import.meta.client && token) {
      try {
        const res = await streamAdvisorChat(
          `${apiV1Base.value}/ai/chat/stream`,
          token,
          { message: trimmed, history },
          () => {
            // JSON-стрим не показываем — только индикатор набора до event done
          }
        )
        if (!hasAssistantPayload(res)) {
          failReply('Советник не вернул ответ. Попробуйте переформулировать вопрос.')
          return
        }
        applyResponse(res.reply?.trim() || res.title || 'Ответ готов', {
          title: res.title,
          blocks: res.blocks,
          actions: res.actions,
          source: res.source ?? 'gemini',
          id: res.id
        })
        return
      } catch {
        /* fallback to POST below */
      }
    }

    try {
      const res = await apiFetch<AiChatResponse>('/ai/chat', {
        method: 'POST',
        body: { message: trimmed, history }
      })
      if (!hasAssistantPayload(res)) {
        failReply('Советник не вернул ответ. Попробуйте переформулировать вопрос.')
        return
      }
      applyResponse(res.reply?.trim() || res.title || 'Ответ готов', {
        title: res.title,
        blocks: res.blocks,
        actions: res.actions,
        source: res.source ?? 'gemini',
        id: res.id
      })
    } catch (e) {
      failReply(formatApiError(e, 'Не удалось получить ответ'))
    }
  }

  async function sendMessage(text: string) {
    const trimmed = text.trim()
    if (!trimmed || typing.value) return

    await initChat()

    error.value = null
    messages.value.push({
      id: newId(),
      role: 'user',
      content: trimmed,
      createdAt: Date.now()
    })
    persist()

    typing.value = true
    const assistantId = newId()

    try {
      await requestReply(trimmed, assistantId)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось получить ответ'
      messages.value = messages.value.filter((m) => m.id !== assistantId)
    } finally {
      typing.value = false
    }
  }

  async function handleAskQuery(ask: string) {
    if (!ask) return
    await initChat()
    const lastUser = [...messages.value].reverse().find((m) => m.role === 'user')
    if (lastUser?.content === ask) return
    await sendMessage(ask)
  }

  return {
    messages,
    typing,
    error,
    initChat,
    resetChat,
    sendMessage,
    handleAskQuery
  }
}
