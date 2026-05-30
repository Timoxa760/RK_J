<script setup lang="ts">
import {
  ONBOARDING_STEP_LABELS,
  ONBOARDING_WIZARD_STEP_COUNT
} from '~/constants/onboardingSteps'
import { surveyStepToSkipKey } from '~/constants/onboardingSkips'

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
  skipSurveyStep,
  selectSurveyMode,
  completeVoiceSurvey,
  syncProfileToApi,
  nextStep,
  prevStep,
  completeOnboarding
} = useOnboarding()

const { diagnosis, loading: diagnosisLoading, loadDiagnosis } = useOnboardingDiagnosis()

const expenseAdded = ref(false)
const voiceProgress = ref<{ current: number; total: number } | null>(null)

onMounted(() => {
  hydrate()
})

watch(step, async (s) => {
  if (s !== 7) return
  await syncProfileToApi()
  await loadDiagnosis()
})

const progressView = computed(() => {
  if (step.value <= 1) return null

  const total = ONBOARDING_WIZARD_STEP_COUNT - 2

  if (step.value === 2) {
    return { current: 1, total, label: ONBOARDING_STEP_LABELS[1] }
  }

  if (draft.value.survey_input_mode === 'voice') {
    if (step.value === 3 && voiceProgress.value) {
      return {
        current: 1 + voiceProgress.value.current,
        total,
        label: ONBOARDING_STEP_LABELS[step.value - 1]
      }
    }
    if (step.value === 7) {
      return { current: total - 1, total, label: ONBOARDING_STEP_LABELS[6] }
    }
    if (step.value === 8) {
      return { current: total, total, label: ONBOARDING_STEP_LABELS[7] }
    }
  }

  if (step.value >= 3 && step.value <= 8) {
    return {
      current: Math.min(step.value - 1, total),
      total,
      label: ONBOARDING_STEP_LABELS[step.value - 1]
    }
  }

  return null
})

const textSurveyStep = computed(() => step.value - 2)

function onSkipTextStep(surveyStep: number) {
  const key = surveyStepToSkipKey(surveyStep)
  if (key) skipSurveyStep(key)
}

function onVoiceProgress(current: number, total: number) {
  voiceProgress.value = { current, total }
}
</script>

<template>
  <OnboardingShell
    :progress="progressView"
    :show-progress="step > 1 && step < ONBOARDING_WIZARD_STEP_COUNT"
    :minimal-header="step === 1"
    :step-labels="[...ONBOARDING_STEP_LABELS].slice(1)"
  >
    <div class="mm-onboarding-survey">
      <Alert v-if="error" variant="destructive" class="mb-4">
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>

      <div :key="`${step}-${draft.survey_input_mode ?? 'none'}`" class="mm-onboarding-survey__step mm-onboarding-survey__step--enter">
        <OnboardingWelcome v-if="step === 1" @start="nextStep" />

        <OnboardingIntro
          v-else-if="step === 2"
          @back="prevStep"
          @select-voice="selectSurveyMode('voice')"
          @select-text="selectSurveyMode('text')"
        />

        <OnboardingVoiceFlow
          v-else-if="step === 3 && draft.survey_input_mode === 'voice'"
          :draft="draft"
          @back="prevStep"
          @complete="completeVoiceSurvey"
          @patch="patchDraft"
          @progress="onVoiceProgress"
        />

        <OnboardingTextFlow
          v-else-if="step >= 3 && step <= 6 && draft.survey_input_mode === 'text'"
          :draft="draft"
          :step="textSurveyStep"
          @back="prevStep"
          @next="nextStep"
          @patch="patchDraft"
          @set-goal-kind="setGoalKind"
          @patch-breakdown="patchBreakdown"
          @add-fixed-expense="addFixedExpense"
          @update-fixed-expense="(i, p) => updateFixedExpense(i, p)"
          @remove-fixed-expense="removeFixedExpense"
          @skip="onSkipTextStep(textSurveyStep)"
        />

        <OnboardingSummary
          v-else-if="step === 7"
          :draft="draft"
          :summary="summary"
          :diagnosis="diagnosis"
          :diagnosis-loading="diagnosisLoading"
          @back="prevStep"
          @next="nextStep"
        />

        <OnboardingFirstAction
          v-else-if="step === 8"
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
