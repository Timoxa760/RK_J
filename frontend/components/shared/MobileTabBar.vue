<script setup lang="ts">
import { BarChart3, CreditCard, Home, Menu, ReceiptText } from 'lucide-vue-next'
import { useSidebar } from '~/components/ui/sidebar/utils'

const route = useRoute()
const { setOpenMobile } = useSidebar()

const tabs = [
  { to: '/dashboard', label: 'Моя картина', icon: Home },
  { to: '/receipts', label: 'Расходы', icon: ReceiptText },
  { to: '/credits', label: 'Кредиты', icon: CreditCard },
  { to: '/analytics', label: 'Прогноз', icon: BarChart3 }
]
</script>

<template>
  <nav
    class="fixed inset-x-0 bottom-0 z-30 border-t bg-background/95 backdrop-blur-md md:hidden mm-safe-bottom"
    aria-label="Основная навигация"
  >
    <div class="grid grid-cols-5">
      <NuxtLink
        v-for="tab in tabs"
        :key="tab.to"
        :to="tab.to"
        class="flex flex-col items-center gap-0.5 px-1 py-2.5 text-[10px] font-medium transition"
        :class="route.path === tab.to ? 'text-primary' : 'text-muted-foreground'"
      >
        <component :is="tab.icon" class="size-5" stroke-width="1.75" />
        <span class="max-w-full truncate">{{ tab.label }}</span>
      </NuxtLink>

      <button
        type="button"
        class="flex flex-col items-center gap-0.5 px-1 py-2.5 text-[10px] font-medium text-muted-foreground"
        aria-label="Открыть меню"
        @click="setOpenMobile(true)"
      >
        <Menu class="size-5" stroke-width="1.75" />
        <span>Ещё</span>
      </button>
    </div>
  </nav>
</template>
