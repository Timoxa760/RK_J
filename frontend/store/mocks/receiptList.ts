import type { ReceiptListItem } from '~/types/api'
import { mockReceiptManual, mockReceiptVoice } from '~/store/mocks/receipts'

/** Демо-лента покупок для страницы /receipts. */
export const mockReceiptListItems: ReceiptListItem[] = [
  {
    id: 'demo-receipt-1',
    store: mockReceiptManual.store,
    amount: mockReceiptManual.amount,
    date: '2026-05-28',
    category: mockReceiptManual.category
  },
  {
    id: 'demo-receipt-2',
    store: mockReceiptVoice.store,
    amount: mockReceiptVoice.total,
    date: '2026-05-27',
    category: mockReceiptVoice.category,
    items: mockReceiptVoice.items
  },
  {
    id: 'demo-receipt-3',
    store: 'Ozon',
    amount: 8400,
    date: '2026-05-25',
    category: 'Развлечения',
    impulse_count: 1
  }
]
