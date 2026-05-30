<script setup lang="ts">
import type { TimeMachineResponse } from '~/types/api'
import { GOALS, formatRub } from '~/constants/productCopy'

const props = defineProps<{
  data: TimeMachineResponse | null
  /** Текущий баланс накоплений из профиля / credits — приоритетнее первой точки timemachine */
  currentSavings?: number | null
}>()

function resolveCurrentSavings(points: TimeMachineResponse['points']): number | null {
  if (props.currentSavings != null && props.currentSavings > 0) {
    return props.currentSavings
  }
  const first = points[0]
  if (first && first.actual > 0) return first.actual
  return null
}

const snapshot = computed(() => {
  const points = props.data?.points ?? []
  const now = resolveCurrentSavings(points)
  if (now == null) return null

  if (!points.length) {
    return {
      now,
      months: null as number | null,
      forecast: null as number | null,
      optimistic: null as number | null,
      diff: 0
    }
  }

  const last = points[points.length - 1]!
  const months = points.length
  const diff = last.optimistic - last.actual

  return {
    months,
    now,
    forecast: last.actual,
    optimistic: last.optimistic,
    diff
  }
})
</script>

<template>
  <div v-if="snapshot" class="mm-simple-chart mm-savings-chart">
    <div class="mm-savings-chart__grid" :class="{ 'sm:grid-cols-1': snapshot.forecast == null }">
      <div class="mm-savings-chart__card">
        <p class="mm-savings-chart__caption">{{ GOALS.savingsCurrentLabel }}</p>
        <p class="mm-savings-chart__amount">{{ formatRub(snapshot.now) }}</p>
      </div>
      <div
        v-if="snapshot.forecast != null && snapshot.months != null"
        class="mm-savings-chart__card mm-savings-chart__card--accent"
      >
        <p class="mm-savings-chart__caption">Через {{ snapshot.months }} мес.</p>
        <p class="mm-savings-chart__amount">{{ formatRub(snapshot.forecast) }}</p>
      </div>
    </div>

    <div v-if="snapshot.forecast != null && snapshot.diff > 0" class="mm-savings-chart__hint">
      {{ GOALS.savingsChartOpportunity(formatRub(snapshot.diff), formatRub(snapshot.optimistic!)) }}
    </div>
    <div v-else-if="snapshot.forecast != null" class="mm-savings-chart__hint">
      {{ GOALS.savingsEven }}
    </div>
  </div>
  <p v-else class="mm-simple-chart__empty">Нет данных о накоплениях</p>
</template>
