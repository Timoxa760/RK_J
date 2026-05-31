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

const API_ERROR_TEXT: Record<string, string> = {
  'invalid phone or password': 'Неверный телефон или пароль',
  'user already exists': 'Этот номер уже зарегистрирован',
  'phone and password required': 'Укажите телефон и пароль',
  'phone required': 'Укажите телефон',
  'password must be at least 8 characters': 'Пароль — минимум 8 символов',
  'invalid or expired code': 'Неверный или просроченный код',
  'registration failed': 'Не удалось зарегистрироваться',
  'request failed': 'Не удалось выполнить запрос'
}

const STATUS_FALLBACKS: Record<number, string> = {
  401: 'Неверный телефон или пароль',
  409: 'Этот номер уже зарегистрирован',
  400: 'Проверьте введённые данные'
}

function extractPayload(error: unknown): ApiErrorPayload | null {
  if (!error || typeof error !== 'object') return null
  const data = (error as { data?: ApiErrorPayload }).data
  if (data && typeof data === 'object') return data
  return null
}

function extractStatusCode(error: unknown): number | undefined {
  if (!error || typeof error !== 'object') return undefined
  const code = (error as { statusCode?: number }).statusCode
  return typeof code === 'number' ? code : undefined
}

function mapApiErrorText(raw: string): string | null {
  const key = raw.trim().toLowerCase()
  return API_ERROR_TEXT[key] ?? null
}

/** Человекочитаемое сообщение из ответа API или $fetch-ошибки. */
export function formatApiError(error: unknown, fallback: string): string {
  const payload = extractPayload(error)
  if (payload?.code && ERROR_MESSAGES[payload.code]) {
    return ERROR_MESSAGES[payload.code]
  }
  if (payload?.error?.trim()) {
    const mapped = mapApiErrorText(payload.error)
    if (mapped) return mapped
    const known = Object.values(ERROR_MESSAGES).find((msg) => msg === payload.error)
    if (known) return known
    if (!looksLikeRawHttpError(payload.error)) {
      return payload.error.trim()
    }
  }
  if (payload?.message?.trim() && !looksLikeRawHttpError(payload.message)) {
    const mapped = mapApiErrorText(payload.message)
    if (mapped) return mapped
    return payload.message.trim()
  }

  const status = extractStatusCode(error)
  if (status != null && STATUS_FALLBACKS[status]) {
    return STATUS_FALLBACKS[status]!
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
