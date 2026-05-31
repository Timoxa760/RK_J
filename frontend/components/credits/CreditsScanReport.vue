<script setup lang="ts">
import { Check, Plus } from 'lucide-vue-next'
import type { CreditScanResponse } from '~/types/api'
import { CREDITS, PROFILE } from '~/constants/productCopy'
import {
  buildCreditScanReport,
  creditFixedExpenseTitle,
  isCreditInFixedExpenses,
  mergeCreditFixedExpense
} from '~/utils/creditReport'

const props = defineProps<{
  scan: CreditScanResponse
}>()

const { enrichedDashboard } = useCredits()
const { profile, loadProfile, saveProfile, syncProfileToApi } = useFinancialProfile()
const { refreshAdvisorContext } = useAdvisorContext()

const addingFixed = ref(false)
const fixedAdded = ref(false)
const fixedError = ref<string | null>(null)

const report = computed(() =>
  buildCreditScanReport(props.scan, {
    monthlyIncome: enrichedDashboard.value?.monthly_income || undefined,
    fixedExpensesTotal: (profile.value.fixed_expenses ?? []).reduce(
      (sum, row) => sum + Math.max(0, row.amount),
      0
    )
  })
)

const fixedTitle = computed(() => creditFixedExpenseTitle(props.scan.parsed.bank))

const alreadyInFixed = computed(() =>
  isCreditInFixedExpenses(profile.value.fixed_expenses, fixedTitle.value, report.value.monthlyPayment)
)

onMounted(() => {
  loadProfile()
})

watch(
  () => props.scan.credit_id,
  () => {
    fixedAdded.value = false
    fixedError.value = null
  }
)

async function addToFixedExpenses() {
  if (report.value.monthlyPayment <= 0) return
  addingFixed.value = true
  fixedError.value = null
  try {
    const next = mergeCreditFixedExpense(
      profile.value.fixed_expenses ?? [],
      fixedTitle.value,
      report.value.monthlyPayment
    )
    saveProfile({
      fixed_expenses: next,
      skipped_expenses: false
    })
    await syncProfileToApi(profile.value)
    await refreshAdvisorContext()
    fixedAdded.value = true
  } catch (e) {
    fixedError.value = e instanceof Error ? e.message : 'Не удалось сохранить'
  } finally {
    addingFixed.value = false
  }
}
</script>

<template>
  <Card class="overflow-hidden" data-demo="credit-scan-report">
    <CardHeader class="gap-2 border-b bg-muted/20 pb-4">
      <p class="text-xs font-medium uppercase tracking-wide text-muted-foreground">
        {{ report.title }}
      </p>
      <CardTitle class="text-xl leading-tight">{{ report.subtitle }}</CardTitle>
      <CardDescription v-if="scan.confidence != null" class="text-xs">
        {{ CREDITS.scanConfidence(scan.confidence) }}
      </CardDescription>
    </CardHeader>

    <CardContent class="space-y-5 pt-5">
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <SharedKeyTextStrip
          v-for="metric in report.metrics"
          :key="metric.label"
          :label="metric.label"
          :value="metric.value"
          :tone="metric.tone"
        />
      </div>

      <p
        v-for="(hint, index) in report.metrics.filter((m) => m.hint).map((m) => m.hint)"
        :key="`hint-${index}`"
        class="text-xs text-muted-foreground"
      >
        {{ hint }}
      </p>

      <section
        class="rounded-xl border bg-muted/30 p-4"
        aria-label="Выводы по договору"
      >
        <h3 class="text-sm font-semibold">Выводы</h3>
        <ul class="mt-3 space-y-2 text-sm leading-relaxed text-muted-foreground">
          <li v-for="(line, index) in report.insights" :key="index" class="flex gap-2">
            <span class="mt-2 size-1.5 shrink-0 rounded-full bg-primary" aria-hidden="true" />
            <span>{{ line }}</span>
          </li>
        </ul>
        <p v-if="report.rateVsMarketLabel" class="mt-3 text-sm font-medium">
          {{ report.rateVsMarketLabel }}
        </p>
      </section>

      <div class="rounded-xl border border-primary/20 bg-primary/5 p-4">
        <p class="text-sm font-medium">{{ PROFILE.fixedExpensesTitle }}</p>
        <p class="mt-1 text-sm text-muted-foreground">
          {{ PROFILE.fixedExpensesHint }}
        </p>
        <div class="mt-4 flex flex-wrap items-center gap-3">
          <Button
            v-if="!alreadyInFixed && !fixedAdded"
            type="button"
            class="gap-2"
            :disabled="addingFixed || report.monthlyPayment <= 0"
            data-demo="credit-add-fixed-expense"
            @click="addToFixedExpenses"
          >
            <Plus class="size-4" />
            {{ addingFixed ? 'Сохраняем…' : `Внести ${report.monthlyPayment.toLocaleString('ru-RU')} ₽/мес` }}
          </Button>
          <div
            v-else
            class="inline-flex items-center gap-2 rounded-lg border border-primary/30 bg-background px-3 py-2 text-sm"
          >
            <Check class="size-4 text-primary" />
            <span>{{ fixedTitle }} уже в обязательных расходах</span>
          </div>
        </div>
        <p v-if="fixedError" class="mt-2 text-sm text-destructive">{{ fixedError }}</p>
      </div>
    </CardContent>
  </Card>
</template>
