<script setup lang="ts">
import { Mic, MicOff, RotateCcw } from 'lucide-vue-next'

const emit = defineEmits<{
  submit: [audio: Blob]
}>()

const props = defineProps<{
  busy?: boolean
}>()

const {
  status,
  errorMessage,
  blob,
  durationMs,
  supported,
  start,
  stop,
  reset
} = useAudioRecorder()

async function toggleRecording() {
  if (props.busy) return
  if (status.value === 'recording') {
    await stop()
    return
  }
  reset()
  await start()
}

async function send() {
  if (props.busy) return
  let audio = blob.value
  if (status.value === 'recording') {
    audio = await stop()
  }
  if (!audio || audio.size === 0) return
  emit('submit', audio)
}

const durationLabel = computed(() => {
  const sec = Math.max(1, Math.round(durationMs.value / 1000))
  return `${sec} сек`
})
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
        :disabled="!supported || busy"
        :aria-pressed="status === 'recording'"
        @click="toggleRecording"
      >
        <MicOff v-if="status === 'recording'" class="size-5" />
        <Mic v-else class="size-5" />
      </Button>
      <div class="min-w-0 flex-1">
        <p class="text-xs font-medium text-muted-foreground">
          {{
            status === 'recording'
              ? 'Поток слушает…'
              : status === 'recorded'
                ? `Запись готова (${durationLabel})`
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
        :disabled="busy"
        aria-label="Записать заново"
        @click="reset"
      >
        <RotateCcw class="size-4" />
      </Button>
    </div>

    <p v-if="errorMessage" class="text-sm text-destructive">{{ errorMessage }}</p>

    <Button
      :disabled="busy || (status !== 'recorded' && status !== 'recording')"
      @click="send"
    >
      {{ busy ? 'Обрабатываем…' : status === 'recording' ? 'Остановить и отправить' : 'Отправить Потоку' }}
    </Button>
  </div>
</template>
