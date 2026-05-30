import type {
  ReceiptManualRequest,
  ReceiptManualResponse,
  ReceiptVoiceResponse
} from '~/types/api'
import {
  appendStoredReceipt,
  receiptFromManual,
  receiptFromVoice
} from '~/utils/receiptListStorage'
import { toReceiptIsoDate } from '~/utils/receiptDate'

export type ReceiptSubmitResult = ReceiptManualResponse | ReceiptVoiceResponse

export function useReceiptSubmit() {
  const { apiFetch } = useApi()
  const toast = useToast()

  const submitting = useState<boolean>('receipt-submitting', () => false)
  const lastResult = useState<ReceiptSubmitResult | null>('receipt-last-result', () => null)

  async function submitManual(payload: {
    store: string
    amount: number
    category: string
    date: string
  }) {
    submitting.value = true
    lastResult.value = null

    const body: ReceiptManualRequest = {
      store: payload.store.trim() || 'Не указан',
      amount: payload.amount,
      category: payload.category,
      date: toReceiptIsoDate(payload.date)
    }

    try {
      lastResult.value = await apiFetch<ReceiptManualResponse>('/receipt/manual', {
        method: 'POST',
        body
      })

      appendStoredReceipt(receiptFromManual(lastResult.value as ReceiptManualResponse))
      toast.show(
        `Расход ${body.amount.toLocaleString('ru-RU')} ₽ учтён · ${body.category}`,
        'success'
      )
      return lastResult.value
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Не удалось сохранить расход'
      toast.show(msg, 'error')
      throw e
    } finally {
      submitting.value = false
    }
  }

  async function submitVoiceAudio(audio: Blob) {
    submitting.value = true
    lastResult.value = null

    try {
      const form = new FormData()
      const ext = audio.type.includes('mp4') ? 'mp4' : 'webm'
      form.append('audio', audio, `recording.${ext}`)

      lastResult.value = await apiFetch<ReceiptVoiceResponse>('/receipt/voice', {
        method: 'POST',
        body: form
      })

      appendStoredReceipt(receiptFromVoice(lastResult.value as ReceiptVoiceResponse))
      const res = lastResult.value as ReceiptVoiceResponse
      toast.show(
        `Поток разобрал: ${res.total.toLocaleString('ru-RU')} ₽ · ${res.store}`,
        'success'
      )
      return lastResult.value
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Не удалось распознать голос'
      toast.show(msg, 'error')
      throw e
    } finally {
      submitting.value = false
    }
  }

  return {
    submitting,
    lastResult,
    submitManual,
    submitVoiceAudio
  }
}
