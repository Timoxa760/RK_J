import type { CategoriesResponse } from '~/types/api'

/** Канонические категории — те же, что у AI/regex при разборе покупок (без новых имён). */
export const STANDARD_EXPENSE_CATEGORIES = [
  'Продукты',
  'Кафе и рестораны',
  'Транспорт',
  'Доставка',
  'Подписки',
  'ЖКХ',
  'Развлечения',
  'Одежда',
  'Здоровье',
  'Прочие расходы'
] as const

export type StandardExpenseCategory = (typeof STANDARD_EXPENSE_CATEGORIES)[number]

const ALIAS_TO_STANDARD: Record<string, StandardExpenseCategory> = {
  прочее: 'Прочие расходы',
  'прочие расходы': 'Прочие расходы',
  прочие: 'Прочие расходы',
  связь: 'Прочие расходы',
  продукты: 'Продукты',
  'кафе и рестораны': 'Кафе и рестораны',
  кафе: 'Кафе и рестораны',
  рестораны: 'Кафе и рестораны',
  транспорт: 'Транспорт',
  доставка: 'Доставка',
  подписки: 'Подписки',
  подписка: 'Подписки',
  жкх: 'ЖКХ',
  коммунальные: 'ЖКХ',
  развлечения: 'Развлечения',
  одежда: 'Одежда',
  здоровье: 'Здоровье',
  кредиты: 'Прочие расходы',
  доход: 'Прочие расходы'
}

/** Приводит название из API/AI к одной из стандартных категорий. */
export function normalizeExpenseCategory(name: string): StandardExpenseCategory {
  const trimmed = name.trim()
  if (!trimmed) return 'Прочие расходы'

  const lower = trimmed.toLowerCase()
  if (ALIAS_TO_STANDARD[lower]) return ALIAS_TO_STANDARD[lower]

  for (const standard of STANDARD_EXPENSE_CATEGORIES) {
    const stdLower = standard.toLowerCase()
    if (lower === stdLower || lower.includes(stdLower) || stdLower.includes(lower)) {
      return standard
    }
  }

  return 'Прочие расходы'
}

export interface UserCategoryOption {
  name: StandardExpenseCategory
  amount: number
  share: number
}

/** Сводит траты пользователя в 5–10 стандартных категорий для симулятора. */
export function buildUserCategoryOptions(
  categories: CategoriesResponse | null,
  maxItems = 8
): UserCategoryOption[] {
  const totals = new Map<StandardExpenseCategory, number>()

  for (const row of categories?.categories ?? []) {
    const amount = row.amount ?? row.total ?? 0
    if (amount <= 0) continue
    const normalized = normalizeExpenseCategory(row.name)
    totals.set(normalized, (totals.get(normalized) ?? 0) + amount)
  }

  const grandTotal = [...totals.values()].reduce((sum, value) => sum + value, 0)
  if (grandTotal <= 0) return []

  let items = [...totals.entries()]
    .map(([name, amount]) => ({
      name,
      amount,
      share: amount / grandTotal
    }))
    .sort((a, b) => b.amount - a.amount)

  if (items.length > maxItems) {
    const head = items.slice(0, maxItems - 1)
    const tailAmount = items.slice(maxItems - 1).reduce((sum, row) => sum + row.amount, 0)
    const other = head.find((row) => row.name === 'Прочие расходы')
    if (other) {
      other.amount += tailAmount
      other.share = other.amount / grandTotal
    } else {
      head.push({
        name: 'Прочие расходы',
        amount: tailAmount,
        share: tailAmount / grandTotal
      })
    }
    items = head
  }

  return items.map((row) => ({
    name: row.name,
    amount: row.amount,
    share: row.amount / grandTotal
  }))
}

export const CATEGORY_EMOJI: Partial<Record<StandardExpenseCategory, string>> = {
  Продукты: '🛒',
  'Кафе и рестораны': '☕',
  Транспорт: '🚌',
  Доставка: '🛵',
  Подписки: '📱',
  ЖКХ: '🏠',
  Развлечения: '🎬',
  Одежда: '👕',
  Здоровье: '💊',
  'Прочие расходы': '📦'
}
