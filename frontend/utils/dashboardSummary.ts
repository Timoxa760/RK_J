import type {
  CategoriesResponse,
  CreditsDashboardResponse,
  FinancialProfile,
  InsightItem,
  TimeMachineResponse
} from '~/types/api'
import { ACTIONS, GOALS, HEALTH, CREDITS, formatRub } from '~/constants/productCopy'
import { percentDti } from '~/utils/apiNormalize'
import { isPlaceholderTimemachine } from '~/utils/dashboardProjections'

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
  goalOpportunityThousands: number | null
  habitInsight: string | null
  goalHint: string
  weeklyAction: string
  behaviorInsight: string | null
  dti: number | null
  dtiTone: HealthTone
}

function sumFixedExpenses(profile: FinancialProfile | null | undefined): number {
  return (profile?.fixed_expenses ?? []).reduce((sum, row) => sum + (row.amount ?? 0), 0)
}

function profileIncome(profile: FinancialProfile | null | undefined): number {
  if (!profile) return 0
  return (profile.active_income ?? 0) + (profile.passive_income ?? 0)
}

/** Агрегаты из профиля, категорий и credits dashboard. */
function sumFromProfile(input: {
  profile: FinancialProfile | null | undefined
  credits: CreditsDashboardResponse | null
  categories?: CategoriesResponse | null
}): { income: number; expenses: number; savings: number } {
  const profileIncomeValue = profileIncome(input.profile)
  const fixedExpenses = sumFixedExpenses(input.profile)
  const variableExpenses = (input.categories?.categories ?? []).reduce(
    (sum, row) => sum + (row.amount ?? row.total ?? 0),
    0
  )
  const income = input.credits?.monthly_income ?? profileIncomeValue
  const expenses =
    fixedExpenses > 0 && variableExpenses > 0
      ? fixedExpenses + variableExpenses
      : variableExpenses > 0
        ? variableExpenses
        : fixedExpenses
  const savings = input.credits?.savings ?? input.profile?.emergency_fund ?? 0
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

export function buildGoalForecast(
  timemachine: TimeMachineResponse | null,
  profile?: FinancialProfile | null
): string {
  const savingsBalance = profile?.emergency_fund ?? 0
  const useTimemachine =
    timemachine?.points.length &&
    !isPlaceholderTimemachine(timemachine, savingsBalance)

  if (!useTimemachine) {
    return buildProfileGoalHint(profile)
  }

  const diff = timemachine!.difference_final ?? timemachine!.delta ?? 0
  const horizon = timemachine!.points.length

  if (diff > 0) {
    const thousands = Math.round(diff / 1000)
    return GOALS.savingsOpportunity(thousands, horizon)
  }

  const last = timemachine!.points[timemachine!.points.length - 1]
  if (last) {
    const formatted = Math.round(last.actual / 1000)
    return GOALS.savingsOnTrack(formatted, horizon)
  }

  return GOALS.savingsEven
}

function buildProfileGoalHint(profile: FinancialProfile | null | undefined): string {
  if (!profile || profile.skipped_goal || (profile.goal_amount ?? 0) < 1000) {
    return GOALS.addExpensesForecast
  }

  const title = profile.goal_title?.trim() || 'Цель'
  const amount = profile.goal_amount ?? 0
  const monthlySaving =
    profile.active_income > 0 ? Math.round(profile.active_income * 0.1) : 0

  if (monthlySaving > 0) {
    const months = Math.ceil(amount / monthlySaving)
    return `Цель «${title}» — ${formatRub(amount)}. При ~${formatRub(monthlySaving)}/мес. ориентир ${months} мес.`
  }

  return `Цель «${title}» — ${formatRub(amount)}.`
}

function goalOpportunityThousands(
  timemachine: TimeMachineResponse | null,
  savingsBalance = 0
): number | null {
  if (!timemachine?.points.length || isPlaceholderTimemachine(timemachine, savingsBalance)) {
    return null
  }
  const diff = timemachine.difference_final ?? timemachine.delta ?? 0
  if (diff <= 0) return null
  return Math.round(diff / 1000)
}

export function buildDashboardSummary(input: {
  profile: FinancialProfile | null | undefined
  timemachine: TimeMachineResponse | null
  credits: CreditsDashboardResponse | null
  topInsight: InsightItem | null
  categories?: CategoriesResponse | null
}): DashboardSummary {
  const { income, expenses, savings } = sumFromProfile({
    profile: input.profile,
    credits: input.credits,
    categories: input.categories
  })

  const monthlyIncome = income
  const monthlyExpenses = expenses
  const cushion = savings
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

  const goalForecast = buildGoalForecast(input.timemachine, input.profile)
  const opportunityThousands = goalOpportunityThousands(
    input.timemachine,
    input.profile?.emergency_fund ?? 0
  )
  const insightText =
    input.topInsight?.description ?? input.topInsight?.body ?? input.topInsight?.title ?? null

  const horizon = input.timemachine?.points?.length ?? 0
  const goalHint =
    opportunityThousands != null && horizon > 0
      ? GOALS.savingsPain(opportunityThousands, horizon)
      : goalForecast

  const weeklyAction = input.topInsight?.title ?? ACTIONS.addPurchaseHint

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
    goalOpportunityThousands: opportunityThousands,
    habitInsight: insightText,
    goalHint,
    weeklyAction,
    behaviorInsight: insightText,
    dti,
    dtiTone: dtiTone(dti)
  }
}
