import type {
  ForecastResponse,
  SimulateScenarioRequest,
  SimulateScenarioResponse,
  TimeMachineResponse
} from '~/types/api'
import { formatScenarioResult } from '~/utils/analyticsNarrative'
import { normalizeForecast, normalizeTimeMachine } from '~/utils/apiNormalize'

export function useAnalytics() {
  const { apiFetch } = useApi()

  const timeMachine = ref<TimeMachineResponse | null>(null)
  const forecast = ref<ForecastResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const scenarioResult = ref<string | null>(null)
  const scenarioSimulation = ref<TimeMachineResponse | null>(null)
  const scenarioLoading = ref(false)

  async function loadAll() {
    loading.value = true
    error.value = null
    try {
      const [tmRaw, fcRaw] = await Promise.all([
        apiFetch<TimeMachineResponse>('/dashboard/timemachine'),
        apiFetch<ForecastResponse>('/forecast?days=7')
      ])
      timeMachine.value = normalizeTimeMachine(tmRaw)
      forecast.value = normalizeForecast(fcRaw)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка загрузки аналитики'
      timeMachine.value = null
      forecast.value = null
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
      const monthlySaving =
        res.scenario?.monthly_saving ?? res.difference_final / (params.months ?? 60)
      scenarioResult.value = formatScenarioResult(
        res.difference_final,
        monthlySaving,
        params.reduction_percent
      )
      scenarioSimulation.value = normalizeTimeMachine(res)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка симуляции'
    } finally {
      scenarioLoading.value = false
    }
  }

  return {
    timeMachine,
    forecast,
    loading,
    error,
    scenarioResult,
    scenarioSimulation,
    scenarioLoading,
    loadAll,
    simulateScenario
  }
}
