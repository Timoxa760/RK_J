import type { AiDiagnosisResponse } from '~/types/api'

/** GET /ai/diagnosis — по docs/api/API_Contract.md (v3). */
export const mockDiagnosis: AiDiagnosisResponse = {
  score: 72,
  grade: 'B',
  indicators: [
    { name: 'Долговая нагрузка', value: 28, norm: '<30', status: 'good' },
    { name: 'Подушка безопасности', value: 4.2, norm: '>3', status: 'good' },
    { name: 'Накопления от дохода', value: 15, norm: '>20', status: 'warning' },
    { name: 'Импульсивные траты', value: 32, norm: '<25', status: 'critical' },
    { name: 'Стабильность доходов', value: 85, norm: '>70', status: 'good' }
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
