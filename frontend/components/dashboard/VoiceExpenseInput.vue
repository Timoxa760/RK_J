<script setup lang="ts">
import { Mic, MicOff, RotateCcw } from 'lucide-vue-next'

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
  if (!res || !('total' in res)) return null
  return `${res.store} — ${res.total.toLocaleString('ru-RU')} ₽ · ${res.category}`
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
    submitError.value = 'Сначала запишите покупку — нажмите микрофон.'
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
      submitError.value = 'Сервис распознавания речи недоступен. Проверьте, что Whisper запущен.'
      return
    }
    submitError.value = e instanceof Error ? e.message : 'Не удалось отправить запись'
  }
}
</script>

<template>
  <div class="space-y-4">
    <p class="text-sm text-muted-foreground">
      Нажмите микрофон и скажите, что купили — Поток разберёт и запишет. Подключать ФНС не обязательно.
    </p>

    <div class="flex items-center gap-3 rounded-lg border bg-muted/50 px-3 py-3">
      <Button
        type="button"
        size="icon"
        :variant="status === 'recording' ? 'default' : 'secondary'"
        class="shrink-0 rounded-full"
        :disabled="!supported || submitting"
        :aria-pressed="status === 'recording'"
        @click="toggleRecording"
      >
        <MicOff v-if="status === 'recording'" class="size-5" />
        <Mic v-else class="size-5" />
      </Button>
      <div class="min-w-0 flex-1">
        <p class="text-xs font-medium text-muted-foreground">
          {{
            submitting
              ? 'Обрабатываем…'
              : status === 'recording'
                ? 'Поток слушает… нажмите ещё раз или «Отправить»'
                : supported
                  ? 'Нажмите, чтобы записать'
                  : 'Запись недоступна в этом браузере'
          }}
        </p>
        <div v-if="status === 'recording'" class="mt-1 flex gap-0.5 mm-voice-demo__bars" aria-hidden="true">
          <span /><span /><span /><span /><span />
        </div>
      </div>
      <Button
        v-if="status === 'recorded'"
        type="button"
        variant="ghost"
        size="icon"
        class="shrink-0"
        :disabled="submitting"
        aria-label="Записать заново"
        @click="reset"
      >
        <RotateCcw class="size-4" />
      </Button>
    </div>

    <p v-if="errorMessage" class="text-sm text-destructive">{{ errorMessage }}</p>
    <p v-if="submitError" class="text-sm text-destructive">{{ submitError }}</p>

    <div
      v-if="addedHint && !submitting"
      class="rounded-lg border border-[color:var(--mm-primary)]/25 bg-[color:var(--mm-primary-soft)]/50 px-3 py-2 text-sm"
    >
      <span class="font-medium text-[color:var(--mm-primary)]">Записали: </span>
      {{ addedHint }}
    </div>

    <Button
      :disabled="submitting || (status !== 'recorded' && status !== 'recording')"
      @click="send"
    >
      {{ submitting ? 'Обрабатываем…' : status === 'recording' ? 'Остановить и отправить' : 'Отправить Потоку' }}
    </Button>
  </div>
</template>
