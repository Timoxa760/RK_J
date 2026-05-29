<script setup lang="ts">
import { use } from 'echarts/core'
import { ScatterChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, VisualMapComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { StoresResponse } from '~/types/api'
import { chartThemeLight } from '~/types/api'

use([CanvasRenderer, ScatterChart, TooltipComponent, GridComponent, VisualMapComponent])

const props = defineProps<{
  data: StoresResponse | null
}>()

const chartTheme = chartThemeLight

function impulseColor(ratio: number) {
  if (ratio < 0.33) return '#e8955f'
  if (ratio < 0.66) return '#c4a574'
  return '#c4847a'
}

const option = computed(() => {
  if (!props.data?.stores.length) return null
  const seriesData = props.data.stores.map((s) => ({
    name: s.name,
    value: [s.avg_check, s.visits, s.total, s.impulse_ratio],
    itemStyle: { color: impulseColor(s.impulse_ratio) }
  }))
  const maxTotal = Math.max(...props.data.stores.map((s) => s.total))

  return {
    backgroundColor: chartTheme.backgroundColor,
    tooltip: {
      formatter: (p: { data: { name: string; value: number[] } }) => {
        const [avg, visits, total, impulse] = p.data.value
        return `${p.data.name}<br/>Средний чек: ${avg} ₽<br/>Покупок: ${visits}<br/>Всего: ${total} ₽<br/>Импульсивность: ${Math.round(impulse * 100)}%`
      }
    },
    grid: { left: 48, right: 24, top: 24, bottom: 48 },
    xAxis: {
      name: 'Средний чек',
      nameLocation: 'middle',
      nameGap: 28,
      axisLine: { lineStyle: { color: chartTheme.axisColor } },
      axisLabel: { color: chartTheme.textStyle.color }
    },
    yAxis: {
      name: 'Покупки',
      axisLine: { lineStyle: { color: chartTheme.axisColor } },
      axisLabel: { color: chartTheme.textStyle.color },
      splitLine: { lineStyle: { color: chartTheme.splitLine } }
    },
    series: [
      {
        type: 'scatter',
        symbolSize: (val: number[]) => Math.max(12, (val[2] / maxTotal) * 60),
        data: seriesData,
        label: {
          show: true,
          formatter: '{b}',
          position: 'top',
          fontSize: 10,
          color: chartTheme.textStyle.color
        }
      }
    ]
  }
})
</script>

<template>
  <ClientOnly>
    <VChart
      v-if="option"
      class="h-full min-h-[220px] w-full sm:min-h-[280px] lg:min-h-[300px]"
      :option="option"
      autoresize
    />
    <p v-else class="flex min-h-[220px] items-center justify-center text-sm text-slate-500 sm:min-h-[280px] lg:min-h-[300px]">
      Нет данных
    </p>
  </ClientOnly>
</template>
