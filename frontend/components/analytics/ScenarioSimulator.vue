<script setup lang="ts">
import { SCENARIO_OPTIONS } from '~/types/api'
import { scenarioResultPrefix } from '~/utils/dashboardCopy'

const props = defineProps<{
  result: string | null
  loading?: boolean
  embedded?: boolean
}>()

const formattedResult = computed(() =>
  props.result ? scenarioResultPrefix(props.result) : null
)

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
  <component :is="embedded ? 'div' : 'Card'" data-demo="scenario-simulator" :class="embedded ? 'space-y-4' : undefined">
    <component :is="embedded ? 'div' : 'CardHeader'">
      <CardTitle class="text-base">А если меньше тратить на…</CardTitle>
      <CardDescription>
        Сократите категорию — увидите, как меняется прогноз накоплений
      </CardDescription>
    </component>
    <component :is="embedded ? 'div' : 'CardContent'" class="space-y-4">
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

      <Alert v-if="formattedResult">
        <AlertTitle>Последствия сценария</AlertTitle>
        <AlertDescription>{{ formattedResult }}</AlertDescription>
      </Alert>
    </component>
  </component>
</template>
