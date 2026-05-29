<script setup lang="ts">
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { TimeMachineResponse } from '~/types/api'
import { chartThemeLight } from '~/types/api'

use([CanvasRenderer, LineChart, TooltipComponent, LegendComponent, GridComponent])

const props = defineProps<{
  data: TimeMachineResponse | null
}>()

const chartTheme = chartThemeLight

const option = computed(() => {
  if (!props.data?.points.length) return null
  const months = props.data.points.map((p) => p.month)
  return {
    backgroundColor: chartTheme.backgroundColor,
    tooltip: { trigger: 'axis' },
    legend: {
      data: ['Факт', 'Оптимистичный'],
      textStyle: { color: chartTheme.textStyle.color }
    },
    grid: { left: 48, right: 24, top: 40, bottom: 32 },
    xAxis: {
      type: 'category',
      data: months,
      axisLabel: { color: chartTheme.textStyle.color, rotate: 45, fontSize: 10 }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: chartTheme.textStyle.color,
        formatter: (v: number) => `${Math.round(v / 1000)}k`
      },
      splitLine: { lineStyle: { color: chartTheme.splitLine } }
    },
    series: [
      {
        name: 'Факт',
        type: 'line',
        smooth: true,
        data: props.data.points.map((p) => p.actual),
        itemStyle: { color: '#e8955f' }
      },
      {
        name: 'Оптимистичный',
        type: 'line',
        smooth: true,
        data: props.data.points.map((p) => p.optimistic),
        lineStyle: { type: 'dashed' },
        itemStyle: { color: '#f5c4a0' }
      }
    ]
  }
})
</script>

<template>
  <ClientOnly>
    <VChart
      v-if="option"
      class="h-full min-h-[220px] max-h-[420px] w-full sm:min-h-[280px] lg:min-h-[300px]"
      :option="option"
      autoresize
    />
    <p v-else class="flex min-h-[220px] items-center justify-center text-sm text-slate-500 sm:min-h-[280px] lg:min-h-[300px]">
      Нет данных
    </p>
  </ClientOnly>
</template>
