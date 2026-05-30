import type {
  AiDiagnosisResponse,
  Goal,
  InsightItem,
  TimeMachineResponse
} from '~/types/api'
import { buildGoalForecast } from '~/utils/dashboardSummary'
import { buildGoalProgressText } from '~/utils/pageNarrative'
import { formatRub } from '~/constants/productCopy'

export interface AdvisorContext {
  diagnosis: AiDiagnosisResponse | null
  topInsight: InsightItem | null
  timemachine: TimeMachineResponse | null
  primaryGoal: Goal | null
  goalForecast: string
}

const QUICK_PROMPTS = [
  'Как сократить доставку?',
  'Когда дойду до цели?',
  'С чего начать экономить?'
] as const

export { QUICK_PROMPTS }

export function goalProgressPercent(goal: Goal): number {
  if (!goal.target_amount) return 0
  return Math.min(
    100,
    goal.progress_percent ?? Math.round((goal.current_amount / goal.target_amount) * 100)
  )
}

export function buildGoalCloserPrompt(goal: Goal): string {
  const pct = goalProgressPercent(goal)
  const remaining = Math.max(0, goal.target_amount - goal.current_amount)
  return `Как быстрее достичь цели «${goal.title}»? Сейчас ${pct}%, осталось ${formatRub(remaining)}.`
}

export function buildGoalCloserLabel(goalTitle: string): string {
  return `Приблизить «${goalTitle}»`
}

export function buildWeeklyActionPrompt(title: string, description?: string): string {
  const detail = description ? ` ${description}` : ''
  return `Расскажите подробнее про шаг «${title}».${detail} Что сделать в первую очередь?`
}

export function buildReceiptsGoalDelayPrompt(impactText: string): string {
  return `Эта покупка отодвигает цель: ${impactText}. Как тратить так, чтобы цель не сдвигалась?`
}

export function buildGettingStartedPrompt(): string {
  return 'С чего начать, чтобы картина денег и советы стали точнее?'
}

export function decodeAskQuery(value: string | string[] | null | undefined): string | null {
  if (!value || Array.isArray(value)) return null
  try {
    return decodeURIComponent(value).trim() || null
  } catch {
    return value.trim() || null
  }
}

function normalize(text: string) {
  return text.trim().toLowerCase()
}

export function buildAdvisorGreeting(_ctx: AdvisorContext): string {
  return 'Привет! Спросите про совет недели или шаги плана.'
}

export function buildAdvisorReply(message: string, ctx: AdvisorContext): string {
  const q = normalize(message)

  if (/план|состав|шаг|что делать/.test(q)) {
    const steps: string[] = ['Вот короткий план на ближайшее время:']
    if (ctx.diagnosis) {
      steps.push(`1. ${ctx.diagnosis.main_action.title} — ${ctx.diagnosis.main_action.description}`)
    }
    if (ctx.topInsight) {
      steps.push(
        `2. ${ctx.topInsight.title}${ctx.topInsight.body ? `: ${ctx.topInsight.body}` : ''}`
      )
    }
    steps.push('3. Записывайте покупки голосом — так картина и советы обновляются быстрее.')
    return steps.join('\n')
  }

  if (/урез|сократ|эконом|меньше трат|куда резать/.test(q)) {
    if (ctx.diagnosis) {
      const saving =
        ctx.diagnosis.main_action.potential_savings > 0
          ? ` Это около ${formatRub(ctx.diagnosis.main_action.potential_savings)} в месяц.`
          : ''
      return `${ctx.diagnosis.main_action.title}: ${ctx.diagnosis.main_action.description}${saving}`
    }
    if (ctx.topInsight) {
      return `${ctx.topInsight.title}. ${ctx.topInsight.body ?? ctx.topInsight.description ?? ''}`
    }
    return 'Добавьте несколько покупок — подскажу, где чаще всего «утекают» деньги.'
  }

  if (/быстрее|приблиз|достич|осталось/.test(q)) {
    const goalLine = buildGoalProgressText(ctx.primaryGoal)
    const forecast = ctx.goalForecast || buildGoalForecast(ctx.timemachine)
    return `${goalLine} ${forecast} Могу подсказать, что урезать в первую очередь — спросите «где урезать».`
  }

  if (/когда|срок|цел|дойду|накоп/.test(q)) {
    const goalLine = buildGoalProgressText(ctx.primaryGoal)
    const forecast = ctx.goalForecast || buildGoalForecast(ctx.timemachine)
    return `${goalLine} ${forecast}`
  }

  if (/привет|здрав|добр/.test(q)) {
    return buildAdvisorGreeting(ctx)
  }

  if (ctx.diagnosis) {
    return `Главное сейчас — ${ctx.diagnosis.main_action.title.toLowerCase()}. ${ctx.diagnosis.main_action.description} Можете спросить «составь план» или «где урезать».`
  }

  if (ctx.topInsight) {
    return `${ctx.topInsight.title}. ${ctx.topInsight.body ?? ctx.topInsight.description ?? 'Добавьте расходы — советы станут точнее.'}`
  }

  return 'Запишите покупку или спросите «составь план», «где урезать» или «когда дойду до цели».'
}
