<script setup lang="ts">
import { buildGoalForecast } from '~/utils/dashboardSummary'
import { buildAnalyticsPageNarrative } from '~/utils/pageNarrative'
import { forecastAnomalyAlert } from '~/utils/analyticsNarrative'

const {
  timeMachine,
  forecast,
  loading,
  error,
  scenarioResult,
  scenarioSimulation,
  scenarioLoading,
  loadAll,
  simulateScenario
} = useAnalytics()

const { insights, topInsight, loading: insightsLoading, fetchInsights } = useInsights()

const scenario = ref<'reduce_delivery' | 'reduce_cafe' | 'reduce_entertainment' | 'custom'>(
  'reduce_cafe'
)
const percent = ref(20)

const pageLoading = computed(() => loading.value || insightsLoading.value)

const goalForecast = computed(() => buildGoalForecast(timeMachine.value))

const pageNarrative = computed(() =>
  buildAnalyticsPageNarrative({
    forecast: forecast.value,
    topInsight: topInsight.value,
    scenarioResult: scenarioResult.value,
    goalForecast: goalForecast.value
  })
)

const anomalyAlert = computed(() => forecastAnomalyAlert(forecast.value))

onMounted(async () => {
  await Promise.all([loadAll(), fetchInsights()])
})

async function runSimulation() {
  await simulateScenario({
    scenario: scenario.value,
    reduction_percent: percent.value,
    months: 60
  })
}
</script>

<template>
  <div class="mx-auto w-full max-w-6xl space-y-6">
    <SharedPageNarrative :narrative="pageNarrative" :loading="pageLoading && !forecast" />

    <Alert v-if="error" variant="destructive">
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <Skeleton v-if="pageLoading && !forecast" class="h-[480px] w-full" />

    <template v-else>
      <Card>
        <CardHeader>
          <CardTitle class="text-base">Прогноз трат (7 дней)</CardTitle>
          <CardDescription v-if="goalForecast">{{ goalForecast }}</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <Alert v-if="anomalyAlert" variant="default">
            <AlertTitle>Необычно много потратили</AlertTitle>
            <AlertDescription>{{ anomalyAlert }}</AlertDescription>
          </Alert>
          <div class="mm-chart-wrap mm-chart-wrap--md">
            <ChartsForecastChart :data="forecast" />
          </div>
        </CardContent>
      </Card>

      <AnalyticsScenarioSimulator
        v-model:scenario="scenario"
        v-model:percent="percent"
        :result="scenarioResult"
        :simulation="scenarioSimulation"
        :loading="scenarioLoading"
        @simulate="runSimulation"
      />

      <Card class="border-dashed">
        <CardContent class="flex flex-col gap-2 py-4 sm:flex-row sm:items-center sm:justify-between">
          <p class="text-sm text-muted-foreground">
            Полная картина «если ничего не менять» — на главной.
          </p>
          <Button variant="outline" size="sm" as-child>
            <NuxtLink to="/dashboard" data-demo="timemachine">Смотреть накопления →</NuxtLink>
          </Button>
        </CardContent>
      </Card>

      <AnalyticsDetective v-if="insights?.insights.length" :insights="insights.insights" />
    </template>
  </div>
</template>
