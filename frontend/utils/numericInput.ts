const DIGIT_NAV_KEYS = new Set([
  'Backspace',
  'Delete',
  'Tab',
  'Escape',
  'Enter',
  'ArrowLeft',
  'ArrowRight',
  'ArrowUp',
  'ArrowDown',
  'Home',
  'End'
])

/** Блокирует ввод нецифровых символов в поле телефона/кода. */
export function allowDigitKeydown(event: KeyboardEvent) {
  if (event.ctrlKey || event.metaKey || event.altKey) return
  if (DIGIT_NAV_KEYS.has(event.key)) return
  if (/^\d$/.test(event.key)) return
  event.preventDefault()
}

export function digitsOnly(value: string, maxLen: number): string {
  return value.replace(/\D/g, '').slice(0, maxLen)
}

export function onDigitPaste(event: ClipboardEvent, apply: (digits: string) => void, maxLen: number) {
  event.preventDefault()
  const text = event.clipboardData?.getData('text') ?? ''
  apply(digitsOnly(text, maxLen))
}
