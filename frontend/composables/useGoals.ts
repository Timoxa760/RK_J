import type { CreateGoalRequest, Goal, GoalsListResponse } from '~/types/api'
import { mockGoals } from '~/store/mocks'
import { currentUserStorageKey } from '~/utils/userStorage'

const GOALS_PREFIX = 'potok:goals'

function readStoredGoals(): Goal[] {
  if (!import.meta.client) return []
  try {
    const raw = localStorage.getItem(currentUserStorageKey(GOALS_PREFIX))
    if (!raw) return []
    const parsed = JSON.parse(raw) as Goal[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function writeStoredGoals(goals: Goal[]) {
  if (!import.meta.client) return
  localStorage.setItem(currentUserStorageKey(GOALS_PREFIX), JSON.stringify(goals))
}

function mergeGoals(remote: Goal[], local: Goal[]): Goal[] {
  const map = new Map<string, Goal>()
  for (const g of local) map.set(g.id, g)
  for (const g of remote) map.set(g.id, g)
  return [...map.values()]
}

function withProgress(goal: Goal): Goal {
  const progress =
    goal.target_amount > 0
      ? Math.min(100, Math.round((goal.current_amount / goal.target_amount) * 100))
      : 0
  return { ...goal, progress_percent: goal.progress_percent ?? progress }
}

export function useGoals() {
  const { apiFetch, apiFetchWithDemo, demoMode } = useApi()

  const goals = ref<Goal[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const error = ref<string | null>(null)

  const primaryGoal = computed(() => goals.value[0] ?? null)

  function applyGoals(next: Goal[]) {
    goals.value = next.map(withProgress)
    writeStoredGoals(goals.value)
  }

  async function fetchGoals() {
    loading.value = true
    error.value = null
    try {
      if (demoMode.value) {
        const remote = await apiFetchWithDemo<GoalsListResponse>('/goals', { goals: mockGoals })
        const merged = mergeGoals(remote.goals ?? mockGoals, readStoredGoals())
        applyGoals(merged.length ? merged : mockGoals)
      } else {
        const remote = await apiFetch<GoalsListResponse>('/goals')
        const merged = mergeGoals(remote.goals ?? [], readStoredGoals())
        applyGoals(merged)
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка загрузки целей'
      const stored = readStoredGoals()
      applyGoals(demoMode.value ? (stored.length ? stored : mockGoals) : stored)
    } finally {
      loading.value = false
    }
  }

  async function createGoal(payload: CreateGoalRequest) {
    saving.value = true
    error.value = null
    try {
      const created = await apiFetch<Goal>('/goals', { method: 'POST', body: payload })
      applyGoals([withProgress(created), ...goals.value])
      return created
    } catch {
      const localGoal: Goal = withProgress({
        id: `local-${Date.now()}`,
        title: payload.title,
        target_amount: payload.target_amount,
        current_amount: 0,
        progress_percent: 0,
        target_date: payload.target_date,
        auto_save_percent: payload.auto_save_percent
      })
      applyGoals([localGoal, ...goals.value])
      return localGoal
    } finally {
      saving.value = false
    }
  }

  function monthsToGoal(goal: Goal, monthlySaving: number): number | null {
    const remaining = goal.target_amount - goal.current_amount
    if (remaining <= 0 || monthlySaving <= 0) return null
    return Math.ceil(remaining / monthlySaving)
  }

  return {
    goals,
    primaryGoal,
    loading,
    saving,
    error,
    fetchGoals,
    createGoal,
    monthsToGoal
  }
}
