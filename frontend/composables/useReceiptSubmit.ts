import type {
  ManualExpenseResponse,
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
import { formatApiError } from '~/utils/apiError'

export type ReceiptSubmitResult = ReceiptManualResponse | ReceiptVoiceResponse

function voiceResultFromExpense(text: string, res: ManualExpenseResponse): ReceiptVoiceResponse {
  const label =
    res.description?.trim() ||
    text.trim().slice(0, 48) ||
    res.category ||
    'Покупка'
  return {
    receipt_id: res.id || `voice-${Date.now()}`,
    store: label,
    items: [
      {
        name: label,
        price: res.amount,
        quantity: 1,
        category: res.category
      }
    ],
    total: res.amount,
    category: res.category,
    confidence: res.parsed ? 0.85 : 0.75
  }
}

export function useReceiptSubmit() {
  const { apiFetch } = useApi()
  const toast = useToast()
  const authStore = useAuthStore()

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
      const msg = formatApiError(e, 'Не удалось сохранить расход')
      toast.show(msg, 'error')
      throw new Error(msg)
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
      const msg = formatApiError(e, 'Не удалось распознать голос')
      toast.show(msg, 'error')
      throw new Error(msg)
    } finally {
      submitting.value = false
    }
  }

  async function submitVoiceTranscript(text: string) {
    submitting.value = true
    lastResult.value = null

    const trimmed = text.trim()
    if (!trimmed) {
      submitting.value = false
      throw new Error('Скажите, что купили, например: «колбаса 300 рублей».')
    }

    if (!/\d/.test(trimmed)) {
      submitting.value = false
      throw new Error('Не удалось найти сумму. Добавьте цену, например: «колбаса 300 рублей».')
    }

    const userId = authStore.user?.phone || authStore.user?.id
    if (!userId) {
      submitting.value = false
      throw new Error('Сначала войдите в аккаунт')
    }

    try {
      const expense = await apiFetch<ManualExpenseResponse>('/expenses/manual', {
        method: 'POST',
        body: {
          user_id: userId,
          raw_text: trimmed,
          source: 'voice'
        }
      })

      lastResult.value = voiceResultFromExpense(trimmed, expense)
      appendStoredReceipt(receiptFromVoice(lastResult.value as ReceiptVoiceResponse))
      const res = lastResult.value as ReceiptVoiceResponse
      toast.show(
        `Поток разобрал: ${res.total.toLocaleString('ru-RU')} ₽ · ${res.category}`,
        'success'
      )
      return lastResult.value
    } catch (e) {
      const msg = formatApiError(e, 'Не удалось разобрать покупку')
      toast.show(msg, 'error')
      throw new Error(msg)
    } finally {
      submitting.value = false
    }
  }

  return {
    submitting,
    lastResult,
    submitManual,
    submitVoiceAudio,
    submitVoiceTranscript
  }
}
