/** Единый копирайт приложения — простой язык, без жаргона. */

export const NAV = {
  dashboard: 'Моя картина',
  dashboardSubtitle: 'Доходы, расходы и запас',
  receipts: 'Расходы',
  analytics: 'Прогноз',
  credits: 'Кредиты',
  creditsTitle: 'Кредиты: насколько спокойно',
  creditsSubtitle: 'Сколько уходит на погашение',
  social: 'Привычки',
  socialSubtitle: 'Задания с друзьями',
  digest: 'Сводка',
  profile: 'Профиль'
} as const

export const ACTIONS = {
  addPurchaseHint: 'Добавьте покупку — цифры обновятся',
  addPurchaseOrReceipt: 'Добавьте покупку или чек — подскажем следующий шаг',
  weeklyActionLabel: 'Что сделать на этой неделе',
  pageSummaryAria: 'Кратко о странице',
  oneTipNow: 'Совет на сейчас',
  allTips: 'Все советы',
  viewCategoriesOnMain: 'Смотреть категории на главной'
} as const

export const HEALTH = {
  stable: 'С деньгами сейчас спокойно',
  caution: 'Запас есть, но лучше не расслабляться',
  attention: 'Финансовое положение требует внимания',
  attentionPrefix: 'На что обратить внимание',
  stabilityGood: 'Устойчиво',
  stabilityMid: 'Средний запас',
  stabilityLow: 'Мало запаса',
  stabilityNeedData: 'Нужны данные',
  reserveMonths: (months: number) =>
    `Запас примерно на ${months} мес. при текущих расходах.`,
  reserveUnknown: 'Укажите расходы — скажем, на сколько месяцев хватит запаса',
  leftAfterExpenses: 'Остаётся после расходов'
} as const

export const CREDITS = {
  paymentsTitle: 'Платежи по кредитам',
  paymentsHint: 'Сколько процентов дохода уходит на погашение. До 35% — обычно комфортно',
  incomeShare: (pct: number) => `На кредиты уходит ${pct}% дохода`,
  stressIncomeDrop: (pct: number) =>
    `Если доход упадёт на 20%, на погашение уйдёт около ${pct}%`,
  stressReserveMonths: (months: number) =>
    `Если доход упадёт на 20%, запаса хватит примерно на ${months.toFixed(1)} мес.`,
  trafficLight: 'Кредиты: насколько спокойно',
  anotherLoan: 'Потяну ли ещё один кредит?',
  highPaymentsRisk:
    'На кредиты уходит много — новый кредит сделает положение тяжелее',
  cushionTitle: 'Запас на чёрный день'
} as const

export const GOALS = {
  addExpensesForecast: 'Добавьте расходы — покажем, когда примерно дойдёте до цели',
  ifSameSpending: 'К чему идёте при тех же тратах',
  habitSavingsHint: (thousands: number) =>
    `Если чуть сократить траты, к цели можно прийти на ${thousands} тыс. ₽ раньше`
} as const

export const PURCHASES = {
  impulseBadge: 'лишнее',
  impulseCount: (n: number) => `${n} лишн.`,
  onEmotion: 'на эмоциях',
  goalDelay: 'Эта покупка отодвинет цель примерно на'
} as const

export const ANALYTICS = {
  headline: 'Что будет, если так же тратить',
  attentionTitle: 'На что обратить внимание',
  moreTips: (n: number) => `Ещё ${n} ${n === 1 ? 'совет' : n < 5 ? 'совета' : 'советов'}`,
  anomaly: 'Необычно много потратили',
  simulator: 'А если меньше тратить на…',
  savingsChart: 'Как пойдут накопления',
  trajectoryOnMain: 'Полная картина «если ничего не менять» — на главной'
} as const

export const DIGEST = {
  loading: 'Сводка за период',
  headline: 'Как складывается месяц',
  badge: 'Сводка',
  mindfulness: 'Траты под контролем',
  leftAfterExpenses: (amount: string) => `После расходов остаётся ${amount}`
} as const

export const SOCIAL = {
  headline: 'Привычки в компании',
  habitsTitle: 'Как у вас с тратами',
  habitsHint: 'Без сравнения сумм — только ваш ритм',
  leaderboard: 'Рейтинг друзей',
  challenges: 'Задания',
  optionalBadge: 'По желанию'
} as const

export const PROFILE = {
  headline: 'Ваши цифры',
  formTitle: 'Ваши цифры',
  formHint: 'Доход и запас нужны для прогноза на главной и в прогнозе',
  emptyModel: 'Пока мало данных — пройдите короткий опрос'
} as const

export const ONBOARDING = {
  shellTitle: 'Ваши цифры',
  shellSubtitle: 'Несколько простых шагов — и всё станет наглядно',
  welcomeBullet1: 'Спросим про доход, запас и цель — простыми словами',
  welcomeBullet2: 'Покажем, когда примерно дойдёте до цели',
  welcomeBullet3: 'Подскажем первый шаг — минута вашего времени',
  summaryTitle: 'Вот ваша картина',
  summarySubtitle: 'Коротко — и что сделать дальше',
  fixedPayments: 'постоянные платежи',
  recording: 'Записываем…'
} as const

export function formatRub(value: number): string {
  return `${value.toLocaleString('ru-RU')} ₽`
}
