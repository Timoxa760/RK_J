<script setup lang="ts">
import type { Component } from 'vue'
import { Landmark, LayoutGrid, MessageCircle, Plus, ReceiptText } from 'lucide-vue-next'
import { ADVISOR } from '~/constants/productCopy'
import { isAppFeatureEnabled } from '~/constants/featureFlags'

type TabItem = {
  to: string
  label: string
  shortLabel?: string
  icon: Component
}

const route = useRoute()
const { show: showAddExpense } = useAddExpenseSheet()

const showCreditsTab = isAppFeatureEnabled('creditsNav')

const leftTabs: TabItem[] = [
  { to: '/dashboard', label: 'План', icon: LayoutGrid },
  { to: '/receipts', label: 'Расходы', icon: ReceiptText }
]

const advisorTab: TabItem = { to: '/advisor', label: 'Советник', icon: MessageCircle }

const creditsTab: TabItem = {
  to: '/credits',
  label: 'Кредиты',
  shortLabel: 'Кредит',
  icon: Landmark
}

const gridClass = computed(() =>
  showCreditsTab ? 'mm-mobile-tab-bar__grid' : 'mm-mobile-tab-bar__grid mm-mobile-tab-bar__grid--compact'
)

function isActive(path: string) {
  return route.path === path
}

function tabLabel(item: TabItem) {
  return item.shortLabel
    ? { full: item.label, short: item.shortLabel }
    : { full: item.label, short: item.label }
}
</script>

<template>
  <nav
    class="mm-mobile-tab-bar fixed inset-x-0 bottom-0 z-30 border-t bg-background/95 backdrop-blur-md md:hidden mm-safe-bottom"
    aria-label="Основная навигация"
  >
    <button
      type="button"
      class="mm-add-purchase-tab-btn mm-mobile-tab-bar__fab"
      data-demo="add-expense"
      :aria-label="ADVISOR.addPurchaseAria"
      @click="showAddExpense"
    >
      <Plus class="size-6 shrink-0" stroke-width="2.5" />
    </button>

    <div :class="gridClass">
      <NuxtLink
        v-for="item in leftTabs"
        :key="item.to"
        :to="item.to"
        class="mm-mobile-tab-bar__item"
        :class="isActive(item.to) ? 'text-primary' : 'text-muted-foreground'"
      >
        <component :is="item.icon" class="size-5 shrink-0" stroke-width="1.75" />
        <span class="max-w-full truncate">
          <span class="mm-mobile-tab-bar__item-label--full">{{ tabLabel(item).full }}</span>
          <span class="mm-mobile-tab-bar__item-label--short">{{ tabLabel(item).short }}</span>
        </span>
      </NuxtLink>

      <div class="mm-mobile-tab-bar__add-slot">
        <span class="mm-mobile-tab-bar__add-label">Добавить</span>
      </div>

      <NuxtLink
        :to="advisorTab.to"
        class="mm-mobile-tab-bar__item"
        :class="isActive(advisorTab.to) ? 'text-primary' : 'text-muted-foreground'"
      >
        <component :is="advisorTab.icon" class="size-5 shrink-0" stroke-width="1.75" />
        <span class="max-w-full truncate">{{ advisorTab.label }}</span>
      </NuxtLink>

      <NuxtLink
        v-if="showCreditsTab"
        :to="creditsTab.to"
        class="mm-mobile-tab-bar__item"
        :class="isActive(creditsTab.to) ? 'text-primary' : 'text-muted-foreground'"
      >
        <component :is="creditsTab.icon" class="size-5 shrink-0" stroke-width="1.75" />
        <span class="max-w-full truncate">
          <span class="mm-mobile-tab-bar__item-label--full">{{ tabLabel(creditsTab).full }}</span>
          <span class="mm-mobile-tab-bar__item-label--short">{{ tabLabel(creditsTab).short }}</span>
        </span>
      </NuxtLink>
    </div>
  </nav>
</template>
