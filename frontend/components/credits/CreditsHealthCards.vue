<script setup lang="ts">
import { CREDITS } from '~/constants/productCopy'
import type { CreditsDashboardResponse } from '~/types/api'
import type { HealthTone } from '~/utils/dashboardSummary'

defineProps<{
  dashboard: CreditsDashboardResponse
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
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2">
    <Card data-demo="credits-dti">
      <CardHeader class="pb-2">
        <CardTitle class="text-base">{{ CREDITS.trafficLight }}</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-3xl font-bold">{{ dashboard.dti }}%</p>
        <p class="mt-1 text-xs text-muted-foreground">{{ CREDITS.incomeShare(dashboard.dti) }}</p>
        <div class="mt-3 h-2 overflow-hidden rounded-full bg-muted">
          <div
            class="h-full rounded-full transition-all"
            :class="dtiBarClass(dtiTone(dashboard.dti))"
            :style="{ width: `${Math.min(100, dashboard.dti)}%` }"
          />
        </div>
        <p
          v-if="dashboard.stress_test_dti != null"
          class="mt-3 text-sm text-muted-foreground"
        >
          {{ CREDITS.stressIncomeDrop(dashboard.stress_test_dti) }}
        </p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="pb-2">
        <CardTitle class="text-base">{{ CREDITS.cushionTitle }}</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-3xl font-bold">{{ dashboard.savings.toLocaleString('ru-RU') }} ₽</p>
        <p
          v-if="dashboard.stress_test_months != null"
          class="mt-2 text-sm text-muted-foreground"
        >
          {{ CREDITS.stressReserveMonths(dashboard.stress_test_months) }}
        </p>
      </CardContent>
    </Card>
  </div>
</template>
