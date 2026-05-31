<script setup lang="ts">
import type { Component } from 'vue'
import { Landmark, LayoutGrid, MessageCircle, Plus, ReceiptText } from 'lucide-vue-next'
import { ADVISOR } from '~/constants/productCopy'
import { isAppFeatureEnabled } from '~/constants/featureFlags'

type TabItem = {
  to: string
  label: string
  icon: Component
}

const route = useRoute()
const { show: showAddExpense } = useAddExpenseSheet()

const showCreditsTab = isAppFeatureEnabled('creditsNav')

/** 2 + «Добавить» + 2 — профиль только в хедере на мобилке. */
const leftTabs: TabItem[] = [
  { to: '/dashboard', label: 'План', icon: LayoutGrid },
  { to: '/receipts', label: 'Расходы', icon: ReceiptText }
]

const rightTabs = computed<TabItem[]>(() => {
  const advisor: TabItem = { to: '/advisor', label: 'Советник', icon: MessageCircle }
  if (showCreditsTab) {
    return [advisor, { to: '/credits', label: 'Кредиты', icon: Landmark }]
  }
  return [advisor]
})

function isActive(path: string) {
  return route.path === path
}
</script>

<template>
  <nav
    class="mm-mobile-tab-bar fixed inset-x-0 bottom-0 z-30 border-t bg-background/95 backdrop-blur-md md:hidden mm-safe-bottom"
    aria-label="Основная навигация"
  >
    <div class="mm-mobile-tab-bar__row flex items-end">
      <div class="mm-mobile-tab-bar__side flex min-w-0 flex-1 justify-evenly">
        <NuxtLink
          v-for="item in leftTabs"
          :key="item.to"
          :to="item.to"
          class="mm-mobile-tab-bar__item"
          :class="isActive(item.to) ? 'text-primary' : 'text-muted-foreground'"
        >
          <component :is="item.icon" class="size-5" stroke-width="1.75" />
          <span class="max-w-full truncate">{{ item.label }}</span>
        </NuxtLink>
      </div>

      <div class="mm-mobile-tab-bar__add-wrap flex w-[4.25rem] shrink-0 flex-col items-center justify-end pb-1">
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

      <div class="mm-mobile-tab-bar__side flex min-w-0 flex-1 justify-evenly">
        <NuxtLink
          v-for="item in rightTabs"
          :key="item.to"
          :to="item.to"
          class="mm-mobile-tab-bar__item"
          :class="isActive(item.to) ? 'text-primary' : 'text-muted-foreground'"
        >
          <component :is="item.icon" class="size-5" stroke-width="1.75" />
          <span class="max-w-full truncate">{{ item.label }}</span>
        </NuxtLink>
        <span
          v-if="rightTabs.length < 2"
          class="mm-mobile-tab-bar__item pointer-events-none invisible"
          aria-hidden="true"
        >
          <span class="size-5" />
          <span class="text-[10px]">.</span>
        </span>
      </div>
    </div>
  </nav>
</template>
