<script setup lang="ts">
import {
  ONBOARDING_STEP_LABELS,
  ONBOARDING_WIZARD_STEP_COUNT
} from '~/constants/onboardingSteps'

const {
  step,
  draft,
  summary,
  finishing,
  error,
  hydrate,
  patchDraft,
  setGoalKind,
  patchBreakdown,
  addFixedExpense,
  updateFixedExpense,
  removeFixedExpense,
  skipExpenses,
  nextStep,
  prevStep,
  completeOnboarding
} = useOnboarding()

const expenseAdded = ref(false)

onMounted(() => {
  hydrate()
})

const progressView = computed(() => {
  if (step.value <= 1) return null
  return {
    current: step.value - 1,
    total: ONBOARDING_WIZARD_STEP_COUNT - 1,
    label: ONBOARDING_STEP_LABELS[step.value - 1]
  }
})

const surveyStep = computed(() => step.value - 1)

function onSurveySecondary() {
  skipExpenses()
  nextStep()
}
</script>

<template>
  <OnboardingShell
    :progress="progressView"
    :show-progress="step > 1"
    :minimal-header="step === 1"
    :step-labels="[...ONBOARDING_STEP_LABELS].slice(1)"
  >
    <div class="mm-onboarding-survey">
      <Alert v-if="error" variant="destructive" class="mb-4">
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>

      <div :key="step" class="mm-onboarding-survey__step mm-onboarding-survey__step--enter">
        <OnboardingWelcome
          v-if="step === 1"
          @start="nextStep"
        />

        <OnboardingTextFlow
          v-else-if="step >= 2 && step <= 5"
          :draft="draft"
          :step="surveyStep"
          @back="prevStep"
          @next="nextStep"
          @patch="patchDraft"
          @set-goal-kind="setGoalKind"
          @patch-breakdown="patchBreakdown"
          @add-fixed-expense="addFixedExpense"
          @update-fixed-expense="(i, p) => updateFixedExpense(i, p)"
          @remove-fixed-expense="removeFixedExpense"
          @skip-expenses="onSurveySecondary"
        />

        <OnboardingSummary
          v-else-if="step === 6"
          :draft="draft"
          :summary="summary"
          @back="prevStep"
          @next="nextStep"
        />

        <OnboardingFirstAction
          v-else-if="step === 7"
          :finishing="finishing"
          :expense-added="expenseAdded"
          @back="prevStep"
          @finish="completeOnboarding"
          @added="expenseAdded = true"
        />
      </div>
    </div>
  </OnboardingShell>
</template>
