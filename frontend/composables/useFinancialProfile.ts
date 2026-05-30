import type { FinancialProfile, ProfilePatchRequest } from '~/types/api'
import { currentUserStorageKey, userStorageKey } from '~/utils/userStorage'

const PROFILE_PREFIX = 'potok:financial-profile'

export const defaultFinancialProfile: FinancialProfile = {
  active_income: 0,
  passive_income: 0,
  emergency_fund: 0,
  emergency_breakdown: { cash: 0, deposit: 0, investments: 0 },
  fixed_expenses: [],
  goal_kind: 'save',
  goal_title: '',
  goal_amount: 0,
  skipped_income: false,
  skipped_cushion: false,
  skipped_goal: false,
  skipped_expenses: false,
  onboarding_completed: false
}

export function readStoredProfile(
  phone?: string | null,
  userId?: string | null
): FinancialProfile {
  if (!import.meta.client) return { ...defaultFinancialProfile }
  const key =
    phone !== undefined || userId !== undefined
      ? userStorageKey(PROFILE_PREFIX, phone, userId)
      : currentUserStorageKey(PROFILE_PREFIX)
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return { ...defaultFinancialProfile }
    return { ...defaultFinancialProfile, ...JSON.parse(raw) }
  } catch {
    return { ...defaultFinancialProfile }
  }
}

function writeStoredProfile(profile: FinancialProfile) {
  if (!import.meta.client) return
  localStorage.setItem(currentUserStorageKey(PROFILE_PREFIX), JSON.stringify(profile))
}

function profileHasGoal(profile: FinancialProfile): boolean {
  return !profile.skipped_goal && (profile.goal_amount ?? 0) >= 1000
}

function profileHasIncome(profile: FinancialProfile): boolean {
  return (profile.active_income ?? 0) > 0 || (profile.passive_income ?? 0) > 0
}

/** Сохраняет цель и доход из localStorage, если API вернул пустой профиль. */
export function mergeProfileFromApi(
  local: FinancialProfile,
  remote: FinancialProfile
): FinancialProfile {
  const merged: FinancialProfile = { ...defaultFinancialProfile, ...remote }

  if (!profileHasGoal(remote) && profileHasGoal(local)) {
    merged.goal_kind = local.goal_kind
    merged.goal_title = local.goal_title
    merged.goal_amount = local.goal_amount
    merged.skipped_goal = local.skipped_goal
  }

  if (!profileHasIncome(remote) && profileHasIncome(local)) {
    merged.active_income = local.active_income
    merged.passive_income = local.passive_income
    merged.skipped_income = local.skipped_income
  }

  if ((merged.emergency_fund ?? 0) <= 0 && (local.emergency_fund ?? 0) > 0) {
    merged.emergency_fund = local.emergency_fund
    merged.emergency_breakdown = local.emergency_breakdown
    merged.skipped_cushion = local.skipped_cushion
  }

  if (
    !(merged.fixed_expenses?.some((row) => row.amount > 0)) &&
    local.fixed_expenses?.some((row) => row.amount > 0)
  ) {
    merged.fixed_expenses = local.fixed_expenses
    merged.skipped_expenses = local.skipped_expenses
  }

  if (local.onboarding_completed && !remote.onboarding_completed) {
    merged.onboarding_completed = true
  }

  return merged
}

function remoteProfileMissingGoal(remote: FinancialProfile): boolean {
  return !profileHasGoal(remote)
}

function isEndpointMissing(error: unknown): boolean {
  const status = (error as { statusCode?: number })?.statusCode
  if (status === 404 || status === 501) return true
  if (error instanceof Error) {
    return /404|501|not found|не найден/i.test(error.message)
  }
  return false
}

export function useFinancialProfile() {
  const { apiFetch, demoMode } = useApi()
  const profile = useState<FinancialProfile>('financial-profile', () => ({
    ...defaultFinancialProfile
  }))

  if (import.meta.client) {
    const bootstrapped = useState('financial-profile-bootstrapped', () => false)
    if (!bootstrapped.value) {
      profile.value = readStoredProfile()
      bootstrapped.value = true
    }
  }

  const totalIncome = computed(() => profile.value.active_income + profile.value.passive_income)

  function loadProfile() {
    profile.value = readStoredProfile()
  }

  function saveProfile(partial: Partial<FinancialProfile>) {
    profile.value = { ...profile.value, ...partial }
    writeStoredProfile(profile.value)
  }

  function resetProfile() {
    profile.value = { ...defaultFinancialProfile }
    writeStoredProfile(profile.value)
  }

  async function syncProfileToApi(source: FinancialProfile) {
    if (demoMode.value) return

    const body: ProfilePatchRequest = {
      active_income: Math.max(0, source.active_income),
      passive_income: Math.max(0, source.passive_income),
      emergency_fund: Math.max(0, source.emergency_fund),
      emergency_breakdown: source.emergency_breakdown,
      fixed_expenses: source.fixed_expenses ?? [],
      goal_kind: source.goal_kind,
      goal_title: source.goal_title,
      goal_amount: Math.max(0, source.goal_amount ?? 0),
      skipped_income: source.skipped_income,
      skipped_cushion: source.skipped_cushion,
      skipped_goal: source.skipped_goal,
      skipped_expenses: source.skipped_expenses,
      survey_input_mode: source.survey_input_mode,
      onboarding_completed: source.onboarding_completed
    }

    try {
      await apiFetch('/users/me/profile', { method: 'PATCH', body })
    } catch (error) {
      if (!isEndpointMissing(error)) throw error
    }
  }

  async function fetchProfileFromApi() {
    if (demoMode.value) return
    const local = readStoredProfile()
    try {
      const remote = await apiFetch<FinancialProfile>('/users/me/profile')
      const merged = mergeProfileFromApi(local, remote)
      profile.value = merged
      writeStoredProfile(merged)

      if (profileHasGoal(local) && remoteProfileMissingGoal(remote)) {
        await syncProfileToApi(merged)
      }
    } catch (error) {
      if (!isEndpointMissing(error)) throw error
    }
  }

  async function markOnboardingCompleteOnApi() {
    if (demoMode.value) return

    try {
      await apiFetch('/users/me/onboarding/complete', {
        method: 'POST',
        body: { onboarding_completed: true }
      })
    } catch (error) {
      if (!isEndpointMissing(error)) throw error
    }
  }

  return {
    profile,
    totalIncome,
    loadProfile,
    saveProfile,
    resetProfile,
    syncProfileToApi,
    fetchProfileFromApi,
    markOnboardingCompleteOnApi
  }
}
