<script setup lang="ts">
const PURCHASE_VOICE_EXAMPLES = [
  'купил продукты на 3500',
  'кофе 450 рублей',
  'кроссовки за 8900'
]

const emit = defineEmits<{
  done: []
}>()

const transcript = ref('')
const listening = ref(false)
const speechSupported = ref(false)
const submitError = ref('')

const { submitting, submitVoiceTranscript, lastResult } = useReceiptSubmit()

const addedHint = computed(() => {
  const res = lastResult.value
  if (!res) return null
  if ('total' in res) {
    return `${res.store} — ${res.total.toLocaleString('ru-RU')} ₽ · ${res.category}`
  }
  return `${res.store} — ${res.amount.toLocaleString('ru-RU')} ₽ · ${res.category}`
})

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

function pickExample(example: string) {
  void submitPhrase(example)
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
  <div class="space-y-4">
    <p class="text-center text-sm leading-relaxed text-[color:var(--mm-text)]">
      Скажите, что купили — как в опросе: нажмите микрофон и говорите.
    </p>

    <div v-if="speechSupported" class="mm-onboarding-voice__hero flex justify-center py-2">
      <OnboardingMicOrb
        :listening="listening"
        :parsing="submitting"
        label="Нажмите микрофон и назовите покупку"
        @click="toggleListen"
      />
    </div>

    <p
      v-else
      class="text-center text-sm text-[color:var(--mm-text-muted)]"
    >
      Голос недоступен в этом браузере — выберите пример:
    </p>

    <p
      v-if="listening && transcript"
      class="mm-onb-transcript-bubble"
    >
      «{{ transcript }}»
    </p>

    <p v-if="submitError" class="text-center text-sm text-destructive">{{ submitError }}</p>

    <div
      v-if="addedHint && !submitting"
      class="rounded-xl border border-[color:var(--mm-primary)]/25 bg-[color:var(--mm-primary-soft)]/60 px-4 py-3 text-center text-sm text-[color:var(--mm-text)]"
    >
      <span class="font-medium text-[color:var(--mm-primary)]">Записали: </span>
      {{ addedHint }}
    </div>

    <div
      v-if="!addedHint && (!speechSupported || submitError)"
      class="mm-onboarding-voice__chips flex flex-wrap justify-center gap-2"
    >
      <button
        v-for="example in PURCHASE_VOICE_EXAMPLES"
        :key="example"
        type="button"
        class="mm-onb-chip"
        :disabled="submitting"
        @click="pickExample(example)"
      >
        {{ example }}
      </button>
    </div>
  </div>
</template>
