<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { PURCHASES } from '~/constants/productCopy'
import type { ReceiptListItem } from '~/types/api'
import { formatGoalImpact } from '~/utils/receiptImpact'
import { buildReceiptsGoalDelayPrompt } from '~/utils/advisorChat'
import { ADVISOR } from '~/constants/productCopy'

const { receipts, selected, selectReceipt, closeDetail, refresh } = useReceiptList()
const { addedVersion, show: showAddExpense } = useAddExpenseSheet()

watch(addedVersion, () => {
  refresh()
})

const totals = computed(() => {
  const total = receipts.value.reduce((sum, r) => sum + r.amount, 0)
  const impulse = receipts.value.reduce((sum, r) => sum + (r.impulse_count ?? 0), 0)
  return { total, impulse, count: receipts.value.length }
})

function impactFor(receipt: ReceiptListItem) {
  return formatGoalImpact(receipt.amount)
}
</script>

<template>
  <div class="mx-auto w-full max-w-4xl space-y-6">
    <header>
      <h1 class="text-xl font-semibold tracking-tight sm:text-2xl">Расходы</h1>
      <p class="mt-1 text-sm text-muted-foreground">Добавляйте покупки — они учтутся в вашей картине.</p>
    </header>

    <Card v-if="receipts.length" class="overflow-hidden" data-demo="receipts">
      <CardHeader class="pb-2">
        <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle class="text-base">Лента покупок</CardTitle>
            <CardDescription>Нажмите на покупку — детализация по позициям</CardDescription>
          </div>
          <p class="text-sm text-muted-foreground">
            {{ totals.count }} шт. · {{ totals.total.toLocaleString('ru-RU') }} ₽
            <span v-if="totals.impulse"> · {{ PURCHASES.impulseCount(totals.impulse) }}</span>
          </p>
        </div>
      </CardHeader>
      <CardContent class="p-0">
        <ul class="divide-y">
          <li
            v-for="r in receipts"
            :key="r.id"
            class="flex cursor-pointer flex-col gap-2 px-4 py-4 text-sm transition-colors hover:bg-muted/50 sm:flex-row sm:items-center sm:justify-between sm:px-6"
            @click="selectReceipt(r)"
          >
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <p class="font-medium">{{ r.store }}</p>
                <Badge v-if="r.impulse_count" variant="secondary" class="text-[10px]">
                  {{ PURCHASES.impulseBadge }}
                </Badge>
              </div>
              <p class="mt-0.5 text-xs text-muted-foreground">
                {{ r.date }}
                <span v-if="r.category"> · {{ r.category }}</span>
              </p>
              <p v-if="impactFor(r)" class="mt-1 text-xs text-primary">{{ impactFor(r) }}</p>
            </div>
            <p class="shrink-0 text-base font-semibold">{{ r.amount.toLocaleString('ru-RU') }} ₽</p>
          </li>
        </ul>
      </CardContent>
    </Card>

    <Card v-else class="text-center">
      <CardContent class="py-12">
        <p class="text-sm text-muted-foreground">Покупок пока нет.</p>
        <Button class="mt-4 mm-add-purchase-btn gap-2" @click="showAddExpense">
          <Plus class="size-4" />
          Добавить
        </Button>
      </CardContent>
    </Card>

    <Dialog :open="Boolean(selected)" @update:open="(v) => !v && closeDetail()">
      <DialogContent v-if="selected" class="max-h-[85vh] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{{ selected.store }}</DialogTitle>
          <DialogDescription>{{ selected.date }}</DialogDescription>
        </DialogHeader>
        <p class="text-xl font-bold">{{ selected.amount.toLocaleString('ru-RU') }} ₽</p>
        <p v-if="impactFor(selected)" class="text-sm text-primary">
          {{ PURCHASES.goalDelay }} {{ impactFor(selected).replace('≈ ', '') }}.
        </p>
        <AdvisorAskButton
          v-if="selected && impactFor(selected)"
          :prompt="buildReceiptsGoalDelayPrompt(impactFor(selected)!)"
          :label="ADVISOR.askGoalDelay"
        />
        <ul v-if="selected.items?.length" class="space-y-2">
          <li
            v-for="item in selected.items"
            :key="item.name"
            class="flex justify-between gap-2 rounded-lg bg-muted px-3 py-2 text-sm"
          >
            <span>
              {{ item.name }}
              <Badge v-if="item.impulse" variant="secondary" class="ml-1 text-[10px]">{{ PURCHASES.impulseBadge }}</Badge>
            </span>
            <span class="shrink-0 text-muted-foreground">
              {{ (item.price * (item.quantity ?? 1)).toLocaleString('ru-RU') }} ₽
            </span>
          </li>
        </ul>
      </DialogContent>
    </Dialog>
  </div>
</template>
