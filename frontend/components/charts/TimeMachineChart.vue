<script setup lang="ts">
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { TimeMachineResponse } from '~/types/api'
import {
  baseGrid,
  brandColors,
  chartAxisLabel,
  chartThemeLight,
  formatAxisMoney,
  formatMonthLabel,
  sparseLabelInterval
} from '~/utils/chartTheme'

use([CanvasRenderer, LineChart, TooltipComponent, LegendComponent, GridComponent])

const props = defineProps<{
  data: TimeMachineResponse | null
  size?: 'sm' | 'md' | 'lg' | 'full'
}>()

const { containerRef, isCompact } = useChartViewport()

const option = computed(() => {
  if (!props.data?.points.length) return null
  const months = props.data.points.map((p) => p.month)
  const grid = baseGrid(isCompact.value)
  grid.bottom = isCompact.value ? 56 : 40
  grid.top = 48

  return {
    backgroundColor: chartThemeLight.backgroundColor,
    tooltip: { trigger: 'axis' },
    legend: {
      data: ['Факт', 'Оптимистичный'],
      top: 0,
      left: 'center',
      textStyle: { ...chartThemeLight.textStyle, fontSize: 11 }
    },
    grid,
    xAxis: {
      type: 'category',
      data: months,
      axisLabel: {
        ...chartAxisLabel(10),
        rotate: isCompact.value ? 0 : 45,
        interval: sparseLabelInterval(months.length, isCompact.value),
        formatter: (value: string) => formatMonthLabel(value)
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        ...chartAxisLabel(10),
        formatter: (v: number) => formatAxisMoney(v)
      },
      splitLine: { lineStyle: { color: chartThemeLight.splitLine } }
    },
    series: [
      {
        name: 'Факт',
        type: 'line',
        smooth: true,
        data: props.data.points.map((p) => p.actual),
        itemStyle: { color: brandColors.primary }
      },
      {
        name: 'Оптимистичный',
        type: 'line',
        smooth: true,
        data: props.data.points.map((p) => p.optimistic),
        lineStyle: { type: 'dashed' },
        itemStyle: { color: brandColors.primaryLight }
      }
    ]
  }
})

</script>

<template>
  <ClientOnly>
    <div ref="containerRef" class="h-full w-full min-h-0">
      <VChart v-if="option" class="h-full w-full" :option="option" autoresize />
      <p
        v-else
        class="flex h-full w-full items-center justify-center text-sm text-[color:var(--mm-text-soft)]"
      >
        Нет данных
      </p>
    </div>
  </ClientOnly>
</template>
