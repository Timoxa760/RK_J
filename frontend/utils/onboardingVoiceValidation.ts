import type { OnboardingDraft, OnboardingParseStep } from '~/types/api'

export function isMeaningfulOnboardingPatch(
  step: OnboardingParseStep | string,
  patch: Partial<OnboardingDraft>
): boolean {
  if (!patch || Object.keys(patch).length === 0) return false
  switch (step) {
    case 'income':
      return (patch.active_income ?? 0) > 0 || (patch.passive_income ?? 0) > 0
    case 'cushion':
      return (patch.emergency_fund ?? 0) > 0
    case 'goal':
      return (patch.goal_amount ?? 0) >= 1000
    case 'expenses':
      return Boolean(patch.skipped_expenses) || (patch.fixed_expenses?.length ?? 0) > 0
    default:
      return Object.keys(patch).length > 0
  }
}
