import type { ChallengeItem, ChallengeType, LeaderboardEntry } from '~/types/api'
import { mockChallenges, mockLeaderboard } from '~/store/mocks'
import { normalizeLeaderboard } from '~/utils/apiNormalize'
import { currentUserStorageKey } from '~/utils/userStorage'

const CHALLENGES_PREFIX = 'potok:challenges'

function readStoredChallenges(): ChallengeItem[] {
  if (!import.meta.client) return []
  try {
    const raw = localStorage.getItem(currentUserStorageKey(CHALLENGES_PREFIX))
    if (!raw) return []
    const parsed = JSON.parse(raw) as ChallengeItem[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function writeStoredChallenges(challenges: ChallengeItem[]) {
  if (!import.meta.client) return
  localStorage.setItem(currentUserStorageKey(CHALLENGES_PREFIX), JSON.stringify(challenges))
}

export function useSocial() {
  const { apiFetch, apiFetchWithDemo, demoMode } = useApi()

  const challenges = ref<ChallengeItem[]>([])
  const leaderboard = ref<LeaderboardEntry[]>([])
  const selectedChallengeId = ref<string | null>(null)
  const inviteToken = ref<string | null>(null)
  const inviteCopied = ref(false)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const newTitle = ref('')
  const newType = ref<ChallengeType>('least_spend')
  const newDuration = ref(7)

  async function loadChallenges() {
    loading.value = true
    error.value = null
    try {
      const stored = readStoredChallenges()
      challenges.value = stored.length ? stored : mockChallenges
      if (!selectedChallengeId.value && challenges.value.length) {
        selectedChallengeId.value = challenges.value[0]?.id ?? null
      }
      const selected = challenges.value.find((c) => c.id === selectedChallengeId.value)
      inviteToken.value = selected?.invite_token ?? null
      await loadLeaderboard()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка загрузки'
      challenges.value = mockChallenges
      await loadLeaderboard()
    } finally {
      loading.value = false
    }
  }

  async function loadLeaderboard() {
    const id = selectedChallengeId.value ?? challenges.value[0]?.id
    if (!id) {
      leaderboard.value = mockLeaderboard
      return
    }
    try {
      const res = await apiFetchWithDemo(`/challenges/${id}/leaderboard`, {
        challenge_id: id,
        type: challenges.value.find((c) => c.id === id)?.type,
        leaderboard: mockLeaderboard.map((r) => ({
          position: r.position,
          username: r.username,
          avatar: r.avatar,
          relative_score: r.relative_score
        }))
      })
      leaderboard.value = normalizeLeaderboard(res)
    } catch {
      leaderboard.value = mockLeaderboard
    }
  }

  async function createChallenge() {
    const title = newTitle.value.trim()
    if (!title) return
    error.value = null
    try {
      if (demoMode.value) {
        const item: ChallengeItem = {
          id: String(Date.now()),
          title,
          type: newType.value,
          participants: 1,
          status: 'active',
          invite_token: `invite-${Date.now()}`
        }
        challenges.value = [item, ...challenges.value]
        inviteToken.value = item.invite_token ?? null
        selectedChallengeId.value = item.id
        newTitle.value = ''
        await loadLeaderboard()
        return
      }
      const created = await apiFetch<ChallengeItem>('/challenges', {
        method: 'POST',
        body: {
          type: newType.value,
          title,
          duration_days: newDuration.value,
          max_participants: 10
        }
      })
      challenges.value = [created, ...challenges.value]
      writeStoredChallenges(challenges.value)
      inviteToken.value = created.invite_token ?? null
      selectedChallengeId.value = created.id
      newTitle.value = ''
      await loadLeaderboard()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось создать задание'
    }
  }

  async function selectChallenge(id: string) {
    selectedChallengeId.value = id
    const ch = challenges.value.find((c) => c.id === id)
    inviteToken.value = ch?.invite_token ?? null
    inviteCopied.value = false
    await loadLeaderboard()
  }

  function buildInviteUrl(token: string) {
    if (import.meta.client) {
      return `${window.location.origin}/social?invite=${encodeURIComponent(token)}`
    }
    return `/social?invite=${encodeURIComponent(token)}`
  }

  async function copyInvite() {
    if (!inviteToken.value) return
    const url = buildInviteUrl(inviteToken.value)
    try {
      await navigator.clipboard.writeText(url)
      inviteCopied.value = true
      setTimeout(() => {
        inviteCopied.value = false
      }, 2000)
    } catch {
      error.value = 'Не удалось скопировать ссылку'
    }
  }

  return {
    challenges,
    leaderboard,
    selectedChallengeId,
    inviteToken,
    inviteCopied,
    loading,
    error,
    newTitle,
    newType,
    newDuration,
    loadChallenges,
    createChallenge,
    selectChallenge,
    copyInvite,
    buildInviteUrl
  }
}
