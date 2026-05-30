<script setup lang="ts">
import { PROFILE } from '~/constants/productCopy'
import { buildProfilePageNarrative } from '~/utils/pageNarrative'
import { needsOnboarding, resetOnboardingForCurrentUser } from '~/composables/useOnboarding'

const { profile, totalIncome, loadProfile } = useFinancialProfile()
const { goals, loading: goalsLoading, fetchGoals } = useGoals()

const pageLoading = computed(() => goalsLoading.value)

const showSurveyPrompt = computed(() => needsOnboarding())

const profileIncomplete = computed(
  () => totalIncome.value <= 0 && profile.value.emergency_fund <= 0
)

async function retakeSurvey() {
  resetOnboardingForCurrentUser()
  await navigateTo('/onboarding')
}

const pageNarrative = computed(() =>
  buildProfilePageNarrative({
    profile: profile.value,
    goals: goals.value
  })
)

onMounted(() => {
  loadProfile()
  void fetchGoals()
})
</script>

<template>
  <div class="mx-auto w-full max-w-4xl space-y-6">
    <SharedPageNarrative :narrative="pageNarrative" :loading="pageLoading && !goals.length" />

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

    <Alert v-else-if="profileIncomplete && !goals.length">
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
    <ProfileGoalsSection />
    <ProfileProvidersSection />
  </div>
</template>
