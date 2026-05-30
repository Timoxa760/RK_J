<script setup lang="ts">
import { use } from 'echarts/core'
import { SankeyChart as SankeyChartType } from 'echarts/charts'
import { TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { SankeyResponse } from '~/types/api'
import { chartThemeLight } from '~/utils/chartTheme'

use([CanvasRenderer, SankeyChartType, TooltipComponent])

const props = defineProps<{
  data: SankeyResponse | null
  size?: 'sm' | 'md' | 'lg' | 'full'
}>()

const { containerRef, isCompact } = useChartViewport()

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
    backgroundColor: chartThemeLight.backgroundColor,
    tooltip: { trigger: 'item', triggerOn: 'mousemove' },
    series: [
      {
        type: 'sankey',
        layout: 'none',
        left: isCompact.value ? 8 : '4%',
        right: isCompact.value ? 8 : '4%',
        top: isCompact.value ? 12 : 20,
        bottom: isCompact.value ? 12 : 20,
        nodeWidth: isCompact.value ? 12 : 20,
        nodeGap: isCompact.value ? 8 : 14,
        layoutIterations: 32,
        emphasis: { focus: 'adjacency' },
        data: nodes,
        links: props.data.links,
        lineStyle: { color: 'gradient', curveness: 0.5 },
        label: {
          color: chartThemeLight.textStyle.color,
          fontSize: isCompact.value ? 9 : 11,
          show: !isCompact.value,
          position: 'inside'
        }
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
