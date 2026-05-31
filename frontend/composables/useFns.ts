import type { FnsConnectionState, FnsSyncResponse, ReceiptListItem } from '~/types/api'
import { fnsReceiptCatalog } from '~/constants/fnsReceiptCatalog'
import { useDashboardStore } from '~/store/dashboardStore'
import { useAuthStore } from '~/store/authStore'
import { appendStoredReceipt } from '~/utils/receiptListStorage'
import { readFnsConnection, writeFnsConnection } from '~/utils/fnsStorage'
import { FNS } from '~/constants/productCopy'

const SEND_CODE_DELAY_MS = 1_100
const VERIFY_CODE_DELAY_MS = 1_400
const CONNECT_DELAY_MS = 900
const RECEIPT_IMPORT_STEP_MS = 100

export type FnsConnectStep = 'phone' | 'code'

function mergeImportedIds(current: string[], ids: string[]): string[] {
  const set = new Set(current)
  for (const id of ids) set.add(id)
  return [...set]
}

function pendingReceipts(importedIds: string[]): ReceiptListItem[] {
  const known = new Set(importedIds)
  return fnsReceiptCatalog.filter((item) => !known.has(item.id))
}

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export function useFns() {
  const { apiFetch } = useApi()
  const authStore = useAuthStore()
  const dashboardStore = useDashboardStore()

  const connection = useState<FnsConnectionState>('fns-connection', () => defaultConnection())
  const syncing = useState('fns-sync-loading', () => false)
  const connectBusy = useState('fns-connect-busy', () => false)
  const connectStep = useState<FnsConnectStep>('fns-connect-step', () => 'phone')
  const connectPhone = useState('fns-connect-phone', () => '')
  const error = useState<string | null>('fns-error', () => null)
  const success = useState<string | null>('fns-success', () => null)

  function defaultConnection(): FnsConnectionState {
    if (import.meta.client) return readFnsConnection()
    return {
      connected: false,
      phone: '',
      auto_sync: true,
      imported_ids: []
    }
  }

  function hydrate() {
    connection.value = readFnsConnection()
  }

  if (import.meta.client) {
    onMounted(hydrate)
  }

  async function pushReceiptToExpenses(item: ReceiptListItem) {
    const userId = authStore.user?.phone || authStore.user?.id
    if (!userId) return

    try {
      await apiFetch('/expenses/manual', {
        method: 'POST',
        body: {
          user_id: userId,
          amount: item.amount,
          category: item.category ?? 'Другое',
          description: item.store,
          date: item.date,
          source: 'fns'
        }
      })
    } catch {
      // Чек остаётся в локальной ленте, даже если API расходов недоступен.
    }
  }

  async function importReceipts(receipts: ReceiptListItem[]): Promise<FnsSyncResponse> {
    const imported: ReceiptListItem[] = []

    for (let index = 0; index < receipts.length; index += 1) {
      const receipt = receipts[index]
      appendStoredReceipt(receipt)
      await pushReceiptToExpenses(receipt)
      imported.push(receipt)
      if (index < receipts.length - 1) {
        await delay(RECEIPT_IMPORT_STEP_MS)
      }
    }

    if (imported.length > 0) {
      await dashboardStore.loadAll({ silent: true })
    }

    return {
      imported: imported.length,
      receipts: imported,
      synced_at: new Date().toISOString()
    }
  }

  function applySyncResult(result: FnsSyncResponse) {
    const next: FnsConnectionState = {
      ...connection.value,
      last_sync_at: result.synced_at,
      imported_ids: mergeImportedIds(
        connection.value.imported_ids,
        result.receipts.map((item) => item.id)
      )
    }
    writeFnsConnection(next)
    connection.value = next
    success.value =
      result.imported > 0 ? FNS.syncImported(result.imported) : FNS.syncUpToDate
  }

  async function syncFromCatalog(): Promise<FnsSyncResponse> {
    return importReceipts(pendingReceipts(connection.value.imported_ids))
  }

  function resetConnectFlow() {
    connectStep.value = 'phone'
    connectPhone.value = ''
    connectBusy.value = false
  }

  function backToPhoneStep() {
    connectStep.value = 'phone'
    error.value = null
  }

  async function sendConnectCode(phone: string): Promise<boolean> {
    error.value = null
    success.value = null
    connectBusy.value = true

    const normalizedPhone = phone.trim()
    if (!normalizedPhone) {
      error.value = FNS.phoneRequired
      connectBusy.value = false
      return false
    }

    try {
      await delay(SEND_CODE_DELAY_MS)
      connectPhone.value = normalizedPhone
      connectStep.value = 'code'
      return true
    } finally {
      connectBusy.value = false
    }
  }

  async function verifyConnectCode(code: string): Promise<FnsSyncResponse | null> {
    error.value = null
    success.value = null
    connectBusy.value = true
    syncing.value = true

    const trimmed = code.trim()
    if (!trimmed) {
      error.value = FNS.codeRequired
      connectBusy.value = false
      syncing.value = false
      return null
    }

    if (!connectPhone.value) {
      error.value = FNS.phoneRequired
      connectStep.value = 'phone'
      connectBusy.value = false
      syncing.value = false
      return null
    }

    try {
      await delay(VERIFY_CODE_DELAY_MS)
      await delay(CONNECT_DELAY_MS)

      const next: FnsConnectionState = {
        ...connection.value,
        connected: true,
        phone: connectPhone.value,
        connected_at: new Date().toISOString(),
        auto_sync: true
      }
      writeFnsConnection(next)
      connection.value = next

      const result = await syncFromCatalog()
      applySyncResult(result)
      success.value = FNS.connected
      resetConnectFlow()
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : FNS.connectFailed
      return null
    } finally {
      connectBusy.value = false
      syncing.value = false
    }
  }

  async function sync(): Promise<FnsSyncResponse | null> {
    if (!connection.value.connected) return null

    error.value = null
    success.value = null
    syncing.value = true

    try {
      const result = await syncFromCatalog()
      applySyncResult(result)
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : FNS.syncFailed
      return null
    } finally {
      syncing.value = false
    }
  }

  function disconnect() {
    const next: FnsConnectionState = {
      ...connection.value,
      connected: false,
      connected_at: undefined,
      last_sync_at: undefined
    }
    writeFnsConnection(next)
    connection.value = next
    success.value = FNS.disconnected
    error.value = null
    resetConnectFlow()
  }

  const importedCount = computed(() => connection.value.imported_ids.length)
  const pendingCount = computed(() =>
    pendingReceipts(connection.value.imported_ids).length
  )
  const connectLoading = computed(() => connectBusy.value || syncing.value)

  return {
    connection,
    syncing,
    connectBusy,
    connectLoading,
    connectStep,
    connectPhone,
    error,
    success,
    importedCount,
    pendingCount,
    hydrate,
    sendConnectCode,
    verifyConnectCode,
    resetConnectFlow,
    backToPhoneStep,
    sync,
    disconnect
  }
}
