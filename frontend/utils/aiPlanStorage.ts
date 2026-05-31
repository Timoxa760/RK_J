import type { AiDiagnosisResponse } from '~/types/api'
import type { FinancialPlan } from '~/utils/financialPlan'
import { currentUserStorageKey } from '~/utils/userStorage'

const AI_PLAN_STORAGE_PREFIX = 'potok:ai-plan'

export interface StoredAiPlan {
  plan: FinancialPlan
  diagnosis: AiDiagnosisResponse | null
  loadedAt: number
}

export function readStoredAiPlan(): StoredAiPlan | null {
  if (!import.meta.client) return null
  try {
    const raw = localStorage.getItem(currentUserStorageKey(AI_PLAN_STORAGE_PREFIX))
    if (!raw) return null
    const parsed = JSON.parse(raw) as StoredAiPlan
    if (!parsed.plan || !parsed.loadedAt) return null
    return parsed
  } catch {
    return null
  }
}

export function writeStoredAiPlan(data: StoredAiPlan) {
  if (!import.meta.client) return
  localStorage.setItem(currentUserStorageKey(AI_PLAN_STORAGE_PREFIX), JSON.stringify(data))
}

export function clearStoredAiPlanForCurrentUser() {
  if (!import.meta.client) return
  localStorage.removeItem(currentUserStorageKey(AI_PLAN_STORAGE_PREFIX))
}
