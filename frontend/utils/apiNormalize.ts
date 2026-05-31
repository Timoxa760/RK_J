import type {
  CategoriesResponse,
  CompareResponse,
  CreditsDashboardResponse,
  ForecastResponse,
  InsightItem,
  InsightsResponse,
  LeaderboardEntry,
  LeaderboardResponse,
  SankeyResponse,
  StoresResponse,
  TimeMachineApiResponse,
  TimeMachineResponse
} from '~/types/api'

export function normalizeSankey(raw: SankeyResponse): SankeyResponse {
  return raw
}

export function normalizeStores(raw: StoresResponse): StoresResponse {
  return {
    stores: raw.stores.map((s, i) => ({
      ...s,
      id: s.id ?? String(i + 1),
      visits: s.visits ?? s.purchases ?? 0
    }))
  }
}

export function normalizeCategories(raw: CategoriesResponse): CategoriesResponse {
  return {
    categories: raw.categories.map((c) => {
      const amount = c.amount ?? c.total ?? 0
      return {
        name: c.name,
        amount,
        share: c.share ?? 0,
        subcategories: (c.subcategories ?? []).map((sub) => ({
          name: sub.name,
          items: sub.items.map((item) => ({
            name: item.name,
            price: item.price,
            quantity: item.quantity,
            amount:
              item.amount ??
              item.total ??
              (item.price ?? 0) * (item.quantity ?? 1)
          }))
        }))
      }
    })
  }
}

export function normalizeCompare(raw: CompareResponse): CompareResponse {
  return {
    months: raw.months.map((m) => ({
      month: m.month ?? m.label ?? '',
      label: m.label ?? m.month,
      categories: m.categories.map((cat) => ({
        name: cat.name,
        amount: cat.amount ?? cat.total ?? 0,
        share: cat.share
      }))
    })),
    insights: raw.insights
  }
}

export function normalizeTimeMachine(
  raw: TimeMachineResponse | TimeMachineApiResponse | Record<string, unknown>
): TimeMachineResponse {
  if ('points' in raw && Array.isArray((raw as TimeMachineResponse).points)) {
    const data = raw as TimeMachineResponse
    return {
      points: data.points,
      delta: data.delta ?? data.difference_final ?? 0,
      difference_final: data.difference_final
    }
  }

  const api = raw as {
    months?: string[]
    real_savings?: number[]
    optimized_savings?: number[]
    difference_final?: number
  }
  const months = api.months ?? []
  const points = months.map((month, i) => ({
    month,
    actual: api.real_savings?.[i] ?? 0,
    optimistic: api.optimized_savings?.[i] ?? 0
  }))
  const delta =
    api.difference_final ??
    (points.length
      ? (points[points.length - 1]?.optimistic ?? 0) - (points[points.length - 1]?.actual ?? 0)
      : 0)

  return { points, delta, difference_final: api.difference_final }
}

export function normalizeForecast(raw: ForecastResponse): ForecastResponse {
  if (raw.dates?.length && raw.forecast?.length) {
    return raw
  }

  const legacyPoints = raw.points as Array<{ month?: string; date?: string; amount: number }> | undefined
  if (legacyPoints?.length) {
    return {
      dates: legacyPoints.map((p) => p.date ?? p.month ?? ''),
      forecast: legacyPoints.map((p) => p.amount),
      upper_bound: legacyPoints.map((p) => Math.round(p.amount * 1.2)),
      lower_bound: legacyPoints.map((p) => Math.round(p.amount * 0.8))
    }
  }

  const apiPoints = (raw as { points?: Array<{ date?: string; amount?: number }> }).points
  if (apiPoints?.length) {
    return {
      dates: apiPoints.map((p) => p.date ?? ''),
      forecast: apiPoints.map((p) => p.amount ?? 0),
      upper_bound: apiPoints.map((p) => Math.round((p.amount ?? 0) * 1.2)),
      lower_bound: apiPoints.map((p) => Math.round((p.amount ?? 0) * 0.8))
    }
  }

  return { dates: [], forecast: [] }
}

export function normalizeInsights(raw: InsightsResponse): InsightsResponse {
  const placeholderTypes = new Set(['subscription', 'duplicate', 'overprice'])
  return {
    insights: (raw.insights ?? [])
      .filter((item) => !placeholderTypes.has(item.type))
      .map((item, i) => normalizeInsight(item, i))
  }
}

export function normalizeInsight(item: InsightItem, index = 0): InsightItem {
  const severity =
    item.severity === 'critical' ? 'warning' : item.severity
  return {
    id: item.id ?? String(index + 1),
    type: item.type,
    title: item.title,
    body: item.body ?? item.description ?? '',
    description: item.description ?? item.body,
    severity,
    amount: item.amount,
    merchant: item.merchant,
    store: item.store
  }
}

export function normalizeCredits(raw: CreditsDashboardResponse): CreditsDashboardResponse {
  const dti = raw.dti <= 1 ? Math.round(raw.dti * 100) : Math.round(raw.dti)
  const stressDti =
    raw.stress_test_dti != null
      ? raw.stress_test_dti <= 1
        ? Math.round(raw.stress_test_dti * 100)
        : Math.round(raw.stress_test_dti)
      : undefined

  return {
    ...raw,
    dti,
    stress_test_dti: stressDti,
    stress_test_months: raw.stress_test_months,
    credits: (raw.credits ?? []).map((c) => ({
      id: c.id,
      name: c.name ?? c.bank ?? 'Кредит',
      bank: c.bank ?? c.name,
      balance: c.balance ?? c.remaining ?? 0,
      remaining: c.remaining ?? c.balance,
      payment: c.payment ?? c.monthly_payment ?? 0,
      monthly_payment: c.monthly_payment ?? c.payment,
      rate: c.rate,
      term_months: c.term_months,
      amount: c.amount,
      next_payment: c.next_payment
    }))
  }
}

export function normalizeLeaderboard(raw: LeaderboardResponse): LeaderboardEntry[] {
  return raw.leaderboard.map((row) => ({
    position: row.position,
    username: row.username,
    avatar: row.avatar,
    relative_score: row.relative_score,
    id: row.id ?? String(row.position),
    display_name: row.display_name ?? row.username,
    rank: row.rank ?? row.position
  }))
}

export function percentDti(value: number): number {
  return value <= 1 ? Math.round(value * 100) : Math.round(value)
}
