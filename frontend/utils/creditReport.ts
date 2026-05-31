import type { CreditScanResponse, CreditsDashboardResponse } from '~/types/api'
import { CREDITS, formatRub } from '~/constants/productCopy'

export interface CreditReportMetric {
  label: string
  value: string
  tone?: 'default' | 'accent' | 'warn' | 'primary' | 'good'
  hint?: string
}

export interface CreditScanReportModel {
  title: string
  subtitle: string
  metrics: CreditReportMetric[]
  insights: string[]
  monthlyPayment: number
  paymentEstimated: boolean
  fixedExpenseTitle: string
  rateVsMarketLabel: string | null
}

export function estimateCreditMonthlyPayment(
  amount: number,
  rate: number,
  termMonths: number
): number {
  if (amount <= 0 || termMonths <= 0) return 0
  if (rate <= 0) return Math.round(amount / termMonths)
  const r = rate / 100 / 12
  const pow = (1 + r) ** termMonths
  if (pow <= 1) return Math.round(amount / termMonths)
  return Math.round((amount * r * pow) / (pow - 1))
}

export function resolveScanPayment(parsed: CreditScanResponse['parsed']): {
  amount: number
  estimated: boolean
} {
  if (parsed.monthly_payment > 0) {
    return {
      amount: Math.round(parsed.monthly_payment),
      estimated: Boolean(parsed.payment_estimated)
    }
  }
  const estimated = estimateCreditMonthlyPayment(
    parsed.amount,
    parsed.rate,
    parsed.term_months
  )
  return { amount: estimated, estimated: estimated > 0 }
}

function formatTermMonths(months: number): string {
  if (months % 12 === 0 && months >= 12) {
    const years = months / 12
    return `${years} ${years === 1 ? 'год' : years < 5 ? 'года' : 'лет'} (${months} мес.)`
  }
  return `${months} мес.`
}

function rateVsMarketLabel(scan: CreditScanResponse): string | null {
  const vs = scan.rate_vs_market
  if (vs === 'above') return CREDITS.scanRateAbove
  if (vs === 'at_or_below') return CREDITS.scanRateAtOrBelow
  if (vs === 'unknown') return CREDITS.scanRateUnknown
  return null
}

export function buildCreditScanReport(
  scan: CreditScanResponse,
  opts?: {
    monthlyIncome?: number
    fixedExpensesTotal?: number
  }
): CreditScanReportModel {
  const { parsed } = scan
  const payment = resolveScanPayment(parsed)
  const totalPaid = payment.amount * parsed.term_months
  const overpay = Math.max(0, Math.round(totalPaid - parsed.amount))
  const dti =
    opts?.monthlyIncome && opts.monthlyIncome > 0
      ? Math.round((payment.amount / opts.monthlyIncome) * 1000) / 10
      : null

  const metrics: CreditReportMetric[] = [
    {
      label: 'Сумма договора',
      value: formatRub(parsed.amount),
      tone: 'primary'
    },
    {
      label: 'Ставка',
      value: `${parsed.rate}% годовых`,
      tone: scan.rate_vs_market === 'above' ? 'warn' : 'default',
      hint:
        scan.benchmark_rate != null
          ? CREDITS.scanBenchmark(scan.benchmark_rate)
          : undefined
    },
    {
      label: 'Срок',
      value: formatTermMonths(parsed.term_months)
    },
    {
      label: 'Платёж в месяц',
      value: formatRub(payment.amount),
      tone: 'accent',
      hint: payment.estimated ? 'Рассчитан по сумме, ставке и сроку' : undefined
    },
    {
      label: 'Переплата за срок',
      value: formatRub(overpay),
      hint: `Всего выплат ~${formatRub(Math.round(totalPaid))}`
    }
  ]

  if (dti != null) {
    metrics.push({
      label: 'Доля дохода',
      value: `${dti}%`,
      tone: dti >= 35 ? 'warn' : dti >= 20 ? 'default' : 'good',
      hint: CREDITS.incomeShare(dti)
    })
  }

  const insights: string[] = []
  if (payment.estimated) {
    insights.push(
      'В договоре нет графика платежей — ежемесячную сумму посчитали по аннуитетной формуле.'
    )
  }
  if (scan.benchmark_rate != null) {
    const delta = Math.round((parsed.rate - scan.benchmark_rate) * 10) / 10
    if (delta > 0.5) {
      insights.push(
        `Ставка ${parsed.rate}% — на ~${delta} п.п. выше среднерыночной (~${scan.benchmark_rate}%).`
      )
    } else if (delta <= 0) {
      insights.push(`Ставка ${parsed.rate}% — на уровне или ниже рынка (~${scan.benchmark_rate}%).`)
    }
  }
  insights.push(
    `За ${parsed.term_months} мес. вернёте около ${formatRub(Math.round(totalPaid))}, из них ${formatRub(overpay)} — проценты.`
  )
  if (opts?.fixedExpensesTotal != null) {
    const nextFixed = (opts.fixedExpensesTotal ?? 0) + payment.amount
    insights.push(
      `Если добавить платёж в «каждый месяц», в прогнозе будет ~${formatRub(nextFixed)}/мес регулярных трат.`
    )
  } else {
    insights.push(
      'Добавьте платёж в «каждый месяц» — так главная и советник учтут кредит в прогнозе.'
    )
  }

  return {
    title: 'Разбор договора',
    subtitle: parsed.bank,
    metrics,
    insights,
    monthlyPayment: payment.amount,
    paymentEstimated: payment.estimated,
    fixedExpenseTitle: `Кредит · ${parsed.bank}`,
    rateVsMarketLabel: rateVsMarketLabel(scan)
  }
}

export function creditFixedExpenseTitle(bank: string): string {
  return `Кредит · ${bank}`
}

export function mergeCreditFixedExpense(
  rows: { title: string; amount: number }[],
  title: string,
  amount: number
): { title: string; amount: number }[] {
  const next = [...rows]
  const index = next.findIndex((row) => row.title === title)
  const item = { title, amount: Math.max(0, Math.round(amount)) }
  if (index >= 0) {
    next[index] = item
  } else {
    next.push(item)
  }
  return next.filter((row) => row.title.trim() || row.amount > 0)
}

export function isCreditInFixedExpenses(
  rows: { title: string; amount: number }[] | undefined,
  title: string,
  amount: number
): boolean {
  const row = rows?.find((item) => item.title === title)
  return Boolean(row && row.amount === Math.round(amount))
}

export function buildCreditsPortfolioSummary(dashboard: CreditsDashboardResponse): string {
  const count = dashboard.credits.length
  const debt = dashboard.total_debt ?? 0
  const payments = dashboard.monthly_payments ?? 0
  if (count === 0) return ''
  return `${count} ${count === 1 ? 'договор' : count < 5 ? 'договора' : 'договоров'} · долг ${formatRub(debt)} · платежи ${formatRub(payments)}/мес`
}
