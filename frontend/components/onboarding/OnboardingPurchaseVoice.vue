<script setup lang="ts">
const emit = defineEmits<{
  done: []
}>()

const {
  status,
  errorMessage,
  blob,
  supported,
  start,
  stop,
  reset
} = useAudioRecorder()

const submitError = ref('')
const { submitting, submitVoiceAudio, lastResult } = useReceiptSubmit()

const addedHint = computed(() => {
  const res = lastResult.value
  if (!res) return null
  if ('total' in res) {
    return `${res.store} — ${res.total.toLocaleString('ru-RU')} ₽ · ${res.category}`
  }
  return `${res.store} — ${res.amount.toLocaleString('ru-RU')} ₽ · ${res.category}`
})

async function toggleRecording() {
  if (submitting.value) return
  submitError.value = ''
  if (status.value === 'recording') {
    const audio = await stop()
    if (audio && audio.size > 0) {
      await submitAudio(audio)
    }
    return
  }
  reset()
  await start()
}

async function send() {
  if (submitting.value) return
  submitError.value = ''
  let audio = blob.value
  if (status.value === 'recording') {
    audio = await stop()
  }
  if (!audio || audio.size === 0) {
    submitError.value = 'Сначала запишите покупку — нажмите на орб.'
    return
  }
  await submitAudio(audio)
}

async function submitAudio(audio: Blob) {
  try {
    await submitVoiceAudio(audio)
    emit('done')
  } catch (e) {
    const status = (e as { statusCode?: number })?.statusCode
    if (status === 503) {
      submitError.value = 'Сервис распознавания речи недоступен. Проверьте Whisper.'
      return
    }
    submitError.value = e instanceof Error ? e.message : 'Не удалось отправить запись'
  }
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex justify-center py-2">
      <button
        type="button"
        class="mm-onb-mic-orb-hit mm-onb-mic-orb-hit--preview border-0 bg-transparent p-0"
        :disabled="!supported || submitting"
        :aria-pressed="status === 'recording'"
        aria-label="Записать покупку голосом"
        @click="toggleRecording"
      >
        <OnboardingMicOrbVisual
          :listening="status === 'recording'"
          compact
          :ambient="status !== 'recording'"
        />
      </button>
    </div>

    <p class="text-center text-xs text-[color:var(--mm-text-muted)]">
      {{
        submitting
          ? 'Распознаём…'
          : status === 'recording'
            ? 'Поток слушает — нажмите ещё раз, чтобы отправить'
            : supported
              ? 'Нажмите на орб, скажите покупку и нажмите ещё раз'
              : 'Запись голоса недоступна в этом браузере'
      }}
    </p>

    <p v-if="errorMessage" class="text-center text-xs text-destructive">{{ errorMessage }}</p>
    <p v-if="submitError" class="text-center text-xs text-destructive">{{ submitError }}</p>

    <div
      v-if="addedHint && !submitting"
      class="rounded-xl border border-[color:var(--mm-primary)]/25 bg-[color:var(--mm-primary-soft)]/60 px-4 py-3 text-center text-sm text-[color:var(--mm-text)]"
    >
      <span class="font-medium text-[color:var(--mm-primary)]">Записали: </span>
      {{ addedHint }}
    </div>

    <div v-if="!addedHint" class="mm-onb-form-panel space-y-3">
      <Button
        type="button"
        class="w-full"
        variant="outline"
        :disabled="submitting || status !== 'recorded'"
        @click="reset"
      >
        Записать заново
      </Button>
      <Button
        type="button"
        class="w-full"
        :disabled="submitting || (status !== 'recorded' && status !== 'recording')"
        @click="send"
      >
        {{ submitting ? 'Обрабатываем…' : status === 'recording' ? 'Остановить и отправить' : 'Добавить покупку' }}
      </Button>
    </div>
  </div>
</template>
