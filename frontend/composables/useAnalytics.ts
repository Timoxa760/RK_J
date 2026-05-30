import type { SimulateScenarioRequest, SimulateScenarioResponse, TimeMachineResponse } from '~/types/api'
import { buildUserCategoryOptions } from '~/constants/expenseCategories'
import { formatScenarioResult } from '~/utils/analyticsNarrative'
import { normalizeTimeMachine } from '~/utils/apiNormalize'
import {
  buildScenarioResult,
  isPlaceholderTimemachine
} from '~/utils/dashboardProjections'
import { useDashboardStore } from '~/store/dashboardStore'

const SCENARIO_ENUM_TO_CATEGORY: Record<
  Exclude<SimulateScenarioRequest['scenario'], 'custom'>,
  string
> = {
  reduce_delivery: 'Доставка',
  reduce_cafe: 'Кафе и рестораны',
  reduce_entertainment: 'Развлечения'
}

function resolveScenarioCategory(
  scenario: SimulateScenarioRequest['scenario'],
  categories: ReturnType<typeof useDashboardStore>['categories']
): string {
  if (scenario === 'custom') {
    return buildUserCategoryOptions(categories)[0]?.name ?? 'Прочие расходы'
  }
  return SCENARIO_ENUM_TO_CATEGORY[scenario]
}

export function useAnalytics() {
  const { apiFetch } = useApi()
  const { profile } = useFinancialProfile()
  const dashboardStore = useDashboardStore()

  const timeMachine = ref<TimeMachineResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const scenarioResult = ref<string | null>(null)
  const scenarioSimulation = ref<TimeMachineResponse | null>(null)
  const scenarioLoading = ref(false)

  async function loadAll() {
    loading.value = true
    error.value = null
    try {
      const tmRaw = await apiFetch<TimeMachineResponse>('/dashboard/timemachine')
      timeMachine.value = normalizeTimeMachine(tmRaw)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка загрузки аналитики'
      timeMachine.value = null
    } finally {
      loading.value = false
    }
  }

  async function simulateScenario(params: {
    scenario: SimulateScenarioRequest['scenario']
    reduction_percent: number
    months?: number
  }) {
    scenarioLoading.value = true
    error.value = null
    scenarioResult.value = null
    scenarioSimulation.value = null

    try {
      const body: SimulateScenarioRequest = {
        scenario: params.scenario,
        reduction_percent: params.reduction_percent,
        months: params.months ?? 60
      }
      const res = await apiFetch<SimulateScenarioResponse>('/scenarios/simulate', {
        method: 'POST',
        body
      })
      const normalized = normalizeTimeMachine(res)
      const savingsBalance = profile.value.emergency_fund ?? 0

      if (isPlaceholderTimemachine(normalized, savingsBalance)) {
        const local = buildScenarioResult({
          profile: profile.value,
          categories: dashboardStore.categories,
          categoryName: resolveScenarioCategory(params.scenario, dashboardStore.categories),
          reductionPercent: params.reduction_percent,
          months: params.months ?? 12
        })
        scenarioResult.value = local.message
        scenarioSimulation.value = local.timemachine
        return
      }

      const monthlySaving =
        res.scenario?.monthly_saving ?? res.difference_final / (params.months ?? 60)
      scenarioResult.value = formatScenarioResult(
        res.difference_final,
        monthlySaving,
        params.reduction_percent
      )
      scenarioSimulation.value = normalized
    } catch {
      const local = buildScenarioResult({
        profile: profile.value,
        categories: dashboardStore.categories,
        categoryName: resolveScenarioCategory(params.scenario, dashboardStore.categories),
        reductionPercent: params.reduction_percent,
        months: params.months ?? 12
      })
      scenarioResult.value = local.message
      scenarioSimulation.value = local.timemachine
    } finally {
      scenarioLoading.value = false
    }
  }

  return {
    timeMachine,
    loading,
    error,
    scenarioResult,
    scenarioSimulation,
    scenarioLoading,
    loadAll,
    simulateScenario
  }
}
