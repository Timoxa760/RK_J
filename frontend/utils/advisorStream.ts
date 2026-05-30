import type { AiChatResponse } from '~/types/api'

export async function streamAdvisorChat(
  apiBase: string,
  token: string,
  body: { message: string; history: Array<{ role: string; content: string }> },
  onDelta: (text: string) => void
): Promise<AiChatResponse> {
  const url = `${apiBase.replace(/\/$/, '')}/api/v1/ai/chat/stream`
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify(body)
  })

  if (!res.ok || !res.body) {
    throw new Error(`stream failed: ${res.status}`)
  }

  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let donePayload: AiChatResponse | null = null

  while (true) {
    const { value, done } = await reader.read()
    if (done) break
    buffer += decoder.decode(value, { stream: true })

    let boundary = buffer.indexOf('\n\n')
    while (boundary >= 0) {
      const block = buffer.slice(0, boundary)
      buffer = buffer.slice(boundary + 2)
      parseSseBlock(block, onDelta, (payload) => {
        donePayload = payload
      })
      boundary = buffer.indexOf('\n\n')
    }
  }

  if (buffer.trim()) {
    parseSseBlock(buffer, onDelta, (payload) => {
      donePayload = payload
    })
  }

  if (!donePayload?.reply) {
    throw new Error('empty stream response')
  }
  return donePayload
}

function parseSseBlock(
  block: string,
  onDelta: (text: string) => void,
  onDone: (payload: AiChatResponse) => void
) {
  const lines = block.split('\n')
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
      if (parsed.text) onDelta(parsed.text)
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
