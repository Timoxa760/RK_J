import type {
  ManualExpenseItem,
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
import { useAuthStore } from '~/store/authStore'
import { formatApiError } from '~/utils/apiError'

export type ReceiptSubmitResult = ReceiptManualResponse | ReceiptVoiceResponse

function expenseRowsFromResponse(res: ManualExpenseResponse): ManualExpenseItem[] {
  if (res.expenses?.length) {
    return res.expenses
  }
  return [
    {
      id: res.id,
      amount: res.amount,
      category: res.category,
      description: res.description
    }
  ]
}

function voiceResultFromExpense(
  text: string,
  row: ManualExpenseItem,
  parsed: boolean,
  index = 0
): ReceiptVoiceResponse {
  const label =
    row.description?.trim() ||
    text.trim().slice(0, 48) ||
    row.category ||
    'Покупка'
  return {
    receipt_id: row.id || `voice-${Date.now()}-${index}`,
    store: label,
    items: [
      {
        name: label,
        price: row.amount,
        quantity: 1,
        category: row.category
      }
    ],
    total: row.amount,
    category: row.category,
    confidence: parsed ? 0.85 : 0.75
  }
}

function voiceResultsFromResponse(text: string, res: ManualExpenseResponse): ReceiptVoiceResponse[] {
  return expenseRowsFromResponse(res).map((row, index) =>
    voiceResultFromExpense(text, row, res.parsed, index)
  )
}

function summarizeVoiceResults(results: ReceiptVoiceResponse[]): ReceiptVoiceResponse {
  if (results.length === 1) {
    return results[0]!
  }
  const total = results.reduce((sum, row) => sum + row.total, 0)
  const categories = [...new Set(results.map((row) => row.category))]
  return {
    receipt_id: results[0]!.receipt_id,
    store: `${results.length} покупки`,
    items: results.flatMap((row) => row.items),
    total,
    category: categories.join(' · '),
    confidence: results[0]!.confidence
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
      throw new Error('Пустая фраза')
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

      const results = voiceResultsFromResponse(trimmed, expense)
      if (!results.length) {
        throw new Error('Не удалось разобрать покупку')
      }
      for (const result of results) {
        appendStoredReceipt(receiptFromVoice(result))
      }
      lastResult.value = summarizeVoiceResults(results)
      const summary = lastResult.value as ReceiptVoiceResponse
      toast.show(
        results.length > 1
          ? `Записали ${results.length} покупки · ${summary.total.toLocaleString('ru-RU')} ₽`
          : `Поток разобрал: ${summary.total.toLocaleString('ru-RU')} ₽ · ${summary.category}`,
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
