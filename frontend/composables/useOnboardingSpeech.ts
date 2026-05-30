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

export function useOnboardingSpeech(onFinalTranscript: (text: string) => void) {
  const transcript = ref('')
  const listening = ref(false)
  const speechSupported = ref(false)

  let recognition: InstanceType<SpeechRecognitionCtor> | null = null

  onMounted(() => {
    if (!import.meta.client) return
    const win = window as Window & { webkitSpeechRecognition?: SpeechRecognitionCtor }
    const SR = win.SpeechRecognition ?? win.webkitSpeechRecognition
    speechSupported.value = Boolean(SR)
    if (!SR) return

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
      if (said) onFinalTranscript(said)
    }
    recognition.onerror = () => {
      listening.value = false
    }
  })

  onUnmounted(() => {
    recognition?.abort()
  })

  function toggleListen(disabled?: boolean) {
    if (!recognition || disabled) return
    if (listening.value) {
      recognition.stop()
      return
    }
    transcript.value = ''
    listening.value = true
    recognition.start()
  }

  function clearTranscript() {
    transcript.value = ''
  }

  return {
    transcript,
    listening,
    speechSupported,
    toggleListen,
    clearTranscript
  }
}
