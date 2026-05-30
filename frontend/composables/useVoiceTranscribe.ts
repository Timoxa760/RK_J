import type { VoiceTranscribeResponse } from '~/types/api'

export function useVoiceTranscribe() {
  const { apiFetch, demoMode } = useApi()
  const transcribing = ref(false)
  const lastTranscript = ref('')

  function clearLastTranscript() {
    lastTranscript.value = ''
  }

  async function transcribeAudio(audio: Blob, demoFallback?: string): Promise<string> {
    transcribing.value = true
    lastTranscript.value = ''
    try {
      if (demoMode.value) {
        await new Promise((r) => setTimeout(r, 700))
        const text = demoFallback?.trim() || '180 тысяч в месяц'
        lastTranscript.value = text
        return text
      }

      const form = new FormData()
      const ext = audio.type.includes('mp4') ? 'mp4' : 'webm'
      form.append('audio', audio, `recording.${ext}`)

      const res = await apiFetch<VoiceTranscribeResponse>('/voice/transcribe', {
        method: 'POST',
        body: form
      })
      const text = res.text.trim()
      lastTranscript.value = text
      return text
    } catch (e) {
      lastTranscript.value = ''
      throw e
    } finally {
      transcribing.value = false
    }
  }

  return {
    transcribing,
    lastTranscript,
    transcribeAudio,
    clearLastTranscript
  }
}
