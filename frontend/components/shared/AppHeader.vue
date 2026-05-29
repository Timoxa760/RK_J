<script setup lang="ts">
import { User } from 'lucide-vue-next'
import { useAuthStore } from '~/store/authStore'

defineProps<{
  title?: string
  subtitle?: string
}>()

const authStore = useAuthStore()
const route = useRoute()

const pageTitles: Record<string, { title: string; subtitle: string }> = {
  '/dashboard': { title: 'Диагноз', subtitle: 'Финансовое состояние' },
  '/receipts': { title: 'Расходы', subtitle: 'Куда уходят деньги' },
  '/analytics': { title: 'Прогноз', subtitle: 'Что будет дальше' },
  '/credits': { title: 'Кредитный светофор', subtitle: 'Долговая нагрузка' },
  '/social': { title: 'Привычки', subtitle: 'Челленджи' },
  '/digest': { title: 'Сводка', subtitle: 'Итоги периода' },
  '/profile': { title: 'Профиль', subtitle: 'Настройки' }
}

const meta = computed(() => pageTitles[route.path] ?? { title: 'Поток', subtitle: '' })

const demoRef = ref<{ start: () => void } | null>(null)

function logout() {
  authStore.logout()
  navigateTo('/')
}

function startDemo() {
  demoRef.value?.start()
}
</script>

<template>
  <header
    class="sticky top-0 z-20 border-b border-[color:var(--mm-border)] bg-white/95 backdrop-blur-sm mm-safe-top projector-header"
  >
    <div class="flex h-14 items-center gap-2 px-3 sm:h-16 sm:gap-3 sm:px-6 lg:px-8">
      <div class="min-w-0 flex-1">
        <p
          v-if="subtitle || meta.subtitle"
          class="hidden truncate text-xs text-[color:var(--mm-text-soft)] sm:block"
        >
          {{ subtitle || meta.subtitle }}
        </p>
        <h1 class="truncate text-base font-semibold text-[color:var(--mm-text)] sm:text-lg projector-title">
          {{ title || meta.title }}
        </h1>
      </div>

      <div class="flex shrink-0 items-center gap-1.5 sm:gap-2">
        <button
          v-if="authStore.isAuthenticated"
          type="button"
          class="rounded-full border border-[color:var(--mm-border)] px-2.5 py-1.5 text-xs text-[color:var(--mm-text-muted)] transition hover:bg-[color:var(--mm-bg-muted)] sm:px-3"
          @click="startDemo"
        >
          Демо
        </button>
        <slot name="actions" />
        <button
          v-if="authStore.isAuthenticated"
          type="button"
          class="mm-pill !gap-1.5 !px-2.5 !py-1.5 text-xs sm:!px-3"
          aria-label="Выйти"
          @click="logout"
        >
          <User class="h-3.5 w-3.5 shrink-0" />
          <span class="hidden sm:inline">Выйти</span>
        </button>
      </div>
    </div>
    <DemoDemoMode ref="demoRef" />
  </header>
</template>
