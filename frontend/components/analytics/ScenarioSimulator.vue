<script setup lang="ts">
import type { TimeMachineResponse } from '~/types/api'
import { SCENARIO_OPTIONS } from '~/types/api'

defineProps<{
  result: string | null
  simulation: TimeMachineResponse | null
  loading?: boolean
}>()

const scenario = defineModel<'reduce_delivery' | 'reduce_cafe' | 'reduce_entertainment' | 'custom'>(
  'scenario',
  { default: 'reduce_cafe' }
)
const percent = defineModel<number>('percent', { default: 20 })

const emit = defineEmits<{
  simulate: []
}>()

const scenarioOptions = SCENARIO_OPTIONS
</script>

<template>
  <Card data-demo="scenario-simulator">
    <CardHeader>
      <CardTitle class="text-base">А если меньше тратить на…</CardTitle>
      <CardDescription>
        Сократите категорию — увидите, как сдвигается траектория накоплений
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-end">
        <div class="min-w-0 flex-1 space-y-2">
          <Label for="scenario-select">Категория</Label>
          <Select v-model="scenario">
            <SelectTrigger id="scenario-select" class="w-full sm:w-[220px]">
              <SelectValue placeholder="Сценарий" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="opt in scenarioOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="space-y-2">
          <Label for="scenario-percent">Сокращение</Label>
          <div class="flex items-center gap-2">
            <Input
              id="scenario-percent"
              v-model.number="percent"
              type="number"
              min="5"
              max="50"
              class="w-full sm:w-24"
            />
            <span class="shrink-0 text-sm text-muted-foreground">%</span>
          </div>
        </div>
        <Button :disabled="loading" @click="emit('simulate')">
          {{ loading ? 'Считаем…' : 'Симулировать' }}
        </Button>
      </div>

      <Alert v-if="result">
        <AlertTitle>Последствия сценария</AlertTitle>
        <AlertDescription>{{ result }}</AlertDescription>
      </Alert>

      <div v-if="simulation" class="mm-chart-wrap mm-chart-wrap--md">
        <ChartsTimeMachineChart :data="simulation" />
      </div>
    </CardContent>
  </Card>
</template>
