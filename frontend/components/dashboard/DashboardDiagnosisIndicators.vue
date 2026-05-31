<script setup lang="ts">
import type { AiDiagnosisIndicator, AiDiagnosisIndicatorStatus } from '~/types/api'
import { ADVISOR } from '~/constants/productCopy'

defineProps<{
  indicators: AiDiagnosisIndicator[]
  loading?: boolean
}>()

const statusLabel: Record<AiDiagnosisIndicatorStatus, string> = {
  good: 'Хорошо',
  warning: 'Стоит улучшить',
  critical: ADVISOR.diagnosisStatusUrgent
}

const statusVariant: Record<
  AiDiagnosisIndicatorStatus,
  'default' | 'secondary' | 'destructive'
> = {
  good: 'default',
  warning: 'secondary',
  critical: 'destructive'
}

function formatValue(name: string, value: number): string {
  if (name.toLowerCase().includes('подушка')) {
    return `${value.toLocaleString('ru-RU', { maximumFractionDigits: 1 })} мес.`
  }
  if (name.toLowerCase().includes('нагрузка') || name.toLowerCase().includes('доход')) {
    return `${value}%`
  }
  return value.toLocaleString('ru-RU')
}

function cardClass(status: AiDiagnosisIndicatorStatus): string {
  if (status === 'critical') return 'border-amber-500/60'
  return ''
}
</script>

<template>
  <section aria-label="Показатели картины">
    <div v-if="loading" class="mm-diagnosis-indicators-grid">
      <Skeleton v-for="i in 5" :key="i" class="h-16 w-full" />
    </div>

    <div v-else class="mm-diagnosis-indicators-grid">
      <Card
        v-for="item in indicators"
        :key="item.name"
        class="flex flex-col"
        :class="cardClass(item.status)"
      >
        <CardHeader class="space-y-1 p-3 pb-1">
          <div class="flex items-start justify-between gap-2">
            <CardDescription class="text-xs leading-snug">{{ item.name }}</CardDescription>
            <Badge :variant="statusVariant[item.status]" class="shrink-0 text-[10px]">
              {{ statusLabel[item.status] }}
            </Badge>
          </div>
          <CardTitle class="text-base tabular-nums">
            {{ formatValue(item.name, item.value) }}
            <span class="text-xs font-normal text-muted-foreground">норма {{ item.norm }}</span>
          </CardTitle>
        </CardHeader>
      </Card>
    </div>
  </section>
</template>
