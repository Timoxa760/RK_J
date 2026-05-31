export interface FinancialReportLoadingStage {
  id: string
  label: string
  hint?: string
  /** Сколько держим этап, пока запрос ещё идёт */
  durationMs?: number
}

export const FINANCIAL_REPORT_LOADING_STAGES: FinancialReportLoadingStage[] = [
  {
    id: 'profile',
    label: 'Читаем профиль',
    hint: 'Доход, цель и подушка',
    durationMs: 2200
  },
  {
    id: 'expenses',
    label: 'Смотрим траты',
    hint: 'Категории и динамика расходов',
    durationMs: 2800
  },
  {
    id: 'obligations',
    label: 'Учитываем обязательства',
    hint: 'Кредиты и фиксированные платежи',
    durationMs: 2400
  },
  {
    id: 'stability',
    label: 'Оцениваем устойчивость',
    hint: 'Запас, прогноз цели и риски',
    durationMs: 3200
  },
  {
    id: 'plan',
    label: 'Формируем план',
    hint: 'Шаги и персональный совет',
    durationMs: 4500
  }
]

export const ONBOARDING_DIAGNOSIS_LOADING_STAGES: FinancialReportLoadingStage[] = [
  {
    id: 'save',
    label: 'Сохраняем ответы',
    hint: 'Фиксируем доход, цель и расходы',
    durationMs: 1600
  },
  {
    id: 'model',
    label: 'Строим модель',
    hint: 'Считаем запас и путь к цели',
    durationMs: 2600
  },
  {
    id: 'action',
    label: 'Подбираем первое действие',
    hint: 'Один конкретный шаг на старт',
    durationMs: 3800
  }
]

export const FINANCIAL_REPORT_LOADING_COPY = {
  title: 'Собираем финансовый отчёт',
  subtitle: 'Обычно это занимает до минуты — идём по шагам',
  waitMore: 'Ещё немного — дожимаем детали',
  onboardingTitle: 'Готовим вашу картину',
  onboardingSubtitle: 'Считаем модель и первое действие'
} as const
