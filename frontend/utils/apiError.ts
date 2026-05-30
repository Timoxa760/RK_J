const KNOWN_API_ERRORS: Record<string, string> = {
  'amount required': 'Не удалось найти сумму. Скажите цену, например: «колбаса 300 рублей».',
  'user_id required': 'Сначала войдите в аккаунт.',
  'invalid json': 'Не удалось прочитать запрос. Попробуйте ещё раз.',
  'save failed': 'Не удалось сохранить покупку. Попробуйте позже.'
}

function isRawHttpMessage(message: string): boolean {
  return /^\[(?:GET|POST|PUT|PATCH|DELETE|HEAD)]/i.test(message) || /^fetch failed/i.test(message)
}

/** Человекочитаемый текст из ошибки $fetch / API. */
export function formatApiError(error: unknown, fallback: string): string {
  if (error && typeof error === 'object') {
    const data = (error as { data?: { error?: string; message?: string } }).data
    const apiError = data?.error?.trim() || data?.message?.trim()
    if (apiError) {
      return KNOWN_API_ERRORS[apiError] ?? apiError
    }
  }

  if (error instanceof Error) {
    const msg = error.message.trim()
    if (KNOWN_API_ERRORS[msg]) return KNOWN_API_ERRORS[msg]
    if (msg && !isRawHttpMessage(msg)) return msg
  }

  return fallback
}
