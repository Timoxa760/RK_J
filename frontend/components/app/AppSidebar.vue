<script setup lang="ts">
import {
  BarChart3,
  CreditCard,
  FileText,
  Home,
  Newspaper,
  PieChart,
  ReceiptText,
  Users
} from 'lucide-vue-next'
import { useAuthStore } from '~/store/authStore'
import { useSidebar } from '~/components/ui/sidebar/utils'

const authStore = useAuthStore()
const route = useRoute()
const { setOpenMobile, isMobile } = useSidebar()

function onNavClick() {
  if (isMobile.value) setOpenMobile(false)
}

const mainNav = [
  { to: '/dashboard', label: 'Моя картина', icon: Home },
  { to: '/receipts', label: 'Расходы', icon: ReceiptText },
  { to: '/analytics', label: 'Прогноз', icon: BarChart3 },
  { to: '/credits', label: 'Кредиты', icon: CreditCard },
  { to: '/social', label: 'Привычки', icon: Users },
  { to: '/digest', label: 'Сводка', icon: Newspaper },
  { to: '/profile', label: 'Профиль', icon: PieChart }
]
</script>

<template>
  <Sidebar collapsible="offcanvas">
    <SidebarHeader class="border-b border-sidebar-border">
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <NuxtLink to="/dashboard" class="flex items-center gap-2" @click="onNavClick">
              <FileText class="size-5 text-primary" />
              <span class="font-semibold">Поток</span>
            </NuxtLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>Навигация</SidebarGroupLabel>
        <SidebarMenu>
          <SidebarMenuItem v-for="item in mainNav" :key="item.to">
            <SidebarMenuButton as-child :is-active="route.path === item.to">
              <NuxtLink :to="item.to" @click="onNavClick">
                <component :is="item.icon" class="size-4" />
                <span>{{ item.label }}</span>
              </NuxtLink>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter v-if="authStore.user" class="border-t border-sidebar-border">
      <p class="truncate px-2 py-1 text-xs text-muted-foreground">
        {{ authStore.user.phone }}
      </p>
    </SidebarFooter>
  </Sidebar>
</template>
