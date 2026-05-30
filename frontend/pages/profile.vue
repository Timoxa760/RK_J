<script setup lang="ts">
import { PROFILE } from '~/constants/productCopy'
import { buildProfilePageNarrative } from '~/utils/pageNarrative'
import { needsOnboarding, resetOnboardingForCurrentUser } from '~/composables/useOnboarding'

const { profile, totalIncome, loadProfile } = useFinancialProfile()
const { primaryGoal } = useGoals()

const showSurveyPrompt = computed(() => needsOnboarding())

const profileIncomplete = computed(
  () =>
    totalIncome.value <= 0 &&
    profile.value.emergency_fund <= 0 &&
    !primaryGoal.value
)

async function retakeSurvey() {
  resetOnboardingForCurrentUser()
  await navigateTo('/onboarding')
}

const pageNarrative = computed(() =>
  buildProfilePageNarrative({
    profile: profile.value,
    goals: primaryGoal.value ? [primaryGoal.value] : []
  })
)

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <div class="mx-auto w-full max-w-4xl space-y-6">
    <SharedPageNarrative :narrative="pageNarrative" />

    <Alert v-if="showSurveyPrompt">
      <AlertTitle>{{ PROFILE.emptyModel }}</AlertTitle>
      <AlertDescription class="space-y-2">
        <p>Укажите доход, запас и цель — тогда прогноз на других экранах станет точнее.</p>
        <div class="flex flex-wrap gap-2">
          <Button size="sm" variant="secondary" @click="retakeSurvey">
            Пройти короткий опрос
          </Button>
          <AdvisorAskButton size="sm" variant="secondary" />
        </div>
      </AlertDescription>
    </Alert>

    <Alert v-else-if="profileIncomplete">
      <AlertTitle>{{ PROFILE.emptyModel }}</AlertTitle>
      <AlertDescription class="space-y-2">
        <p>Укажите доход, запас и цель — тогда прогноз на других экранах станет точнее.</p>
        <div class="flex flex-wrap gap-2">
          <Button as-child size="sm" variant="secondary">
            <NuxtLink to="/onboarding">Пройти короткий опрос</NuxtLink>
          </Button>
          <AdvisorAskButton size="sm" variant="secondary" />
        </div>
      </AlertDescription>
    </Alert>

    <ProfileFinancialForm />
    <ProfileProvidersSection />
  </div>
</template>
