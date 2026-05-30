import type { CategoriesResponse, ForecastResponse, TimeMachineResponse } from '~/types/api'
import { formatRub } from '~/constants/productCopy'
import { normalizeForecast } from '~/utils/apiNormalize'

export function formatChartDay(date: string): string {
  const match = date.match(/^(\d{4})-(\d{2})-(\d{2})/)
  if (!match) return date
  const months = [
    'янв',
    'фев',
    'мар',
    'апр',
    'май',
    'июн',
    'июл',
    'авг',
    'сен',
    'окт',
    'ноя',
    'дек'
  ]
  const day = Number(match[3])
  const month = months[Number(match[2]) - 1] ?? match[2]
  return `${day} ${month}`
}

export function buildForecastSummary(data: ForecastResponse | null): string {
  const normalized = data ? normalizeForecast(data) : null
  if (!normalized?.forecast.length) {
    return 'Добавьте расходы — покажем, сколько примерно уйдёт на ближайшие дни.'
  }

  const values = normalized.forecast
  const avg = Math.round(values.reduce((sum, v) => sum + v, 0) / values.length)
  const max = Math.max(...values)
  const maxIndex = values.indexOf(max)
  const maxDay = normalized.dates[maxIndex] ? formatChartDay(normalized.dates[maxIndex]) : ''

  if (maxDay) {
    return `В среднем ${formatRub(avg)} в день. Больше всего — ${formatRub(max)} (${maxDay}).`
  }

  return `В среднем ${formatRub(avg)} в день.`
}

export function buildSavingsSummary(data: TimeMachineResponse | null): string {
  if (!data?.points.length) {
    return 'Пока мало данных — добавьте покупки, и мы покажем, как растут накопления.'
  }

  const first = data.points[0]!
  const last = data.points[data.points.length - 1]!
  const diff = last.optimistic - last.actual
  const months = data.points.length

  if (diff > 0) {
    return `За ${months} мес. на счёте может быть ${formatRub(last.actual)}. Если чуть экономить — ${formatRub(last.optimistic)} (разница ${formatRub(diff)}).`
  }

  return `Через ${months} мес. на счёте может быть около ${formatRub(last.actual)}.`
}

export function buildCategoriesSummary(data: CategoriesResponse | null): string {
  if (!data?.categories.length) {
    return 'Категории появятся после первых покупок.'
  }

  const sorted = [...data.categories].sort(
    (a, b) => (b.amount ?? b.total ?? 0) - (a.amount ?? a.total ?? 0)
  )
  const top = sorted[0]!
  const topAmount = top.amount ?? top.total ?? 0
  const total = sorted.reduce((sum, c) => sum + (c.amount ?? c.total ?? 0), 0)
  const share = total > 0 ? Math.round((topAmount / total) * 100) : 0

  return `Больше всего уходит на «${top.name}» — ${formatRub(topAmount)} (${share}% всех трат).`
}
