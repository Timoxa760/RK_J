export function estimateGoalImpactDays(amount: number, dailySaving = 1_200): number {
  if (amount <= 0 || dailySaving <= 0) return 0
  return Math.max(1, Math.round(amount / dailySaving))
}

export function formatGoalImpact(amount: number): string {
  const days = estimateGoalImpactDays(amount)
  if (!days) return ''
  if (days === 1) return '≈ 1 день к цели'
  if (days < 7) return `≈ ${days} дн. к цели`
  const weeks = Math.round(days / 7)
  return weeks === 1 ? '≈ 1 нед. к цели' : `≈ ${weeks} нед. к цели`
}
