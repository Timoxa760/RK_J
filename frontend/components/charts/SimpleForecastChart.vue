<script setup lang="ts">
import type { ForecastResponse } from '~/types/api'
import { brandColors } from '~/utils/chartTheme'
import { formatChartDay } from '~/utils/chartSummaries'
import { normalizeForecast } from '~/utils/apiNormalize'
import { formatRub } from '~/constants/productCopy'

const props = defineProps<{
  data: ForecastResponse | null
}>()

const rows = computed(() => {
  const normalized = props.data ? normalizeForecast(props.data) : null
  if (!normalized?.forecast.length) return []

  const max = Math.max(...normalized.forecast, 1)

  return normalized.forecast.map((amount, i) => ({
    label: formatChartDay(normalized.dates[i] ?? `День ${i + 1}`),
    amount,
    width: Math.max(12, Math.round((amount / max) * 100))
  }))
})
</script>

<template>
  <ul v-if="rows.length" class="mm-simple-chart space-y-3" role="list">
    <li v-for="row in rows" :key="row.label" class="mm-simple-chart__row">
      <div class="mm-simple-chart__label-row">
        <span class="mm-simple-chart__label">{{ row.label }}</span>
        <span class="mm-simple-chart__value">{{ formatRub(row.amount) }}</span>
      </div>
      <div class="mm-simple-chart__track" aria-hidden="true">
        <div
          class="mm-simple-chart__bar"
          :style="{ width: `${row.width}%`, backgroundColor: brandColors.primary }"
        />
      </div>
    </li>
  </ul>
  <p v-else class="mm-simple-chart__empty">Нет данных для прогноза</p>
</template>
