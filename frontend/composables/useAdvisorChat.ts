import type { AdvisorChatAction, AiChatHistoryResponse, AiChatResponse } from '~/types/api'
import { useAuthStore } from '~/store/authStore'
import { currentUserStorageKey } from '~/utils/userStorage'
import {
  buildAdvisorGreeting,
  buildAdvisorReply,
  historyToTurns,
  toApiHistory,
  type AdvisorContext
} from '~/utils/advisorChat'
import { streamAdvisorChat } from '~/utils/advisorStream'

const CHAT_PREFIX = 'potok:advisor-chat'

export interface ChatTurn {
  id: string
  role: 'user' | 'assistant'
  content: string
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

export function useAdvisorChat(getContext: () => AdvisorContext) {
  const { apiFetch, apiV1Base } = useApi()
  const authStore = useAuthStore()

  const messages = useState<ChatTurn[]>('advisor-chat-messages', () => [])
  const typing = useState('advisor-chat-typing', () => false)
  const error = useState<string | null>('advisor-chat-error', () => null)
  const initialized = useState('advisor-chat-initialized', () => false)

  function persist() {
    writeStoredChat(messages.value)
  }

  function seedGreeting() {
    if (messages.value.length) return
    messages.value = [
      {
        id: newId(),
        role: 'assistant',
        content: buildAdvisorGreeting(getContext()),
        createdAt: Date.now(),
        source: 'local'
      }
    ]
    persist()
  }

  async function loadServerHistory() {
    try {
      const res = await apiFetch<AiChatHistoryResponse>('/ai/chat/history?limit=50')
      if (res.messages?.length) {
        messages.value = historyToTurns(res.messages)
        persist()
        return
      }
    } catch {
      /* fallback to local */
    }
    messages.value = readStoredChat()
  }

  async function initChat() {
    if (initialized.value) return
    await loadServerHistory()
    seedGreeting()
    initialized.value = true
  }

  async function resetChat() {
    try {
      await apiFetch('/ai/chat/history', { method: 'DELETE' })
    } catch {
      /* local reset still works */
    }
    messages.value = []
    seedGreeting()
    persist()
  }

  async function requestReply(trimmed: string, assistantId: string) {
    const ctx = getContext()
    const history = toApiHistory(messages.value)

    const applyResponse = (reply: string, meta?: Partial<ChatTurn>) => {
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
        let streamed = ''
        const res = await streamAdvisorChat(
          `${apiV1Base.value}/ai/chat/stream`,
          token,
          { message: trimmed, history },
          (delta) => {
            streamed += delta
            const idx = messages.value.findIndex((m) => m.id === assistantId)
            if (idx >= 0) {
              messages.value[idx] = {
                ...messages.value[idx]!,
                content: streamed,
                streaming: true
              }
            }
          }
        )
        applyResponse(res.reply, {
          actions: res.actions,
          source: res.source ?? 'gemini',
          id: res.id
        })
        return
      } catch {
        /* fallback below */
      }
    }

    try {
      const res = await apiFetch<AiChatResponse>('/ai/chat', {
        method: 'POST',
        body: { message: trimmed, history }
      })
      applyResponse(res.reply?.trim() || buildAdvisorReply(trimmed, ctx), {
        actions: res.actions,
        source: res.source ?? 'gemini',
        id: res.id
      })
    } catch {
      applyResponse(buildAdvisorReply(trimmed, ctx), { source: 'local' })
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
