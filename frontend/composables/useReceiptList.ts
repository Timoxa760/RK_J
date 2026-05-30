import type { ReceiptListItem } from '~/types/api'
import { mockReceiptListItems } from '~/store/mocks/receiptList'
import { readStoredReceipts } from '~/utils/receiptListStorage'

export function useReceiptList() {
  const { demoMode } = useApi()

  const receipts = ref<ReceiptListItem[]>([])
  const selected = ref<ReceiptListItem | null>(null)

  function refresh() {
    const stored = readStoredReceipts()
    if (stored.length) {
      receipts.value = stored
      return
    }
    receipts.value = demoMode.value ? [...mockReceiptListItems] : []
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
    selectReceipt,
    closeDetail,
    refresh
  }
}
