<script setup lang="ts">
import type { HabitIndex } from '~/utils/habitIndex'

defineProps<{
  habitIndex: HabitIndex
  loading?: boolean
}>()

const toneClass: Record<HabitIndex['tone'], string> = {
  good: 'text-emerald-600',
  warn: 'text-amber-600',
  risk: 'text-red-600'
}

const barClass: Record<HabitIndex['tone'], string> = {
  good: 'bg-emerald-500',
  warn: 'bg-amber-500',
  risk: 'bg-red-500'
}
</script>

<template>
  <Card data-demo="habit-index">
    <CardHeader>
      <CardTitle class="text-base">Как у вас с тратами</CardTitle>
      <CardDescription>Без сравнения сумм — только ваш ритм</CardDescription>
    </CardHeader>
    <CardContent>
      <Skeleton v-if="loading" class="h-20 w-full" />

      <div v-else class="space-y-4">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div class="space-y-1">
            <p class="text-3xl font-semibold tabular-nums" :class="toneClass[habitIndex.tone]">
              {{ habitIndex.score > 0 ? habitIndex.score : '—' }}
              <span v-if="habitIndex.score > 0" class="text-lg font-normal text-muted-foreground">
                / 100
              </span>
            </p>
            <p class="text-sm font-medium">{{ habitIndex.label }}</p>
          </div>

          <div class="w-full sm:max-w-xs">
            <div class="h-2 overflow-hidden rounded-full bg-muted">
              <div
                class="h-full rounded-full transition-all"
                :class="barClass[habitIndex.tone]"
                :style="{ width: `${habitIndex.score}%` }"
              />
            </div>
          </div>
        </div>

        <p v-if="habitIndex.insight" class="text-sm text-muted-foreground">
          {{ habitIndex.insight }}
        </p>
        <p v-else class="text-sm text-muted-foreground">
          {{ habitIndex.challengeHint }}
        </p>

        <p class="text-xs text-muted-foreground">
          Подробнее в
          <NuxtLink to="/receipts" class="text-primary underline">ленте расходов</NuxtLink>
          и
          <NuxtLink to="/analytics" class="text-primary underline">аналитике</NuxtLink>.
        </p>
      </div>
    </CardContent>
  </Card>
</template>
