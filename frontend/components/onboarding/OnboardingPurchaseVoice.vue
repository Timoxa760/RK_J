<script setup lang="ts">
const emit = defineEmits<{
  done: []
}>()

const transcript = ref('')
const listening = ref(false)
const speechSupported = ref(false)
const submitError = ref('')

const { submitting, submitVoiceTranscript } = useReceiptSubmit()

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

function toggleListen() {
  if (!recognition || submitting.value) return
  if (listening.value) {
    recognition.stop()
    return
  }
  submitError.value = ''
  transcript.value = ''
  listening.value = true
  recognition.start()
}

async function submitPhrase(rawText: string) {
  const text = rawText.trim()
  if (!text || submitting.value) return
  submitError.value = ''
  try {
    await submitVoiceTranscript(text)
    emit('done')
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Не удалось добавить покупку'
  }
}

onMounted(() => {
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
      if (said) void submitPhrase(said)
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
  <div class="mm-onboarding-voice__stage space-y-4">
    <p class="mm-onboarding-voice__prompt">
      Нажмите микрофон и скажите, что купили — Поток разберёт и запишет.
    </p>

    <div v-if="speechSupported" class="mm-onboarding-voice__hero">
      <OnboardingMicOrb
        :listening="listening"
        :parsing="submitting"
        label="Нажмите микрофон и ответьте вслух"
        @click="toggleListen"
      />
    </div>

    <p
      v-else
      class="max-w-sm text-sm text-[color:var(--mm-text-muted)]"
    >
      Голос недоступен в этом браузере — переключитесь на вкладку «Вручную».
    </p>

    <p
      v-if="listening && transcript"
      class="mm-onb-transcript-bubble"
    >
      «{{ transcript }}»
    </p>

    <p v-if="submitError" class="text-sm text-destructive">{{ submitError }}</p>
  </div>
</template>
