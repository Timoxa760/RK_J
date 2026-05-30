<script setup lang="ts">
import { CreditCard, LayoutGrid, Menu, Plus, ReceiptText } from 'lucide-vue-next'
import { ADVISOR } from '~/constants/productCopy'
import { useSidebar } from '~/components/ui/sidebar/utils'
import { isAppFeatureEnabled } from '~/constants/featureFlags'

const route = useRoute()
const { setOpenMobile } = useSidebar()
const { show: showAddExpense } = useAddExpenseSheet()

const showCredits = computed(() => isAppFeatureEnabled('creditsNav'))
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
        v-if="showCredits"
        to="/credits"
        class="mm-mobile-tab-bar__item"
        :class="route.path === '/credits' ? 'text-primary' : 'text-muted-foreground'"
      >
        <CreditCard class="size-5" stroke-width="1.75" />
        <span class="max-w-full truncate">Кредиты</span>
      </NuxtLink>
      <div v-else class="mm-mobile-tab-bar__item" aria-hidden="true" />

      <button
        type="button"
        class="mm-mobile-tab-bar__item text-muted-foreground"
        aria-label="Открыть меню"
        @click="setOpenMobile(true)"
      >
        <Menu class="size-5" stroke-width="1.75" />
        <span>Ещё</span>
      </button>
    </div>
  </nav>
</template>
