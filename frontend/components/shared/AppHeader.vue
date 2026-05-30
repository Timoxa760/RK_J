<script setup lang="ts">
import { User } from 'lucide-vue-next'
import { NAV } from '~/constants/productCopy'
import { useAuthStore } from '~/store/authStore'

defineProps<{
  title?: string
  subtitle?: string
}>()

const authStore = useAuthStore()
const route = useRoute()

const pageTitles: Record<string, { title: string; subtitle: string }> = {
  '/dashboard': { title: NAV.dashboard, subtitle: NAV.dashboardSubtitle },
  '/receipts': { title: NAV.receipts, subtitle: 'Куда уходят деньги' },
  '/credits': { title: NAV.creditsTitle, subtitle: NAV.creditsSubtitle },
  '/social': { title: NAV.social, subtitle: NAV.socialSubtitle },
  '/digest': { title: NAV.digest, subtitle: 'Итоги периода' },
  '/profile': { title: NAV.profile, subtitle: 'Настройки' }
}

const meta = computed(() => pageTitles[route.path] ?? { title: 'Поток', subtitle: '' })

function logout() {
  authStore.logout()
  navigateTo('/')
}
</script>

<template>
  <header
    class="mm-app-shell-header sticky top-0 z-20 flex h-14 shrink-0 items-center gap-2 border-b bg-background/95 px-3 backdrop-blur-sm sm:h-16 sm:gap-3 sm:px-6 mm-safe-top md:border-b-0"
  >
    <div class="flex min-w-0 shrink-0 items-center gap-2 md:gap-0">
      <SidebarTrigger class="md:hidden" />
      <NuxtLink
        to="/dashboard"
        class="mm-app-header-brand shrink-0 md:hidden"
        aria-label="Поток — на главную"
      >
        <SharedHeroFlowWord variant="brand" />
      </NuxtLink>
    </div>
    <div class="min-w-0 flex-1">
      <p
        v-if="subtitle || meta.subtitle"
        class="hidden truncate text-xs text-muted-foreground sm:block"
      >
        {{ subtitle || meta.subtitle }}
      </p>
      <h1 class="truncate text-base font-semibold sm:text-lg">
        {{ title || meta.title }}
      </h1>
    </div>

    <div class="flex shrink-0 items-center gap-1.5 sm:gap-2">
      <slot name="actions" />
      <Button
        v-if="authStore.isAuthenticated"
        variant="outline"
        size="sm"
        class="gap-1.5"
        @click="logout"
      >
        <User class="size-3.5 shrink-0" />
        <span class="hidden sm:inline">Выйти</span>
      </Button>
    </div>
  </header>
</template>
