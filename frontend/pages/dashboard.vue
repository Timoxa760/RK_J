<script setup lang="ts">
import { buildDashboardSummary, hasCreditsData } from '~/utils/dashboardSummary'
import { narrativeFromDiagnosis, narrativeFromDashboardSummary } from '~/utils/pageNarrative'
import { buildFinancialPlan } from '~/utils/financialPlan'
import { mockInsights } from '~/store/mocks'
import { normalizeInsights } from '~/utils/apiNormalize'
import { buildCategoriesSummary } from '~/utils/chartSummaries'
import type { FinancialPlan } from '~/utils/financialPlan'
import type { InsightsResponse } from '~/types/api'

const { categories, timemachine, loading, error, loadAll, retry } = useDashboard()
const { dashboard: credits, loading: creditsLoading, fetchDashboard } = useCredits()
const { diagnosis, loading: diagnosisLoading, fetchDiagnosis } = useDiagnosis()
const { profile, loadProfile } = useFinancialProfile()
const { primaryGoal, fetchGoals } = useGoals()
const { insights, topInsight, loading: insightsLoading, fetchInsights } = useInsights()
const {
  forecast,
  loading: forecastLoading,
  error: forecastError,
  scenarioResult,
  scenarioLoading,
  loadAll: loadForecast,
  simulateScenario
} = useAnalytics()

const { refreshAdvisorContext } = useAdvisorContext()

const { apiFetchWithDemo, demoMode } = useApi()
const { addedVersion } = useAddExpenseSheet()

const insightsData = ref<InsightsResponse | null>(null)

const scenario = ref<'reduce_delivery' | 'reduce_cafe' | 'reduce_entertainment' | 'custom'>(
  'reduce_cafe'
)
const percent = ref(20)

async function loadInsightCard() {
  try {
    const raw = await apiFetchWithDemo('/insights', mockInsights)
    insightsData.value = normalizeInsights(raw)
  } catch {
    if (demoMode.value) {
      insightsData.value = normalizeInsights(mockInsights)
    }
  }
}

const summary = computed(() =>
  buildDashboardSummary({
    profile: profile.value,
    timemachine: timemachine.value,
    credits: credits.value,
    topInsight: topInsight.value ?? insightsData.value?.insights?.[0] ?? null
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

const topInsightRef = computed(
  () => topInsight.value ?? insightsData.value?.insights?.[0] ?? null
)

const pageNarrative = computed(() => {
  const insight = topInsightRef.value
  if (diagnosis.value) {
    return narrativeFromDiagnosis(diagnosis.value, summary.value, insight)
  }
  return narrativeFromDashboardSummary(summary.value, insight)
})

const narrativeLoading = computed(
  () =>
    !initialLoadDone.value &&
    (loading.value ||
      creditsLoading.value ||
      insightsLoading.value ||
      diagnosisLoading.value ||
      forecastLoading.value)
)

const planLoading = computed(
  () => !initialLoadDone.value && narrativeLoading.value
)

const showCredits = computed(() => hasCreditsData(credits.value))

const plan = ref<FinancialPlan | null>(null)
const planRefreshing = ref(false)

function rebuildPlan() {
  planRefreshing.value = true
  plan.value = buildFinancialPlan({
    primaryGoal: primaryGoal.value,
    summary: summary.value,
    timemachine: timemachine.value,
    diagnosis: diagnosis.value,
    topInsight: topInsight.value ?? insightsData.value?.insights?.[0] ?? null
  })
  planRefreshing.value = false
}

const categoriesSummary = computed(() => buildCategoriesSummary(categories.value))

const chartsRefreshing = ref(false)
const initialLoadDone = ref(false)
const chartsLoading = computed(
  () => chartsRefreshing.value || (loading.value && !categories.value && !timemachine.value)
)

const allInsights = computed(() => insights.value?.insights ?? insightsData.value?.insights ?? [])

async function refreshData(options?: { soft?: boolean }) {
  if (options?.soft) chartsRefreshing.value = true
  loadProfile()
  await Promise.all([
    loadAll({ silent: options?.soft }),
    fetchDashboard(),
    fetchInsights(),
    loadInsightCard(),
    fetchDiagnosis(),
    loadForecast(),
    fetchGoals()
  ])
  rebuildPlan()
  await refreshAdvisorContext({ silent: true })
  if (options?.soft) chartsRefreshing.value = false
}

async function refreshAll() {
  await refreshData()
  initialLoadDone.value = true
}

onMounted(async () => {
  loadProfile()
  await refreshAll()
})

watch(addedVersion, () => {
  refreshData({ soft: true })
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
  <div class="mx-auto w-full max-w-none space-y-5 lg:space-y-6">
    <SharedPageNarrative :narrative="pageNarrative" :loading="narrativeLoading" />

    <AdvisorFinancialPlanCard
      mega
      :plan="plan"
      :summary="displaySummary"
      :diagnosis="diagnosis"
      :diagnosis-loading="diagnosisLoading && !diagnosis"
      :categories="categories"
      :forecast="forecast"
      :timemachine="timemachine"
      :categories-summary="categoriesSummary"
      :charts-loading="chartsLoading"
      :credits="credits"
      :credits-loading="creditsLoading"
      :show-credits="showCredits"
      :dti-tone="summary.dtiTone"
      :insights="allInsights"
      v-model:scenario="scenario"
      v-model:percent="percent"
      :scenario-result="scenarioResult"
      :scenario-loading="scenarioLoading"
      :loading="planLoading || planRefreshing"
      @refresh="rebuildPlan"
      @simulate="runSimulation"
    />

    <Alert v-if="error" variant="destructive">
      <AlertTitle>Не удалось загрузить данные</AlertTitle>
      <AlertDescription class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <span>{{ error }}</span>
        <Button variant="outline" size="sm" class="shrink-0 border-destructive/30" @click="retry">
          Повторить
        </Button>
      </AlertDescription>
    </Alert>

    <Alert v-if="forecastError" variant="destructive">
      <AlertDescription class="text-base">{{ forecastError }}</AlertDescription>
    </Alert>
  </div>
</template>
