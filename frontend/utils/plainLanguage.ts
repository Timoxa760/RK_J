/** Простые подписи вместо финансового жаргона в UI. */

const INDICATOR_LABELS: Array<{ match: RegExp; label: string }> = [
  { match: /emergency|эмердженси|emergenc/i, label: 'Запас на чёрный день' },
  { match: /runway|run.?way/i, label: 'Запас на чёрный день' },
  { match: /dti|debt.?to.?income/i, label: 'Платежи по кредитам' },
  { match: /подушк|cushion|reserve/i, label: 'Запас на чёрный день' },
  { match: /накоплен.*доход|savings.?rate|saving/i, label: 'Сколько откладываете' },
  { match: /импульс|impulse/i, label: 'Покупки на эмоциях' },
  { match: /стабильност.*доход|income.?stab/i, label: 'Насколько стабилен доход' },
  { match: /долгов.*нагруз|debt.?load/i, label: 'Платежи по кредитам' },
  { match: /mindfulness|осознан/i, label: 'Контроль трат' },
  { match: /cashflow|cash.?flow|свободн.*поток/i, label: 'Остаётся после трат' }
]

export function humanizeIndicatorName(name: string): string {
  const trimmed = name.trim()
  if (!trimmed) return trimmed

  for (const { match, label } of INDICATOR_LABELS) {
    if (match.test(trimmed)) return label
  }

  return trimmed
    .replace(/\s*\(\s*мес\.?\s*\)\s*/gi, '')
    .replace(/\s*\(\s*months?\s*\)\s*/gi, '')
    .replace(/\s*\(\s*%?\s*\)\s*/g, '')
    .trim()
}

export function humanizeIndicatorNorm(norm: string, indicatorName?: string): string {
  const n = norm.trim()
  if (!n) return n

  const name = (indicatorName ?? '').toLowerCase()
  const isMonths =
    name.includes('подуш') ||
    name.includes('запас') ||
    name.includes('emergency') ||
    name.includes('эмердженси') ||
    name.includes('runway') ||
    n.includes('мес')

  if (isMonths) {
    if (n.startsWith('>')) return `обычно от ${n.slice(1).trim()} мес.`
    if (n.startsWith('<')) return `обычно до ${n.slice(1).trim()} мес.`
    return `обычно ${n}`
  }

  if (n.endsWith('%') || n.includes('%')) return `обычно ${n}`
  if (n.startsWith('>')) return `обычно от ${n.slice(1).trim()}%`
  if (n.startsWith('<')) return `обычно до ${n.slice(1).trim()}%`
  return `обычно ${n}`
}

export function formatIndicatorValue(name: string, value: number): string {
  const key = name.toLowerCase()

  if (
    key.includes('подуш') ||
    key.includes('запас') ||
    key.includes('emergency') ||
    key.includes('эмердженси') ||
    key.includes('runway') ||
    key.includes('мес')
  ) {
    return `${value.toLocaleString('ru-RU', { maximumFractionDigits: 1 })} мес.`
  }

  if (
    key.includes('нагруз') ||
    key.includes('доход') ||
    key.includes('dti') ||
    key.includes('накоплен') ||
    key.includes('импульс') ||
    key.includes('стабиль') ||
    key.includes('%')
  ) {
    return `${value}%`
  }

  return value.toLocaleString('ru-RU')
}
