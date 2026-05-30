<script setup lang="ts">
import { CREDITS } from '~/constants/productCopy'
import type { CreditsDashboardResponse } from '~/types/api'
import type { HealthTone } from '~/utils/dashboardSummary'
import { buildCreditsIntro } from '~/utils/dashboardCopy'

const props = defineProps<{
  credits: CreditsDashboardResponse | null
  dtiTone: HealthTone
  loading?: boolean
  embedded?: boolean
}>()

const introLine = computed(() =>
  props.credits ? buildCreditsIntro(props.credits.dti) : null
)

function dtiBarClass(tone: HealthTone) {
  if (tone === 'good') return 'bg-primary'
  if (tone === 'warn') return 'bg-amber-500'
  return 'bg-destructive'
}
</script>

<template>
  <component :is="embedded ? 'div' : 'Card'" data-demo="credits-dti" :class="embedded ? 'space-y-3' : undefined">
    <component :is="embedded ? 'div' : 'CardHeader'" :class="embedded ? undefined : 'pb-2'">
      <div class="flex items-start justify-between gap-3">
        <div>
          <CardTitle class="text-base">{{ CREDITS.paymentsTitle }}</CardTitle>
          <CardDescription v-if="!embedded">{{ CREDITS.paymentsHint }}</CardDescription>
        </div>
        <Button variant="link" class="h-auto shrink-0 p-0 text-xs" as-child>
          <NuxtLink to="/credits">Подробнее →</NuxtLink>
        </Button>
      </div>
    </component>
    <component :is="embedded ? 'div' : 'CardContent'">
      <Skeleton v-if="loading && !credits" class="h-14 w-full" />
      <template v-else-if="credits">
        <p v-if="embedded && introLine" class="text-sm font-medium">{{ introLine }}</p>
        <p class="text-3xl font-bold">{{ credits.dti }}%</p>
        <p class="mt-1 text-xs text-muted-foreground">
          Это {{ credits.dti }} ₽ из каждых 100 ₽ дохода на погашение кредитов
        </p>
        <div class="mt-3 h-2 overflow-hidden rounded-full bg-muted">
          <div
            class="h-full rounded-full transition-all"
            :class="dtiBarClass(dtiTone)"
            :style="{ width: `${Math.min(100, credits.dti)}%` }"
          />
        </div>
        <div class="mt-3 grid gap-2 text-sm text-muted-foreground sm:grid-cols-2">
          <p v-if="credits.stress_test_months != null">
            {{ CREDITS.stressReserveMonths(credits.stress_test_months) }}
          </p>
          <p v-if="credits.monthly_payments">
            Платежи: {{ credits.monthly_payments.toLocaleString('ru-RU') }} ₽/мес
          </p>
        </div>
      </template>
    </component>
  </component>
</template>
