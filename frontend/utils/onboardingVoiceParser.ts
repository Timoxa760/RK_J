/**
 * Fallback-парсер RU-фраз для онбординга (доход, подушка, цель, обязательства).
 * На `back` планируется `POST /api/v1/onboarding/parse` (user-service / ai-processor).
 * В prod вызывается через `useOnboardingParse`; здесь — demo и откат при 404.
 */
import type { FixedExpense, GoalKind, EmergencyFundBreakdown } from '~/types/api'
import { GOAL_KIND_LABELS } from '~/constants/onboardingGoals'

function normalizeText(text: string) {
  return text.toLowerCase().replace(/\s+/g, ' ').trim()
}

export function isSkipAnswer(text: string) {
  return /пропуст|нет обязательн|не знаю|позже|ничего|без обязательн/i.test(text)
}

/** Извлекает суммы в рублях из свободной фразы. */
export function extractRubles(text: string): number[] {
  const normalized = normalizeText(text)
  const found: number[] = []

  const thousandRe =
    /(\d[\d\s.,]*)\s*(?:тысяч(?:и|ей|а)?|тыс(?:яч(?:и|ей|a)?|и|ч)?|тыщ(?:a|и|ей|у)?|т(?:р|р\.|\.р\.?)|k(?=\s|$|[^a-z0-9])|к(?=\s|$|[^0-9а-яё]))/gi
  let match: RegExpExecArray | null
  while ((match = thousandRe.exec(normalized)) !== null) {
    const base = Number(match[1].replace(/\s/g, '').replace(',', '.'))
    if (!Number.isNaN(base)) found.push(Math.round(base * 1000))
  }

  const rubRe = /(\d[\d\s.,]{2,})\s*(?:₽|руб|р\b)/gi
  while ((match = rubRe.exec(normalized)) !== null) {
    const value = Number(match[1].replace(/\s/g, '').replace(',', '.'))
    if (!Number.isNaN(value)) found.push(Math.round(value))
  }

  const bareRe = /\b(\d{4,7})\b/g
  while ((match = bareRe.exec(normalized)) !== null) {
    const value = Number(match[1])
    if (!Number.isNaN(value)) found.push(value)
  }

  return [...new Set(found)].filter((n) => n > 0)
}

export function parseIncomeAnswer(text: string) {
  const nums = extractRubles(text)
  const passiveHint = /пассив|аренд|дивиденд|процент|подработк/i.test(text)

  if (nums.length >= 2) {
    return { active_income: nums[0], passive_income: nums[1] }
  }
  if (nums.length === 1 && passiveHint) {
    return { active_income: 0, passive_income: nums[0] }
  }
  if (nums.length === 1) {
    return { active_income: nums[0], passive_income: 0 }
  }
  return {}
}

export function parseCushionAnswer(text: string) {
  const breakdown = parseEmergencyBreakdownLocal(text)
  if (breakdown.cash > 0 || breakdown.deposit > 0 || breakdown.investments > 0) {
    const total = breakdown.cash + breakdown.deposit + breakdown.investments
    return { emergency_fund: total, emergency_breakdown: breakdown }
  }
  const nums = extractRubles(text)
  if (!nums.length) return {}
  return { emergency_fund: nums[0] }
}

function parseEmergencyBreakdownLocal(text: string): EmergencyFundBreakdown {
  const out = { cash: 0, deposit: 0, investments: 0 }
  const normalized = normalizeText(text)
  const segments = splitVoiceSegmentsLocal(normalized)

  const cashKw = ['налич', 'налик', 'наличк', 'кэш', 'cash', 'на руках', 'дома']
  const depKw = ['вклад', 'депозит', 'счет', 'счёт', 'банк', 'сбер', 'вкладах']
  const invKw = ['инвест', 'акци', 'облига', 'брокер', 'фонд', 'крипт']

  for (const seg of segments) {
    const nums = extractRubles(seg)
    if (!nums.length) continue
    const amount = nums[0]
    const n = normalizeText(seg)
    if (cashKw.some((k) => n.includes(k))) out.cash += amount
    else if (depKw.some((k) => n.includes(k))) out.deposit += amount
    else if (invKw.some((k) => n.includes(k))) out.investments += amount
  }
  return out
}

function splitVoiceSegmentsLocal(text: string): string[] {
  const withCommas = text.replace(/;/g, ',').replace(/\s+и\s+/g, ',')
  const parts = withCommas
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)

  if (parts.length > 1) return parts

  const chunks = text.match(
    /\d[\d\s.,]*\s*(?:тысяч(?:и|ей|a)?|тыс(?:яч(?:и|ей|a)?|и|ч)?|тыщ(?:a|и|ей|u)?|т(?:р|р\.|\.р\.?)?|₽|руб(?:лей|ля)?)[^.]*?(?=,\s*\d|\s+и\s+\d|$)/gi
  )
  if (chunks?.length) return chunks.map((c) => c.trim())

  return [text]
}

export function parseGoalAnswer(text: string) {
  const normalized = normalizeText(text)
  const nums = extractRubles(text)
  let goal_kind: GoalKind = 'save'
  let goal_title = GOAL_KIND_LABELS.save

  if (/подушк|резерв|запас/i.test(normalized)) {
    goal_kind = 'cushion'
    goal_title = GOAL_KIND_LABELS.cushion
  } else if (/отпуск|путешеств/i.test(normalized)) {
    goal_kind = 'save'
    goal_title = 'Отпуск'
  } else if (/квартир|машин|авто|ремонт|ипотек|покупк/i.test(normalized)) {
    goal_kind = 'purchase'
    goal_title = GOAL_KIND_LABELS.purchase
  } else if (/накоп/i.test(normalized)) {
    goal_kind = 'save'
    goal_title = GOAL_KIND_LABELS.save
  }

  const titleMatch = normalized.match(/цель\s+(.+?)(?:\s+на|\s+\d|$)/i)
  if (titleMatch?.[1]) {
    goal_kind = 'other'
    goal_title = titleMatch[1].trim().replace(/^\W+|\W+$/g, '') || goal_title
  }

  return {
    goal_kind,
    goal_title,
    goal_amount: nums[0] ?? 0
  }
}

export function parseExpensesAnswer(text: string): {
  skipped: boolean
  items: FixedExpense[]
} {
  if (isSkipAnswer(text)) {
    return { skipped: true, items: [] }
  }

  const nums = extractRubles(text)
  const normalized = normalizeText(text)
  const items: FixedExpense[] = []

  if (/аренд|квартир/i.test(normalized) && nums[0]) {
    items.push({ title: 'Аренда', amount: nums[0] })
  }
  if (/кредит|ипотек/i.test(normalized)) {
    const amount = nums[items.length] ?? nums[0] ?? 0
    if (amount) items.push({ title: 'Кредит', amount })
  }
  if (/парковк/i.test(normalized)) {
    const amount = nums[items.length] ?? 0
    if (amount) items.push({ title: 'Парковка', amount })
  }

  if (!items.length && nums[0]) {
    items.push({ title: 'Обязательный платёж', amount: nums[0] })
  }

  return { skipped: false, items }
}
