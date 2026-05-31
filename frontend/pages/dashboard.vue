<script setup lang="ts">
import { buildDashboardSummary, hasCreditsData } from '~/utils/dashboardSummary'
import { narrativeFromDiagnosis, narrativeFromDashboardSummary } from '~/utils/pageNarrative'
import { buildCategoriesSummary } from '~/utils/chartSummaries'
import { resolveSavingsTimemachine } from '~/utils/dashboardProjections'

const { categories, timemachine, loading, error, loadAll, retry } = useDashboard()
const { dashboard: credits, loading: creditsLoading, fetchDashboard } = useCredits()
const { profile, loadProfile, fetchProfileFromApi } = useFinancialProfile()
const { insights, topInsight, loading: insightsLoading, fetchInsights } = useInsights()
const { plan, diagnosisFromPlan, loading: aiPlanLoading, fetchPlan } = useAiPlan()

const { refreshAdvisorContext } = useAdvisorContext()

const { addedVersion } = useAddExpenseSheet()

const selectedCategory = ref('')
const percent = ref(20)

const projectedTimemachine = computed(() =>
  resolveSavingsTimemachine(timemachine.value, profile.value, categories.value)
)

const summary = computed(() =>
  buildDashboardSummary({
    profile: profile.value,
    timemachine: projectedTimemachine.value,
    credits: credits.value,
    topInsight: topInsight.value,
    categories: categories.value
  })
)

const displayDiagnosis = computed(() => diagnosisFromPlan.value)

const displaySummary = computed(() => {
  const base = summary.value
  const d = displayDiagnosis.value
  if (!d) return base
  return {
    ...base,
    weeklyAction: d.main_action.title,
    behaviorInsight: d.main_action.description
  }
})

const pageNarrative = computed(() => {
  const insight = topInsight.value
  if (displayDiagnosis.value) {
    return narrativeFromDiagnosis(displayDiagnosis.value, summary.value, insight)
  }
  return narrativeFromDashboardSummary(summary.value, insight)
})

const narrativeLoading = computed(
  () =>
    !initialLoadDone.value &&
    (loading.value || creditsLoading.value || insightsLoading.value || aiPlanLoading.value)
)

const planLoading = computed(
  () => !initialLoadDone.value && (narrativeLoading.value || aiPlanLoading.value)
)

const showCredits = computed(() => hasCreditsData(credits.value))

const planRefreshing = ref(false)

async function rebuildPlan() {
  planRefreshing.value = true
  await fetchPlan({
    summary: summary.value,
    timemachine: projectedTimemachine.value,
    topInsight: topInsight.value
  })
  planRefreshing.value = false
}

const categoriesSummary = computed(() => buildCategoriesSummary(categories.value))

const chartsRefreshing = ref(false)
const initialLoadDone = ref(false)
const chartsLoading = computed(
  () => chartsRefreshing.value || (loading.value && !categories.value && !timemachine.value)
)

const allInsights = computed(() => insights.value?.insights ?? [])

async function refreshData(options?: { soft?: boolean }) {
  if (options?.soft) chartsRefreshing.value = true
  loadProfile()
  await Promise.all([loadAll({ silent: options?.soft }), fetchDashboard(), fetchInsights()])
  await rebuildPlan()
  await refreshAdvisorContext({ silent: true })
  if (options?.soft) chartsRefreshing.value = false
}

async function refreshAll() {
  await refreshData()
  initialLoadDone.value = true
}

onMounted(async () => {
  loadProfile()
  await fetchProfileFromApi()
  await refreshAll()
})

watch(addedVersion, () => {
  refreshData({ soft: true })
})
</script>

<template>
  <div class="mx-auto w-full max-w-none space-y-4 sm:space-y-5">
    <SharedPageNarrative :narrative="pageNarrative" :loading="narrativeLoading">
      <template #aside>
        <DashboardMindfulnessScore
          :diagnosis="displayDiagnosis"
          :loading="aiPlanLoading && !displayDiagnosis"
        />
      </template>
    </SharedPageNarrative>

    <AdvisorFinancialPlanCard
      mega
      :plan="plan"
      :summary="displaySummary"
      :diagnosis="displayDiagnosis"
      :diagnosis-loading="aiPlanLoading && !displayDiagnosis"
      :categories="categories"
      :profile="profile"
      :timemachine="projectedTimemachine"
      :categories-summary="categoriesSummary"
      :charts-loading="chartsLoading"
      :credits="credits"
      :credits-loading="creditsLoading"
      :show-credits="showCredits"
      :dti-tone="summary.dtiTone"
      :insights="allInsights"
      v-model:selected-category="selectedCategory"
      v-model:percent="percent"
      :loading="planLoading || planRefreshing"
      @refresh="rebuildPlan"
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
  </div>
</template>
