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
const parseError = ref('')
const lastAnswer = ref('')

const question = computed(() => getVoiceQuestionForSurveyStep(props.surveyStep))

const { transcript, listening, speechSupported, toggleListen, clearTranscript } =
  useOnboardingSpeech((text) => {
    void submitAnswer(text)
  })

async function submitAnswer(rawText: string) {
  if (!question.value || parsing.value) return
  const text = rawText.trim()
  if (!text) return

  parseError.value = ''
  lastAnswer.value = text
  try {
    const patch = await parseStep(question.value.id, text)
    if (Object.keys(patch).length === 0) {
      parseError.value = 'Не удалось разобрать ответ. Попробуйте иначе или введите вручную.'
      return
    }
    emit('patch', patch)
  } catch {
    parseError.value = 'Не удалось разобрать ответ. Повторите или введите вручную.'
  }
}

function pickExample(example: string) {
  clearTranscript()
  void submitAnswer(example)
}

const resultHint = computed(() => formatVoiceResult(props.surveyStep, props.draft))

function formatRub(n: number) {
  return n.toLocaleString('ru-RU')
}

function formatVoiceResult(step: number, draft: OnboardingDraft): string | null {
  switch (step) {
    case 1:
      if (draft.active_income <= 0 && draft.passive_income <= 0) return null
      return `Основной доход ${formatRub(draft.active_income)} ₽/мес · ещё ${formatRub(draft.passive_income)} ₽/мес`
    case 2:
      if (draft.emergency_fund <= 0) return null
      return `Запас ${formatRub(draft.emergency_fund)} ₽`
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
        class="mm-onb-mic-orb-hit mm-onb-mic-orb-hit--preview border-0 bg-transparent p-0"
        :disabled="!speechSupported || parsing"
        :aria-pressed="listening"
        :aria-label="listening ? 'Остановить запись' : 'Нажмите и ответьте вслух'"
        @click="toggleListen(parsing)"
      >
        <OnboardingMicOrbVisual
          :listening="listening"
          :parsing="parsing"
          compact
          :ambient="!listening && !parsing"
        />
      </button>
    </div>

    <p class="text-center text-xs text-[color:var(--mm-text-muted)]">
      {{
        parsing
          ? 'Записываем…'
          : listening
            ? 'Поток слушает'
            : speechSupported
              ? 'Нажмите на орб и расскажите'
              : 'Голос недоступен — выберите пример ниже'
      }}
    </p>

    <p v-if="listening && transcript" class="mm-onb-transcript-bubble">
      «{{ transcript }}»
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

    <div v-if="!speechSupported || parseError" class="flex flex-wrap justify-center gap-2">
      <button
        v-for="example in question.examples"
        :key="example"
        type="button"
        class="mm-onb-chip"
        :disabled="parsing"
        @click="pickExample(example)"
      >
        {{ example }}
      </button>
    </div>

    <p v-if="lastAnswer && resultHint" class="text-center text-xs text-[color:var(--mm-text-muted)]">
      Можно перейти на «Вручную» и поправить цифры
    </p>
  </div>
</template>
