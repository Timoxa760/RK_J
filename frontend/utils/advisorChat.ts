import type {
  AiChatMessage,
  AiChatResponse,
  AiDiagnosisResponse,
  CategoriesResponse,
  Goal,
  InsightItem,
  TimeMachineResponse
} from '~/types/api'
import { buildUserCategoryOptions } from '~/constants/expenseCategories'
import { formatRub } from '~/constants/productCopy'
import { parseAdvisorStoredContent } from '~/utils/advisorStructured'

export interface AdvisorContext {
  diagnosis: AiDiagnosisResponse | null
  topInsight: InsightItem | null
  timemachine: TimeMachineResponse | null
  primaryGoal: Goal | null
  goalForecast: string
  categories: CategoriesResponse | null
}

const DEFAULT_QUICK_PROMPTS = [
  'С чего начать экономить?',
  'Когда дойду до цели?',
  'Составь план'
] as const

export function buildDynamicQuickPrompts(ctx: AdvisorContext): string[] {
  const prompts: string[] = []
  const seen = new Set<string>()

  const add = (text: string) => {
    const line = text.trim()
    if (!line || seen.has(line)) return
    seen.add(line)
    prompts.push(line)
  }

  const action = ctx.diagnosis?.main_action
  if (action?.title) {
    add(`Что делать с «${action.title}»?`)
  }

  const weakIndicator = ctx.diagnosis?.indicators?.find(
    (row) => row.status === 'critical' || row.status === 'warning'
  )
  if (weakIndicator) {
    add(`Как улучшить «${weakIndicator.name}»?`)
  }

  if (ctx.topInsight?.title) {
    add(`Почему «${ctx.topInsight.title}»?`)
  }

  if (ctx.primaryGoal?.title) {
    const remaining = Math.max(0, ctx.primaryGoal.target_amount - ctx.primaryGoal.current_amount)
    if (remaining > 0) {
      add(`Когда накоплю на «${ctx.primaryGoal.title}»?`)
    } else {
      add(`Как приблизить «${ctx.primaryGoal.title}»?`)
    }
  }

  const topCategory = buildUserCategoryOptions(ctx.categories)[0]
  if (topCategory && topCategory.amount > 0) {
    add(`Как урезать «${topCategory.name}»?`)
  }

  if (ctx.goalForecast && ctx.primaryGoal) {
    add('Составь план до цели')
  }

  add('Где больше всего уходит денег?')
  add('С чего начать улучшать картину?')
  add(DEFAULT_QUICK_PROMPTS[2]!)

  return prompts.slice(0, 3)
}

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

export function historyToTurns(
  rows: Array<{
    id: string
    role: 'user' | 'assistant'
    content: string
    actions?: AiChatResponse['actions']
    source?: string
    created_at: number
  }>
) {
  return rows.map((row) => {
    const parsed = parseAdvisorStoredContent(row.content)
    return {
      id: row.id,
      role: row.role,
      content: parsed.plain,
      title: parsed.title,
      blocks: parsed.blocks,
      createdAt: row.created_at,
      actions: row.actions,
      source: (row.source as 'gemini' | 'heuristic' | undefined) ?? undefined
    }
  })
}

export function toApiHistory(
  messages: Array<{ role: 'user' | 'assistant'; content: string }>
): AiChatMessage[] {
  return messages
    .filter((m) => m.role === 'user' || m.role === 'assistant')
    .slice(-10)
    .map((m) => ({ role: m.role, content: m.content }))
}
