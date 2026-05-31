<script setup lang="ts">
import { LayoutGrid, MessageCircle, PieChart, Plus, ReceiptText } from 'lucide-vue-next'
import { ADVISOR } from '~/constants/productCopy'

const route = useRoute()
const { show: showAddExpense } = useAddExpenseSheet()
</script>

<template>
  <nav
    class="mm-mobile-tab-bar fixed inset-x-0 bottom-0 z-30 border-t bg-background/95 backdrop-blur-md md:hidden mm-safe-bottom"
    aria-label="Основная навигация"
  >
    <div class="grid grid-cols-5 items-end">
      <NuxtLink
        to="/dashboard"
        class="mm-mobile-tab-bar__item"
        :class="route.path === '/dashboard' ? 'text-primary' : 'text-muted-foreground'"
      >
        <LayoutGrid class="size-5" stroke-width="1.75" />
        <span class="max-w-full truncate">План</span>
      </NuxtLink>

      <NuxtLink
        to="/receipts"
        class="mm-mobile-tab-bar__item"
        :class="route.path === '/receipts' ? 'text-primary' : 'text-muted-foreground'"
      >
        <ReceiptText class="size-5" stroke-width="1.75" />
        <span class="max-w-full truncate">Расходы</span>
      </NuxtLink>

      <div class="mm-mobile-tab-bar__add-wrap flex flex-col items-center justify-end pb-1">
        <button
          type="button"
          class="mm-add-purchase-tab-btn"
          data-demo="add-expense"
          :aria-label="ADVISOR.addPurchaseAria"
          @click="showAddExpense"
        >
          <Plus class="size-6 shrink-0" stroke-width="2.5" />
        </button>
        <span class="mt-1 max-w-full truncate text-[10px] font-medium text-muted-foreground">
          Добавить
        </span>
      </div>

      <NuxtLink
        to="/advisor"
        class="mm-mobile-tab-bar__item"
        :class="route.path === '/advisor' ? 'text-primary' : 'text-muted-foreground'"
      >
        <MessageCircle class="size-5" stroke-width="1.75" />
        <span class="max-w-full truncate">Советник</span>
      </NuxtLink>

      <NuxtLink
        to="/profile"
        class="mm-mobile-tab-bar__item"
        :class="route.path === '/profile' ? 'text-primary' : 'text-muted-foreground'"
      >
        <PieChart class="size-5" stroke-width="1.75" />
        <span class="max-w-full truncate">Профиль</span>
      </NuxtLink>
    </div>
  </nav>
</template>
