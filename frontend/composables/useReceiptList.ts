import type { ReceiptListItem, ReceiptsListResponse } from '~/types/api'
import { readStoredReceipts, removeStoredReceipt } from '~/utils/receiptListStorage'

function mergeReceiptLists(remote: ReceiptListItem[], local: ReceiptListItem[]): ReceiptListItem[] {
  const seen = new Set(remote.map((item) => item.id))
  const pending = local.filter((item) => item.id && !seen.has(item.id))
  if (!pending.length) return remote
  return [...pending, ...remote]
}

export function useReceiptList() {
  const { apiFetch } = useApi()

  const receipts = ref<ReceiptListItem[]>([])
  const selected = ref<ReceiptListItem | null>(null)
  const loading = ref(false)
  const deleting = ref(false)
  const error = ref<string | null>(null)

  async function refresh() {
    loading.value = true
    error.value = null

    try {
      const remote = await apiFetch<ReceiptsListResponse>('/receipts')
      const fromApi = remote.receipts ?? []
      receipts.value = mergeReceiptLists(fromApi, readStoredReceipts())
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

  async function deleteReceipt(id: string) {
    deleting.value = true
    error.value = null

    try {
      await apiFetch(`/receipts/${encodeURIComponent(id)}`, { method: 'DELETE' })
      removeStoredReceipt(id)
    } catch (e) {
      const hadLocal = readStoredReceipts().some((item) => item.id === id)
      if (hadLocal) {
        removeStoredReceipt(id)
      } else {
        error.value = e instanceof Error ? e.message : 'Не удалось удалить покупку'
        deleting.value = false
        return
      }
    }

    if (selected.value?.id === id) {
      closeDetail()
    }

    await refresh()
    deleting.value = false
  }

  onMounted(refresh)

  return {
    receipts,
    selected,
    loading,
    deleting,
    error,
    selectReceipt,
    closeDetail,
    deleteReceipt,
    refresh
  }
}
