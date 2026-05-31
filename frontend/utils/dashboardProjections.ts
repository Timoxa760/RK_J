import type {
  CategoriesResponse,
  FinancialProfile,
  ForecastResponse,
  TimeMachineResponse
} from '~/types/api'
import { buildUserCategoryOptions } from '~/constants/expenseCategories'
import { normalizeForecast, normalizeTimeMachine } from '~/utils/apiNormalize'

function sumCategories(categories: CategoriesResponse | null | undefined): number {
  return (categories?.categories ?? []).reduce(
    (sum, row) => sum + (row.amount ?? row.total ?? 0),
    0
  )
}

function sumFixedExpenses(profile: FinancialProfile | null | undefined): number {
  return (profile?.fixed_expenses ?? []).reduce((sum, row) => sum + (row.amount ?? 0), 0)
}

function profileIncome(profile: FinancialProfile | null | undefined): number {
  if (!profile || profile.skipped_income) return 0
  return (profile.active_income ?? 0) + (profile.passive_income ?? 0)
}

function incomeKnown(profile: FinancialProfile | null | undefined): boolean {
  return !profile?.skipped_income && profileIncome(profile) > 0
}

function monthSpend(profile: FinancialProfile | null | undefined, categories: CategoriesResponse | null) {
  const variable = sumCategories(categories)
  const fixed = sumFixedExpenses(profile)
  if (fixed > 0 && variable > 0) return fixed + variable
  if (variable > 0) return variable
  return fixed
}

/** Demo GET /forecast: 9800 + i×120 за 7 дней */
export function isPlaceholderForecast(raw: ForecastResponse | null | undefined): boolean {
  if (!raw) return false
  const normalized = normalizeForecast(raw)
  const values = normalized.forecast
  if (values.length !== 7) return false
  return values.every((value, index) => value === 9800 + index * 120)
}

/** Demo POST /scenarios/simulate или битый timemachine с базой 500k */
export function isPlaceholderTimemachine(
  raw: TimeMachineResponse | null | undefined,
  savingsBalance = 0
): boolean {
  if (!raw?.points.length) return false

  const first = raw.points[0]?.actual ?? 0
  if (first >= 499_000 && first <= 501_000) {
    if (raw.points.length > 1) {
      const step = raw.points[1]!.actual - raw.points[0]!.actual
      if (Math.abs(step - 12_000) < 1) return true
    }
    return raw.points.length <= 60
  }

  // PG раньше отдавал накопительную сумму трат, а не накоплений
  if (savingsBalance > 0 && raw.points.length <= 12) {
    const last = raw.points[raw.points.length - 1]?.actual ?? 0
    if (last > savingsBalance * 1.2 && last > 10_000) {
      const monotonic = raw.points.every((point, index) => {
        if (index === 0) return true
        return point.actual >= raw.points[index - 1]!.actual
      })
      if (monotonic && last !== savingsBalance) return true
    }
  }

  return false
}

/** Прогноз дневных трат на N дней из профиля и категорий текущего месяца */
export function buildSpendForecast(
  profile: FinancialProfile | null | undefined,
  categories: CategoriesResponse | null,
  days = 7
): ForecastResponse {
  const monthly = monthSpend(profile, categories)
  if (monthly <= 0) {
    return { dates: [], forecast: [] }
  }

  const now = new Date()
  const dayOfMonth = Math.max(now.getDate(), 1)
  const variable = sumCategories(categories)
  const dailyRate =
    variable > 0 ? variable / dayOfMonth : monthly / new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()

  const dates: string[] = []
  const forecast: number[] = []

  for (let i = 0; i < days; i++) {
    const date = new Date(now)
    date.setDate(date.getDate() + i)
    dates.push(date.toISOString().slice(0, 10))
    forecast.push(Math.max(0, Math.round(dailyRate)))
  }

  return { dates, forecast }
}

/** Прогноз накоплений: подушка + (доход − траты) × месяцы */
export function buildSavingsTimemachine(
  profile: FinancialProfile | null | undefined,
  categories: CategoriesResponse | null,
  months = 6
): TimeMachineResponse {
  const income = profileIncome(profile)
  const expenses = monthSpend(profile, categories)
  const monthlySaving = income - expenses
  const start = Math.max(0, profile?.emergency_fund ?? 0)

  if (start <= 0 && income <= 0 && expenses <= 0) {
    return { points: [], delta: 0, difference_final: 0 }
  }

  const points: TimeMachineResponse['points'] = []
  let actual = start
  let optimistic = start
  const now = new Date()

  for (let i = 1; i <= months; i++) {
    actual += monthlySaving
    optimistic += monthlySaving > 0 ? Math.round(monthlySaving * 1.1) : monthlySaving

    const date = new Date(now.getFullYear(), now.getMonth() + i, 1)
    const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
    points.push({
      month,
      actual: Math.max(0, Math.round(actual)),
      optimistic: Math.max(0, Math.round(optimistic))
    })
  }

  const diff = points.length ? points[points.length - 1]!.optimistic - points[points.length - 1]!.actual : 0

  return {
    points,
    delta: diff,
    difference_final: diff
  }
}

export function resolveSpendForecast(
  api: ForecastResponse | null | undefined,
  profile: FinancialProfile | null | undefined,
  categories: CategoriesResponse | null,
  days = 7
): ForecastResponse | null {
  if (api && !isPlaceholderForecast(api)) {
    const normalized = normalizeForecast(api)
    if (normalized.forecast.length) return normalized
  }
  const built = buildSpendForecast(profile, categories, days)
  return built.forecast.length ? built : null
}

export function resolveSavingsTimemachine(
  api: TimeMachineResponse | null | undefined,
  profile: FinancialProfile | null | undefined,
  categories: CategoriesResponse | null,
  months = 6
): TimeMachineResponse | null {
  const savingsBalance = profile?.emergency_fund ?? 0
  if (api?.points.length && !isPlaceholderTimemachine(api, savingsBalance)) {
    return normalizeTimeMachine(api)
  }
  const built = buildSavingsTimemachine(profile, categories, months)
  return built.points.length ? built : null
}

/** Локальный расчёт эффекта сценария без demo API */
export interface ScenarioPreview {
  categoryName: string
  reductionPercent: number
  months: number
  /** Траты по выбранной категории за текущий месяц (из API, нормализованные) */
  categorySpend: number
  monthlySaving: number
  totalGain: number
  /** Подушка / emergency_fund из профиля */
  currentBalance: number
  freeCashflow: number
  baselineEnd: number
  optimizedEnd: number
  /** Доход известен — можно прогнозировать подушку, не только экономию по категории */
  incomeKnown: boolean
  /** Прирост подушки за период без сокращения выбранной категории */
  baselineGain: number
  /** Прирост с учётом сокращения (baselineGain + totalGain) */
  optimizedGain: number
  hasData: boolean
  message: string
}

function categorySpendByName(
  categories: CategoriesResponse | null,
  categoryName: string
): number {
  if (!categoryName) return 0
  const options = buildUserCategoryOptions(categories)
  return options.find((row) => row.name === categoryName)?.amount ?? 0
}

function totalCategorySpend(categories: CategoriesResponse | null): number {
  return sumCategories(categories)
}

function monthlyFreeCashflow(
  profile: FinancialProfile | null | undefined,
  categories: CategoriesResponse | null
): number {
  return profileIncome(profile) - monthSpend(profile, categories)
}

export function buildScenarioPreview(input: {
  profile: FinancialProfile | null | undefined
  categories: CategoriesResponse | null
  categoryName: string
  reductionPercent: number
  months?: number
}): ScenarioPreview {
  const months = input.months ?? 12
  const categoryOptions = buildUserCategoryOptions(input.categories)
  const categoryName =
    categoryOptions.find((row) => row.name === input.categoryName)?.name ??
    input.categoryName
  const categorySpend = categorySpendByName(input.categories, categoryName)
  const currentBalance = Math.max(0, input.profile?.emergency_fund ?? 0)
  const freeCashflow = monthlyFreeCashflow(input.profile, input.categories)
  const totalSpend = totalCategorySpend(input.categories)
  const monthlySaving = Math.round(categorySpend * (input.reductionPercent / 100))
  const totalGain = monthlySaving * months

  const knownIncome = incomeKnown(input.profile)
  const baselineEnd = Math.max(
    0,
    Math.round(currentBalance + freeCashflow * months)
  )
  const optimizedEnd = Math.max(
    0,
    Math.round(currentBalance + (freeCashflow + monthlySaving) * months)
  )
  const baselineGain = knownIncome ? Math.max(0, Math.round(freeCashflow * months)) : 0
  const optimizedGain = baselineGain + totalGain

  if (totalSpend <= 0) {
    return {
      categoryName,
      reductionPercent: input.reductionPercent,
      months,
      categorySpend: 0,
      monthlySaving: 0,
      totalGain: 0,
      currentBalance,
      freeCashflow,
      baselineEnd: currentBalance,
      optimizedEnd: currentBalance,
      incomeKnown: knownIncome,
      baselineGain: 0,
      optimizedGain: 0,
      hasData: false,
      message: 'Добавьте покупки — покажем эффект от сокращения трат.'
    }
  }

  if (categorySpend <= 0) {
    return {
      categoryName,
      reductionPercent: input.reductionPercent,
      months,
      categorySpend: 0,
      monthlySaving: 0,
      totalGain: 0,
      currentBalance,
      freeCashflow,
      baselineEnd: currentBalance,
      optimizedEnd: currentBalance,
      incomeKnown: knownIncome,
      baselineGain: 0,
      optimizedGain: 0,
      hasData: false,
      message: `За этот месяц нет трат в «${categoryName}». Выберите категорию, где уже есть расходы.`
    }
  }

  const message = knownIncome
    ? `Сократите «${categoryName}» на ${input.reductionPercent}% — +${monthlySaving.toLocaleString('ru-RU')} ₽/мес к свободным деньгам (сейчас ${categorySpend.toLocaleString('ru-RU')} ₽/мес в категории).`
    : `Сократите «${categoryName}» на ${input.reductionPercent}% — освободите ~${monthlySaving.toLocaleString('ru-RU')} ₽/мес. Укажите доход в профиле для прогноза подушки.`

  return {
    categoryName,
    reductionPercent: input.reductionPercent,
    months,
    categorySpend,
    monthlySaving,
    totalGain,
    currentBalance,
    freeCashflow,
    baselineEnd,
    optimizedEnd,
    incomeKnown: knownIncome,
    baselineGain,
    optimizedGain,
    hasData: true,
    message
  }
}

export function buildScenarioResult(input: {
  profile: FinancialProfile | null | undefined
  categories: CategoriesResponse | null
  categoryName: string
  reductionPercent: number
  months?: number
}): { message: string; timemachine: TimeMachineResponse | null; preview: ScenarioPreview } {
  const preview = buildScenarioPreview(input)
  const months = input.months ?? 12

  if (!preview.hasData) {
    return { message: preview.message, timemachine: null, preview }
  }

  const freeCashflow = preview.freeCashflow
  const monthlySaving = preview.monthlySaving
  const start = preview.currentBalance
  const points: TimeMachineResponse['points'] = []
  let baseline = start
  let optimized = start
  const now = new Date()

  for (let i = 1; i <= Math.min(months, 12); i++) {
    baseline += freeCashflow
    optimized += freeCashflow + monthlySaving
    const date = new Date(now.getFullYear(), now.getMonth() + i, 1)
    points.push({
      month: `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`,
      actual: Math.max(0, Math.round(baseline)),
      optimistic: Math.max(0, Math.round(optimized))
    })
  }

  return {
    message: preview.message,
    preview,
    timemachine: {
      points,
      delta: preview.totalGain,
      difference_final: preview.totalGain
    }
  }
}
