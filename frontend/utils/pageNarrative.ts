import type {
  AiDiagnosisResponse,
  CreditsDashboardResponse,
  DigestResponse,
  FinancialProfile,
  ForecastResponse,
  Goal,
  InsightItem,
  ReceiptListItem
} from '~/types/api'
import {
  ACTIONS,
  ANALYTICS,
  CREDITS,
  DIGEST,
  HEALTH,
  NAV,
  PROFILE,
  SOCIAL,
  formatRub
} from '~/constants/productCopy'
import type { DashboardSummary, HealthTone } from '~/utils/dashboardSummary'
import type { HabitIndex } from '~/utils/habitIndex'
import {
  adviceFromDiagnosis,
  buildAdviceBlock,
  buildAdviceHint,
  buildDashboardContextFacts,
  type DashboardContextFact
} from '~/utils/dashboardCopy'

export interface PageNarrativeBlock {
  headline: string
  paragraphs: string[]
  /** @deprecated Плитки контекста в contextFacts */
  contextLine?: string
  /** Оценка, доход, траты, запас — отдельные плитки */
  contextFacts?: DashboardContextFact[]
  /** Hook-цифра для совета недели */
  goalOpportunityThousands?: number | null
  /** @deprecated hero не рендерит; данные в contextLine / плане */
  statsLine?: string
  /** @deprecated hero не рендерит; таб «Возможность» в плане */
  forecastLine?: string
  /** @deprecated Встроено в adviceHint при необходимости */
  attentionLine?: string
  healthEmoji?: '🟢' | '🟡' | '🔴'
  healthTone?: HealthTone
  badgeLabel?: string
  weeklyAction?: string
  adviceHint?: string
  callout?: string
}

export function buildGoalProgressText(goal: Goal | null, monthlySaving?: number): string {
  if (!goal) {
    return 'Добавьте цель в профиле — покажем, сколько осталось и когда примерно дойдёте.'
  }

  const remaining = Math.max(0, goal.target_amount - goal.current_amount)
  if (remaining <= 0) {
    return `Цель «${goal.title}» достигнута — можно поставить следующую.`
  }

  const pct = goal.progress_percent ?? Math.round((goal.current_amount / goal.target_amount) * 100)
  let text = `Цель «${goal.title}»: ${pct}% · осталось ${formatRub(remaining)}.`

  if (monthlySaving && monthlySaving > 0) {
    const months = Math.ceil(remaining / monthlySaving)
    text += ` При тех же тратах — примерно ${months} мес.`
  } else if (goal.target_date) {
    text += ` Срок: ${goal.target_date}.`
  }

  return text
}

export function buildWeeklyAction(
  digest: DigestResponse | null,
  summary: Pick<DashboardSummary, 'weeklyAction'> | null,
  topInsight: InsightItem | null
): string {
  if (digest?.ai_advice) return digest.ai_advice
  if (summary?.weeklyAction) return summary.weeklyAction
  if (topInsight?.title) return topInsight.title
  return ACTIONS.addPurchaseOrReceipt
}

export function buildDigestPageNarrative(input: {
  digest: DigestResponse | null
  primaryGoal: Goal | null
  weeklyAction?: string
}): PageNarrativeBlock {
  const { digest, primaryGoal, weeklyAction } = input

  if (!digest) {
    return {
      headline: DIGEST.loading,
      paragraphs: ['Собираем цифры за месяц…'],
      badgeLabel: DIGEST.badge
    }
  }

  const free = digest.total_income - digest.total_spent
  const freeText =
    free >= 0
      ? DIGEST.leftAfterExpenses(`+${formatRub(free)}`)
      : `Расходы больше дохода на ${formatRub(Math.abs(free))}.`
  const picture = `За ${digest.period.from} — ${digest.period.to}: доход ${formatRub(digest.total_income)}, траты ${formatRub(digest.total_spent)}. ${freeText} Отложили ${formatRub(digest.saved)}.`

  return {
    headline: DIGEST.headline,
    paragraphs: [picture, buildGoalProgressText(primaryGoal, digest.saved)],
    weeklyAction: weeklyAction ?? digest.ai_advice,
    callout: digest.insights_summary,
    badgeLabel: DIGEST.badge
  }
}

export function buildReceiptsPageNarrative(input: {
  receipts: ReceiptListItem[]
  topInsight: InsightItem | null
}): PageNarrativeBlock {
  const { receipts, topInsight } = input
  const impulseTotal = receipts.reduce((sum, r) => sum + (r.impulse_count ?? 0), 0)
  const total = receipts.reduce((sum, r) => sum + r.amount, 0)

  const paragraphs: string[] = []
  if (receipts.length) {
    paragraphs.push(
      `В списке ${receipts.length} покупок на ${formatRub(total)}.` +
        (impulseTotal ? ` Покупок «на эмоциях»: ${impulseTotal}.` : '')
    )
  } else {
    paragraphs.push('Покупок пока нет — нажмите «Добавить».')
  }

  if (topInsight && receipts.length) {
    paragraphs.push(topInsight.body ?? topInsight.description ?? topInsight.title)
  } else if (impulseTotal > 0) {
    paragraphs.push('Часть покупок — «лишние» — может отодвинуть цель на несколько недель.')
  }

  return {
    headline: 'Куда уходят деньги',
    paragraphs,
    weeklyAction:
      receipts.length && topInsight?.title ? topInsight.title : ACTIONS.addPurchaseHint,
    badgeLabel: NAV.receipts
  }
}

export function buildAnalyticsPageNarrative(input: {
  forecast: ForecastResponse | null
  topInsight: InsightItem | null
  scenarioResult: string | null
  goalForecast?: string
}): PageNarrativeBlock {
  const { forecast, topInsight, scenarioResult, goalForecast } = input
  const paragraphs: string[] = []

  if (goalForecast) {
    paragraphs.push(goalForecast)
  } else {
    paragraphs.push('Здесь видно, как ваши траты влияют на срок до цели.')
  }

  if (forecast?.forecast?.length) {
    const avg = Math.round(
      forecast.forecast.reduce((a, b) => a + b, 0) / forecast.forecast.length
    )
    const max = Math.max(...forecast.forecast)
    paragraphs.push(
      `На ближайшие ${forecast.forecast.length} дней в среднем ~${formatRub(avg)} в день` +
        (max > avg * 1.25 ? ` — в один день может быть до ${formatRub(max)}.` : '.')
    )
  }

  if (scenarioResult) {
    paragraphs.push(scenarioResult)
  } else if (topInsight) {
    paragraphs.push(topInsight.body ?? topInsight.description ?? topInsight.title)
  }

  return {
    headline: ANALYTICS.headline,
    paragraphs,
    weeklyAction: topInsight?.title,
    badgeLabel: NAV.analytics
  }
}

export function buildCreditsPageNarrative(credits: CreditsDashboardResponse | null): PageNarrativeBlock {
  if (!credits) {
    return {
      headline: CREDITS.paymentsTitle,
      paragraphs: ['Считаем долги и запас…'],
      badgeLabel: NAV.credits
    }
  }

  const dti = credits.dti
  let healthTone: HealthTone = 'good'
  let healthEmoji: '🟢' | '🟡' | '🔴' = '🟢'
  if (dti >= 50) {
    healthTone = 'risk'
    healthEmoji = '🔴'
  } else if (dti >= 35) {
    healthTone = 'warn'
    healthEmoji = '🟡'
  }

  const paragraphs = [
    CREDITS.incomeShare(dti) +
      (credits.stress_test_months != null
        ? ` ${CREDITS.stressReserveMonths(credits.stress_test_months)}`
        : ''),
    'Ниже — ипотека: примерный платёж, «сейчас или подождать» и сравнение банков.'
  ]

  return {
    headline: CREDITS.anotherLoan,
    paragraphs,
    healthEmoji,
    healthTone,
    badgeLabel: CREDITS.trafficLight,
    weeklyAction:
      dti >= 35
        ? 'Перед новым кредитом — посмотрите ипотечный разбор и комфортный платёж.'
        : undefined
  }
}

export function buildProfilePageNarrative(input: {
  profile: FinancialProfile
  goals: Goal[]
}): PageNarrativeBlock {
  const { profile, goals } = input
  const income = profile.active_income + profile.passive_income
  const paragraphs: string[] = []

  if (income > 0 || profile.emergency_fund > 0) {
    paragraphs.push(
      `Доход ${formatRub(income)}/мес` +
        (profile.emergency_fund ? `, запас ${formatRub(profile.emergency_fund)}.` : '.')
    )
  } else {
    paragraphs.push('Укажите доход и запас — так прогноз цели станет точнее.')
  }

  if (goals.length) {
    paragraphs.push(buildGoalProgressText(goals[0], profile.active_income * 0.1))
  } else {
    paragraphs.push('Поставьте первую цель — отпуск, запас или крупная покупка.')
  }

  return {
    headline: PROFILE.headline,
    paragraphs,
    badgeLabel: NAV.profile,
    weeklyAction: goals.length ? undefined : 'Создайте цель — покажем, когда примерно дойдёте.'
  }
}

export function buildSocialPageNarrative(input?: { habitIndex?: HabitIndex | null }): PageNarrativeBlock {
  const paragraphs = [
    'Можно соревноваться с друзьями — без показа сумм, только ваши привычки.',
    'Чем спокойнее траты, тем выше оценка — можно закрепить это вместе.'
  ]

  const habit = input?.habitIndex
  if (habit?.insight) {
    paragraphs.push(habit.insight)
  } else if (habit && habit.score > 0) {
    paragraphs.push(`Сейчас оценка ${habit.score} из 100 — ${habit.label.toLowerCase()}.`)
  }

  return {
    headline: SOCIAL.headline,
    paragraphs,
    badgeLabel: SOCIAL.optionalBadge,
    weeklyAction: habit?.challengeHint
  }
}

export function narrativeFromDashboardSummary(
  summary: DashboardSummary,
  topInsight?: InsightItem | null
): PageNarrativeBlock {
  const advice = buildAdviceBlock({
    diagnosis: null,
    topInsight: topInsight ?? null,
    summary
  })

  return {
    headline: summary.healthHeadline,
    contextFacts: buildDashboardContextFacts({ diagnosis: null, summary }),
    goalOpportunityThousands: summary.goalOpportunityThousands,
    paragraphs: [],
    healthEmoji: summary.healthEmoji,
    healthTone: summary.healthTone,
    weeklyAction: advice.weeklyAction,
    adviceHint: advice.adviceHint,
    badgeLabel: NAV.dashboard
  }
}

function healthFromDiagnosisScore(score: number): {
  headline: string
  tone: HealthTone
  emoji: '🟢' | '🟡' | '🔴'
} {
  if (score >= 70) {
    return { headline: HEALTH.stable, tone: 'good', emoji: '🟢' }
  }
  if (score >= 50) {
    return { headline: HEALTH.caution, tone: 'warn', emoji: '🟡' }
  }
  return { headline: HEALTH.attention, tone: 'risk', emoji: '🔴' }
}

export function narrativeFromDiagnosis(
  diagnosis: AiDiagnosisResponse,
  summary?: DashboardSummary | null,
  topInsight?: InsightItem | null
): PageNarrativeBlock {
  const health = healthFromDiagnosisScore(diagnosis.score)

  const advice = summary
    ? buildAdviceBlock({ diagnosis, topInsight: topInsight ?? null, summary })
    : {
        weeklyAction: adviceFromDiagnosis(diagnosis.main_action.title),
        adviceHint: buildAdviceHint({
          diagnosis,
          summary: summary ?? null,
          weeklyAction: adviceFromDiagnosis(diagnosis.main_action.title)
        })
      }

  return {
    headline: `${health.headline} · ${diagnosis.score} из 100`,
    contextFacts: summary
      ? buildDashboardContextFacts({ diagnosis, summary })
      : [{ id: 'grade', label: 'Оценка', value: diagnosis.grade, tone: 'accent' }],
    goalOpportunityThousands: summary?.goalOpportunityThousands ?? null,
    paragraphs: [],
    healthEmoji: health.emoji,
    healthTone: health.tone,
    weeklyAction: advice.weeklyAction,
    adviceHint: advice.adviceHint,
    badgeLabel: `Оценка ${diagnosis.grade}`,
    callout:
      diagnosis.next_check_days > 0
        ? `Обновим картину через ${diagnosis.next_check_days} дн.`
        : undefined
  }
}
