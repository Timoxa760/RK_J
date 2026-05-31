/** Единый копирайт приложения — простой язык, без жаргона. */

export const NAV = {
  dashboard: 'План',
  dashboardSubtitle: 'Советы, план и картина денег',
  advisor: 'Советник',
  advisorSubtitle: 'Спросите про план, цель и траты',
  receipts: 'Расходы',
  analytics: 'Советник',
  credits: 'Кредиты',
  creditsTitle: 'Кредиты: насколько спокойно',
  creditsSubtitle: 'Скан договора PDF',
  social: 'Привычки',
  socialSubtitle: 'Задания с друзьями',
  profile: 'Профиль'
} as const

export const ACTIONS = {
  addPurchaseHint: 'Добавьте покупку — цифры обновятся',
  addPurchaseOrReceipt: 'Добавьте покупку — подскажем следующий шаг',
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
  cushionTitle: 'Запас на чёрный день',
  scanConfidence: (pct: number) =>
    `Уверенность распознавания: ${Math.round(pct <= 1 ? pct * 100 : pct)}%`,
  scanBenchmark: (rate: number) => `Средняя ставка по рынку: ~${rate}%`,
  scanRateAbove: 'Ставка выше рынка',
  scanRateAtOrBelow: 'Ставка на уровне или ниже рынка',
  scanRateUnknown: 'Сравнение со рынком недоступно',
  scanReportTitle: 'Разбор договора',
  scanAddFixed: 'Внести в обязательные расходы',
  scanFixedDone: 'Платёж учтён в обязательных расходах',
  dtiNeedIncome: 'Платёж по договору есть — укажите доход, чтобы показать долю в процентах.',
  dtiIncomeCta: 'Указать доход в профиле',
  monthlyPaymentsLine: (amount: number) =>
    `Суммарный платёж по договорам: ${amount.toLocaleString('ru-RU')} ₽/мес`,
  cushionNeedData: 'Укажите подушку в профиле — покажем, на сколько месяцев хватит запаса.',
  deleteCredit: 'Удалить из списка'
} as const

export const GOALS = {
  addExpensesForecast: 'Добавьте расходы — покажем, сколько можете отложить',
  opportunityLabel: 'Ваша возможность',
  opportunityAmount: (thousands: number) => `${thousands} тыс. ₽`,
  savingsOpportunity: (thousands: number, months: number) =>
    `Можете получить ${thousands} тыс. ₽ за ${months} мес., если чуть сократите траты.`,
  savingsPain: (thousands: number, months: number) =>
    `Сейчас упускаете эту сумму — без изменений за ${months} мес. накопите на ${thousands} тыс. ₽ меньше.`,
  savingsOnTrack: (thousands: number, months: number) =>
    `При тех же тратах через ${months} мес. на счёте может быть ~${thousands} тыс. ₽ — темп нормальный.`,
  savingsEven: 'Накопления идут ровно — один небольшой шаг ускорит путь к цели.',
  savingsCurrentLabel: 'Текущие накопления',
  savingsChartOpportunity: (gain: string, total: string) =>
    `Можете получить ${gain}, если чуть экономить — на счёте будет ${total}.`,
  habitSavingsHint: (thousands: number) =>
    `Сократите траты — и получите ещё ~${thousands} тыс. ₽ к запасу.`
} as const

export const PURCHASES = {
  impulseBadge: 'лишнее',
  impulseCount: (n: number) => `${n} лишн.`,
  onEmotion: 'на эмоциях',
  goalDelay: 'Эта покупка отодвинет цель примерно на',
  deleteReceipt: 'Удалить покупку',
  deleteConfirm: 'Удалить эту покупку? Её не будет в картине расходов.',
  fnsBadge: 'ФНС'
} as const

export const ANALYTICS = {
  headline: 'Что будет, если так же тратить',
  attentionTitle: 'На что обратить внимание',
  moreTips: (n: number) => `Ещё ${n} ${n === 1 ? 'совет' : n < 5 ? 'совета' : 'советов'}`,
  anomaly: 'Необычно много потратили',
  simulator: 'А если меньше тратить на…',
  savingsChart: 'Как пойдут накопления',
  trajectoryOnMain: 'Прогноз накоплений — ниже на этой странице'
} as const

export const ADVISOR = {
  chatTitle: 'Советник',
  chatHint: 'Вопрос снизу или свой текст — ответит советник',
  chatReset: 'Очистить диалог',
  chatLocalReply: 'Локальный ответ',
  chatWeeklyHint: 'Совет недели',
  planTitle: 'Ваш план',
  planTitleMega: 'Финансовый план',
  planHint: 'Цель, срок и три ближайших шага',
  planHintMega: 'Цель, диагностика, деньги и сценарии — всё в одном месте',
  planTabSteps: 'Шаги',
  planTabOpportunity: 'Возможность',
  planTabDiagnosis: 'Диагностика',
  planTabMoney: 'Картина денег',
  planTabCredits: 'Кредиты',
  planTabExplore: 'Что если',
  mindfulnessTitle: 'Траты под контролем',
  weeklyAdviceTitle: 'Совет недели',
  weeklyAdviceHint: 'Один шаг из вашего финансового плана',
  weeklyAdviceHintShort: 'Первый шаг — в вашем плане ниже',
  planStepPrimaryLabel: 'Главное на неделю',
  planTabOpportunityHint: 'Сколько можете получить, если чуть изменить траты',
  planTabDiagnosisHint: 'Показатели, из которых складывается оценка',
  planTabMoneyHint: 'Куда уходят деньги и как растут накопления',
  planTabCreditsHint: 'Доля дохода на погашение кредитов',
  planTabExploreHint: 'Почему совет именно такой и что будет, если изменить траты',
  diagnosisStatusUrgent: 'Срочно',
  addPurchaseLabel: 'Добавить покупку',
  addPurchaseAria: 'Голосом или вручную — совет станет точнее',
  chatPlaceholder: 'Напишите вопрос…',
  chatEmptyHint: 'Выберите вопрос ниже или напишите свой — ответит советник',
  chatHistoryToday: 'Сегодня',
  chatSourceAi: 'Ответ ИИ',
  chatSourceHeuristic: 'Базовый ответ',
  chatSourceLocal: 'Локальный ответ',
  askCloser: (goalTitle: string) => `Приблизить «${goalTitle}»`,
  askAboutAction: 'Спросить советника',
  askGettingStarted: 'Спросить, с чего начать',
  askGoalDelay: 'Как не отодвигать цель?',
  forecastTitle: 'Прогноз трат (7 дней)',
  savingsTitle: 'Как пойдут накопления',
  categoriesTitle: 'Категории'
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
  formTitle: 'Финансовая модель',
  formHint: 'Доход, запас и цель — основа прогноза на главной и в советнике',
  emptyModel: 'Пока мало данных — пройдите короткий опрос',
  fixedExpensesTitle: 'Обязательные расходы',
  fixedExpensesHint:
    'Аренда, кредиты, связь — учитываются в прогнозе. Можно добавить, если пропустили на опросе.',
  addFixedExpense: 'Добавить платёж'
} as const

export const FNS = {
  title: 'Чеки ФНС',
  hint: 'Подключите «Мой налог» — покупки с кассы сами попадут в расходы и аналитику.',
  phoneLabel: 'Телефон из «Мой налог»',
  connectDialog: 'Войдите тем же номером, что в приложении «Мой налог». Мы отправим код подтверждения.',
  codeLabel: 'Код из СМС',
  codePlaceholder: '6 цифр',
  sendCode: 'Получить код',
  sendingCode: 'Отправляем код…',
  codeSent: (phone: string) => `Код отправлен на ${phone}. Введите его ниже.`,
  codeRequired: 'Введите код из СМС',
  confirmCode: 'Подтвердить',
  verifying: 'Проверяем код…',
  changePhone: 'Другой номер',
  resendCode: 'Отправить код ещё раз',
  connect: 'Подключить ФНС',
  connecting: 'Подключаем…',
  connected: 'ФНС подключена — чеки добавлены в расходы',
  disconnected: 'ФНС отключена',
  disconnect: 'Отключить',
  sync: 'Обновить чеки',
  syncing: 'Синхронизация…',
  syncImported: (count: number) =>
    `Добавили ${count} ${count === 1 ? 'чек' : count < 5 ? 'чека' : 'чеков'} в расходы`,
  syncUpToDate: 'Новых чеков нет — всё уже в ленте',
  syncFailed: 'Не удалось синхронизировать чеки',
  connectFailed: 'Не удалось подключить ФНС',
  phoneRequired: 'Укажите телефон',
  statusConnected: 'Подключено',
  statusDisconnected: 'Не подключено',
  importedLine: (count: number) => `В расходах: ${count} чеков из ФНС`,
  lastSync: (when: string) => `Последняя синхронизация: ${when}`,
  pendingLine: (count: number) =>
    `Ещё ${count} ${count === 1 ? 'новый чек' : count < 5 ? 'новых чека' : 'новых чеков'} — нажмите «Обновить чеки».`,
  openReceipts: 'Открыть ленту расходов'
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
