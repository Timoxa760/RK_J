<script setup lang="ts">
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { CompareResponse } from '~/types/api'
import { chartThemeLight } from '~/types/api'

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent])

const props = defineProps<{
  data: CompareResponse | null
}>()

const colors = ['#e8955f', '#f0a66b', '#f5c4a0', '#d4824a', '#f5dcc8']

const option = computed(() => {
  if (!props.data?.months.length) return null
  const [prev, curr] = props.data.months
  if (!prev || !curr) return null

  return {
    backgroundColor: chartThemeLight.backgroundColor,
    tooltip: { trigger: 'item' },
    legend: {
      bottom: 0,
      textStyle: { color: chartThemeLight.textStyle.color, fontSize: 11 }
    },
    series: [
      {
        name: prev.month,
        type: 'pie',
        radius: ['36%', '52%'],
        center: ['28%', '45%'],
        data: prev.categories.map((c, i) => ({
          name: c.name,
          value: c.amount,
          itemStyle: { color: colors[i % colors.length] }
        })),
        label: { fontSize: 10, color: chartThemeLight.textStyle.color }
      },
      {
        name: curr.month,
        type: 'pie',
        radius: ['36%', '52%'],
        center: ['72%', '45%'],
        data: curr.categories.map((c, i) => ({
          name: c.name,
          value: c.amount,
          itemStyle: { color: colors[i % colors.length] }
        })),
        label: { fontSize: 10, color: chartThemeLight.textStyle.color }
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
