import type { ReceiptListItem, ReceiptsListResponse } from '~/types/api'
import { readStoredReceipts } from '~/utils/receiptListStorage'

export function useReceiptList() {
  const { apiFetch } = useApi()

  const receipts = ref<ReceiptListItem[]>([])
  const selected = ref<ReceiptListItem | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function refresh() {
    loading.value = true
    error.value = null

    try {
      const remote = await apiFetch<ReceiptsListResponse>('/receipts')
      receipts.value = remote.receipts ?? []
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить расходы'
      receipts.value = readStoredReceipts()
    } finally {
      loading.value = false
    }
  }

  function selectReceipt(item: ReceiptListItem) {
    selected.value = item
  }

  function closeDetail() {
    selected.value = null
  }

  onMounted(refresh)

  return {
    receipts,
    selected,
    loading,
    error,
    selectReceipt,
    closeDetail,
    refresh
  }
}
