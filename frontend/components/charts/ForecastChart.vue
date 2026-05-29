<script setup lang="ts">
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { ForecastResponse } from '~/types/api'
import { chartThemeLight } from '~/types/api'

use([CanvasRenderer, BarChart, TooltipComponent, GridComponent])

const props = defineProps<{
  data: ForecastResponse | null
}>()

const chartTheme = chartThemeLight

const option = computed(() => {
  if (!props.data?.points.length) return null
  return {
    backgroundColor: chartTheme.backgroundColor,
    tooltip: { trigger: 'axis' },
    grid: { left: 48, right: 16, top: 24, bottom: 32 },
    xAxis: {
      type: 'category',
      data: props.data.points.map((p) => p.month),
      axisLabel: { color: chartTheme.textStyle.color }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: chartTheme.textStyle.color },
      splitLine: { lineStyle: { color: chartTheme.splitLine } }
    },
    series: [
      {
        type: 'bar',
        data: props.data.points.map((p) => p.amount),
        itemStyle: { color: '#e8955f', borderRadius: [4, 4, 0, 0] }
      }
    ]
  }
})
</script>

<template>
  <ClientOnly>
    <VChart
      v-if="option"
      class="h-full min-h-[220px] max-h-[420px] w-full sm:min-h-[260px] lg:min-h-[280px]"
      :option="option"
      autoresize
    />
    <p v-else class="flex min-h-[220px] items-center justify-center text-sm text-slate-500 sm:min-h-[260px] lg:min-h-[280px]">
      Нет данных
    </p>
  </ClientOnly>
</template>
