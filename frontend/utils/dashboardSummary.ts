import type {
  CompareResponse,
  CreditsDashboardResponse,
  InsightItem,
  SankeyResponse,
  StoresResponse,
  TimeMachineResponse
} from '~/types/api'
import { ACTIONS, GOALS, HEALTH, CREDITS } from '~/constants/productCopy'
import { percentDti } from '~/utils/apiNormalize'

export type HealthTone = 'good' | 'warn' | 'risk'

export interface DashboardSummary {
  income: number
  expenses: number
  freeCashflow: number
  savingsBalance: number
  runwayMonths: number | null
  stabilityLabel: string
  stabilityTone: HealthTone
  healthTone: HealthTone
  healthEmoji: '🟢' | '🟡' | '🔴'
  healthHeadline: string
  mainRisk: string | null
  goalForecast: string
  habitInsight: string | null
  goalHint: string
  weeklyAction: string
  behaviorInsight: string | null
  dti: number | null
  dtiTone: HealthTone
}

const SAVINGS_TARGETS = new Set(['Накопления'])

/** Агрегаты из GET /dashboard/sankey — только links, как в API-контракте. */
function sumFromSankey(sankey: SankeyResponse | null): {
  income: number
  expenses: number
  savings: number
} {
  if (!sankey?.links?.length) {
    return { income: 0, expenses: 0, savings: 0 }
  }

  let savings = 0
  let expenses = 0

  for (const link of sankey.links) {
    if (SAVINGS_TARGETS.has(link.target)) {
      savings += link.value
    } else {
      expenses += link.value
    }
  }

  const incomeNode = sankey.nodes.find((n) => n.name === 'Зарплата')
  const fromLinks = sankey.links
    .filter((l) => l.source === 'Зарплата')
    .reduce((sum, l) => sum + l.value, 0)
  const income = incomeNode?.value ?? (fromLinks || savings + expenses)

  return { income, expenses, savings }
}

/** Есть ли данные GET /credits/dashboard для отображения блока. */
export function hasCreditsData(credits: CreditsDashboardResponse | null | undefined): boolean {
  if (!credits) return false
  return (
    (credits.credits?.length ?? 0) > 0 ||
    (credits.total_debt ?? 0) > 0 ||
    (credits.monthly_payments ?? 0) > 0
  )
}

function healthFromRunway(months: number | null): { label: string; tone: HealthTone } {
  if (months == null || months <= 0) {
    return { label: HEALTH.stabilityNeedData, tone: 'warn' }
  }
  if (months >= 6) return { label: HEALTH.stabilityGood, tone: 'good' }
  if (months >= 3) return { label: HEALTH.stabilityMid, tone: 'warn' }
  return { label: HEALTH.stabilityLow, tone: 'risk' }
}

function dtiTone(dti: number | null): HealthTone {
  if (dti == null) return 'warn'
  if (dti < 35) return 'good'
  if (dti < 50) return 'warn'
  return 'risk'
}

function compositeHealth(input: {
  runwayMonths: number | null
  dti: number | null
  freeCashflow: number
  stressMonths: number | null | undefined
}): { tone: HealthTone; emoji: '🟢' | '🟡' | '🔴'; headline: string } {
  const risks: HealthTone[] = []

  if (input.freeCashflow < 0) risks.push('risk')
  if (input.runwayMonths != null && input.runwayMonths < 3) risks.push('risk')
  if (input.dti != null && input.dti >= 50) risks.push('risk')
  if (input.stressMonths != null && input.stressMonths < 3) risks.push('risk')

  if (risks.includes('risk')) {
    return { tone: 'risk', emoji: '🔴', headline: HEALTH.attention }
  }

  const warns: HealthTone[] = []
  if (input.runwayMonths != null && input.runwayMonths < 6) warns.push('warn')
  if (input.dti != null && input.dti >= 35) warns.push('warn')
  if (input.freeCashflow > 0 && input.freeCashflow < 15000) warns.push('warn')

  if (warns.length) {
    return { tone: 'warn', emoji: '🟡', headline: HEALTH.caution }
  }

  return { tone: 'good', emoji: '🟢', headline: HEALTH.stable }
}

function buildMainRisk(input: {
  dti: number | null
  runwayMonths: number | null
  freeCashflow: number
  stressMonths: number | null | undefined
  topInsight: InsightItem | null
}): string | null {
  if (input.topInsight?.description ?? input.topInsight?.body) {
    return input.topInsight.description ?? input.topInsight.body ?? null
  }
  if (input.freeCashflow < 0) {
    return 'Тратите больше, чем получаете — запас будет уменьшаться.'
  }
  if (input.runwayMonths != null && input.runwayMonths < 3) {
    return `Запаса хватит примерно на ${input.runwayMonths} мес. — лучше иметь подушку побольше.`
  }
  if (input.dti != null && input.dti >= 50) {
    return CREDITS.highPaymentsRisk
  }
  if (input.stressMonths != null && input.stressMonths < 4) {
    return `Если доход пропадёт, запаса хватит примерно на ${input.stressMonths.toFixed(1)} мес.`
  }
  return null
}

export function buildGoalForecast(timemachine: TimeMachineResponse | null): string {
  if (!timemachine?.points?.length) {
    return GOALS.addExpensesForecast
  }

  const diff = timemachine.difference_final ?? timemachine.delta ?? 0
  const horizon = timemachine.points.length

  if (diff > 0) {
    const thousands = Math.round(diff / 1000)
    return `Если ничего не менять, за ${horizon} мес. накопите примерно на ${thousands} тыс. ₽ меньше, чем если чуть сократить траты.`
  }

  const last = timemachine.points[timemachine.points.length - 1]
  if (last) {
    const formatted = Math.round(last.actual / 1000)
    return `При тех же тратах через ${horizon} мес. на счёте может быть около ${formatted} тыс. ₽.`
  }

  return 'Накопления идут ровно — можно чуть ускорить путь к цели.'
}

/** Текст из GET /dashboard/compare → insights.biggest_change (без пересчёта). */
function compareInsightText(compare: CompareResponse | null): string | null {
  const change = compare?.insights?.biggest_change
  if (!change) return null
  const sign = change.delta > 0 ? '+' : ''
  const pctSign = change.delta_percent > 0 ? '+' : ''
  return `«${change.category}»: ${sign}${change.delta.toLocaleString('ru-RU')} ₽ (${pctSign}${change.delta_percent}%)`
}

export function buildDashboardSummary(input: {
  sankey: SankeyResponse | null
  compare: CompareResponse | null
  timemachine: TimeMachineResponse | null
  stores: StoresResponse | null
  credits: CreditsDashboardResponse | null
  topInsight: InsightItem | null
}): DashboardSummary {
  const { income: sankeyIncome, expenses: sankeyExpenses, savings: sankeySavings } =
    sumFromSankey(input.sankey)

  const monthlyIncome = input.credits?.monthly_income ?? sankeyIncome
  const monthlyExpenses = sankeyExpenses || 0
  const cushion = input.credits?.savings ?? sankeySavings
  const runwayMonths =
    monthlyExpenses > 0 && cushion > 0 ? Math.round((cushion / monthlyExpenses) * 10) / 10 : null

  const dtiRaw = input.credits?.dti
  const dti = dtiRaw != null ? percentDti(dtiRaw) : null
  const { label: stabilityLabel, tone: stabilityTone } = healthFromRunway(runwayMonths)
  const freeCashflow = monthlyIncome - monthlyExpenses
  const stressMonths = input.credits?.stress_test_months

  const { tone: healthTone, emoji: healthEmoji, headline: healthHeadline } = compositeHealth({
    runwayMonths,
    dti,
    freeCashflow,
    stressMonths
  })

  const goalForecast = buildGoalForecast(input.timemachine)
  const compareText = compareInsightText(input.compare)

  const diff = input.timemachine?.difference_final ?? input.timemachine?.delta ?? 0
  const goalHint =
    diff > 0
      ? GOALS.habitSavingsHint(Math.round(diff / 1000))
      : goalForecast

  const weeklyAction =
    input.topInsight?.title ??
    (input.compare?.insights?.biggest_change
      ? `Обратите внимание на «${input.compare.insights.biggest_change.category}» — траты выросли по сравнению с прошлым месяцем.`
      : ACTIONS.addPurchaseHint)

  const mainRisk = buildMainRisk({
    dti,
    runwayMonths,
    freeCashflow,
    stressMonths,
    topInsight: input.topInsight
  })

  return {
    income: monthlyIncome,
    expenses: monthlyExpenses,
    freeCashflow,
    savingsBalance: cushion,
    runwayMonths,
    stabilityLabel,
    stabilityTone,
    healthTone,
    healthEmoji,
    healthHeadline,
    mainRisk,
    goalForecast,
    habitInsight: compareText,
    goalHint,
    weeklyAction,
    behaviorInsight: input.topInsight?.description ?? input.topInsight?.body ?? compareText,
    dti,
    dtiTone: dtiTone(dti)
  }
}
