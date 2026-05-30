import type { ReceiptListItem, ReceiptManualResponse, ReceiptVoiceResponse } from '~/types/api'
import { currentUserStorageKey } from '~/utils/userStorage'

const RECEIPTS_PREFIX = 'potok:receipt-list'

export function readStoredReceipts(): ReceiptListItem[] {
  if (!import.meta.client) return []
  try {
    const raw = localStorage.getItem(currentUserStorageKey(RECEIPTS_PREFIX))
    if (!raw) return []
    const parsed = JSON.parse(raw) as ReceiptListItem[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

export function writeStoredReceipts(receipts: ReceiptListItem[]) {
  if (!import.meta.client) return
  localStorage.setItem(currentUserStorageKey(RECEIPTS_PREFIX), JSON.stringify(receipts))
}

export function appendStoredReceipt(item: ReceiptListItem) {
  const next = [item, ...readStoredReceipts()]
  writeStoredReceipts(next)
  return next
}

export function receiptFromManual(res: ReceiptManualResponse): ReceiptListItem {
  return {
    id: res.receipt_id,
    store: res.store,
    amount: res.amount,
    date: res.date.slice(0, 10),
    category: res.category
  }
}

export function receiptFromVoice(res: ReceiptVoiceResponse): ReceiptListItem {
  return {
    id: res.receipt_id,
    store: res.store,
    amount: res.total,
    date: new Date().toISOString().slice(0, 10),
    category: res.category,
    items: res.items
  }
}
