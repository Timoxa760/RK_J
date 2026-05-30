<script setup lang="ts">
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { CategoryDetail, CategoriesResponse } from '~/types/api'
import { chartThemeLight } from '~/utils/chartTheme'

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent])

const props = defineProps<{
  data: CategoriesResponse | null
  size?: 'sm' | 'md' | 'lg' | 'full'
}>()

const { containerRef, isCompact, isMedium } = useChartViewport()
const selected = ref<CategoryDetail | null>(null)
const detailOpen = computed({
  get: () => selected.value != null,
  set: (open: boolean) => {
    if (!open) selected.value = null
  }
})

const colors = chartThemeLight.colors

const option = computed(() => {
  if (!props.data?.categories.length) return null
  const compact = isCompact.value || isMedium.value

  return {
    backgroundColor: chartThemeLight.backgroundColor,
    tooltip: { trigger: 'item' },
    legend: compact
      ? {
          type: 'scroll',
          orient: 'horizontal',
          bottom: 4,
          left: 'center',
          width: '92%',
          itemGap: 8,
          textStyle: { ...chartThemeLight.textStyle, fontSize: 10 }
        }
      : {
          orient: 'vertical',
          right: 0,
          top: 'middle',
          textStyle: { ...chartThemeLight.textStyle, fontSize: 11 }
        },
    series: [
      {
        type: 'pie',
        radius: compact ? ['30%', '46%'] : ['42%', '68%'],
        center: compact ? ['50%', '38%'] : ['38%', '50%'],
        data: props.data.categories.map((c, i) => ({
          name: c.name,
          value: c.amount ?? c.total ?? 0,
          itemStyle: { color: colors[i % colors.length] }
        })),
        label: { show: false },
        emphasis: {
          itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.15)' }
        }
      }
    ]
  }
})

function onChartClick(params: { name?: string }) {
  if (!params.name || !props.data) return
  const cat = props.data.categories.find((c) => c.name === params.name)
  if (cat) selected.value = cat
}

</script>

<template>
  <ClientOnly>
    <div ref="containerRef" class="h-full w-full min-h-0">
      <VChart
        v-if="option"
        class="h-full w-full"
        :option="option"
        autoresize
        @click="onChartClick"
      />
      <p
        v-else
        class="flex h-full w-full items-center justify-center text-sm text-[color:var(--mm-text-soft)]"
      >
        Нет данных
      </p>
    </div>
  </ClientOnly>

  <Dialog v-model:open="detailOpen">
    <DialogContent v-if="selected" class="max-h-[80vh] overflow-y-auto sm:max-w-md">
      <DialogHeader>
        <DialogTitle>{{ selected.name }}</DialogTitle>
        <DialogDescription>
          {{ (selected.amount ?? selected.total ?? 0).toLocaleString('ru-RU') }} ₽
          <template v-if="selected.share"> · {{ Math.round(selected.share * 100) }}%</template>
        </DialogDescription>
      </DialogHeader>
      <div v-for="sub in selected.subcategories" :key="sub.name" class="mt-2">
        <p class="text-sm font-medium">{{ sub.name }}</p>
        <ul class="mt-2 space-y-1">
          <li
            v-for="item in sub.items"
            :key="item.name"
            class="flex justify-between text-sm text-muted-foreground"
          >
            <span>{{ item.name }}</span>
            <span>{{ (item.amount ?? (item.price ?? 0) * (item.quantity ?? 1)).toLocaleString('ru-RU') }} ₽</span>
          </li>
        </ul>
      </div>
      <p v-if="!selected.subcategories.length" class="text-sm text-muted-foreground">
        Детализация недоступна
      </p>
    </DialogContent>
  </Dialog>
</template>
