<script setup lang="ts">
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
    <div class="flex justify-center py-2">
      <button
        type="button"
        class="mm-onb-mic-orb-hit mm-onb-mic-orb-hit--preview border-0 bg-transparent p-0"
        :disabled="!supported || busy"
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
        status === 'recording'
          ? 'Поток слушает — нажмите ещё раз, чтобы остановить'
          : status === 'recorded'
            ? `Запись готова (${durationLabel})`
            : supported
              ? 'Нажмите на орб, скажите покупку и отправьте'
              : 'Запись голоса недоступна в этом браузере'
      }}
    </p>

    <p v-if="errorMessage" class="text-center text-xs text-destructive">{{ errorMessage }}</p>

    <div class="mm-onb-form-panel space-y-3">
      <Button
        type="button"
        class="w-full"
        variant="outline"
        :disabled="busy || status !== 'recorded'"
        @click="reset"
      >
        Записать заново
      </Button>
      <Button
        type="button"
        class="w-full"
        :disabled="busy || (status !== 'recorded' && status !== 'recording')"
        @click="send"
      >
        {{ busy ? 'Обрабатываем…' : status === 'recording' ? 'Остановить и отправить' : 'Добавить покупку' }}
      </Button>
    </div>
  </div>
</template>
