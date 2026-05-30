import type {
  ChallengeItem,
  DigestResponse,
  ForecastResponse,
  Goal,
  LeaderboardEntry,
  MortgageBreakdownResponse
} from '~/types/api'

export {
  mockSankey,
  mockStores,
  mockCategories,
  mockCompare,
  mockTimeMachine,
  mockInsights,
  mockCredits
} from '~/store/mocks/dashboard'

export { mockReceiptManual, mockReceiptVoice, mockReceiptFnsScan } from '~/store/mocks/receipts'

export { mockReceiptListItems } from '~/store/mocks/receiptList'

export { mockDiagnosis } from '~/store/mocks/diagnosis'

export { mockChatReplies } from '~/store/mocks/advisor'

export const mockForecast: ForecastResponse = {
  dates: ['2026-05-28', '2026-05-29', '2026-05-30', '2026-05-31'],
  forecast: [5200, 5100, 5300, 4950],
  upper_bound: [6240, 6120, 6360, 5940],
  lower_bound: [4160, 4080, 4240, 3960]
}

export const mockChallenges: ChallengeItem[] = [
  {
    id: '1',
    type: 'least_spend',
    title: 'Неделя без доставки',
    participants: 124,
    status: 'active',
    invite_token: 'invite-week-no-delivery'
  },
  {
    id: '2',
    type: 'streak',
    title: 'Кофе ≤ 3 раза',
    participants: 89,
    status: 'active',
    invite_token: 'invite-coffee-streak'
  }
]

export const mockLeaderboard: LeaderboardEntry[] = [
  { position: 1, username: 'Анна', relative_score: 0, display_name: 'Анна', rank: 1 },
  { position: 2, username: 'Иван', relative_score: 0.35, display_name: 'Иван', rank: 2 },
  { position: 3, username: 'Вы', relative_score: 0.52, display_name: 'Вы', rank: 3 }
]

export const mockDigest: DigestResponse = {
  period: { from: '2026-04-01', to: '2026-04-30' },
  total_spent: 145000,
  total_income: 180000,
  saved: 35000,
  by_category: [
    { name: 'Продукты', total: 52000, percent: 35.9, trend: '+8.3%' },
    { name: 'Кафе', total: 28000, percent: 19.3, trend: '+12%' }
  ],
  word_cloud: ['молоко', 'латте', 'хлеб', 'сыр', 'такси', 'доставка'],
  top_stores: [
    { name: 'Пятёрочка', total: 9100, visits: 14 },
    { name: 'Ozon', total: 8400, visits: 4 }
  ],
  mindfulness_rating: 72,
  ai_advice: 'Попробуйте сократить доставку — это около 9 000 ₽ в месяц',
  insights_summary: 'Найдено 2 скрытые подписки и 3 переплаты'
}

export const mockGoals: Goal[] = [
  {
    id: 'goal-1',
    title: 'Отпуск',
    target_amount: 150000,
    current_amount: 42000,
    progress_percent: 28,
    target_date: '2026-12-01',
    auto_save_percent: 10
  }
]

export const mockMortgageBreakdown: MortgageBreakdownResponse = {
  approval_level: 'medium',
  approval_reason:
    'Доход стабильный, на кредиты уходит умеренная доля. Запаса хватит около 4 месяцев — банк может попросить подтверждение.',
  safe_mortgage_amount: 9_000_000,
  comfortable_payment: 52_000,
  load_risk:
    'После платежа подушка сократится до ~3,2 мес. — ситуация станет более хрупкой, но управляемой.',
  scenario_now:
    'Платёж ~52 000 ₽ в месяц — до цели «Отпуск» пойдёте примерно на 5 месяцев дольше.',
  scenario_wait:
    'Если накопить ещё 180 000 ₽ и снизить долю на кредиты, через 8 месяцев условия могут стать мягче.',
  wait_months: 8,
  banks: [
    {
      id: 'sber',
      bank: 'Сбер',
      rate: 18.2,
      monthly_payment: 54_800,
      total_overpayment: 4_720_000,
      term_months: 240
    },
    {
      id: 'tinkoff',
      bank: 'Т-Банк',
      rate: 17.9,
      monthly_payment: 53_600,
      total_overpayment: 4_560_000,
      term_months: 240
    },
    {
      id: 'vtb',
      bank: 'ВТБ',
      rate: 18.5,
      monthly_payment: 55_400,
      total_overpayment: 4_880_000,
      term_months: 240
    },
    {
      id: 'alfa',
      bank: 'Альфа-Банк',
      rate: 18.0,
      monthly_payment: 54_100,
      total_overpayment: 4_640_000,
      term_months: 240
    }
  ],
  optimal_bank_id: 'tinkoff'
}
