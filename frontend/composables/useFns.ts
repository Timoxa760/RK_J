import type {
  FnsTicketRequest,
  FnsTicketResponse,
  ReceiptFnsScanRequest,
  ReceiptFnsScanResponse
} from '~/types/api'
import { mockReceiptFnsScan } from '~/store/mocks/receipts'
import { parseFnsQr } from '~/utils/fnsQr'

export function useFns() {
  const { apiFetch, demoMode } = useApi()
  const toast = useToast()

  const loading = ref(false)

  async function scanReceipt(body: ReceiptFnsScanRequest) {
    loading.value = true
    try {
      if (demoMode.value) {
        await new Promise((r) => setTimeout(r, 800))
        const res = {
          ...mockReceiptFnsScan,
          receipt_id: `demo-fns-${Date.now()}`
        } satisfies ReceiptFnsScanResponse
        toast.show(
          `Чек принят: ${res.store}, ${res.total.toLocaleString('ru-RU')} ₽`,
          'success'
        )
        return res
      }

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
      if (demoMode.value) {
        await new Promise((r) => setTimeout(r, 800))
        toast.show('Чек ФНС принят — появится в ленте после обработки', 'success')
        return { success: true, receipt_id: 'demo-fns-ticket-1' } satisfies FnsTicketResponse
      }
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

  /** Парсит QR → scan; иначе fallback на /fns/ticket. */
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
