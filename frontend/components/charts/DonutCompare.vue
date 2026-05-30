<script setup lang="ts">
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { CompareResponse } from '~/types/api'
import { chartThemeLight } from '~/utils/chartTheme'

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent])

const props = defineProps<{
  data: CompareResponse | null
  size?: 'sm' | 'md' | 'lg' | 'full'
}>()

const { containerRef, isCompact } = useChartViewport()
const colors = chartThemeLight.colors

const option = computed(() => {
  if (!props.data?.months.length) return null
  const [prev, curr] = props.data.months
  if (!prev || !curr) return null

  const prevLabel = prev.label ?? prev.month ?? 'Пред.'
  const currLabel = curr.label ?? curr.month ?? 'Тек.'

  if (isCompact.value) {
    return {
      backgroundColor: chartThemeLight.backgroundColor,
      tooltip: { trigger: 'item' },
      legend: {
        type: 'scroll',
        bottom: 0,
        left: 'center',
        textStyle: { color: chartThemeLight.textStyle.color, fontSize: 10 }
      },
      graphic: [
        {
          type: 'text',
          left: 'center',
          top: '4%',
          style: { text: prevLabel, fill: chartThemeLight.textStyle.color, fontSize: 11, fontWeight: 600 }
        },
        {
          type: 'text',
          left: 'center',
          top: '52%',
          style: { text: currLabel, fill: chartThemeLight.textStyle.color, fontSize: 11, fontWeight: 600 }
        }
      ],
      series: [
        {
          name: prevLabel,
          type: 'pie',
          radius: ['28%', '42%'],
          center: ['50%', '28%'],
          data: prev.categories.map((c, i) => ({
            name: c.name,
            value: c.amount ?? c.total ?? 0,
            itemStyle: { color: colors[i % colors.length] }
          })),
          label: { show: false }
        },
        {
          name: currLabel,
          type: 'pie',
          radius: ['28%', '42%'],
          center: ['50%', '76%'],
          data: curr.categories.map((c, i) => ({
            name: c.name,
            value: c.amount ?? c.total ?? 0,
            itemStyle: { color: colors[i % colors.length] }
          })),
          label: { show: false }
        }
      ]
    }
  }

  return {
    backgroundColor: chartThemeLight.backgroundColor,
    tooltip: { trigger: 'item' },
    legend: {
      type: 'scroll',
      bottom: 8,
      left: 'center',
      textStyle: { color: chartThemeLight.textStyle.color, fontSize: 11 }
    },
    graphic: [
      {
        type: 'text',
        left: '22%',
        top: '8%',
        style: { text: prevLabel, fill: chartThemeLight.textStyle.color, fontSize: 12, fontWeight: 600, textAlign: 'center' }
      },
      {
        type: 'text',
        left: '72%',
        top: '8%',
        style: { text: currLabel, fill: chartThemeLight.textStyle.color, fontSize: 12, fontWeight: 600, textAlign: 'center' }
      }
    ],
    series: [
      {
        name: prevLabel,
        type: 'pie',
        radius: ['32%', '48%'],
        center: ['28%', '48%'],
        data: prev.categories.map((c, i) => ({
          name: c.name,
          value: c.amount ?? c.total ?? 0,
          itemStyle: { color: colors[i % colors.length] }
        })),
        label: { show: false }
      },
      {
        name: currLabel,
        type: 'pie',
        radius: ['32%', '48%'],
        center: ['72%', '48%'],
        data: curr.categories.map((c, i) => ({
          name: c.name,
          value: c.amount ?? c.total ?? 0,
          itemStyle: { color: colors[i % colors.length] }
        })),
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
