import type { OnboardingDraft, OnboardingParseStep } from '~/types/api'
import {
  parseCushionAnswer,
  parseExpensesAnswer,
  parseGoalAnswer,
  parseIncomeAnswer
} from '~/utils/onboardingVoiceParser'

interface OnboardingParseApiResponse {
  parsed: boolean
  step: OnboardingParseStep
  patch: Partial<OnboardingDraft>
  message?: string
}

/** Парсер опроса онбординга: API + локальный fallback/обогащение. */
export function useOnboardingParse() {
  const { apiFetch, demoMode } = useApi()
  const parsing = ref(false)

  async function parseStep(
    step: OnboardingParseStep,
    rawText: string
  ): Promise<Partial<OnboardingDraft>> {
    parsing.value = true
    try {
      const local = parseOnboardingStepLocal(step, rawText)

      if (!demoMode.value) {
        try {
          const res = await apiFetch<OnboardingParseApiResponse>('/onboarding/parse', {
            method: 'POST',
            body: { step, raw_text: rawText.trim(), locale: 'ru' }
          })
          if (res.parsed && res.patch) {
            return mergeParsePatch(step, res.patch, local)
          }
        } catch {
          /* fallback на локальный regex */
        }
      }

      return local
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

/** Объединяет ответ API с локальным regex — локальный разбор часто полнее для «Запаса». */
function mergeParsePatch(
  step: OnboardingParseStep,
  apiPatch: Partial<OnboardingDraft>,
  localPatch: Partial<OnboardingDraft>
): Partial<OnboardingDraft> {
  if (step === 'cushion') {
    const lb = localPatch.emergency_breakdown
    const localHasBreakdown =
      Boolean(lb) &&
      ((lb?.cash ?? 0) > 0 || (lb?.deposit ?? 0) > 0 || (lb?.investments ?? 0) > 0)
    const apiFund = apiPatch.emergency_fund ?? 0
    const localFund = localPatch.emergency_fund ?? 0

    if (localHasBreakdown && localFund >= apiFund) {
      return {
        ...apiPatch,
        emergency_fund: localFund,
        emergency_breakdown: lb
      }
    }
  }

  if (step === 'income') {
    const apiIncome = (apiPatch.active_income ?? 0) + (apiPatch.passive_income ?? 0)
    const localIncome = (localPatch.active_income ?? 0) + (localPatch.passive_income ?? 0)
    if (localIncome > apiIncome) {
      return { ...apiPatch, ...localPatch }
    }
  }

  return apiPatch
}

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
