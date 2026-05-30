import type { AiChatResponse } from '~/types/api'

const STREAM_TIMEOUT_MS = 25_000

export async function streamAdvisorChat(
  url: string,
  token: string,
  body: { message: string; history: Array<{ role: string; content: string }> },
  onDelta: (text: string) => void
): Promise<AiChatResponse> {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), STREAM_TIMEOUT_MS)

  let res: Response
  try {
    res = await fetch(url, {
      method: 'POST',
      headers: {
        Accept: 'text/event-stream',
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(body),
      signal: controller.signal
    })
  } finally {
    clearTimeout(timeout)
  }

  if (!res.ok || !res.body) {
    throw new Error(`stream failed: ${res.status}`)
  }

  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let donePayload: AiChatResponse | null = null
  let streamed = ''

  const readTimeout = setTimeout(() => controller.abort(), STREAM_TIMEOUT_MS)

  try {
    while (true) {
      const { value, done } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true }).replace(/\r\n/g, '\n')

      let boundary = buffer.indexOf('\n\n')
      while (boundary >= 0) {
        const block = buffer.slice(0, boundary)
        buffer = buffer.slice(boundary + 2)
        parseSseBlock(block, onDelta, (payload) => {
          donePayload = payload
        }, (delta) => {
          streamed += delta
        })
        boundary = buffer.indexOf('\n\n')
      }
    }

    buffer = buffer.replace(/\r\n/g, '\n')
    if (buffer.trim()) {
      parseSseBlock(buffer, onDelta, (payload) => {
        donePayload = payload
      }, (delta) => {
        streamed += delta
      })
    }
  } finally {
    clearTimeout(readTimeout)
    reader.cancel().catch(() => {})
  }

  if (donePayload?.reply) {
    return donePayload
  }
  if (streamed.trim()) {
    return { reply: streamed.trim(), source: 'heuristic' }
  }
  throw new Error('empty stream response')
}

function parseSseBlock(
  block: string,
  onDelta: (text: string) => void,
  onDone: (payload: AiChatResponse) => void,
  onStreamed?: (delta: string) => void
) {
  const normalized = block.trim()
  if (!normalized) return

  const lines = normalized.split('\n')
  let event = 'message'
  let data = ''
  for (const line of lines) {
    if (line.startsWith('event:')) {
      event = line.slice(6).trim()
    } else if (line.startsWith('data:')) {
      data = line.slice(5).trim()
    }
  }
  if (!data) return

  if (event === 'delta') {
    try {
      const parsed = JSON.parse(data) as { text?: string }
      if (parsed.text) {
        onStreamed?.(parsed.text)
        onDelta(parsed.text)
      }
    } catch {
      /* ignore */
    }
    return
  }

  if (event === 'done') {
    try {
      onDone(JSON.parse(data) as AiChatResponse)
    } catch {
      /* ignore */
    }
  }
}
