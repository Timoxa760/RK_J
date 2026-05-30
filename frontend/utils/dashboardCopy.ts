import type { AiDiagnosisIndicator, AiDiagnosisResponse, InsightItem } from '~/types/api'
import { ADVISOR, CREDITS, formatRub } from '~/constants/productCopy'
import type { DashboardSummary } from '~/utils/dashboardSummary'

const INDICATOR_PAIN: Record<string, string> = {
  'импульсивные траты':
    'Импульсивные траты выше нормы — съедают запас быстрее, чем кажется.',
  'накопления от дохода':
    'Накопления от дохода ниже нормы — цель может сдвинуться.',
  'долговая нагрузка':
    'Долговая нагрузка высокая — новые траты могут быть тяжелее.',
  'подушка безопасности':
    'Подушка безопасности ниже нормы — запас на чёрный день тонкий.',
  'стабильность доходов': 'Доход нестабилен — запас важнее обычного.'
}

function indicatorPainLine(indicator: AiDiagnosisIndicator): string {
  const key = indicator.name.toLowerCase()
  for (const [pattern, message] of Object.entries(INDICATOR_PAIN)) {
    if (key.includes(pattern)) return message
  }
  if (indicator.status === 'critical') {
    return `${indicator.name} — выше нормы, стоит скорректировать траты.`
  }
  return `${indicator.name} — ниже нормы, цель может сдвинуться.`
}

export function adviceFromDiagnosis(title: string): string {
  return title
}

export function adviceFromInsight(insight: InsightItem): string {
  const amount = insight.amount ?? 0
  if (amount > 0) {
    return `Проверьте: ${insight.title.toLowerCase()} — ${formatRub(amount)}/мес можно вернуть.`
  }
  return `Проверьте: ${insight.title.toLowerCase()}.`
}

export function adviceEmpty(): string {
  return 'Добавьте одну покупку — соберём первый совет за минуту.'
}

export function adviceHintWithSavings(amount: number): string {
  return `Экономия ~${formatRub(amount)}/мес · шаги ниже в плане`
}

export function adviceHintDefault(): string {
  return ADVISOR.weeklyAdviceHintShort
}

function normalizeForCompare(text: string): string {
  return text.toLowerCase().replace(/\s+/g, ' ').trim()
}

function hintsAreDuplicate(weeklyAction: string, attention: string): boolean {
  const a = normalizeForCompare(weeklyAction)
  const b = normalizeForCompare(attention)
  if (!a || !b) return false
  return a.includes(b.slice(0, 24)) || b.includes(a.slice(0, 24))
}

export function buildAdviceHint(input: {
  diagnosis: AiDiagnosisResponse | null | undefined
  summary?: Pick<DashboardSummary, 'mainRisk'> | null
  weeklyAction?: string
}): string {
  const savings = input.diagnosis?.main_action.potential_savings ?? 0
  if (savings > 0) return adviceHintWithSavings(savings)

  const attention = buildAttentionLine({
    diagnosis: input.diagnosis,
    summary: input.summary
  })
  if (
    attention &&
    input.weeklyAction &&
    !hintsAreDuplicate(input.weeklyAction, attention)
  ) {
    return attention
  }

  return adviceHintDefault()
}

export function buildAdviceBlock(input: {
  diagnosis: AiDiagnosisResponse | null | undefined
  topInsight: InsightItem | null | undefined
  summary: DashboardSummary
}): { weeklyAction: string; adviceHint: string } {
  let weeklyAction: string
  if (input.diagnosis?.main_action.title) {
    weeklyAction = adviceFromDiagnosis(input.diagnosis.main_action.title)
  } else if (input.topInsight?.title) {
    weeklyAction = adviceFromInsight(input.topInsight)
  } else if (input.summary.weeklyAction && input.summary.weeklyAction !== adviceEmpty()) {
    weeklyAction = input.summary.weeklyAction
  } else {
    weeklyAction = adviceEmpty()
  }

  return {
    weeklyAction,
    adviceHint: buildAdviceHint({
      diagnosis: input.diagnosis,
      summary: input.summary,
      weeklyAction
    })
  }
}

export type DashboardContextFactTone = 'default' | 'accent' | 'warn'

export interface DashboardContextFact {
  id: string
  label: string
  value: string
  tone?: DashboardContextFactTone
}

function formatCashflowValue(freeCashflow: number): string {
  if (freeCashflow < 0) {
    return `−${formatRub(Math.abs(freeCashflow))}`
  }
  return formatRub(freeCashflow)
}

/** Плитки «оценка / доход / траты» вместо одной строки через · */
export function buildDashboardContextFacts(input: {
  diagnosis: AiDiagnosisResponse | null | undefined
  summary: DashboardSummary
}): DashboardContextFact[] {
  const { summary, diagnosis } = input

  if (summary.income <= 0 && summary.expenses <= 0) {
    return [
      {
        id: 'need-data',
        label: 'Данные',
        value: 'Укажите доход и траты в профиле',
        tone: 'warn'
      }
    ]
  }

  const facts: DashboardContextFact[] = []

  if (diagnosis) {
    facts.push({
      id: 'grade',
      label: 'Оценка',
      value: diagnosis.grade,
      tone: 'accent'
    })
  } else if (summary.healthHeadline) {
    facts.push({
      id: 'health',
      label: 'Состояние',
      value: summary.healthEmoji
        ? `${summary.healthEmoji} ${summary.healthHeadline}`
        : summary.healthHeadline,
      tone: summary.healthTone === 'good' ? 'default' : 'warn'
    })
  }

  if (summary.income > 0) {
    facts.push({
      id: 'income',
      label: 'Доход',
      value: `${formatRub(summary.income)}/мес`
    })
  }

  if (summary.expenses > 0) {
    facts.push({
      id: 'expenses',
      label: 'Траты',
      value: `${formatRub(summary.expenses)}/мес`
    })
  } else {
    facts.push({
      id: 'expenses-missing',
      label: 'Траты',
      value: 'Добавьте в профиле',
      tone: 'warn'
    })
  }

  const hasExpenseBaseline = summary.expenses > 0 && summary.income > 0
  if (hasExpenseBaseline) {
    facts.push({
      id: 'cashflow',
      label: 'После расходов',
      value: formatCashflowValue(summary.freeCashflow),
      tone: summary.freeCashflow < 0 ? 'warn' : 'accent'
    })
  }

  if (summary.runwayMonths != null && summary.runwayMonths > 0) {
    const rounded = Math.round(summary.runwayMonths * 10) / 10
    facts.push({
      id: 'runway',
      label: 'Запас',
      value: `${rounded} мес.`
    })
  }

  return facts
}

/** @deprecated Используйте buildDashboardContextFacts */
export function buildDashboardContextLine(input: {
  diagnosis: AiDiagnosisResponse | null | undefined
  summary: DashboardSummary
}): string {
  return buildDashboardContextFacts(input)
    .map((f) => `${f.label} ${f.value}`)
    .join(' · ')
}

export function buildAttentionLine(input: {
  diagnosis: AiDiagnosisResponse | null | undefined
  summary?: Pick<DashboardSummary, 'mainRisk'> | null
}): string | undefined {
  const indicators = input.diagnosis?.indicators ?? []
  const critical = indicators.find((i) => i.status === 'critical')
  if (critical) return indicatorPainLine(critical)

  const warning = indicators.find((i) => i.status === 'warning')
  if (warning) return indicatorPainLine(warning)

  if (input.summary?.mainRisk) return input.summary.mainRisk

  return undefined
}

export function buildDiagnosisIntro(diagnosis: AiDiagnosisResponse): string {
  const count = diagnosis.indicators.length
  const word =
    count === 1 ? 'показатель' : count >= 2 && count <= 4 ? 'показателя' : 'показателей'
  return `${count} ${word} · оценка ${diagnosis.grade} · ${diagnosis.score} из 100`
}

export function buildCreditsIntro(dti: number): string {
  const pct = Math.round(dti)
  const comfort = pct < 35 ? 'комфортно' : pct < 50 ? 'на пределе' : 'тесно'
  return `${CREDITS.incomeShare(pct)} — ${comfort}`
}

export function scenarioResultPrefix(result: string): string {
  if (result.startsWith('Вы можете') || result.startsWith('Можете')) return result
  return `Вы можете получить: ${result}`
}
