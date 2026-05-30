<script setup lang="ts">
import type { InsightItem } from '~/types/api'

const props = defineProps<{
  insights: InsightItem[]
}>()

function bodyText(item: InsightItem) {
  return item.body ?? item.description ?? ''
}

const variantMap: Record<string, 'default' | 'destructive'> = {
  info: 'default',
  warning: 'default',
  success: 'default',
  critical: 'destructive'
}
</script>

<template>
  <ul class="space-y-3" data-demo="insights">
    <li v-for="(item, index) in props.insights" :key="item.id ?? index">
      <Alert :variant="variantMap[item.severity] ?? 'default'">
        <AlertTitle>{{ item.title }}</AlertTitle>
        <AlertDescription>
          <p>{{ bodyText(item) }}</p>
          <p v-if="item.amount" class="mt-1 text-xs opacity-80">
            {{ item.amount.toLocaleString('ru-RU') }} ₽
          </p>
        </AlertDescription>
      </Alert>
    </li>
  </ul>
</template>
