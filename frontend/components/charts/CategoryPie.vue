<script setup lang="ts">
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { LegendComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { CategoryDetail, CategoriesResponse } from '~/types/api'
import { chartThemeLight } from '~/types/api'

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent])

const props = defineProps<{
  data: CategoriesResponse | null
}>()

const chartTheme = chartThemeLight
const selected = ref<CategoryDetail | null>(null)

const colors = ['#e8955f', '#f0a66b', '#f5c4a0', '#d4824a', '#f5dcc8', '#c9773f']

const option = computed(() => {
  if (!props.data?.categories.length) return null
  return {
    backgroundColor: chartTheme.backgroundColor,
    tooltip: { trigger: 'item' },
    legend: {
      orient: 'vertical',
      right: 8,
      top: 'center',
      textStyle: { color: chartTheme.textStyle.color, fontSize: 11 }
    },
    series: [
      {
        type: 'pie',
        radius: ['42%', '68%'],
        center: ['40%', '50%'],
        data: props.data.categories.map((c, i) => ({
          name: c.name,
          value: c.amount,
          itemStyle: { color: colors[i % colors.length] }
        })),
        label: { color: chartTheme.textStyle.color },
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

function closeModal() {
  selected.value = null
}
</script>

<template>
  <ClientOnly>
    <VChart
      v-if="option"
      class="h-full min-h-[220px] w-full sm:min-h-[280px] lg:min-h-[300px]"
      :option="option"
      autoresize
      @click="onChartClick"
    />
    <p v-else class="flex min-h-[220px] items-center justify-center text-sm text-slate-500 sm:min-h-[280px] lg:min-h-[300px]">
      Нет данных
    </p>
  </ClientOnly>

  <Teleport to="body">
    <div
      v-if="selected"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4"
      @click.self="closeModal"
    >
      <div class="max-h-[80vh] w-full max-w-md overflow-y-auto rounded-xl border border-slate-200 bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h3 class="text-lg font-semibold text-slate-900">{{ selected.name }}</h3>
            <p class="text-sm text-slate-500">
              {{ selected.amount.toLocaleString('ru-RU') }} ₽ · {{ Math.round(selected.share * 100) }}%
            </p>
          </div>
          <button
            type="button"
            class="text-slate-400 hover:text-slate-600"
            aria-label="Закрыть"
            @click="closeModal"
          >
            ✕
          </button>
        </div>
        <div v-for="sub in selected.subcategories" :key="sub.name" class="mt-4">
          <p class="text-sm font-medium text-slate-700">{{ sub.name }}</p>
          <ul class="mt-2 space-y-1">
            <li
              v-for="item in sub.items"
              :key="item.name"
              class="flex justify-between text-sm text-slate-600"
            >
              <span>{{ item.name }}</span>
              <span>{{ item.amount.toLocaleString('ru-RU') }} ₽</span>
            </li>
          </ul>
        </div>
        <p v-if="!selected.subcategories.length" class="mt-4 text-sm text-slate-500">
          Детализация недоступна
        </p>
      </div>
    </div>
  </Teleport>
</template>
