import type { AiDiagnosisResponse } from '~/types/api'

/** GET /ai/diagnosis — по docs/api/API_Contract.md (v3). */
export const mockDiagnosis: AiDiagnosisResponse = {
  score: 72,
  grade: 'B',
  indicators: [
    { name: 'Платежи по кредитам', value: 28, norm: '<30', status: 'good' },
    { name: 'Запас на чёрный день', value: 4.2, norm: '>3', status: 'good' },
    { name: 'Сколько откладываете', value: 15, norm: '>20', status: 'warning' },
    { name: 'Покупки на эмоциях', value: 32, norm: '<25', status: 'critical' },
    { name: 'Насколько стабилен доход', value: 85, norm: '>70', status: 'good' }
  ],
  main_action: {
    title: 'Сократите доставку еды',
    description:
      'Вы тратите 9 000 ₽ в месяц на доставку. Готовьте дома 3 раза в неделю — это сэкономит 4 500 ₽ в месяц.',
    potential_savings: 4500,
    difficulty: 'easy'
  },
  next_check_days: 30
}
