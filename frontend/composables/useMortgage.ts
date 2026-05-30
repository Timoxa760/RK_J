import type {
  MortgageAnalyzeRequest,
  MortgageApprovalLevel,
  MortgageBreakdownResponse
} from '~/types/api'
import { mockMortgageBreakdown } from '~/store/mocks'

function approvalFromDti(dti: number): MortgageApprovalLevel {
  if (dti < 30) return 'high'
  if (dti < 45) return 'medium'
  return 'low'
}

function scaleBreakdown(
  base: MortgageBreakdownResponse,
  input: MortgageAnalyzeRequest
): MortgageBreakdownResponse {
  const income = input.monthly_income ?? 180_000
  const savings = input.savings ?? 340_000
  const amount = input.mortgage_amount
  const dti = input.existing_dti ?? 28
  const paymentRatio = amount / 12_000_000
  const estimatedPayment = Math.round(base.comfortable_payment * paymentRatio)
  const newDti = Math.min(65, Math.round(dti + paymentRatio * 12))
  const level = approvalFromDti(newDti)

  const reasons: Record<MortgageApprovalLevel, string> = {
    high: `Доход ${income.toLocaleString('ru-RU')} ₽/мес, на кредиты уйдёт ~${newDti}% — шанс одобрения выглядит хорошим.`,
    medium: `На кредиты уйдёт ~${newDti}% дохода, запас ${savings.toLocaleString('ru-RU')} ₽ — банк может попросить подтверждения.`,
    low: `На кредиты уйдёт ~${newDti}% — ипотека на ${amount.toLocaleString('ru-RU')} ₽ выглядит рискованной.`
  }

  return {
    ...base,
    approval_level: level,
    approval_reason: reasons[level],
    safe_mortgage_amount: Math.round(income * 50),
    comfortable_payment: estimatedPayment,
    load_risk:
      savings > estimatedPayment * 4
        ? `После платежа запаса хватит ~${(savings / estimatedPayment).toFixed(1)} мес. — тесновато, но терпимо.`
        : 'После платежа запас станет маленьким — лучше подумать ещё раз.',
    scenario_now: `Платёж ~${estimatedPayment.toLocaleString('ru-RU')} ₽/мес — до цели будете идти медленнее.`,
    scenario_wait: `Через ${base.wait_months} мес. при накоплении ~${Math.round(estimatedPayment * 3).toLocaleString('ru-RU')} ₽ условия могут стать мягче.`,
    banks: base.banks.map((b) => ({
      ...b,
      monthly_payment: Math.round(b.monthly_payment * paymentRatio),
      total_overpayment: Math.round(b.total_overpayment * paymentRatio)
    }))
  }
}

export function useMortgage() {
  const breakdown = ref<MortgageBreakdownResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function analyze(input: MortgageAnalyzeRequest) {
    loading.value = true
    error.value = null
    try {
      await new Promise((r) => setTimeout(r, 300))
      breakdown.value = scaleBreakdown(mockMortgageBreakdown, input)
      return breakdown.value
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Не удалось рассчитать разбор'
      breakdown.value = scaleBreakdown(mockMortgageBreakdown, input)
      return breakdown.value
    } finally {
      loading.value = false
    }
  }

  function reset() {
    breakdown.value = null
    error.value = null
  }

  return {
    breakdown,
    loading,
    error,
    analyze,
    reset
  }
}
