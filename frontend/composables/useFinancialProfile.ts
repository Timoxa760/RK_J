import type { FinancialProfile, ProfilePatchRequest } from '~/types/api'
import { currentUserStorageKey } from '~/utils/userStorage'

const PROFILE_PREFIX = 'potok:financial-profile'

const defaultProfile: FinancialProfile = {
  active_income: 0,
  passive_income: 0,
  emergency_fund: 0,
  emergency_breakdown: { cash: 0, deposit: 0, investments: 0 },
  fixed_expenses: [],
  onboarding_completed: false
}

function readStoredProfile(): FinancialProfile {
  if (!import.meta.client) return { ...defaultProfile }
  try {
    const raw = localStorage.getItem(currentUserStorageKey(PROFILE_PREFIX))
    if (!raw) return { ...defaultProfile }
    return { ...defaultProfile, ...JSON.parse(raw) }
  } catch {
    return { ...defaultProfile }
  }
}

function writeStoredProfile(profile: FinancialProfile) {
  if (!import.meta.client) return
  localStorage.setItem(currentUserStorageKey(PROFILE_PREFIX), JSON.stringify(profile))
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
  const profile = ref<FinancialProfile>({ ...defaultProfile })

  const totalIncome = computed(() => profile.value.active_income + profile.value.passive_income)

  function loadProfile() {
    profile.value = readStoredProfile()
  }

  function saveProfile(partial: Partial<FinancialProfile>) {
    profile.value = { ...profile.value, ...partial }
    writeStoredProfile(profile.value)
  }

  function resetProfile() {
    profile.value = { ...defaultProfile }
    writeStoredProfile(profile.value)
  }

  async function syncProfileToApi(source: FinancialProfile) {
    if (demoMode.value) return

    const body: ProfilePatchRequest = {
      active_income: Math.max(0, source.active_income),
      passive_income: Math.max(0, source.passive_income),
      emergency_fund: Math.max(0, source.emergency_fund),
      emergency_breakdown: source.emergency_breakdown,
      fixed_expenses: source.fixed_expenses ?? []
    }

    try {
      await apiFetch('/users/me/profile', { method: 'PATCH', body })
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
    markOnboardingCompleteOnApi
  }
}
