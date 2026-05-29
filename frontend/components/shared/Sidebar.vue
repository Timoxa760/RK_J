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

defineProps<{
  mobileOpen?: boolean
}>()

defineEmits<{
  close: []
}>()

const authStore = useAuthStore()

const mainNav = [
  { to: '/dashboard', label: 'Диагноз', icon: Home },
  { to: '/receipts', label: 'Расходы', icon: ReceiptText },
  { to: '/analytics', label: 'Прогноз', icon: BarChart3 },
  { to: '/credits', label: 'Кредиты', icon: CreditCard },
  { to: '/social', label: 'Привычки', icon: Users },
  { to: '/digest', label: 'Сводка', icon: Newspaper },
  { to: '/profile', label: 'Профиль', icon: PieChart }
]
</script>

<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 flex w-[min(20rem,88vw)] flex-col border-r border-[color:var(--mm-border)] bg-white transition-transform duration-200 mm-safe-top mm-safe-bottom lg:z-30 lg:w-64 lg:translate-x-0"
    :class="mobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'"
  >
    <div class="flex h-16 items-center border-b border-[color:var(--mm-border-subtle)] px-4">
      <NuxtLink
        to="/dashboard"
        class="flex items-center gap-2 font-semibold text-[color:var(--mm-text)]"
        @click="$emit('close')"
      >
        <FileText class="h-5 w-5 text-[color:var(--mm-primary)]" stroke-width="1.75" />
        Поток
      </NuxtLink>
    </div>

    <nav class="flex-1 overflow-y-auto p-3">
      <NuxtLink
        v-for="item in mainNav"
        :key="item.to"
        :to="item.to"
        class="mb-0.5 flex items-center gap-3 rounded-2xl px-3 py-2.5 text-sm text-[color:var(--mm-text-muted)] transition hover:bg-[color:var(--mm-bg-muted)] hover:text-[color:var(--mm-text)]"
        active-class="!bg-[color:var(--mm-primary-soft)] !font-medium !text-[color:var(--mm-text)]"
        @click="$emit('close')"
      >
        <component :is="item.icon" class="h-4 w-4" stroke-width="1.75" />
        {{ item.label }}
      </NuxtLink>
    </nav>

    <div
      v-if="authStore.user"
      class="border-t border-[color:var(--mm-border-subtle)] p-3 text-xs text-[color:var(--mm-text-soft)]"
    >
      {{ authStore.user.phone }}
    </div>
  </aside>

  <div
    v-if="mobileOpen"
    class="fixed inset-0 z-30 bg-black/25 backdrop-blur-[2px] lg:hidden"
    aria-hidden="true"
    @click="$emit('close')"
  />
</template>
