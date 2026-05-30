<script setup lang="ts">
import { ONBOARDING_SKIP_HINTS, ONBOARDING_SKIP_LABEL, skipPatchForSurveyStep } from '~/constants/onboardingSkips'
import { ONBOARDING_VOICE_QUESTIONS } from '~/constants/onboardingSurvey'
import type { OnboardingDraft } from '~/types/api'

defineProps<{
  draft: OnboardingDraft
}>()

const emit = defineEmits<{
  back: []
  complete: []
  patch: [partial: Partial<OnboardingDraft>]
  progress: [current: number, total: number]
}>()

const questionIndex = ref(0)
const history = ref<{ title: string; answer: string }[]>([])

const currentQuestion = computed(() => ONBOARDING_VOICE_QUESTIONS[questionIndex.value]!)
const surveyStep = computed(() => questionIndex.value + 1)
const isLast = computed(() => questionIndex.value >= ONBOARDING_VOICE_QUESTIONS.length - 1)
const stepLabel = computed(
  () => `Вопрос ${questionIndex.value + 1} из ${ONBOARDING_VOICE_QUESTIONS.length}`
)

function syncProgress() {
  emit('progress', questionIndex.value + 1, ONBOARDING_VOICE_QUESTIONS.length)
}

function onPatch(partial: Partial<OnboardingDraft>) {
  if (partial.emergency_breakdown) {
    const b = partial.emergency_breakdown
    emit('patch', {
      ...partial,
      emergency_fund: partial.emergency_fund ?? b.cash + b.deposit + b.investments
    })
  } else {
    emit('patch', partial)
  }
  if (partial.goal_kind) {
    emit('patch', {
      goal_kind: partial.goal_kind,
      goal_title: partial.goal_title
    })
  }
}

function onAnswered(answer: string) {
  history.value.push({ title: currentQuestion.value.title, answer })
  if (isLast.value) {
    setTimeout(() => emit('complete'), 350)
    return
  }
  questionIndex.value += 1
  syncProgress()
}

function skipQuestion() {
  const stepId = currentQuestion.value.id
  emit('patch', skipPatchForSurveyStep(stepId))
  history.value.push({ title: currentQuestion.value.title, answer: 'Пропущено' })

  if (isLast.value) {
    setTimeout(() => emit('complete'), 350)
    return
  }

  questionIndex.value += 1
  syncProgress()
}

onMounted(() => {
  syncProgress()
})
</script>

<template>
  <OnboardingStepShell
    :title="currentQuestion.title"
    :description="stepLabel"
    show-back
    hide-next
    :secondary-action="{ label: ONBOARDING_SKIP_LABEL }"
    :skip-hint="ONBOARDING_SKIP_HINTS[currentQuestion.id]"
    @back="emit('back')"
    @secondary="skipQuestion"
  >
    <div v-if="history.length" class="mm-onb-timeline mb-6">
      <div
        v-for="(item, i) in history"
        :key="i"
        class="mm-onb-timeline__item"
      >
        <span class="mm-onb-timeline__dot" aria-hidden="true">✓</span>
        <div class="mm-onb-timeline__body">
          <p class="mm-onb-timeline__title">{{ item.title }}</p>
          <p class="mm-onb-timeline__answer">{{ item.answer }}</p>
        </div>
      </div>
    </div>

    <OnboardingSurveyVoiceBlock
      :survey-step="surveyStep"
      :draft="draft"
      @patch="onPatch"
      @answered="onAnswered"
    />
  </OnboardingStepShell>
</template>
