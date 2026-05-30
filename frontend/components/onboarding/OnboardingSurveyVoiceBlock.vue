<script setup lang="ts">
import type { OnboardingDraft } from '~/types/api'
import { getVoiceQuestionForSurveyStep } from '~/constants/onboardingSurvey'

const props = defineProps<{
  surveyStep: number
  draft: OnboardingDraft
}>()

const emit = defineEmits<{
  patch: [partial: Partial<OnboardingDraft>]
}>()

const { parsing, parseStep } = useOnboardingParse()
const { transcribing, lastTranscript, transcribeAudio, clearLastTranscript } = useVoiceTranscribe()
const parseError = ref('')

const question = computed(() => getVoiceQuestionForSurveyStep(props.surveyStep))

const {
  status,
  errorMessage,
  supported,
  start,
  stop,
  reset
} = useAudioRecorder()

const busy = computed(() => parsing.value || transcribing.value)

async function toggleRecording() {
  if (busy.value) return
  if (status.value === 'recording') {
    const audio = await stop()
    if (!audio || audio.size === 0) return
    parseError.value = ''
    try {
      const text = await transcribeAudio(audio, question.value?.examples[0])
      await submitAnswer(text)
    } catch (e) {
      clearLastTranscript()
      const status = (e as { statusCode?: number })?.statusCode
      if (status === 503) {
        parseError.value =
          'Сервис распознавания речи недоступен. Запустите Whisper или выберите пример ниже.'
      } else {
        parseError.value = e instanceof Error ? e.message : 'Не удалось распознать голос'
      }
    }
    return
  }
  parseError.value = ''
  clearLastTranscript()
  reset()
  await start()
}

async function submitAnswer(rawText: string) {
  if (!question.value) return
  const text = rawText.trim()
  if (!text) {
    parseError.value = 'Не расслышали — попробуйте ещё раз или выберите пример.'
    return
  }

  parseError.value = ''
  try {
    const patch = await parseStep(question.value.id, text)
    if (!isMeaningfulPatch(question.value.id, patch)) {
      parseError.value = 'Не удалось разобрать ответ. Попробуйте иначе или введите вручную.'
      return
    }
    emit('patch', patch)
  } catch {
    parseError.value = 'Не удалось разобрать ответ. Повторите или введите вручную.'
  }
}

function isMeaningfulPatch(
  step: string,
  patch: Partial<OnboardingDraft>
): boolean {
  if (!patch || Object.keys(patch).length === 0) return false
  switch (step) {
    case 'income':
      return (patch.active_income ?? 0) > 0 || (patch.passive_income ?? 0) > 0
    case 'cushion':
      return (patch.emergency_fund ?? 0) > 0
    case 'goal':
      return (patch.goal_amount ?? 0) >= 1000
    case 'expenses':
      return Boolean(patch.skipped_expenses) || (patch.fixed_expenses?.length ?? 0) > 0
    default:
      return Object.keys(patch).length > 0
  }
}

function pickExample(example: string) {
  reset()
  void submitAnswer(example)
}

const resultHint = computed(() => formatVoiceResult(props.surveyStep, props.draft))

const displayTranscript = computed(() => {
  if (lastTranscript.value) return lastTranscript.value
  return ''
})

function formatRub(n: number) {
  return n.toLocaleString('ru-RU')
}

function formatVoiceResult(step: number, draft: OnboardingDraft): string | null {
  switch (step) {
    case 1:
      if (draft.active_income <= 0 && draft.passive_income <= 0) return null
      return `Основной доход ${formatRub(draft.active_income)} ₽/мес · ещё ${formatRub(draft.passive_income)} ₽/мес`
    case 2: {
      if (draft.emergency_fund <= 0) return null
      const b = draft.emergency_breakdown
      const parts: string[] = []
      if (b.cash > 0) parts.push(`наличные ${formatRub(b.cash)} ₽`)
      if (b.deposit > 0) parts.push(`вклад ${formatRub(b.deposit)} ₽`)
      if (b.investments > 0) parts.push(`инвестиции ${formatRub(b.investments)} ₽`)
      if (parts.length > 0) {
        return `Запас ${formatRub(draft.emergency_fund)} ₽ (${parts.join(' · ')})`
      }
      return `Запас ${formatRub(draft.emergency_fund)} ₽`
    }
    case 3:
      if (draft.goal_amount < 1000) return null
      return `${draft.goal_title} — ${formatRub(draft.goal_amount)} ₽`
    case 4:
      if (draft.skipped_expenses) return 'Постоянные платежи пропущены — учтём ваши обычные траты'
      if (!draft.fixed_expenses.length) return null
      return draft.fixed_expenses
        .map((e) => `${e.title || 'Платёж'} ${formatRub(e.amount)} ₽`)
        .join(' · ')
    default:
      return null
  }
}
</script>

<template>
  <div v-if="question" class="mm-onb-survey-voice space-y-4">
    <p class="text-sm leading-relaxed text-[color:var(--mm-text)]">
      {{ question.prompt }}
    </p>

    <div class="flex justify-center py-1">
      <button
        type="button"
        class="mm-onb-mic-orb-hit border-0 bg-transparent p-0"
        :class="{ 'mm-onb-mic-orb-hit--listen': status === 'recording' && !busy }"
        :disabled="!supported || busy"
        :aria-pressed="status === 'recording'"
        :aria-label="status === 'recording' ? 'Остановить и отправить' : 'Нажмите и ответьте вслух'"
        @click="toggleRecording"
      >
        <OnboardingMicOrbVisual
          :listening="status === 'recording'"
          :parsing="busy"
          compact
          :ambient="status !== 'recording' && !busy"
        />
      </button>
    </div>

    <p class="text-center text-xs text-[color:var(--mm-text-muted)]">
      {{
        busy
          ? 'Распознаём…'
          : status === 'recording'
            ? 'Поток слушает — нажмите ещё раз, чтобы отправить'
            : supported
              ? 'Нажмите на орб, скажите ответ и нажмите ещё раз'
              : 'Голос недоступен — выберите пример ниже'
      }}
    </p>

    <p v-if="errorMessage" class="text-center text-xs text-destructive">
      {{ errorMessage }}
    </p>

    <p v-if="displayTranscript && !parseError" class="mm-onb-transcript-bubble">
      «{{ displayTranscript }}»
    </p>

    <p v-if="parseError" class="text-center text-sm text-destructive">
      {{ parseError }}
    </p>

    <div
      v-if="resultHint && !parseError"
      class="rounded-xl border border-[color:var(--mm-primary)]/25 bg-[color:var(--mm-primary-soft)]/60 px-4 py-3 text-center text-sm text-[color:var(--mm-text)]"
    >
      <span class="font-medium text-[color:var(--mm-primary)]">Записали: </span>
      {{ resultHint }}
    </div>

    <div v-if="!supported || parseError" class="flex flex-wrap justify-center gap-2">
      <button
        v-for="example in question.examples"
        :key="example"
        type="button"
        class="mm-onb-chip"
        :disabled="busy"
        @click="pickExample(example)"
      >
        {{ example }}
      </button>
    </div>

    <p v-if="resultHint" class="text-center text-xs text-[color:var(--mm-text-muted)]">
      Можно перейти на «Вручную» и поправить цифры
    </p>
  </div>
</template>
