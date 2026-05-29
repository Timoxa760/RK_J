<script setup lang="ts">
import type { ReceiptItem } from '~/types/api'

const receipts = ref<ReceiptItem[]>([
  { id: '1', store: 'Пятёрочка', amount: 1240, date: '2026-05-28', category: 'Продукты' },
  { id: '2', store: 'Starbucks', amount: 420, date: '2026-05-28', category: 'Кафе' },
  { id: '3', store: 'Ozon', amount: 3890, date: '2026-05-27', category: 'Покупки' }
])
</script>

<template>
  <section class="mm-card overflow-hidden" data-demo="receipts">
    <div class="border-b border-[color:var(--mm-border-subtle)] px-4 py-3 sm:px-6 sm:py-4">
      <h2 class="text-sm font-semibold text-[color:var(--mm-text)]">Расходы</h2>
    </div>
    <ul class="divide-y divide-[color:var(--mm-border-subtle)]">
      <li
        v-for="r in receipts"
        :key="r.id"
        class="flex flex-col gap-1 px-4 py-4 text-sm sm:flex-row sm:items-center sm:justify-between sm:px-6"
      >
        <div class="min-w-0">
          <p class="font-medium text-[color:var(--mm-text)]">{{ r.store }}</p>
          <p class="text-xs text-[color:var(--mm-text-soft)]">{{ r.date }} · {{ r.category }}</p>
        </div>
        <p class="shrink-0 font-semibold text-[color:var(--mm-text)]">{{ r.amount.toLocaleString('ru-RU') }} ₽</p>
      </li>
    </ul>
  </section>
</template>
