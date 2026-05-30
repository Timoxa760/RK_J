export interface ApiErrorPayload {
  error?: string
  code?: string
  message?: string
}

const ERROR_MESSAGES: Record<string, string> = {
  amount_required: 'Не услышали сумму. Скажите, например: «колбаса 300 рублей».',
  user_id_required: 'Сначала войдите в аккаунт.',
  save_failed: 'Не удалось сохранить покупку. Попробуйте ещё раз.',
  invalid_json: 'Не удалось разобрать запрос. Попробуйте ещё раз.'
}

function extractPayload(error: unknown): ApiErrorPayload | null {
  if (!error || typeof error !== 'object') return null
  const data = (error as { data?: ApiErrorPayload }).data
  if (data && typeof data === 'object') return data
  return null
}

/** Человекочитаемое сообщение из ответа API или $fetch-ошибки. */
export function formatApiError(error: unknown, fallback: string): string {
  const payload = extractPayload(error)
  if (payload?.code && ERROR_MESSAGES[payload.code]) {
    return ERROR_MESSAGES[payload.code]
  }
  if (payload?.error?.trim()) {
    const known = Object.values(ERROR_MESSAGES).find((msg) => msg === payload.error)
    if (known) return known
    if (!looksLikeRawHttpError(payload.error)) {
      return payload.error.trim()
    }
  }
  if (payload?.message?.trim() && !looksLikeRawHttpError(payload.message)) {
    return payload.message.trim()
  }

  if (error instanceof Error && error.message.trim()) {
    if (ERROR_MESSAGES[error.message]) {
      return ERROR_MESSAGES[error.message]
    }
    if (!looksLikeRawHttpError(error.message)) {
      return error.message.trim()
    }
  }

  return fallback
}

function looksLikeRawHttpError(text: string): boolean {
  const value = text.trim()
  return (
    /^\[?(?:GET|POST|PUT|PATCH|DELETE|HEAD)\]/i.test(value) ||
    /^HTTP/i.test(value) ||
    /^\d{3}\s/.test(value) ||
    value.includes('Bad Gateway') ||
    value.includes('Failed to fetch')
  )
}
