<script setup lang="ts">
import { RefreshCw } from 'lucide-vue-next'
import type {
  AiDiagnosisResponse,
  CategoriesResponse,
  CreditsDashboardResponse,
  FinancialProfile,
  InsightItem,
  TimeMachineResponse
} from '~/types/api'
import type { DashboardSummary, HealthTone } from '~/utils/dashboardSummary'
import type { FinancialPlan } from '~/utils/financialPlan'
import type { PageNarrativeBlock } from '~/utils/pageNarrative'
import { ADVISOR, GOALS, HEALTH } from '~/constants/productCopy'
import { buildDiagnosisIntro } from '~/utils/dashboardCopy'

const props = withDefaults(
  defineProps<{
    plan: FinancialPlan | null
    summary: DashboardSummary | null
    diagnosis: AiDiagnosisResponse | null
    diagnosisLoading?: boolean
    categories: CategoriesResponse | null
    profile: FinancialProfile | null
    timemachine: TimeMachineResponse | null
    categoriesSummary: string
    chartsLoading?: boolean
    credits: CreditsDashboardResponse | null
    creditsLoading?: boolean
    showCredits?: boolean
    dtiTone?: HealthTone
    insights: InsightItem[]
    loading?: boolean
    overviewLoading?: boolean
    narrative?: PageNarrativeBlock | null
    /** Полный план на дашборде — все блоки подряд */
    mega?: boolean
  }>(),
  { mega: false }
)

const selectedCategory = defineModel<string>('selectedCategory', { default: '' })
const percent = defineModel<number>('percent', { default: 20 })

const emit = defineEmits<{
  refresh: []
}>()

const planTab = ref('steps')

const diagnosisIntro = computed(() =>
  props.diagnosis ? buildDiagnosisIntro(props.diagnosis) : ''
)

const opportunityBadge = computed(() => {
  const k = props.summary?.goalOpportunityThousands
  if (!k) return null
  return GOALS.opportunityAmount(k).replace(' ₽', '')
})

const narrativeResolved = computed(() => {
  const n = props.narrative
  return {
    headline: n?.headline ?? '',
    goalOpportunityThousands: n?.goalOpportunityThousands ?? null,
    weeklyAction: n?.weeklyAction ?? props.summary?.weeklyAction ?? '',
    adviceHint: n?.adviceHint ?? ADVISOR.weeklyAdviceHintShort,
    incomeDisplay: n?.incomeDisplay ?? null,
    expensesDisplay: n?.expensesDisplay ?? null,
    expensesWarn: n?.expensesWarn ?? false,
    callout: n?.callout ?? null
  }
})

const showMoneyRow = computed(
  () => Boolean(narrativeResolved.value.incomeDisplay || narrativeResolved.value.expensesDisplay)
)

const showOverview = computed(
  () =>
    props.mega &&
    Boolean(
      props.summary ||
        narrativeResolved.value.weeklyAction ||
        showMoneyRow.value ||
        props.diagnosis ||
        props.diagnosisLoading
    )
)
</script>

<template>
  <Card
    id="financial-plan"
    data-demo="financial-plan"
    class="mm-tier-2 w-full"
    :class="{ 'mm-financial-plan-mega': mega }"
  >
    <CardHeader
      class="flex flex-row items-start justify-between space-y-0"
      :class="mega ? 'mm-financial-plan-mega__card-header' : 'gap-3 p-4 pb-2 sm:p-5 sm:pb-3'"
    >
      <div :class="mega ? 'mm-financial-plan-mega__card-header-inner' : 'space-y-2'">
        <CardTitle :class="mega ? 'mm-financial-plan-mega__card-title' : 'text-xl font-semibold sm:text-2xl'">
          {{ mega ? ADVISOR.planTitleMega : ADVISOR.planTitle }}
        </CardTitle>
        <CardDescription :class="mega ? 'mm-financial-plan-mega__card-desc' : 'text-base'">
          {{ mega ? ADVISOR.planHintMega : ADVISOR.planHint }}
        </CardDescription>
      </div>
      <Button
        type="button"
        variant="outline"
        size="sm"
        class="shrink-0 gap-1.5"
        :disabled="loading"
        @click="emit('refresh')"
      >
        <RefreshCw class="size-3.5" :class="loading ? 'animate-spin' : ''" />
        Обновить
      </Button>
    </CardHeader>
    <CardContent :class="mega ? 'mm-financial-plan-mega__content' : 'space-y-4 p-4 pt-0 sm:p-5 sm:pt-0'">
      <SharedFinancialReportLoading
        v-if="loading && !plan"
        :active="loading"
        :compact="!mega"
      />

      <template v-else-if="plan">
        <section v-if="showOverview" class="mm-financial-plan-mega__section mm-financial-plan-mega__overview">
          <div class="mm-financial-plan-mega__section-head">
            <h3 class="mm-financial-plan-mega__heading">Обзор</h3>
            <p v-if="narrativeResolved.headline" class="mm-financial-plan-mega__hint">
              {{ narrativeResolved.headline }}
            </p>
            <p v-else-if="narrativeResolved.callout" class="mm-financial-plan-mega__hint">
              {{ narrativeResolved.callout }}
            </p>
          </div>

          <div class="mm-financial-plan-mega__section-body space-y-4">
            <div class="mm-narrative-hero mm-financial-plan-mega__narrative">
              <div class="mm-narrative-hero__top">
                <aside
                  v-if="narrativeResolved.weeklyAction"
                  class="mm-narrative-hero__advice"
                  aria-label="Совет недели"
                  data-demo="narrative"
                >
                  <div class="mm-narrative-hero__advice-meta">
                    <span class="mm-narrative-hero__advice-badge">{{ ADVISOR.weeklyAdviceTitle }}</span>
                    <span
                      v-if="narrativeResolved.goalOpportunityThousands"
                      class="mm-narrative-hero__advice-hook"
                    >
                      {{ GOALS.opportunityAmount(narrativeResolved.goalOpportunityThousands) }}
                    </span>
                  </div>
                  <p class="mm-narrative-hero__advice-text">{{ narrativeResolved.weeklyAction }}</p>
                  <p class="mm-narrative-hero__advice-hint">{{ narrativeResolved.adviceHint }}</p>
                </aside>

                <div class="mm-narrative-hero__aside">
                  <DashboardMindfulnessScore
                    :diagnosis="diagnosis"
                    :loading="diagnosisLoading && !diagnosis"
                  />
                </div>
              </div>

              <div
                v-if="showMoneyRow"
                class="mm-narrative-hero__money-row"
                aria-label="Доход и траты"
              >
                <div v-if="narrativeResolved.incomeDisplay" class="mm-narrative-hero__money-card">
                  <span class="mm-narrative-hero__money-label">Доход</span>
                  <span class="mm-narrative-hero__money-value">{{ narrativeResolved.incomeDisplay }}</span>
                </div>
                <div v-if="narrativeResolved.expensesDisplay" class="mm-narrative-hero__money-card">
                  <span class="mm-narrative-hero__money-label">Траты</span>
                  <span
                    class="mm-narrative-hero__money-value"
                    :class="{ 'text-amber-800': narrativeResolved.expensesWarn }"
                  >
                    {{ narrativeResolved.expensesDisplay }}
                  </span>
                </div>
              </div>
            </div>

            <DashboardMetricsGrid
              v-if="summary"
              embedded
              :summary="summary"
              :loading="overviewLoading"
            />
          </div>
        </section>

        <div :class="mega ? 'mm-financial-plan-mega__hero' : 'space-y-3'">
          <p :class="mega ? 'mm-financial-plan-mega__lead' : 'text-lg font-semibold leading-snug'">
            {{ plan.goalTitle }}
          </p>
          <div :class="mega ? 'mm-financial-plan-mega__meta' : 'flex flex-wrap gap-x-6 gap-y-2 text-sm text-muted-foreground'">
            <p>{{ plan.goalProgress }}</p>
            <p v-if="plan.runwayText">{{ plan.runwayText }}</p>
            <p v-if="plan.freeCashflowText">{{ plan.freeCashflowText }}</p>
          </div>
        </div>

        <!-- Mega: все разделы подряд -->
        <template v-if="mega">
          <section class="mm-financial-plan-mega__section">
            <div class="mm-financial-plan-mega__section-head">
              <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabSteps }}</h3>
            </div>
            <ol class="mm-financial-plan-mega__section-body mm-financial-plan-mega__steps">
              <li
                v-for="(step, i) in plan.steps"
                :key="`mega-step-${step.title}-${i}`"
                class="mm-financial-plan-mega__step"
              >
                <span
                  class="mm-plan-step-badge flex size-10 shrink-0 items-center justify-center rounded-full bg-primary/15 text-sm font-bold text-primary sm:size-11 sm:text-base"
                >
                  {{ i + 1 }}
                </span>
                <div class="mm-financial-plan-mega__step-copy">
                  <span v-if="i === 0" class="text-xs font-medium text-primary sm:text-sm">
                    {{ ADVISOR.planStepPrimaryLabel }}
                  </span>
                  <p class="mm-financial-plan-mega__step-title">{{ step.title }}</p>
                  <p v-if="step.description" class="mm-financial-plan-mega__step-desc">
                    {{ step.description }}
                  </p>
                  <AdvisorAskButton
                    v-if="i === 0"
                    class="mt-0.5"
                    :insight-title="step.title"
                    :insight-description="step.description"
                  />
                </div>
              </li>
            </ol>
          </section>

          <section class="mm-financial-plan-mega__section">
            <div class="mm-financial-plan-mega__section-head">
              <h3 class="mm-financial-plan-mega__heading">
                {{ ADVISOR.planTabOpportunity }}
                <span v-if="opportunityBadge" class="ml-2 text-base font-bold text-emerald-700 sm:text-lg">
                  {{ opportunityBadge }}
                </span>
              </h3>
              <p class="mm-financial-plan-mega__hint">{{ ADVISOR.planTabOpportunityHint }}</p>
            </div>
            <div v-if="summary" class="mm-financial-plan-mega__section-body mm-financial-plan-mega__opportunity">
              <p v-if="summary.goalOpportunityThousands" class="mm-financial-plan-mega__opportunity-value">
                {{ GOALS.opportunityAmount(summary.goalOpportunityThousands) }}
              </p>
              <p class="mm-financial-plan-mega__opportunity-text">{{ summary.goalForecast }}</p>
              <p
                v-if="summary.goalHint !== summary.goalForecast"
                class="mm-financial-plan-mega__opportunity-muted"
              >
                {{ summary.goalHint }}
              </p>
              <p v-if="summary.runwayMonths != null" class="mm-financial-plan-mega__opportunity-muted">
                {{ HEALTH.reserveMonths(summary.runwayMonths) }}
              </p>
              <AdvisorAskButton
                v-if="summary.goalOpportunityThousands"
                :prompt="`Что даёт ${GOALS.opportunityAmount(summary.goalOpportunityThousands)} и как приблизить цель?`"
                :label="ADVISOR.askAboutAction"
              />
            </div>
          </section>

          <section
            v-if="diagnosis?.indicators?.length || diagnosisLoading"
            class="mm-financial-plan-mega__section"
          >
            <div class="mm-financial-plan-mega__section-head">
              <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabDiagnosis }}</h3>
              <p v-if="diagnosisIntro" class="mm-financial-plan-mega__hint">{{ diagnosisIntro }}</p>
            </div>
            <div class="mm-financial-plan-mega__section-body">
              <DashboardDiagnosisIndicators
                :indicators="diagnosis?.indicators ?? []"
                :loading="diagnosisLoading && !diagnosis"
              />
              <AdvisorAskButton
                v-if="diagnosis?.main_action"
                :insight-title="diagnosis.main_action.title"
                :insight-description="diagnosis.main_action.description"
              />
            </div>
          </section>

          <section class="mm-financial-plan-mega__section">
            <div class="mm-financial-plan-mega__section-head">
              <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabMoney }}</h3>
              <p class="mm-financial-plan-mega__hint">{{ ADVISOR.planTabMoneyHint }}</p>
            </div>
            <div class="mm-financial-plan-mega__section-body">
              <DashboardMoneyPicture
                embedded
                :categories="categories"
                :timemachine="timemachine"
                :categories-summary="categoriesSummary"
                :current-savings="summary?.savingsBalance ?? null"
                :loading="chartsLoading"
              />
              <AdvisorAskButton
                prompt="Где больше всего уходит денег и что урезать в первую очередь?"
                label="Спросить про траты"
              />
            </div>
          </section>

          <section
            v-if="showCredits || creditsLoading"
            class="mm-financial-plan-mega__section"
          >
            <div class="mm-financial-plan-mega__section-head">
              <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabCredits }}</h3>
              <p class="mm-financial-plan-mega__hint">{{ ADVISOR.planTabCreditsHint }}</p>
            </div>
            <div class="mm-financial-plan-mega__section-body">
              <DashboardCreditsSnapshot
                embedded
                :credits="credits"
                :dti-tone="dtiTone ?? 'warn'"
                :loading="creditsLoading"
              />
              <AdvisorAskButton
                prompt="Стоит ли рефинансировать кредит при моей ставке?"
                label="Спросить про кредит"
              />
            </div>
          </section>

          <section class="mm-financial-plan-mega__section">
            <div class="mm-financial-plan-mega__section-head">
              <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabExplore }}</h3>
              <p class="mm-financial-plan-mega__hint">{{ ADVISOR.planTabExploreHint }}</p>
            </div>
            <div class="mm-financial-plan-mega__section-body gap-4">
              <AnalyticsScenarioSimulator
                v-model:selected-category="selectedCategory"
                v-model:percent="percent"
                embedded
                :profile="profile"
                :categories="categories"
              />
              <AnalyticsDetective v-if="insights.length" embedded :insights="insights" />
            </div>
          </section>
        </template>

        <!-- Компакт: вкладки -->
        <Tabs v-else v-model="planTab" class="w-full">
          <TabsList class="mm-tier-3 h-auto w-full flex-wrap justify-start gap-1 bg-muted/50 p-1">
            <TabsTrigger value="steps" class="text-sm">
              {{ ADVISOR.planTabSteps }}
            </TabsTrigger>
            <TabsTrigger value="opportunity" class="gap-1.5 text-sm">
              {{ ADVISOR.planTabOpportunity }}
              <span v-if="opportunityBadge" class="text-xs font-semibold text-emerald-700">
                {{ opportunityBadge }}
              </span>
            </TabsTrigger>
            <TabsTrigger
              v-if="diagnosis?.indicators?.length || diagnosisLoading"
              value="diagnosis"
              class="text-sm"
            >
              {{ ADVISOR.planTabDiagnosis }}
            </TabsTrigger>
            <TabsTrigger value="money" class="text-sm">
              {{ ADVISOR.planTabMoney }}
            </TabsTrigger>
            <TabsTrigger v-if="showCredits || creditsLoading" value="credits" class="text-sm">
              {{ ADVISOR.planTabCredits }}
            </TabsTrigger>
            <TabsTrigger value="explore" class="text-sm">
              {{ ADVISOR.planTabExplore }}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="steps" class="mt-4 space-y-3">
            <ol class="grid gap-3 md:grid-cols-3">
              <li
                v-for="(step, i) in plan.steps"
                :key="`${step.title}-${i}`"
                class="flex items-start gap-3 rounded-xl border bg-muted/20 p-3"
              >
                <span
                  class="mm-plan-step-badge flex size-9 shrink-0 items-center justify-center rounded-full bg-primary/15 text-sm font-bold text-primary"
                >
                  {{ i + 1 }}
                </span>
                <div class="min-w-0 space-y-1">
                  <span v-if="i === 0" class="text-xs font-medium text-primary">
                    {{ ADVISOR.planStepPrimaryLabel }}
                  </span>
                  <p class="text-base font-medium leading-snug">{{ step.title }}</p>
                  <p v-if="step.description" class="text-sm leading-relaxed text-muted-foreground">
                    {{ step.description }}
                  </p>
                  <AdvisorAskButton
                    v-if="i === 0"
                    class="mt-2"
                    :insight-title="step.title"
                    :insight-description="step.description"
                  />
                </div>
              </li>
            </ol>
          </TabsContent>

          <TabsContent value="opportunity" class="mt-4 space-y-3">
            <p class="text-sm text-muted-foreground">{{ ADVISOR.planTabOpportunityHint }}</p>
            <div v-if="summary" class="mm-tier-3 space-y-3 rounded-xl bg-muted/30 p-4">
              <p
                v-if="summary.goalOpportunityThousands"
                class="text-2xl font-bold text-emerald-700"
              >
                {{ GOALS.opportunityAmount(summary.goalOpportunityThousands) }}
              </p>
              <p class="text-base leading-relaxed">{{ summary.goalForecast }}</p>
              <p
                v-if="summary.goalHint !== summary.goalForecast"
                class="text-sm leading-relaxed text-muted-foreground"
              >
                {{ summary.goalHint }}
              </p>
              <p v-if="summary.runwayMonths != null" class="text-sm text-muted-foreground">
                {{ HEALTH.reserveMonths(summary.runwayMonths) }}
              </p>
              <AdvisorAskButton
                v-if="summary.goalOpportunityThousands"
                :prompt="`Что даёт ${GOALS.opportunityAmount(summary.goalOpportunityThousands)} и как приблизить цель?`"
                :label="ADVISOR.askAboutAction"
              />
            </div>
          </TabsContent>

          <TabsContent value="diagnosis" class="mt-4 space-y-3">
            <p v-if="diagnosisIntro" class="text-sm text-muted-foreground">{{ diagnosisIntro }}</p>
            <DashboardDiagnosisIndicators
              :indicators="diagnosis?.indicators ?? []"
              :loading="diagnosisLoading && !diagnosis"
            />
            <AdvisorAskButton
              v-if="diagnosis?.main_action"
              :insight-title="diagnosis.main_action.title"
              :insight-description="diagnosis.main_action.description"
            />
          </TabsContent>

          <TabsContent value="money" class="mt-4 space-y-3">
            <p class="text-sm text-muted-foreground">{{ ADVISOR.planTabMoneyHint }}</p>
            <DashboardMoneyPicture
              embedded
              :categories="categories"
              :timemachine="timemachine"
              :categories-summary="categoriesSummary"
              :current-savings="summary?.savingsBalance ?? null"
              :loading="chartsLoading"
            />
            <AdvisorAskButton
              prompt="Где больше всего уходит денег и что урезать в первую очередь?"
              label="Спросить про траты"
            />
          </TabsContent>

          <TabsContent value="credits" class="mt-4 space-y-3">
            <p class="text-sm text-muted-foreground">{{ ADVISOR.planTabCreditsHint }}</p>
            <DashboardCreditsSnapshot
              embedded
              :credits="credits"
              :dti-tone="dtiTone ?? 'warn'"
              :loading="creditsLoading"
            />
            <AdvisorAskButton
              prompt="Стоит ли рефинансировать кредит при моей ставке?"
              label="Спросить про кредит"
            />
          </TabsContent>

          <TabsContent value="explore" class="mt-4 space-y-4">
            <p class="text-sm text-muted-foreground">{{ ADVISOR.planTabExploreHint }}</p>
            <AnalyticsScenarioSimulator
              v-model:selected-category="selectedCategory"
              v-model:percent="percent"
              embedded
              :profile="profile"
              :categories="categories"
            />
            <AnalyticsDetective v-if="insights.length" embedded :insights="insights" />
          </TabsContent>
        </Tabs>
      </template>

      <p v-else :class="mega ? 'text-lg text-muted-foreground' : 'text-base text-muted-foreground'">
        Пройдите опрос в профиле — соберём персональный план.
      </p>
    </CardContent>
  </Card>
</template>
