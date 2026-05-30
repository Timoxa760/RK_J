import type { FinancialProfile, GoalKind, OnboardingDraft } from '~/types/api'
import { GOAL_KIND_LABELS } from '~/constants/onboardingGoals'
import { currentUserStorageKey, userStorageKey } from '~/utils/userStorage'

export { GOAL_KIND_LABELS } from '~/constants/onboardingGoals'

export const ONBOARDING_COMPLETED_PREFIX = 'potok:onboarding-completed'
export const ONBOARDING_DRAFT_PREFIX = 'potok:onboarding-draft'
import { ONBOARDING_WIZARD_STEP_COUNT } from '~/constants/onboardingSteps'

export const ONBOARDING_STEP_COUNT = ONBOARDING_WIZARD_STEP_COUNT

export function defaultOnboardingDraft(): OnboardingDraft {
  return {
    active_income: 0,
    passive_income: 0,
    emergency_fund: 0,
    emergency_breakdown: { cash: 0, deposit: 0, investments: 0 },
    goal_kind: 'save',
    goal_title: GOAL_KIND_LABELS.save,
    goal_amount: 0,
    fixed_expenses: [],
    skipped_expenses: false
  }
}

function completedKey(phone?: string | null, userId?: string | null) {
  return userStorageKey(ONBOARDING_COMPLETED_PREFIX, phone, userId)
}

function draftKey(phone?: string | null, userId?: string | null) {
  return userStorageKey(ONBOARDING_DRAFT_PREFIX, phone, userId)
}

export function readOnboardingCompleted(phone?: string | null, userId?: string | null): boolean {
  if (!import.meta.client) return false

  const key =
    phone !== undefined || userId !== undefined
      ? completedKey(phone, userId)
      : currentUserStorageKey(ONBOARDING_COMPLETED_PREFIX)

  return localStorage.getItem(key) === 'true'
}

export function writeOnboardingCompleted(
  value: boolean,
  phone?: string | null,
  userId?: string | null
) {
  if (!import.meta.client) return
  const key =
    phone !== undefined || userId !== undefined
      ? completedKey(phone, userId)
      : currentUserStorageKey(ONBOARDING_COMPLETED_PREFIX)

  if (value) {
    localStorage.setItem(key, 'true')
  } else {
    localStorage.removeItem(key)
  }
}

export function readOnboardingDraft(phone?: string | null, userId?: string | null): OnboardingDraft {
  if (!import.meta.client) return defaultOnboardingDraft()
  const key =
    phone !== undefined || userId !== undefined
      ? draftKey(phone, userId)
      : currentUserStorageKey(ONBOARDING_DRAFT_PREFIX)

  try {
    const raw = localStorage.getItem(key)
    if (!raw) return defaultOnboardingDraft()
    return { ...defaultOnboardingDraft(), ...JSON.parse(raw) }
  } catch {
    return defaultOnboardingDraft()
  }
}

export function writeOnboardingDraft(draft: OnboardingDraft, phone?: string | null, userId?: string | null) {
  if (!import.meta.client) return
  const key =
    phone !== undefined || userId !== undefined
      ? draftKey(phone, userId)
      : currentUserStorageKey(ONBOARDING_DRAFT_PREFIX)
  localStorage.setItem(key, JSON.stringify(draft))
}

export function clearOnboardingDraft(phone?: string | null, userId?: string | null) {
  if (!import.meta.client) return
  const key =
    phone !== undefined || userId !== undefined
      ? draftKey(phone, userId)
      : currentUserStorageKey(ONBOARDING_DRAFT_PREFIX)
  localStorage.removeItem(key)
}

export function goalTitleForKind(kind: GoalKind, custom?: string): string {
  if (kind === 'other') return custom?.trim() || GOAL_KIND_LABELS.other
  return GOAL_KIND_LABELS[kind]
}

export function buildOnboardingSummary(draft: OnboardingDraft) {
  const income = draft.active_income + draft.passive_income
  const monthlySaving = draft.active_income > 0 ? Math.round(draft.active_income * 0.1) : 0
  const fixedTotal = draft.skipped_expenses
    ? 0
    : draft.fixed_expenses.reduce((sum, item) => sum + item.amount, 0)
  const freeCashflow = income - fixedTotal

  let goalForecast = 'Поставьте сумму цели — покажем примерный срок.'
  if (draft.goal_amount > 0 && monthlySaving > 0) {
    const months = Math.ceil(draft.goal_amount / monthlySaving)
    goalForecast = `При текущем поведении до «${draft.goal_title}» примерно ${months} мес.`
  } else if (draft.goal_amount > 0) {
    goalForecast = `Цель «${draft.goal_title}» — ${draft.goal_amount.toLocaleString('ru-RU')} ₽.`
  }

  const runwayMonths =
    fixedTotal > 0 && draft.emergency_fund > 0
      ? Math.round((draft.emergency_fund / fixedTotal) * 10) / 10
      : null

  return {
    income,
    monthlySaving,
    fixedTotal,
    freeCashflow,
    goalForecast,
    runwayMonths
  }
}

export function useOnboarding() {
  const step = ref(1)
  const draft = ref<OnboardingDraft>(defaultOnboardingDraft())
  const finishing = ref(false)
  const error = ref<string | null>(null)

  const { saveProfile, loadProfile, syncProfileToApi, markOnboardingCompleteOnApi } =
    useFinancialProfile()
  const { createGoal } = useGoals()

  const summary = computed(() => buildOnboardingSummary(draft.value))

  function hydrate() {
    draft.value = readOnboardingDraft()
    step.value = 1
  }

  function persistDraft() {
    writeOnboardingDraft(draft.value)
  }

  function patchDraft(partial: Partial<OnboardingDraft>) {
    draft.value = { ...draft.value, ...partial }
    persistDraft()
  }

  function setGoalKind(kind: GoalKind) {
    patchDraft({
      goal_kind: kind,
      goal_title: goalTitleForKind(kind, draft.value.goal_title)
    })
  }

  function syncEmergencyTotalFromBreakdown() {
    const breakdown = draft.value.emergency_breakdown
    patchDraft({
      emergency_fund: breakdown.cash + breakdown.deposit + breakdown.investments
    })
  }

  function patchBreakdown(partial: Partial<OnboardingDraft['emergency_breakdown']>) {
    const emergency_breakdown = { ...draft.value.emergency_breakdown, ...partial }
    patchDraft({
      emergency_breakdown,
      emergency_fund:
        emergency_breakdown.cash + emergency_breakdown.deposit + emergency_breakdown.investments
    })
  }

  function addFixedExpense() {
    patchDraft({
      fixed_expenses: [...draft.value.fixed_expenses, { title: '', amount: 0 }]
    })
  }

  function updateFixedExpense(index: number, partial: Partial<{ title: string; amount: number }>) {
    const next = draft.value.fixed_expenses.map((item, i) =>
      i === index ? { ...item, ...partial } : item
    )
    patchDraft({ fixed_expenses: next })
  }

  function removeFixedExpense(index: number) {
    patchDraft({
      fixed_expenses: draft.value.fixed_expenses.filter((_, i) => i !== index)
    })
  }

  function skipExpenses() {
    patchDraft({ skipped_expenses: true, fixed_expenses: [] })
  }

  function canProceed(current = step.value): boolean {
    switch (current) {
      case 1:
        return true
      case 2:
        return draft.value.active_income > 0 || draft.value.passive_income > 0
      case 3:
        return draft.value.emergency_fund >= 0
      case 4:
        return draft.value.goal_amount >= 1000
      case 5:
        return true
      case 6:
        return true
      case 7:
        return true
      default:
        return false
    }
  }

  function nextStep() {
    if (!canProceed()) return
    step.value = Math.min(ONBOARDING_STEP_COUNT, step.value + 1)
    persistDraft()
  }

  function prevStep() {
    step.value = Math.max(1, step.value - 1)
  }

  function isComplete() {
    return readOnboardingCompleted()
  }

  async function completeOnboarding() {
    finishing.value = true
    error.value = null

    try {
      const current = draft.value
      const profilePayload: FinancialProfile = {
        active_income: Math.max(0, current.active_income),
        passive_income: Math.max(0, current.passive_income),
        emergency_fund: Math.max(0, current.emergency_fund),
        emergency_breakdown: current.emergency_breakdown,
        fixed_expenses: current.skipped_expenses ? [] : current.fixed_expenses,
        onboarding_completed: true
      }

      saveProfile(profilePayload)
      loadProfile()

      await syncProfileToApi(profilePayload)

      if (current.goal_amount >= 1000) {
        await createGoal({
          title: current.goal_title,
          target_amount: current.goal_amount,
          auto_save_percent: 10
        })
      }

      await markOnboardingCompleteOnApi()

      writeOnboardingCompleted(true)
      clearOnboardingDraft()

      await navigateTo('/dashboard')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось завершить опрос'
    } finally {
      finishing.value = false
    }
  }

  return {
    step,
    draft,
    summary,
    finishing,
    error,
    hydrate,
    persistDraft,
    patchDraft,
    setGoalKind,
    syncEmergencyTotalFromBreakdown,
    patchBreakdown,
    addFixedExpense,
    updateFixedExpense,
    removeFixedExpense,
    skipExpenses,
    canProceed,
    nextStep,
    prevStep,
    isComplete,
    completeOnboarding
  }
}
