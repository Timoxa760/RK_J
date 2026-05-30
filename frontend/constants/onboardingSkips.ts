import type { OnboardingDraft, OnboardingParseStep } from '~/types/api'

export type OnboardingSurveySkipKey = OnboardingParseStep

export const ONBOARDING_SKIP_HINTS: Record<OnboardingSurveySkipKey, string> = {
  income: 'Можно пропустить. Укажете доход в профиле — прогноз и советы станут точнее.',
  cushion: 'Можно пропустить. Добавите запас позже — покажем, на сколько месяцев хватит денег.',
  goal: 'Можно пропустить. Поставите цель в профиле — покажем, когда примерно дойдёте.',
  expenses: 'Можно пропустить. Добавите покупки — сами поймём ваши траты.'
}

export const ONBOARDING_SKIP_LABEL = 'Пропустить'

export function skipPatchForSurveyStep(step: OnboardingSurveySkipKey): Partial<OnboardingDraft> {
  switch (step) {
    case 'income':
      return { skipped_income: true, active_income: 0, passive_income: 0 }
    case 'cushion':
      return {
        skipped_cushion: true,
        emergency_fund: 0,
        emergency_breakdown: { cash: 0, deposit: 0, investments: 0 }
      }
    case 'goal':
      return { skipped_goal: true, goal_amount: 0 }
    case 'expenses':
      return { skipped_expenses: true, fixed_expenses: [] }
  }
}

export function surveyStepToSkipKey(step: number): OnboardingSurveySkipKey | null {
  const map: Record<number, OnboardingSurveySkipKey> = {
    1: 'income',
    2: 'cushion',
    3: 'goal',
    4: 'expenses'
  }
  return map[step] ?? null
}

export function draftHasSurveyProgress(draft: OnboardingDraft): boolean {
  const income = draft.active_income + draft.passive_income
  return (
    income > 0 ||
    draft.emergency_fund > 0 ||
    draft.goal_amount >= 1000 ||
    draft.fixed_expenses.some((item) => item.amount > 0) ||
    draft.skipped_income === true ||
    draft.skipped_cushion === true ||
    draft.skipped_goal === true ||
    draft.skipped_expenses === true
  )
}
