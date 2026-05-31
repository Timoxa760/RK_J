<script setup lang="ts">
import { Trash2 } from 'lucide-vue-next'
import type { CreditsDashboardResponse } from '~/types/api'
import { CREDITS, formatRub } from '~/constants/productCopy'

defineProps<{
  dashboard: CreditsDashboardResponse
}>()

const { deleteCredit, deleting } = useCredits()
</script>

<template>
  <Card data-demo="credits-list">
    <CardHeader class="pb-2">
      <CardTitle class="text-base">Ваши договоры</CardTitle>
      <CardDescription>
        {{ dashboard.credits.length }}
        {{
          dashboard.credits.length === 1
            ? 'кредит'
            : dashboard.credits.length < 5
              ? 'кредита'
              : 'кредитов'
        }}
        · платежи {{ formatRub(dashboard.monthly_payments ?? 0) }}/мес
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-3">
      <article
        v-for="credit in dashboard.credits"
        :key="credit.id"
        class="flex items-start justify-between gap-3 rounded-xl border bg-muted/20 p-4"
      >
        <div class="min-w-0 space-y-1">
          <p class="font-medium leading-snug">{{ credit.bank || credit.name || 'Кредит' }}</p>
          <p class="text-sm text-muted-foreground">
            {{ formatRub(credit.amount ?? credit.remaining ?? 0) }}
            <span v-if="credit.rate"> · {{ credit.rate }}%</span>
            <span v-if="credit.term_months"> · {{ credit.term_months }} мес.</span>
          </p>
          <p class="text-sm tabular-nums">
            {{ formatRub(credit.monthly_payment ?? credit.payment ?? 0) }}/мес
          </p>
        </div>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          class="shrink-0 text-muted-foreground hover:text-destructive"
          :disabled="deleting"
          :aria-label="CREDITS.deleteCredit"
          @click="deleteCredit(credit.id)"
        >
          <Trash2 class="size-4" />
        </Button>
      </article>
    </CardContent>
  </Card>
</template>
