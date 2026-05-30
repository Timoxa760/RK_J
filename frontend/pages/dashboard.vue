<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { buildDashboardSummary, hasCreditsData } from '~/utils/dashboardSummary'
import { narrativeFromDiagnosis, narrativeFromDashboardSummary } from '~/utils/pageNarrative'
import { mockInsights } from '~/store/mocks'
import { normalizeInsights } from '~/utils/apiNormalize'
import type { InsightsResponse } from '~/types/api'

const { sankey, stores, categories, compare, timemachine, loading, error, loadAll, retry } =
  useDashboard()

const { dashboard: credits, loading: creditsLoading, fetchDashboard } = useCredits()

const { diagnosis, loading: diagnosisLoading, fetchDiagnosis } = useDiagnosis()

const { apiFetchWithDemo, demoMode } = useApi()
const addOpen = ref(false)

const insights = ref<InsightsResponse | null>(null)
const insightsLoading = ref(false)

async function loadInsights() {
  insightsLoading.value = true
  try {
    const raw = await apiFetchWithDemo('/insights', mockInsights)
    insights.value = normalizeInsights(raw)
  } catch {
    if (demoMode.value) {
      insights.value = normalizeInsights(mockInsights)
    }
  } finally {
    insightsLoading.value = false
  }
}

const summary = computed(() =>
  buildDashboardSummary({
    sankey: sankey.value,
    compare: compare.value,
    timemachine: timemachine.value,
    stores: stores.value,
    credits: credits.value,
    topInsight: insights.value?.insights?.[0] ?? null
  })
)

const displaySummary = computed(() => {
  const base = summary.value
  const d = diagnosis.value
  if (!d) return base
  return {
    ...base,
    weeklyAction: d.main_action.title,
    behaviorInsight: d.main_action.description
  }
})

const pageNarrative = computed(() => {
  if (diagnosis.value) {
    return narrativeFromDiagnosis(diagnosis.value, summary.value)
  }
  const block = narrativeFromDashboardSummary(summary.value)
  return { ...block, weeklyAction: undefined }
})

const narrativeLoading = computed(
  () => loading.value || creditsLoading.value || insightsLoading.value || diagnosisLoading.value
)

const showCredits = computed(() => hasCreditsData(credits.value))

const compareDescription = computed(() => {
  const change = compare.value?.insights?.biggest_change
  if (!change) return undefined
  const sign = change.delta > 0 ? '+' : ''
  const pctSign = change.delta_percent > 0 ? '+' : ''
  return `${change.category}: ${sign}${change.delta.toLocaleString('ru-RU')} ₽ (${pctSign}${change.delta_percent}%)`
})

onMounted(async () => {
  await Promise.all([loadAll(), fetchDashboard(), loadInsights(), fetchDiagnosis()])
})

function onExpenseAdded() {
  loadAll()
  fetchDashboard()
  loadInsights()
  fetchDiagnosis()
}
</script>

<template>
  <div class="mx-auto w-full max-w-6xl mm-page-shell">
    <SharedPageNarrative :narrative="pageNarrative" :loading="narrativeLoading" />

    <DashboardDiagnosisIndicators
      v-if="diagnosis?.indicators?.length || diagnosisLoading"
      class="mt-4"
      :indicators="diagnosis?.indicators ?? []"
      :loading="diagnosisLoading && !diagnosis"
    />

    <div
      class="mt-4 flex flex-col gap-4 rounded-xl border-2 border-primary/30 bg-primary/10 p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between sm:gap-6 sm:p-5"
      data-demo="add-expense"
    >
      <div class="min-w-0 space-y-1">
        <p class="text-base font-semibold text-foreground">Добавить покупку</p>
        <p class="text-sm text-muted-foreground">
          Запишите голосом — Поток разберёт и обновит картину
        </p>
      </div>
      <Button
        size="lg"
        class="h-12 w-full shrink-0 gap-2 px-8 text-base font-semibold sm:w-auto"
        @click="addOpen = true"
      >
        <Plus class="size-5" />
        Добавить
      </Button>
    </div>

    <DashboardMetricsGrid :summary="displaySummary" :loading="narrativeLoading" />

    <ClientOnly>
      <DashboardAddExpenseSheet v-model:open="addOpen" @added="onExpenseAdded" />
    </ClientOnly>

    <Alert v-if="error" variant="destructive">
      <AlertTitle>Не удалось загрузить данные</AlertTitle>
      <AlertDescription class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <span>{{ error }}</span>
        <Button variant="outline" size="sm" class="shrink-0 border-destructive/30" @click="retry">
          Повторить
        </Button>
      </AlertDescription>
    </Alert>

    <DashboardCreditsSnapshot
      v-if="showCredits || creditsLoading"
      :credits="credits"
      :dti-tone="summary.dtiTone"
      :loading="creditsLoading"
    />

    <div class="grid w-full grid-cols-1 items-stretch gap-4 lg:grid-cols-2 lg:gap-6">
      <ChartsChartCard
        title="Откуда деньги и куда уходят"
        description="Откуда приходят деньги и куда уходят"
        size="full"
        col-span="2"
        :loading="loading && !sankey"
      >
        <ChartsSankeyChart :data="sankey" size="full" />
      </ChartsChartCard>

      <ChartsChartCard
        title="По магазинам"
        description="Где вы чаще покупаете — по чекам и вашим записям"
        size="md"
        :loading="loading && !stores"
      >
        <ChartsBubbleChart :data="stores" size="md" />
      </ChartsChartCard>

      <ChartsChartCard
        title="Категории"
        description="Клик по сектору — детализация"
        size="md"
        :loading="loading && !categories"
        data-demo="categories"
      >
        <ChartsCategoryPie :data="categories" size="md" />
      </ChartsChartCard>

      <ChartsChartCard
        title="Как пойдут накопления"
        :description="summary.goalForecast"
        size="full"
        col-span="2"
        :loading="loading && !timemachine"
        data-demo="timemachine"
      >
        <ChartsTimeMachineChart :data="timemachine" size="full" />
      </ChartsChartCard>

      <ChartsChartCard
        title="Сравнение месяцев"
        :description="compareDescription"
        size="md"
        col-span="2"
        :loading="loading && !compare"
      >
        <ChartsDonutCompare :data="compare" size="md" />
      </ChartsChartCard>
    </div>

    <Card v-if="insights?.insights?.length" data-demo="insights">
      <CardHeader>
        <CardTitle class="text-base">Совет на сейчас</CardTitle>
        <CardDescription>Одно простое действие</CardDescription>
      </CardHeader>
      <CardContent>
        <InsightsPanel :insights="insights.insights.slice(0, 1)" />
        <Button variant="link" class="mt-3 h-auto p-0" as-child>
          <NuxtLink to="/analytics">Все советы →</NuxtLink>
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
