<script setup lang="ts">
import type { InsightItem } from '~/types/api'

defineProps<{
  insights: InsightItem[]
}>()

const severityClass: Record<InsightItem['severity'], string> = {
  info: 'border-[color:var(--mm-primary-muted)] bg-[color:var(--mm-primary-soft)] text-[color:var(--mm-primary-hover)]',
  warning: 'border-amber-200/80 bg-amber-50/80 text-amber-900',
  success: 'border-[color:var(--mm-primary-muted)] bg-[color:var(--mm-primary-soft)] text-[color:var(--mm-primary-hover)]'
}
</script>

<template>
  <ul class="space-y-3" data-demo="insights">
    <li
      v-for="item in insights"
      :key="item.id"
      class="rounded-lg border p-4"
      :class="severityClass[item.severity]"
    >
      <p class="font-medium">{{ item.title }}</p>
      <p class="mt-1 text-sm opacity-90">{{ item.body }}</p>
    </li>
  </ul>
</template>
