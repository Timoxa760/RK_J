<script setup lang="ts">
const PURCHASE_VOICE_EXAMPLES = [
  'купил продукты на 3500',
  'кофе 450 рублей',
  'кроссовки за 8900'
]

const emit = defineEmits<{
  done: []
}>()

const {
  status,
  errorMessage,
  supported,
  start,
  stop,
  reset
} = useAudioRecorder()

const submitError = ref('')
const parseError = ref('')
const { transcribing, lastTranscript, transcribeAudio, clearLastTranscript } = useVoiceTranscribe()
const { submitting, submitVoiceTranscript, lastResult } = useReceiptSubmit()

const busy = computed(() => submitting.value || transcribing.value)

const addedHint = computed(() => {
  const res = lastResult.value
  if (!res) return null
  if ('total' in res) {
    return `${res.store} — ${res.total.toLocaleString('ru-RU')} ₽ · ${res.category}`
  }
  return `${res.store} — ${res.amount.toLocaleString('ru-RU')} ₽ · ${res.category}`
})

async function toggleRecording() {
  if (busy.value) return
  if (status.value === 'recording') {
    const audio = await stop()
    if (!audio || audio.size === 0) return
    await processAudio(audio)
    return
  }
  submitError.value = ''
  parseError.value = ''
  clearLastTranscript()
  reset()
  await start()
}

async function processAudio(audio: Blob) {
  submitError.value = ''
  parseError.value = ''
  try {
    const text = await transcribeAudio(audio, PURCHASE_VOICE_EXAMPLES[0])
    await submitTranscript(text)
  } catch (e) {
    clearLastTranscript()
    const code = (e as { statusCode?: number })?.statusCode
    if (code === 503) {
      parseError.value =
        'Сервис распознавания речи недоступен. Запустите Whisper или выберите пример ниже.'
    } else {
      parseError.value = e instanceof Error ? e.message : 'Не удалось распознать голос'
    }
  }
}

async function submitTranscript(text: string) {
  const trimmed = text.trim()
  if (!trimmed) {
    parseError.value = 'Не расслышали — попробуйте ещё раз или выберите пример.'
    return
  }
  try {
    await submitVoiceTranscript(trimmed)
    emit('done')
  } catch (e) {
    const code = (e as { statusCode?: number })?.statusCode
    if (code === 503) {
      submitError.value = 'Сервис обработки речи недоступен.'
    } else {
      submitError.value = e instanceof Error ? e.message : 'Не удалось добавить покупку'
    }
  }
}

function pickExample(example: string) {
  reset()
  clearLastTranscript()
  void submitTranscript(example)
}
</script>

<template>
  <div class="space-y-4">
    <p class="text-center text-sm leading-relaxed text-[color:var(--mm-text)]">
      Скажите, что купили — Поток разберёт сумму и категорию.
    </p>

    <div class="flex justify-center py-2">
      <button
        type="button"
        class="mm-onb-mic-orb-hit mm-onb-mic-orb-hit--preview border-0 bg-transparent p-0"
        :class="{
          'mm-onb-mic-orb-hit--listen': status === 'recording' && !busy,
          'mm-onb-mic-orb-hit--parse': busy && status !== 'recording'
        }"
        :disabled="!supported || busy"
        :aria-pressed="status === 'recording'"
        :aria-label="status === 'recording' ? 'Остановить и отправить' : 'Нажмите и скажите покупку'"
        @click="toggleRecording"
      >
        <OnboardingMicOrbVisual
          :listening="status === 'recording'"
          :parsing="busy"
          compact
          gentle
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
              ? 'Нажмите на орб, скажите покупку и нажмите ещё раз'
              : 'Голос недоступен — выберите пример ниже'
      }}
    </p>

    <p v-if="errorMessage" class="text-center text-xs text-destructive">{{ errorMessage }}</p>
    <p v-if="submitError" class="text-center text-xs text-destructive">{{ submitError }}</p>

    <p v-if="lastTranscript && !parseError" class="mm-onb-transcript-bubble">
      «{{ lastTranscript }}»
    </p>

    <p v-if="parseError" class="text-center text-sm text-destructive">{{ parseError }}</p>

    <div
      v-if="addedHint && !busy"
      class="rounded-xl border border-[color:var(--mm-primary)]/25 bg-[color:var(--mm-primary-soft)]/60 px-4 py-3 text-center text-sm text-[color:var(--mm-text)]"
    >
      <span class="font-medium text-[color:var(--mm-primary)]">Записали: </span>
      {{ addedHint }}
    </div>

    <div v-if="!addedHint && (!supported || parseError)" class="flex flex-wrap justify-center gap-2">
      <button
        v-for="example in PURCHASE_VOICE_EXAMPLES"
        :key="example"
        type="button"
        class="mm-onb-chip"
        :disabled="busy"
        @click="pickExample(example)"
      >
        {{ example }}
      </button>
    </div>
  </div>
</template>
