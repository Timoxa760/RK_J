<script setup lang="ts">
import { ANALYTICS } from '~/constants/productCopy'
import type { InsightItem } from '~/types/api'

const props = defineProps<{
  insights: InsightItem[]
}>()

const expanded = ref(false)

const primary = computed(() => props.insights[0] ?? null)
const rest = computed(() => props.insights.slice(1))

function bodyText(item: InsightItem) {
  return item.body ?? item.description ?? ''
}
</script>

<template>
  <Card data-demo="insights">
    <CardHeader>
      <CardTitle class="text-base">На что обратить внимание</CardTitle>
      <CardDescription>Причины, не цифры — что влияет на цель</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <template v-if="primary">
        <div class="rounded-lg border border-primary/20 bg-primary/5 p-4">
          <p class="text-xs font-medium text-primary">Главный вывод</p>
          <p class="mt-1 font-medium">{{ primary.title }}</p>
          <p v-if="bodyText(primary)" class="mt-2 text-sm text-muted-foreground">
            {{ bodyText(primary) }}
          </p>
        </div>

        <div v-if="rest.length">
          <Button variant="ghost" size="sm" class="h-auto px-0" @click="expanded = !expanded">
            {{ expanded ? 'Скрыть' : ANALYTICS.moreTips(rest.length) }}
          </Button>
          <InsightsPanel v-if="expanded" class="mt-3" :insights="rest" />
        </div>
      </template>

      <p v-else class="text-sm text-muted-foreground">
        Добавьте расходы — появятся паттерны и мягкие подсказки.
      </p>
    </CardContent>
  </Card>
</template>
