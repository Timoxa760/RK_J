<script setup lang="ts">
import { use } from 'echarts/core'
import { SankeyChart as SankeyChartType } from 'echarts/charts'
import { TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { SankeyResponse } from '~/types/api'
import { chartThemeLight } from '~/types/api'

use([CanvasRenderer, SankeyChartType, TooltipComponent])

defineProps<{
  data: SankeyResponse | null
}>()

const chartTheme = chartThemeLight

const option = computed(() => {
  if (!props.data) return null
  const nodes = props.data.nodes.map((n) => ({
    name: n.name,
    itemStyle: {
      color:
        n.category === 'income'
          ? '#e8955f'
          : n.category === 'savings'
            ? '#f0a66b'
            : '#f5c4a0'
    }
  }))
  return {
    backgroundColor: chartTheme.backgroundColor,
    tooltip: { trigger: 'item', triggerOn: 'mousemove' },
    series: [
      {
        type: 'sankey',
        layout: 'none',
        emphasis: { focus: 'adjacency' },
        data: nodes,
        links: props.data.links,
        lineStyle: { color: 'gradient', curveness: 0.5 },
        label: { color: chartTheme.textStyle.color, fontSize: 11 }
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
