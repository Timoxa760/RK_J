import type { ChallengeItem, LeaderboardEntry } from '~/types/api'
import { useAuthStore } from '~/store/authStore'

const mockChallenges: ChallengeItem[] = [
  { id: '1', title: 'Неделя без доставки', participants: 124 },
  { id: '2', title: 'Кофе ≤ 3 раза', participants: 89 }
]

const mockLeaderboard: LeaderboardEntry[] = [
  { id: '1', display_name: 'Участник A', relative_score: 95, rank: 1 },
  { id: '2', display_name: 'Участник B', relative_score: 88, rank: 2 },
  { id: '3', display_name: 'Вы', relative_score: 76, rank: 3 }
]

export function useSocial() {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const challenges = ref<ChallengeItem[]>([])
  const leaderboard = ref<LeaderboardEntry[]>([])
  const loading = ref(false)
  const newTitle = ref('')

  async function fetchJson<T>(path: string, mock: T): Promise<T> {
    if (config.public.demoMode) return mock
    try {
      return await $fetch<T>(path, {
        baseURL: config.public.apiBase,
        headers: authStore.token
          ? { Authorization: `Bearer ${authStore.token}` }
          : undefined
      })
    } catch {
      return mock
    }
  }

  async function loadAll() {
    loading.value = true
    try {
      const c = await fetchJson<{ challenges: ChallengeItem[] }>('/api/v1/social/challenges', {
        challenges: mockChallenges
      })
      const l = await fetchJson<{ entries: LeaderboardEntry[] }>('/api/v1/social/leaderboard', {
        entries: mockLeaderboard
      })
      challenges.value = c.challenges
      leaderboard.value = l.entries
    } finally {
      loading.value = false
    }
  }

  async function createChallenge(title: string) {
    if (!title.trim()) return
    challenges.value = [{ id: String(Date.now()), title, participants: 1 }, ...challenges.value]
    newTitle.value = ''
    await loadAll()
  }

  return { challenges, leaderboard, loading, newTitle, loadAll, createChallenge }
}
