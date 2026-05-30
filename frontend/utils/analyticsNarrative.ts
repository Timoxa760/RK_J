import type { ForecastResponse } from '~/types/api'

export function forecastAnomalyAlert(forecast: ForecastResponse | null): string | null {
  if (!forecast?.forecast?.length) return null

  const values = forecast.forecast
  const avg = values.reduce((sum, v) => sum + v, 0) / values.length
  const max = Math.max(...values)
  const maxIndex = values.indexOf(max)
  const date = forecast.dates?.[maxIndex]

  if (max <= avg * 1.2) return null

  const dateHint = date ? ` (${date})` : ''
  return `Трата может вырасти до ${max.toLocaleString('ru-RU')} ₽${dateHint} — заметно выше среднего ~${Math.round(avg).toLocaleString('ru-RU')} ₽/день.`
}

export function formatScenarioResult(
  differenceFinal: number,
  monthlySaving: number,
  reductionPercent: number
): string {
  const monthsCloser =
    monthlySaving > 0 ? Math.max(1, Math.round(differenceFinal / (monthlySaving * 6))) : null

  let text = `При сокращении на ${reductionPercent}% экономия ~${Math.round(monthlySaving).toLocaleString('ru-RU')} ₽/мес — за период прогноза +${differenceFinal.toLocaleString('ru-RU')} ₽ к накоплениям.`

  if (monthsCloser) {
    text += ` Цель может стать ближе примерно на ${monthsCloser} мес.`
  }

  return text
}
