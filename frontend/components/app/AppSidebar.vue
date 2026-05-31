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
import { useSidebar } from '~/components/ui/sidebar/utils'
import { isAppFeatureEnabled } from '~/constants/featureFlags'

const authStore = useAuthStore()
const route = useRoute()
const { setOpenMobile, isMobile } = useSidebar()
const { show: showAddExpense } = useAddExpenseSheet()

async function onNavClick(to: string) {
  if (isMobile.value) setOpenMobile(false)
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
  <Sidebar collapsible="offcanvas">
    <SidebarHeader class="mm-app-shell-sidebar-header shrink-0 border-0 p-0">
      <NuxtLink
        to="/dashboard"
        class="mm-sidebar-brand-link"
        aria-label="Поток — на главную"
        @click="onNavClick('/dashboard')"
      >
        <SharedHeroFlowWord variant="brand" brand-align="center" />
      </NuxtLink>
    </SidebarHeader>

    <div class="mm-add-purchase-sidebar hidden shrink-0 px-2 py-3 md:block">
      <Button
        type="button"
        class="mm-add-purchase-btn w-full gap-2"
        data-demo="add-expense"
        :aria-label="ADVISOR.addPurchaseAria"
        @click="showAddExpense"
      >
        <Plus class="size-4 shrink-0" stroke-width="2.25" />
        {{ ADVISOR.addPurchaseLabel }}
      </Button>
    </div>

    <SidebarContent class="gap-0 py-2">
      <SidebarGroup>
        <SidebarGroupLabel>Навигация</SidebarGroupLabel>
        <SidebarMenu>
          <SidebarMenuItem v-for="item in mainNav" :key="item.to">
            <SidebarMenuButton as-child :is-active="route.path === item.to">
              <NuxtLink :to="item.to" @click="onNavClick(item.to)">
                <component :is="item.icon" class="size-4" />
                <span>{{ item.label }}</span>
              </NuxtLink>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter v-if="authStore.user" class="shrink-0 border-t border-sidebar-border">
      <p class="truncate px-2 py-1 text-xs text-muted-foreground">
        {{ authStore.user.phone }}
      </p>
    </SidebarFooter>
  </Sidebar>
</template>
