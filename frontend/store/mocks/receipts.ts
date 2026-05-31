import type { ReceiptManualResponse, ReceiptVoiceResponse } from '~/types/api'

/** Mock-данные receipt API — по docs/api/API_Contract.md (v3). */

export const mockReceiptManual: ReceiptManualResponse = {
  receipt_id: 'demo-receipt-manual-1',
  store: 'Пятёрочка',
  amount: 1032.5,
  category: 'Продукты',
  date: '2026-05-30T14:32:00Z',
  status: 'saved'
}

export const mockReceiptVoice: ReceiptVoiceResponse = {
  receipt_id: 'demo-receipt-voice-1',
  store: 'Пятёрочка',
  items: [
    { name: 'Молоко', price: 89.9, quantity: 1 },
    { name: 'Хлеб', price: 45.5, quantity: 1 }
  ],
  total: 135.4,
  category: 'Продукты',
  confidence: 0.92
}
