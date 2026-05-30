<script setup lang="ts">
import { ANALYTICS } from '~/constants/productCopy'
import type { InsightItem } from '~/types/api'

const props = defineProps<{
  insights: InsightItem[]
  embedded?: boolean
}>()

const expanded = ref(false)

const primary = computed(() => props.insights[0] ?? null)
const rest = computed(() => props.insights.slice(1))

function bodyText(item: InsightItem) {
  return item.body ?? item.description ?? ''
}
</script>

<template>
  <component :is="embedded ? 'div' : 'Card'" data-demo="insights" :class="embedded ? 'space-y-4' : undefined">
    <component :is="embedded ? 'div' : 'CardHeader'">
      <component :is="embedded ? 'h3' : 'CardTitle'" class="text-base font-semibold">
        {{ embedded ? 'Почему совет именно такой' : ANALYTICS.attentionTitle }}
      </component>
      <CardDescription v-if="!embedded" class="text-base">
        Причины, не цифры — что влияет на цель
      </CardDescription>
    </component>
    <component :is="embedded ? 'div' : 'CardContent'" class="space-y-4">
      <template v-if="primary">
        <div class="rounded-lg border border-border bg-muted/30 p-4">
          <p class="text-sm font-medium text-muted-foreground">Главный вывод</p>
          <p class="mt-1 text-base font-semibold">{{ primary.title }}</p>
          <p v-if="bodyText(primary)" class="mt-2 text-sm leading-relaxed text-muted-foreground">
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
    </component>
  </component>
</template>
