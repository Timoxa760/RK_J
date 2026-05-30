import type { GoalKind } from '~/types/api'

export const GOAL_KIND_LABELS: Record<GoalKind, string> = {
  save: 'Накопить сумму',
  purchase: 'Крупная покупка',
  cushion: 'Подушка безопасности',
  other: 'Другое'
}
