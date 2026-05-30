import type { OnboardingDraft, OnboardingParseStep } from '~/types/api'
import {
  parseCushionAnswer,
  parseExpensesAnswer,
  parseGoalAnswer,
  parseIncomeAnswer
} from '~/utils/onboardingVoiceParser'

/** Локальный разбор опроса онбординга (API v3 — без POST /onboarding/parse). */
export function parseOnboardingStepLocal(
  step: OnboardingParseStep,
  rawText: string
): Partial<OnboardingDraft> {
  const text = rawText.trim()
  if (!text) return {}

  switch (step) {
    case 'income':
      return parseIncomeAnswer(text)
    case 'cushion':
      return parseCushionAnswer(text)
    case 'goal':
      return parseGoalAnswer(text)
    case 'expenses': {
      const parsed = parseExpensesAnswer(text)
      return {
        fixed_expenses: parsed.items,
        skipped_expenses: parsed.skipped
      }
    }
    default:
      return {}
  }
}

export function useOnboardingParse() {
  const parsing = ref(false)

  async function parseStep(
    step: OnboardingParseStep,
    rawText: string
  ): Promise<Partial<OnboardingDraft>> {
    parsing.value = true
    try {
      return parseOnboardingStepLocal(step, rawText)
    } finally {
      parsing.value = false
    }
  }

  return {
    parsing,
    parseStep,
    parseStepLocal: parseOnboardingStepLocal
  }
}
