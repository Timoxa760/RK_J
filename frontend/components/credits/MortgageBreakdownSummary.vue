<script setup lang="ts">
import type { MortgageApprovalLevel, MortgageBreakdownResponse } from '~/types/api'

defineProps<{
  breakdown: MortgageBreakdownResponse
}>()

const approvalMeta: Record<
  MortgageApprovalLevel,
  { label: string; emoji: string; variant: 'default' | 'secondary' | 'destructive' }
> = {
  high: { label: 'Высокий шанс одобрения', emoji: '🟢', variant: 'default' },
  medium: { label: 'Средний шанс — нужны уточнения', emoji: '🟡', variant: 'secondary' },
  low: { label: 'Низкий шанс — платежи съедят много дохода', emoji: '🔴', variant: 'destructive' }
}
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <Card class="sm:col-span-2 lg:col-span-1">
      <CardHeader class="pb-2">
        <CardDescription>Шанс одобрения</CardDescription>
        <CardTitle class="flex items-center gap-2 text-base">
          <span aria-hidden="true">{{ approvalMeta[breakdown.approval_level].emoji }}</span>
          {{ approvalMeta[breakdown.approval_level].label }}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-sm leading-relaxed text-muted-foreground">{{ breakdown.approval_reason }}</p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="pb-2">
        <CardDescription>Безопасный предел</CardDescription>
        <CardTitle class="text-lg">
          {{ breakdown.safe_mortgage_amount.toLocaleString('ru-RU') }} ₽
        </CardTitle>
      </CardHeader>
      <CardContent class="text-sm text-muted-foreground">
        Комфортный платёж ≈ {{ breakdown.comfortable_payment.toLocaleString('ru-RU') }} ₽/мес
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="pb-2">
        <CardDescription>После ипотеки</CardDescription>
        <CardTitle class="text-base font-medium leading-snug">Как изменится запас</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-sm leading-relaxed text-muted-foreground">{{ breakdown.load_risk }}</p>
      </CardContent>
    </Card>
  </div>
</template>
