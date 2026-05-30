<script setup lang="ts">
import { GOALS, HEALTH } from '~/constants/productCopy'
import type { DashboardSummary } from '~/utils/dashboardSummary'

defineProps<{
  summary: DashboardSummary
  loading?: boolean
}>()

const toneVariant = {
  good: 'default' as const,
  warn: 'secondary' as const,
  risk: 'destructive' as const
}
</script>

<template>
  <section aria-label="Ключевые показатели">
    <div v-if="loading" class="grid w-full grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <Skeleton v-for="i in 4" :key="i" class="h-28 w-full" />
    </div>

    <div v-else class="grid w-full grid-cols-1 items-stretch gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <Card class="flex h-full flex-col">
        <CardHeader class="pb-2">
          <CardDescription>Сейчас</CardDescription>
          <CardTitle class="mm-heading-stretch text-lg">
            {{ summary.income.toLocaleString('ru-RU') }} ₽
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-1 text-sm text-muted-foreground">
          <p>расходы ≈ {{ summary.expenses.toLocaleString('ru-RU') }} ₽/мес</p>
          <p :class="summary.freeCashflow >= 0 ? 'text-emerald-700' : 'text-amber-800'">
            {{ HEALTH.leftAfterExpenses }}:
            {{ summary.freeCashflow >= 0 ? '+' : '' }}{{ summary.freeCashflow.toLocaleString('ru-RU') }} ₽
          </p>
          <p>остаток на счёте ≈ {{ summary.savingsBalance.toLocaleString('ru-RU') }} ₽</p>
        </CardContent>
      </Card>

      <Card class="flex h-full flex-col">
        <CardHeader class="pb-2">
          <CardDescription>{{ GOALS.opportunityLabel }}</CardDescription>
          <CardTitle
            v-if="summary.goalOpportunityThousands"
            class="mm-heading-stretch text-lg text-emerald-700"
          >
            {{ GOALS.opportunityAmount(summary.goalOpportunityThousands) }}
          </CardTitle>
          <CardTitle v-else class="text-base font-medium leading-snug">
            Прогноз накоплений
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-2 text-sm text-muted-foreground">
          <p>{{ summary.goalForecast }}</p>
          <p v-if="summary.goalHint !== summary.goalForecast" class="text-xs">{{ summary.goalHint }}</p>
        </CardContent>
      </Card>

      <Card class="flex h-full flex-col">
        <CardHeader class="pb-2">
          <CardDescription>Запас на чёрный день</CardDescription>
          <Badge :variant="toneVariant[summary.stabilityTone]" class="w-fit">
            {{ summary.stabilityLabel }}
          </Badge>
        </CardHeader>
        <CardContent class="text-sm text-muted-foreground">
          <template v-if="summary.runwayMonths != null">
            {{ HEALTH.reserveMonths(summary.runwayMonths) }}
          </template>
          <template v-else>{{ HEALTH.reserveUnknown }}</template>
        </CardContent>
      </Card>

      <Card class="flex h-full flex-col border-primary/25 bg-primary/5">
        <CardHeader class="pb-2">
          <CardDescription class="text-primary">Что сделать</CardDescription>
          <CardTitle class="mm-heading-stretch text-base leading-snug">
            {{ summary.weeklyAction }}
          </CardTitle>
        </CardHeader>
        <CardContent v-if="summary.behaviorInsight">
          <p class="text-xs leading-relaxed text-muted-foreground">{{ summary.behaviorInsight }}</p>
        </CardContent>
      </Card>
    </div>
  </section>
</template>
