<script setup lang="ts">
import type { AiDiagnosisResponse } from '~/types/api'
import { ADVISOR } from '~/constants/productCopy'

const props = defineProps<{
  diagnosis: AiDiagnosisResponse | null
  loading?: boolean
}>()

const summary = computed(() => {
  const action = props.diagnosis?.main_action
  if (!action?.description) return null
  return action.description
})
</script>

<template>
  <Card class="mm-mindfulness-score h-full" data-demo="mindfulness-score">
    <CardHeader class="space-y-1 pb-2">
      <CardTitle class="text-base">{{ ADVISOR.mindfulnessTitle }}</CardTitle>
      <CardDescription v-if="summary" class="line-clamp-3 text-sm leading-relaxed">
        {{ summary }}
      </CardDescription>
    </CardHeader>
    <CardContent class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <Skeleton v-if="loading" class="h-10 w-24 sm:ml-auto" />
      <div v-else-if="diagnosis" class="sm:ml-auto sm:text-right">
        <p class="text-3xl font-bold tabular-nums leading-none text-primary sm:text-4xl">
          {{ diagnosis.score }}/100
        </p>
        <p v-if="diagnosis.grade" class="mt-1 text-sm text-muted-foreground">
          Оценка {{ diagnosis.grade }}
        </p>
      </div>
      <p v-else class="text-sm leading-relaxed text-muted-foreground">
        Обновите план — посчитаем оценку по вашим тратам и цели.
      </p>
    </CardContent>
  </Card>
</template>
