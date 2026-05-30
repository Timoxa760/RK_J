import type { FinancialProfile, Goal } from '~/types/api'

/** Цель из профиля (без goal-service). */
export function goalFromProfile(profile: FinancialProfile): Goal | null {
  if (profile.skipped_goal || (profile.goal_amount ?? 0) < 1000) {
    return null
  }
  const title = profile.goal_title?.trim() || 'Финансовая цель'
  const target = profile.goal_amount ?? 0
  return {
    id: 'profile-goal',
    title,
    target_amount: target,
    current_amount: 0,
    progress_percent: 0
  }
}

/** @deprecated Используйте goalFromProfile / useProfileGoal. Совместимость API. */
export function useGoals() {
  const { profile, loadProfile } = useFinancialProfile()

  const primaryGoal = computed(() => goalFromProfile(profile.value))
  const goals = computed(() => (primaryGoal.value ? [primaryGoal.value] : []))
  const loading = ref(false)
  const saving = ref(false)
  const error = ref<string | null>(null)

  async function fetchGoals() {
    loadProfile()
  }

  async function createGoal(_payload: unknown) {
    /* цель сохраняется через PATCH profile / онбординг */
  }

  return {
    goals,
    primaryGoal,
    loading,
    saving,
    error,
    fetchGoals,
    createGoal
  }
}

export function useProfileGoal() {
  const { profile } = useFinancialProfile()
  return computed(() => goalFromProfile(profile.value))
}
