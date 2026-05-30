import type {
  ForecastResponse,
  SimulateScenarioRequest,
  SimulateScenarioResponse,
  TimeMachineResponse
} from '~/types/api'
import { mockForecast, mockTimeMachine } from '~/store/mocks'
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
      timeMachine.value = normalizeTimeMachine(mockTimeMachine)
      forecast.value = normalizeForecast(mockForecast)
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

    const applyResult = (
      simulation: TimeMachineResponse,
      differenceFinal: number,
      monthlySaving: number
    ) => {
      scenarioResult.value = formatScenarioResult(
        differenceFinal,
        monthlySaving,
        params.reduction_percent
      )
      scenarioSimulation.value = simulation
    }

    const applyFallback = () => {
      const monthlySaving = Math.round(4_500 * (params.reduction_percent / 20))
      const differenceFinal = monthlySaving * 80
      const base = normalizeTimeMachine(mockTimeMachine)
      applyResult(
        {
          ...base,
          points: base.points.map((p, i) => ({
            ...p,
            optimistic: p.optimistic + i * Math.round(monthlySaving * 0.8)
          })),
          delta: differenceFinal,
          difference_final: differenceFinal
        },
        differenceFinal,
        monthlySaving
      )
    }

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
      applyResult(normalizeTimeMachine(res), res.difference_final, monthlySaving)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Ошибка симуляции'
      applyFallback()
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
