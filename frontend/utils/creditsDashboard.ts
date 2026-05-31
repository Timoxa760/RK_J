import type { CreditItem, CreditsDashboardResponse, FinancialProfile } from '~/types/api'
import { percentDti } from '~/utils/apiNormalize'
import { estimateCreditMonthlyPayment } from '~/utils/creditReport'

export interface EnrichedCreditsDashboard extends CreditsDashboardResponse {
  /** false — доход неизвестен, процент DTI показывать нельзя */
  dti_available: boolean
}

function profileIncome(profile: FinancialProfile | null | undefined): number {
  if (!profile || profile.skipped_income) return 0
  return Math.max(0, (profile.active_income ?? 0) + (profile.passive_income ?? 0))
}

function profileSavings(profile: FinancialProfile | null | undefined): number {
  if (!profile || profile.skipped_cushion) return 0
  return Math.max(0, profile.emergency_fund ?? 0)
}

export function resolveCreditMonthlyPayment(credit: CreditItem): number {
  const direct = credit.monthly_payment ?? credit.payment ?? 0
  if (direct > 0) return Math.round(direct)
  const amount = credit.amount ?? credit.remaining ?? credit.balance ?? 0
  const rate = credit.rate ?? 0
  const term = credit.term_months ?? 0
  if (amount > 0 && term > 0) {
    return estimateCreditMonthlyPayment(amount, rate, term)
  }
  return 0
}

export function sumCreditMonthlyPayments(credits: CreditItem[] | undefined): number {
  return (credits ?? []).reduce((sum, credit) => sum + resolveCreditMonthlyPayment(credit), 0)
}

/** Дополняет dashboard доходом из профиля и пересчитывает DTI/запас. */
export function enrichCreditsDashboard(
  dashboard: CreditsDashboardResponse,
  profile?: FinancialProfile | null
): EnrichedCreditsDashboard {
  const credits = (dashboard.credits ?? []).map((credit) => {
    const monthly_payment = resolveCreditMonthlyPayment(credit)
    return {
      ...credit,
      monthly_payment,
      payment: monthly_payment
    }
  })

  const monthlyPayments =
    sumCreditMonthlyPayments(credits) || Math.max(0, dashboard.monthly_payments ?? 0)

  const incomeFromApi = Math.max(0, dashboard.monthly_income ?? 0)
  const incomeFromProfile = profileIncome(profile)
  const monthlyIncome = incomeFromApi > 0 ? incomeFromApi : incomeFromProfile

  const savingsFromApi = Math.max(0, dashboard.savings ?? 0)
  const savingsFromProfile = profileSavings(profile)
  const savings = savingsFromApi > 0 ? savingsFromApi : savingsFromProfile

  const dtiAvailable = monthlyPayments > 0 && monthlyIncome > 0
  const dti = dtiAvailable
    ? percentDti(monthlyPayments / monthlyIncome)
    : percentDti(dashboard.dti ?? 0)

  const stressMonths =
    monthlyPayments > 0 && savings > 0
      ? Math.round((savings / monthlyPayments) * 10) / 10
      : dashboard.stress_test_months

  return {
    ...dashboard,
    credits,
    monthly_payments: monthlyPayments,
    monthly_income: monthlyIncome,
    savings,
    dti: dtiAvailable ? dti : 0,
    dti_available: dtiAvailable,
    stress_test_months: stressMonths
  }
}
