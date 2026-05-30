<script setup lang="ts">
import type { CategoriesResponse } from '~/types/api'
import { chartThemeLight } from '~/utils/chartTheme'
import { formatRub } from '~/constants/productCopy'

const props = defineProps<{
  data: CategoriesResponse | null
}>()

const rows = computed(() => {
  if (!props.data?.categories.length) return []

  const items = props.data.categories.map((c) => ({
    name: c.name,
    amount: c.amount ?? c.total ?? 0
  }))

  const total = items.reduce((sum, c) => sum + c.amount, 0) || 1
  const max = Math.max(...items.map((c) => c.amount), 1)

  return items
    .sort((a, b) => b.amount - a.amount)
    .map((item, i) => ({
      ...item,
      share: Math.round((item.amount / total) * 100),
      width: Math.max(12, Math.round((item.amount / max) * 100)),
      color: chartThemeLight.colors[i % chartThemeLight.colors.length]
    }))
})
</script>

<template>
  <ul v-if="rows.length" class="mm-simple-chart space-y-4" role="list">
    <li v-for="row in rows" :key="row.name" class="mm-simple-chart__row">
      <div class="mm-simple-chart__label-row">
        <span class="mm-simple-chart__label">{{ row.name }}</span>
        <span class="mm-simple-chart__value">
          {{ formatRub(row.amount) }}
          <span class="mm-simple-chart__share">· {{ row.share }}%</span>
        </span>
      </div>
      <div class="mm-simple-chart__track" aria-hidden="true">
        <div
          class="mm-simple-chart__bar"
          :style="{ width: `${row.width}%`, backgroundColor: row.color }"
        />
      </div>
    </li>
  </ul>
  <p v-else class="mm-simple-chart__empty">Категории появятся после покупок</p>
</template>
