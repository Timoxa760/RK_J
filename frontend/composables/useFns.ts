import type {
  FnsTicketRequest,
  FnsTicketResponse,
  ReceiptFnsScanRequest,
  ReceiptFnsScanResponse
} from '~/types/api'
import { parseFnsQr } from '~/utils/fnsQr'

export function useFns() {
  const { apiFetch } = useApi()
  const toast = useToast()

  const loading = ref(false)

  async function scanReceipt(body: ReceiptFnsScanRequest) {
    loading.value = true
    try {
      const res = await apiFetch<ReceiptFnsScanResponse>('/receipt/fns/scan', {
        method: 'POST',
        body
      })
      toast.show(
        `Чек проверен: ${res.store}, ${res.total.toLocaleString('ru-RU')} ₽`,
        'success'
      )
      return res
    } catch (e) {
      toast.show(
        e instanceof Error ? e.message : 'ФНС временно недоступна — попробуйте позже',
        'error'
      )
      throw e
    } finally {
      loading.value = false
    }
  }

  async function submitTicket(qr: string) {
    loading.value = true
    const body: FnsTicketRequest = { qr: qr.trim() }

    try {
      const res = await apiFetch<FnsTicketResponse>('/fns/ticket', {
        method: 'POST',
        body
      })
      toast.show('Чек проверен через ФНС', 'success')
      return res
    } catch (e) {
      toast.show(
        e instanceof Error ? e.message : 'ФНС временно недоступна — попробуйте позже',
        'error'
      )
      throw e
    } finally {
      loading.value = false
    }
  }

  async function submitQr(qr: string) {
    const parsed = parseFnsQr(qr)
    if (parsed) {
      return scanReceipt(parsed)
    }
    return submitTicket(qr)
  }

  return {
    loading,
    scanReceipt,
    submitTicket,
    submitQr
  }
}
