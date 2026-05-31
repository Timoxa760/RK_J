<script setup lang="ts">
import {
  CreditCard,
  LayoutGrid,
  MessageCircle,
  PieChart,
  Plus,
  ReceiptText
} from 'lucide-vue-next'
import { ADVISOR } from '~/constants/productCopy'
import { useAuthStore } from '~/store/authStore'
import { isAppFeatureEnabled } from '~/constants/featureFlags'

const authStore = useAuthStore()
const route = useRoute()
const shellReady = useShellReady()
const { show: showAddExpense } = useAddExpenseSheet()

async function onNavClick(to: string) {
  if (route.path !== to) {
    await navigateTo(to)
  }
}

const mainNav = computed(() =>
  [
    { to: '/dashboard', label: 'План', icon: LayoutGrid },
    { to: '/advisor', label: 'Советник', icon: MessageCircle },
    { to: '/receipts', label: 'Расходы', icon: ReceiptText },
    isAppFeatureEnabled('creditsNav')
      ? { to: '/credits', label: 'Кредиты', icon: CreditCard }
      : null,
    { to: '/profile', label: 'Профиль', icon: PieChart }
  ].filter(Boolean) as Array<{ to: string; label: string; icon: typeof LayoutGrid }>
)
</script>

<template>
  <aside class="mm-app-sidebar" aria-label="Навигация">
    <header class="mm-app-shell-sidebar-header shrink-0">
      <NuxtLink
        to="/dashboard"
        class="mm-sidebar-brand-link"
        aria-label="Поток — на главную"
        @click="onNavClick('/dashboard')"
      >
        <SharedHeroFlowWord variant="brand" brand-align="center" />
      </NuxtLink>
    </header>

    <div class="mm-add-purchase-sidebar shrink-0 px-3 py-3">
      <Button
        type="button"
        class="mm-add-purchase-btn w-full min-w-0 gap-2"
        data-demo="add-expense"
        :aria-label="ADVISOR.addPurchaseAria"
        @click="showAddExpense"
      >
        <Plus class="size-4 shrink-0" stroke-width="2.25" />
        <span class="truncate">{{ ADVISOR.addPurchaseLabel }}</span>
      </Button>
    </div>

    <nav class="mm-app-sidebar-nav min-h-0 flex-1 overflow-y-auto overflow-x-hidden px-2 py-2">
      <p class="mm-app-sidebar-label">Навигация</p>
      <ul class="space-y-1">
        <li v-for="item in mainNav" :key="item.to">
          <NuxtLink
            :to="item.to"
            class="mm-app-sidebar-link"
            :class="{ 'mm-app-sidebar-link--active': route.path === item.to }"
            @click="onNavClick(item.to)"
          >
            <component :is="item.icon" class="size-4 shrink-0" stroke-width="1.75" />
            <span class="truncate">{{ item.label }}</span>
          </NuxtLink>
        </li>
      </ul>
    </nav>

    <footer
      v-if="shellReady && authStore.user"
      class="shrink-0 border-t border-sidebar-border px-3 py-2"
    >
      <p class="truncate text-xs text-muted-foreground">
        {{ authStore.user.phone }}
      </p>
    </footer>
  </aside>
</template>
