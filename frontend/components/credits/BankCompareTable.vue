<script setup lang="ts">
import type { MortgageBreakdownResponse } from '~/types/api'

const props = defineProps<{
  breakdown: MortgageBreakdownResponse
}>()

const optimal = computed(() =>
  props.breakdown.banks.find((b) => b.id === props.breakdown.optimal_bank_id)
)
</script>

<template>
  <Card data-demo="bank-compare">
    <CardHeader>
      <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <CardTitle class="text-base">Сравнение банков</CardTitle>
          <CardDescription>Платёж, переплата и итоговая стоимость — без скрытых столбцов</CardDescription>
        </div>
        <Badge v-if="optimal" variant="default" class="w-fit shrink-0">
          Оптимально: {{ optimal.bank }}
        </Badge>
      </div>
    </CardHeader>
    <CardContent class="overflow-x-auto">
      <table class="w-full min-w-[520px] text-sm">
        <thead>
          <tr class="border-b text-left text-muted-foreground">
            <th class="pb-2 pr-4 font-medium">Банк</th>
            <th class="pb-2 pr-4 font-medium">Ставка</th>
            <th class="pb-2 pr-4 font-medium">Платёж/мес</th>
            <th class="pb-2 pr-4 font-medium">Переплата</th>
            <th class="pb-2 font-medium">Срок</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="bank in breakdown.banks"
            :key="bank.id"
            class="border-b last:border-0"
            :class="bank.id === breakdown.optimal_bank_id ? 'bg-primary/5' : ''"
          >
            <td class="py-3 pr-4 font-medium">
              {{ bank.bank }}
              <Badge
                v-if="bank.id === breakdown.optimal_bank_id"
                variant="secondary"
                class="ml-1 text-[10px]"
              >
                лучший баланс
              </Badge>
            </td>
            <td class="py-3 pr-4 text-muted-foreground">{{ bank.rate }}%</td>
            <td class="py-3 pr-4">{{ bank.monthly_payment.toLocaleString('ru-RU') }} ₽</td>
            <td class="py-3 pr-4 text-muted-foreground">
              {{ bank.total_overpayment.toLocaleString('ru-RU') }} ₽
            </td>
            <td class="py-3 text-muted-foreground">{{ bank.term_months }} мес.</td>
          </tr>
        </tbody>
      </table>

      <p v-if="optimal" class="mt-4 text-sm text-muted-foreground">
        {{ optimal.bank }} даёт баланс риска, платежа и переплаты — при текущей модели это выглядит
        оптимальным вариантом.
      </p>
    </CardContent>
  </Card>
</template>
