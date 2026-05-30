export type AudioRecorderStatus = 'idle' | 'recording' | 'recorded' | 'unsupported' | 'denied'

export function useAudioRecorder() {
  const status = ref<AudioRecorderStatus>('idle')
  const errorMessage = ref<string | null>(null)
  const blob = ref<Blob | null>(null)
  const durationMs = ref(0)

  let mediaRecorder: MediaRecorder | null = null
  let chunks: Blob[] = []
  let stream: MediaStream | null = null
  let startedAt = 0

  const supported = computed(() => {
    if (!import.meta.client) return false
    return (
      typeof navigator !== 'undefined' &&
      Boolean(navigator.mediaDevices?.getUserMedia) &&
      typeof MediaRecorder !== 'undefined'
    )
  })

  function cleanup() {
    if (mediaRecorder && mediaRecorder.state !== 'inactive') {
      try {
        mediaRecorder.stop()
      } catch {
        /* ignore */
      }
    }
    stream?.getTracks().forEach((t) => t.stop())
    stream = null
    mediaRecorder = null
    chunks = []
  }

  function reset() {
    cleanup()
    blob.value = null
    durationMs.value = 0
    errorMessage.value = null
    status.value = supported.value ? 'idle' : 'unsupported'
  }

  async function start() {
    errorMessage.value = null
    blob.value = null
    durationMs.value = 0

    if (!supported.value) {
      status.value = 'unsupported'
      errorMessage.value = 'Запись голоса недоступна в этом браузере'
      return
    }

    try {
      stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      const mimeType = MediaRecorder.isTypeSupported('audio/webm')
        ? 'audio/webm'
        : MediaRecorder.isTypeSupported('audio/mp4')
          ? 'audio/mp4'
          : ''

      mediaRecorder = mimeType
        ? new MediaRecorder(stream, { mimeType })
        : new MediaRecorder(stream)

      chunks = []
      mediaRecorder.ondataavailable = (e) => {
        if (e.data.size > 0) chunks.push(e.data)
      }

      startedAt = Date.now()
      mediaRecorder.start()
      status.value = 'recording'
    } catch {
      cleanup()
      status.value = 'denied'
      errorMessage.value = 'Нет доступа к микрофону. Разрешите запись в настройках браузера.'
    }
  }

  function stop(): Promise<Blob | null> {
    return new Promise((resolve) => {
      if (!mediaRecorder || status.value !== 'recording') {
        resolve(blob.value)
        return
      }

      mediaRecorder.onstop = () => {
        durationMs.value = Date.now() - startedAt
        const type = mediaRecorder?.mimeType || 'audio/webm'
        const recorded = chunks.length ? new Blob(chunks, { type }) : null
        blob.value = recorded
        status.value = recorded ? 'recorded' : 'idle'
        stream?.getTracks().forEach((t) => t.stop())
        stream = null
        mediaRecorder = null
        chunks = []
        resolve(recorded)
      }

      try {
        mediaRecorder.stop()
      } catch {
        status.value = 'idle'
        resolve(null)
      }
    })
  }

  onMounted(() => {
    if (!import.meta.client) return
    status.value = supported.value ? 'idle' : 'unsupported'
  })

  onUnmounted(() => {
    cleanup()
  })

  return {
    status,
    errorMessage,
    blob,
    durationMs,
    supported,
    start,
    stop,
    reset
  }
}
