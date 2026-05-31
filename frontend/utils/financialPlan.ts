import type {
  AiDiagnosisResponse,
  Goal,
  InsightItem,
  TimeMachineResponse
} from '~/types/api'
import type { DashboardSummary } from '~/utils/dashboardSummary'
import { buildGoalProgressText } from '~/utils/pageNarrative'
import { formatRub } from '~/constants/productCopy'

export interface FinancialPlanStep {
  title: string
  description: string
}

export interface FinancialPlan {
  goalTitle: string
  goalProgress: string
  steps: FinancialPlanStep[]
  runwayText: string | null
  freeCashflowText: string | null
  updatedAt: number
}

const GOAL_STEP_TITLE = 'Проверьте прогресс цели'

/** Подставляет цель из профиля, если API-план её не знает. */
export function applyProfileGoalToPlan(plan: FinancialPlan, primaryGoal: Goal | null): FinancialPlan {
  if (!primaryGoal) return plan

  const goalProgress = buildGoalProgressText(primaryGoal)
  const steps = plan.steps.map((step) =>
    step.title === GOAL_STEP_TITLE ? { ...step, description: goalProgress } : step
  )

  return {
    ...plan,
    goalTitle: primaryGoal.title,
    goalProgress,
    steps
  }
}

export function buildFinancialPlan(input: {
  primaryGoal: Goal | null
  summary: DashboardSummary
  timemachine: TimeMachineResponse | null
  diagnosis: AiDiagnosisResponse | null
  topInsight: InsightItem | null
}): FinancialPlan {
  const goalTitle = input.primaryGoal?.title ?? 'Финансовая цель'
  const goalProgress = buildGoalProgressText(input.primaryGoal)

  const steps: FinancialPlanStep[] = []

  if (input.diagnosis?.main_action) {
    steps.push({
      title: input.diagnosis.main_action.title,
      description: input.diagnosis.main_action.description
    })
  }

  if (input.topInsight) {
    const insightTitle = input.topInsight.title
    if (!steps.some((s) => s.title === insightTitle)) {
      steps.push({
        title: insightTitle,
        description: input.topInsight.body ?? input.topInsight.description ?? ''
      })
    }
  }

  const goalStepDescription = buildGoalProgressText(input.primaryGoal)
  if (!steps.some((s) => s.description === goalStepDescription)) {
    steps.push({
      title: 'Проверьте прогресс цели',
      description: goalStepDescription
    })
  }

  while (steps.length < 3) {
    steps.push({
      title: 'Добавляйте покупки',
      description: 'Записывайте траты голосом — советы станут точнее.'
    })
  }

  const runwayText =
    input.summary.runwayMonths != null && input.summary.expenses > 0
      ? `Запас примерно на ${input.summary.runwayMonths} мес. при текущих расходах.`
      : null

  const freeCashflowText =
    input.summary.income > 0 && input.summary.expenses > 0
      ? `После расходов остаётся ${formatRub(input.summary.freeCashflow)}/мес.`
      : null

  return {
    goalTitle,
    goalProgress,
    steps: steps.slice(0, 3),
    runwayText,
    freeCashflowText,
    updatedAt: Date.now()
  }
}
