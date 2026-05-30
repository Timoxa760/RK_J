import type { AiChatMessage } from '~/types/api'
import { currentUserStorageKey } from '~/utils/userStorage'
import {
  buildAdvisorGreeting,
  buildAdvisorReply,
  type AdvisorContext
} from '~/utils/advisorChat'

const CHAT_PREFIX = 'potok:advisor-chat'

export interface ChatTurn {
  id: string
  role: 'user' | 'assistant'
  content: string
  createdAt: number
}

function readStoredChat(): ChatTurn[] {
  if (!import.meta.client) return []
  try {
    const raw = localStorage.getItem(currentUserStorageKey(CHAT_PREFIX))
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

function readSessionChat(): ChatTurn[] {
  if (!import.meta.client) return []
  try {
    const raw = sessionStorage.getItem(currentUserStorageKey(CHAT_PREFIX))
    if (!raw) return readStoredChat()
    const parsed = JSON.parse(raw) as ChatTurn[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function newId() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

export function useAdvisorChat(getContext: () => AdvisorContext) {
  const { apiFetch } = useApi()

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
        createdAt: Date.now()
      }
    ]
    persist()
  }

  function initChat() {
    if (initialized.value) return
    messages.value = readSessionChat()
    seedGreeting()
    initialized.value = true
  }

  function resetChat() {
    messages.value = []
    seedGreeting()
    persist()
  }

  async function sendMessage(text: string) {
    const trimmed = text.trim()
    if (!trimmed || typing.value) return

    error.value = null
    messages.value.push({
      id: newId(),
      role: 'user',
      content: trimmed,
      createdAt: Date.now()
    })
    persist()

    typing.value = true

    try {
      let reply: string
      const ctx = getContext()
      const history: AiChatMessage[] = messages.value
        .filter((m) => m.role === 'user' || m.role === 'assistant')
        .slice(-10)
        .map((m) => ({ role: m.role, content: m.content }))

      try {
        const res = await apiFetch<{ reply: string }>('/ai/chat', {
          method: 'POST',
          body: { message: trimmed, history }
        })
        reply = res.reply?.trim() || buildAdvisorReply(trimmed, ctx)
      } catch {
        reply = buildAdvisorReply(trimmed, ctx)
      }

      messages.value.push({
        id: newId(),
        role: 'assistant',
        content: reply,
        createdAt: Date.now()
      })
      persist()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось получить ответ'
    } finally {
      typing.value = false
    }
  }

  async function handleAskQuery(ask: string) {
    if (!ask) return
    initChat()
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
