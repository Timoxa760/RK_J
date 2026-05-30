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
    /** Полный план на дашборде — все блоки подряд */
    mega?: boolean
  }>(),
  { mega: false }
)

const scenario = defineModel<'reduce_delivery' | 'reduce_cafe' | 'reduce_entertainment' | 'custom'>(
  'scenario',
  { default: 'reduce_cafe' }
)
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
</script>

<template>
  <Card
    id="financial-plan"
    data-demo="financial-plan"
    class="mm-tier-2 w-full"
    :class="{ 'mm-financial-plan-mega': mega }"
  >
    <CardHeader
      class="flex flex-row items-start justify-between gap-3 space-y-0"
      :class="mega ? 'p-5 pb-3 sm:p-8 sm:pb-4' : 'p-4 pb-2 sm:p-5 sm:pb-3'"
    >
      <div class="space-y-2">
        <CardTitle :class="mega ? 'text-2xl font-bold sm:text-3xl' : 'text-xl font-semibold sm:text-2xl'">
          {{ mega ? ADVISOR.planTitleMega : ADVISOR.planTitle }}
        </CardTitle>
        <CardDescription :class="mega ? 'max-w-3xl text-base sm:text-lg' : 'text-base'">
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
    <CardContent :class="mega ? 'space-y-10 p-5 pt-0 sm:p-8 sm:pt-0' : 'space-y-4 p-4 pt-0 sm:p-5 sm:pt-0'">
      <Skeleton v-if="loading && !plan" :class="mega ? 'h-48 w-full rounded-2xl' : 'h-36 w-full rounded-xl'" />

      <template v-else-if="plan">
        <div
          class="space-y-3"
          :class="mega ? 'mm-financial-plan-mega__hero rounded-2xl border bg-muted/25 p-5 sm:p-6' : ''"
        >
          <p :class="mega ? 'text-2xl font-bold leading-tight sm:text-3xl' : 'text-lg font-semibold leading-snug'">
            {{ plan.goalTitle }}
          </p>
          <div
            class="flex flex-wrap gap-x-6 gap-y-2 text-muted-foreground"
            :class="mega ? 'text-base sm:text-lg' : 'text-sm'"
          >
            <p>{{ plan.goalProgress }}</p>
            <p v-if="plan.runwayText">{{ plan.runwayText }}</p>
            <p v-if="plan.freeCashflowText">{{ plan.freeCashflowText }}</p>
          </div>
          <p
            v-if="mega && summary?.weeklyAction"
            class="border-t border-border/60 pt-4 text-base leading-relaxed sm:text-lg"
          >
            <span class="font-semibold text-foreground">Главное на неделю: </span>
            {{ summary.weeklyAction }}
          </p>
        </div>

        <!-- Mega: все разделы подряд -->
        <template v-if="mega">
          <section class="mm-financial-plan-mega__section space-y-4">
            <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabSteps }}</h3>
            <ol class="grid gap-4 lg:grid-cols-1">
              <li
                v-for="(step, i) in plan.steps"
                :key="`mega-step-${step.title}-${i}`"
                class="flex items-start gap-4 rounded-2xl border bg-muted/20 p-4 sm:p-5"
              >
                <span
                  class="mm-plan-step-badge flex size-11 shrink-0 items-center justify-center rounded-full bg-primary/15 text-base font-bold text-primary sm:size-12 sm:text-lg"
                >
                  {{ i + 1 }}
                </span>
                <div class="min-w-0 space-y-2">
                  <span v-if="i === 0" class="text-sm font-medium text-primary">
                    {{ ADVISOR.planStepPrimaryLabel }}
                  </span>
                  <p class="text-lg font-semibold leading-snug sm:text-xl">{{ step.title }}</p>
                  <p v-if="step.description" class="text-base leading-relaxed text-muted-foreground">
                    {{ step.description }}
                  </p>
                  <AdvisorAskButton
                    v-if="i === 0"
                    class="mt-1"
                    :insight-title="step.title"
                    :insight-description="step.description"
                  />
                </div>
              </li>
            </ol>
          </section>

          <section class="mm-financial-plan-mega__section space-y-4">
            <h3 class="mm-financial-plan-mega__heading">
              {{ ADVISOR.planTabOpportunity }}
              <span v-if="opportunityBadge" class="ml-2 text-lg font-bold text-emerald-700">
                {{ opportunityBadge }}
              </span>
            </h3>
            <p class="text-base text-muted-foreground">{{ ADVISOR.planTabOpportunityHint }}</p>
            <div v-if="summary" class="mm-tier-3 space-y-4 rounded-2xl bg-muted/30 p-5 sm:p-6">
              <p
                v-if="summary.goalOpportunityThousands"
                class="text-3xl font-bold text-emerald-700 sm:text-4xl"
              >
                {{ GOALS.opportunityAmount(summary.goalOpportunityThousands) }}
              </p>
              <p class="text-lg leading-relaxed">{{ summary.goalForecast }}</p>
              <p
                v-if="summary.goalHint !== summary.goalForecast"
                class="text-base leading-relaxed text-muted-foreground"
              >
                {{ summary.goalHint }}
              </p>
              <p v-if="summary.runwayMonths != null" class="text-base text-muted-foreground">
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
            class="mm-financial-plan-mega__section space-y-4"
          >
            <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabDiagnosis }}</h3>
            <p v-if="diagnosisIntro" class="text-base text-muted-foreground">{{ diagnosisIntro }}</p>
            <DashboardDiagnosisIndicators
              :indicators="diagnosis?.indicators ?? []"
              :loading="diagnosisLoading && !diagnosis"
            />
          </section>

          <section class="mm-financial-plan-mega__section space-y-4">
            <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabMoney }}</h3>
            <p class="text-base text-muted-foreground">{{ ADVISOR.planTabMoneyHint }}</p>
            <DashboardMoneyPicture
              embedded
              :categories="categories"
              :timemachine="timemachine"
              :categories-summary="categoriesSummary"
              :current-savings="summary?.savingsBalance ?? null"
              :loading="chartsLoading"
            />
          </section>

          <section
            v-if="showCredits || creditsLoading"
            class="mm-financial-plan-mega__section space-y-4"
          >
            <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabCredits }}</h3>
            <p class="text-base text-muted-foreground">{{ ADVISOR.planTabCreditsHint }}</p>
            <DashboardCreditsSnapshot
              embedded
              :credits="credits"
              :dti-tone="dtiTone ?? 'warn'"
              :loading="creditsLoading"
            />
          </section>

          <section class="mm-financial-plan-mega__section space-y-5">
            <h3 class="mm-financial-plan-mega__heading">{{ ADVISOR.planTabExplore }}</h3>
            <p class="text-base text-muted-foreground">{{ ADVISOR.planTabExploreHint }}</p>
            <AnalyticsScenarioSimulator
              v-model:scenario="scenario"
              v-model:percent="percent"
              embedded
              :profile="profile"
              :categories="categories"
            />
            <AnalyticsDetective v-if="insights.length" embedded :insights="insights" />
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
          </TabsContent>

          <TabsContent value="credits" class="mt-4 space-y-3">
            <p class="text-sm text-muted-foreground">{{ ADVISOR.planTabCreditsHint }}</p>
            <DashboardCreditsSnapshot
              embedded
              :credits="credits"
              :dti-tone="dtiTone ?? 'warn'"
              :loading="creditsLoading"
            />
          </TabsContent>

          <TabsContent value="explore" class="mt-4 space-y-4">
            <p class="text-sm text-muted-foreground">{{ ADVISOR.planTabExploreHint }}</p>
            <AnalyticsScenarioSimulator
              v-model:scenario="scenario"
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
