<script setup lang="ts">
import { buildSocialPageNarrative } from '~/utils/pageNarrative'
import { buildHabitIndex } from '~/utils/habitIndex'
import { mockStores } from '~/store/mocks'
import { normalizeStores } from '~/utils/apiNormalize'
import type { StoresResponse } from '~/types/api'

const {
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
} = useSocial()

const { apiFetchWithDemo, demoMode } = useApi()
const stores = ref<StoresResponse | null>(null)
const storesLoading = ref(false)

const habitIndex = computed(() => buildHabitIndex(stores.value))

const pageNarrative = computed(() =>
  buildSocialPageNarrative({ habitIndex: habitIndex.value })
)

const pageLoading = computed(() => loading.value && !challenges.value.length)

const inviteUrl = computed(() =>
  inviteToken.value ? buildInviteUrl(inviteToken.value) : ''
)

const selectedChallenge = computed(() =>
  challenges.value.find((c) => c.id === selectedChallengeId.value) ?? null
)

async function loadStores() {
  storesLoading.value = true
  try {
    const raw = await apiFetchWithDemo('/dashboard/stores', mockStores)
    stores.value = normalizeStores(raw)
  } catch {
    if (demoMode.value) {
      stores.value = normalizeStores(mockStores)
    }
  } finally {
    storesLoading.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadChallenges(), loadStores()])
})
</script>

<template>
  <div class="mx-auto w-full max-w-4xl space-y-6">
    <Alert v-if="error" variant="destructive">
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <SharedPageNarrative :narrative="pageNarrative" :loading="pageLoading" />

    <SocialHabitIndexCard :habit-index="habitIndex" :loading="storesLoading" />

    <template v-if="!pageLoading">
      <SocialChallengesSection
        v-model:new-title="newTitle"
        v-model:new-type="newType"
        v-model:new-duration="newDuration"
        :challenges="challenges"
        :selected-challenge-id="selectedChallengeId"
        :invite-url="inviteUrl"
        :invite-copied="inviteCopied"
        :loading="loading"
        @create="createChallenge"
        @select="selectChallenge"
        @copy-invite="copyInvite"
      />

      <SocialLeaderboardSection
        :leaderboard="leaderboard"
        :loading="loading && !leaderboard.length"
        :challenge-title="selectedChallenge?.title"
      />
    </template>
  </div>
</template>
