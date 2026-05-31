<script setup lang="ts">
import { CREDITS, formatRub } from '~/constants/productCopy'
import type { EnrichedCreditsDashboard } from '~/utils/creditsDashboard'
import type { HealthTone } from '~/utils/dashboardSummary'

const props = defineProps<{
  dashboard: EnrichedCreditsDashboard
}>()

function dtiTone(dti: number): HealthTone {
  if (dti < 35) return 'good'
  if (dti < 50) return 'warn'
  return 'risk'
}

function dtiBarClass(tone: HealthTone) {
  if (tone === 'good') return 'bg-primary'
  if (tone === 'warn') return 'bg-amber-500'
  return 'bg-destructive'
}

const tone = computed(() => (props.dashboard.dti_available ? dtiTone(props.dashboard.dti) : 'warn'))
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2">
    <Card data-demo="credits-dti">
      <CardHeader class="pb-2">
        <CardTitle class="text-base">{{ CREDITS.trafficLight }}</CardTitle>
        <CardDescription v-if="dashboard.dti_available">{{ CREDITS.paymentsHint }}</CardDescription>
      </CardHeader>
      <CardContent>
        <template v-if="dashboard.dti_available">
          <p class="text-3xl font-bold tabular-nums">{{ dashboard.dti }}%</p>
          <p class="mt-1 text-xs text-muted-foreground">{{ CREDITS.incomeShare(dashboard.dti) }}</p>
          <div class="mt-3 h-2 overflow-hidden rounded-full bg-muted">
            <div
              class="h-full rounded-full transition-all"
              :class="dtiBarClass(tone)"
              :style="{ width: `${Math.min(100, dashboard.dti)}%` }"
            />
          </div>
        </template>

        <template v-else>
          <p class="text-3xl font-bold tabular-nums">
            {{ formatRub(dashboard.monthly_payments ?? 0) }}
            <span class="text-lg font-semibold text-muted-foreground">/мес</span>
          </p>
          <p class="mt-2 text-sm text-muted-foreground">
            {{ CREDITS.dtiNeedIncome }}
          </p>
          <Button variant="link" class="mt-1 h-auto p-0 text-sm" as-child>
            <NuxtLink to="/profile">{{ CREDITS.dtiIncomeCta }}</NuxtLink>
          </Button>
        </template>

        <p v-if="dashboard.monthly_payments" class="mt-3 text-sm text-muted-foreground">
          {{ CREDITS.monthlyPaymentsLine(dashboard.monthly_payments) }}
        </p>
        <p
          v-if="dashboard.stress_test_months != null && dashboard.dti_available"
          class="mt-2 text-sm text-muted-foreground"
        >
          {{ CREDITS.stressReserveMonths(dashboard.stress_test_months) }}
        </p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="pb-2">
        <CardTitle class="text-base">{{ CREDITS.cushionTitle }}</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-3xl font-bold tabular-nums">{{ dashboard.savings.toLocaleString('ru-RU') }} ₽</p>
        <p
          v-if="dashboard.stress_test_months != null"
          class="mt-2 text-sm text-muted-foreground"
        >
          {{ CREDITS.stressReserveMonths(dashboard.stress_test_months) }}
        </p>
        <p
          v-else-if="!dashboard.savings"
          class="mt-2 text-sm text-muted-foreground"
        >
          {{ CREDITS.cushionNeedData }}
        </p>
      </CardContent>
    </Card>
  </div>
</template>
