<script setup lang="ts">
import type { Goal } from '~/types/api'
import { buildGoalProgressText } from '~/utils/pageNarrative'

const props = defineProps<{
  goal: Goal | null
  monthlySaving?: number
}>()

const progressText = computed(() => buildGoalProgressText(props.goal, props.monthlySaving))

const progressPercent = computed(() => {
  if (!props.goal?.target_amount) return 0
  return Math.min(
    100,
    props.goal.progress_percent ??
      Math.round((props.goal.current_amount / props.goal.target_amount) * 100)
  )
})
</script>

<template>
  <Card data-demo="digest-goal">
    <CardHeader class="pb-2">
      <CardDescription>Прогресс к цели</CardDescription>
      <CardTitle class="text-base font-medium leading-snug">
        {{ goal?.title ?? 'Цель не задана' }}
      </CardTitle>
    </CardHeader>
    <CardContent class="space-y-3">
      <template v-if="goal">
        <div class="h-2 overflow-hidden rounded-full bg-muted">
          <div
            class="h-full rounded-full bg-primary transition-all"
            :style="{ width: `${progressPercent}%` }"
          />
        </div>
        <p class="text-sm text-muted-foreground">{{ progressText }}</p>
        <AdvisorAskButton
          v-if="progressPercent < 100"
          :goal="goal"
        />
      </template>
      <template v-else>
        <p class="text-sm text-muted-foreground">
          Поставьте цель в профиле — покажем срок и прогресс накоплений.
        </p>
        <Button variant="outline" size="sm" as-child>
          <NuxtLink to="/profile">Добавить цель</NuxtLink>
        </Button>
      </template>
    </CardContent>
  </Card>
</template>
