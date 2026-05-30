<script setup lang="ts">
import { ONBOARDING_SKIP_HINTS, ONBOARDING_SKIP_LABEL, skipPatchForSurveyStep } from '~/constants/onboardingSkips'
import { ONBOARDING_VOICE_QUESTIONS } from '~/constants/onboardingSurvey'
import type { OnboardingDraft, OnboardingParseStep } from '~/types/api'

defineProps<{
  draft: OnboardingDraft
}>()

const emit = defineEmits<{
  back: []
  complete: []
  patch: [partial: Partial<OnboardingDraft>]
  progress: [current: number, total: number]
}>()

const { parsing, parseStep } = useOnboardingParse()

const questionIndex = ref(0)
const transcript = ref('')
const listening = ref(false)
const speechSupported = ref(false)
const history = ref<{ title: string; answer: string }[]>([])
const parseError = ref('')

type SpeechRecognitionCtor = new () => {
  lang: string
  interimResults: boolean
  continuous: boolean
  start: () => void
  stop: () => void
  abort: () => void
  onresult: ((event: { results: { [i: number]: { [j: number]: { transcript: string } } } }) => void) | null
  onend: (() => void) | null
  onerror: (() => void) | null
}

let recognition: InstanceType<SpeechRecognitionCtor> | null = null

const currentQuestion = computed(() => ONBOARDING_VOICE_QUESTIONS[questionIndex.value]!)
const isLast = computed(() => questionIndex.value >= ONBOARDING_VOICE_QUESTIONS.length - 1)
const stepLabel = computed(
  () => `Вопрос ${questionIndex.value + 1} из ${ONBOARDING_VOICE_QUESTIONS.length}`
)

function syncProgress() {
  emit('progress', questionIndex.value + 1, ONBOARDING_VOICE_QUESTIONS.length)
}

async function submitAnswer(rawText: string) {
  const text = rawText.trim()
  if (!text || parsing.value) return

  parseError.value = ''
  const step = currentQuestion.value.id as OnboardingParseStep
  try {
    const patch = await parseStep(step, text)
    emit('patch', patch)
  } catch {
    parseError.value = 'Не удалось разобрать ответ. Повторите вслух или выберите вариант ниже.'
    return
  }

  history.value.push({ title: currentQuestion.value.title, answer: text })
  transcript.value = ''

  if (isLast.value) {
    setTimeout(() => emit('complete'), 350)
    return
  }

  questionIndex.value += 1
  syncProgress()
}

function toggleListen() {
  if (!recognition || parsing.value) return
  if (listening.value) {
    recognition.stop()
    return
  }
  parseError.value = ''
  transcript.value = ''
  listening.value = true
  recognition.start()
}

function pickChip(example: string) {
  submitAnswer(example)
}

function skipQuestion() {
  const stepId = currentQuestion.value.id
  emit('patch', skipPatchForSurveyStep(stepId))
  history.value.push({ title: currentQuestion.value.title, answer: 'Пропущено' })
  transcript.value = ''
  parseError.value = ''

  if (isLast.value) {
    setTimeout(() => emit('complete'), 350)
    return
  }

  questionIndex.value += 1
  syncProgress()
}

onMounted(() => {
  syncProgress()

  if (!import.meta.client) return
  const win = window as Window & { webkitSpeechRecognition?: SpeechRecognitionCtor }
  const SR = win.SpeechRecognition ?? win.webkitSpeechRecognition
  speechSupported.value = Boolean(SR)
  if (SR) {
    recognition = new SR()
    recognition.lang = 'ru-RU'
    recognition.interimResults = true
    recognition.continuous = false
    recognition.onresult = (event) => {
      let chunk = ''
      for (let i = 0; i < event.results.length; i++) {
        chunk += event.results[i]?.[0]?.transcript ?? ''
      }
      transcript.value = chunk.trim()
    }
    recognition.onend = () => {
      const wasListening = listening.value
      listening.value = false
      if (!wasListening) return
      const said = transcript.value.trim()
      if (said) submitAnswer(said)
    }
    recognition.onerror = () => {
      listening.value = false
    }
  }
})

onUnmounted(() => {
  recognition?.abort()
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

    <div class="mm-onboarding-voice__stage">
      <p class="mm-onboarding-voice__prompt">
        {{ currentQuestion.prompt }}
      </p>

      <div v-if="speechSupported" class="mm-onboarding-voice__hero">
        <OnboardingMicOrb
          :listening="listening"
          :parsing="parsing"
          label="Нажмите микрофон и ответьте вслух"
          @click="toggleListen"
        />
      </div>

      <p
        v-else
        class="max-w-sm text-sm text-[color:var(--mm-text-muted)]"
      >
        Голос недоступен в этом браузере — выберите готовый ответ:
      </p>

      <p
        v-if="listening && transcript"
        class="mm-onb-transcript-bubble"
      >
        «{{ transcript }}»
      </p>

      <p v-if="parseError" class="text-sm text-destructive">{{ parseError }}</p>
    </div>

    <div
      v-if="!speechSupported || parseError"
      class="mm-onboarding-voice__chips"
    >
      <button
        v-for="example in currentQuestion.examples"
        :key="example"
        type="button"
        class="mm-onb-chip"
        :disabled="parsing"
        @click="pickChip(example)"
      >
        {{ example }}
      </button>
    </div>
  </OnboardingStepShell>
</template>
