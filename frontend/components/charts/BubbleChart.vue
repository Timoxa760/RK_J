<script setup lang="ts">
import { use } from 'echarts/core'
import { ScatterChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { StoresResponse } from '~/types/api'
import { baseGrid, brandColors, chartAxisLabel, chartThemeLight } from '~/utils/chartTheme'

use([CanvasRenderer, ScatterChart, TooltipComponent, GridComponent])

const props = defineProps<{
  data: StoresResponse | null
  size?: 'sm' | 'md' | 'lg' | 'full'
}>()

const { containerRef, isCompact } = useChartViewport()

function impulseColor(ratio: number) {
  if (ratio < 0.33) return brandColors.primary
  if (ratio < 0.66) return brandColors.primaryHover
  return brandColors.primaryDeep
}

const option = computed(() => {
  if (!props.data?.stores.length) return null
  const seriesData = props.data.stores.map((s) => ({
    name: s.name,
    value: [s.avg_check, s.visits ?? s.purchases ?? 0, s.total, s.impulse_ratio],
    itemStyle: { color: impulseColor(s.impulse_ratio) }
  }))
  const maxTotal = Math.max(...props.data.stores.map((s) => s.total))
  const grid = baseGrid(isCompact.value)
  grid.top = isCompact.value ? 16 : 24
  grid.bottom = isCompact.value ? 56 : 56
  grid.left = isCompact.value ? 36 : 48

  return {
    backgroundColor: chartThemeLight.backgroundColor,
    tooltip: {
      formatter: (p: { data: { name: string; value: number[] } }) => {
        const [avg, visits, total, impulse] = p.data.value
        return `${p.data.name}<br/>Средний чек: ${avg} ₽<br/>Покупок: ${visits}<br/>Всего: ${total} ₽<br/>Покупок «на эмоциях»: ${Math.round((impulse ?? 0) * 100)}%`
      }
    },
    grid,
    xAxis: {
      name: isCompact.value ? '' : 'Средний чек',
      nameLocation: 'middle',
      nameGap: 28,
      axisLine: { lineStyle: { color: chartThemeLight.axisColor } },
      axisLabel: chartAxisLabel(10)
    },
    yAxis: {
      name: isCompact.value ? '' : 'Покупки',
      nameGap: 12,
      axisLine: { lineStyle: { color: chartThemeLight.axisColor } },
      axisLabel: chartAxisLabel(10),
      splitLine: { lineStyle: { color: chartThemeLight.splitLine } }
    },
    series: [
      {
        type: 'scatter',
        symbolSize: (val: number[]) => Math.max(12, (val[2]! / maxTotal) * 60),
        data: seriesData,
        label: { show: false }
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
