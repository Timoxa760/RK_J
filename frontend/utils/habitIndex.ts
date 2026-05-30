import type { StoresResponse } from '~/types/api'

export type HabitTone = 'good' | 'warn' | 'risk'

export interface HabitIndex {
  score: number
  label: string
  tone: HabitTone
  insight: string | null
  challengeHint: string
}

export function buildHabitIndex(stores: StoresResponse | null): HabitIndex {
  if (!stores?.stores?.length) {
    return {
      score: 0,
      label: 'Нужны данные',
      tone: 'warn',
      insight: null,
      challengeHint: 'Добавьте расходы или чеки — покажем, как у вас с тратами.'
    }
  }

  const totalSpend = stores.stores.reduce((sum, store) => sum + store.total, 0)
  const impulseWeighted = stores.stores.reduce(
    (sum, store) => sum + store.impulse_ratio * store.total,
    0
  )
  const avgImpulse = totalSpend > 0 ? impulseWeighted / totalSpend : 0
  const score = Math.round(Math.max(25, Math.min(92, 100 - avgImpulse * 70)))

  let label: string
  let tone: HabitTone
  if (score >= 75) {
    label = 'Спокойно'
    tone = 'good'
  } else if (score >= 55) {
    label = 'Смешанно'
    tone = 'warn'
  } else {
    label = 'Много лишних трат'
    tone = 'risk'
  }

  const top = [...stores.stores].sort((a, b) => b.impulse_ratio - a.impulse_ratio)[0]
  let insight: string | null = null
  if (top && top.impulse_ratio >= 0.45) {
    const pct = Math.round(top.impulse_ratio * 100)
    insight = `Около ${pct}% покупок в «${top.name}» — на эмоциях. Задание с друзьями поможет закрепить привычку.`
  } else {
    insight = 'Лишних покупок мало — можно соревноваться в днях подряд без срыва.'
  }

  const challengeHint =
    score >= 75
      ? 'Предложите друзьям «дней подряд без срыва» — закрепите хороший ритм.'
      : 'Попробуйте «Меньше трат» или «Кофе ≤ 3 раза» — без сравнения сумм.'

  return { score, label, tone, insight, challengeHint }
}
