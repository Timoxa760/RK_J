<script setup lang="ts">
import { use } from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { ForecastResponse } from '~/types/api'
import { baseGrid, brandColors, chartAxisLabel, chartFontFamily, chartThemeLight } from '~/utils/chartTheme'
import { normalizeForecast } from '~/utils/apiNormalize'

use([CanvasRenderer, BarChart, LineChart, TooltipComponent, LegendComponent, GridComponent])

const props = defineProps<{
  data: ForecastResponse | null
}>()

const { containerRef, isCompact } = useChartViewport()

const normalized = computed(() => (props.data ? normalizeForecast(props.data) : null))

const option = computed(() => {
  const data = normalized.value
  if (!data?.dates.length) return null

  const hasBounds = Boolean(data.upper_bound?.length && data.lower_bound?.length)
  const grid = baseGrid(isCompact.value)
  grid.top = 12
  grid.bottom = hasBounds ? 52 : 36

  const series: Record<string, unknown>[] = [
    {
      name: 'Прогноз',
      type: 'bar',
      data: data.forecast,
      itemStyle: { color: brandColors.primary, borderRadius: [4, 4, 0, 0] },
      z: 2
    }
  ]

  if (hasBounds) {
    series.push(
      {
        name: 'Обычно до',
        type: 'line',
        data: data.upper_bound,
        smooth: true,
        lineStyle: { type: 'dashed', color: brandColors.primaryLight, width: 1.5 },
        symbol: 'none',
        z: 1
      },
      {
        name: 'Обычно от',
        type: 'line',
        data: data.lower_bound,
        smooth: true,
        lineStyle: { type: 'dashed', color: brandColors.primaryHover, width: 1.5 },
        symbol: 'none',
        z: 1
      }
    )
  }

  return {
    backgroundColor: chartThemeLight.backgroundColor,
    tooltip: { trigger: 'axis' },
    legend: hasBounds
      ? {
          bottom: 0,
          left: 'center',
          itemGap: 20,
          itemWidth: 14,
          textStyle: { ...chartThemeLight.textStyle, fontSize: 11 }
        }
      : { show: false },
    grid,
    xAxis: {
      type: 'category',
      data: data.dates,
      axisLabel: chartAxisLabel(10)
    },
    yAxis: {
      type: 'value',
      axisLabel: chartAxisLabel(10),
      splitLine: { lineStyle: { color: chartThemeLight.splitLine } }
    },
    series
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
