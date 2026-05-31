<script setup lang="ts">
import { LogOut } from 'lucide-vue-next'
import { NAV } from '~/constants/productCopy'
import { useAuth } from '~/composables/useAuth'

defineProps<{
  title?: string
  subtitle?: string
}>()

const { authStore, logout: signOut } = useAuth()
const route = useRoute()

const pageTitles: Record<string, { title: string; subtitle: string }> = {
  '/dashboard': { title: NAV.dashboard, subtitle: NAV.dashboardSubtitle },
  '/advisor': { title: NAV.advisor, subtitle: NAV.advisorSubtitle },
  '/receipts': { title: NAV.receipts, subtitle: 'Куда уходят деньги' },
  '/credits': { title: NAV.creditsTitle, subtitle: NAV.creditsSubtitle },
  '/profile': { title: NAV.profile, subtitle: 'Доход, запас и цель' }
}

const meta = computed(() => pageTitles[route.path] ?? { title: 'Поток', subtitle: '' })

function logout() {
  signOut()
  navigateTo('/')
}
</script>

<template>
  <header
    class="mm-app-shell-header z-20 flex h-14 shrink-0 items-center gap-2 border-b bg-background/95 px-3 backdrop-blur-sm sm:h-16 sm:gap-3 sm:px-6 mm-safe-top md:border-b-0"
  >
    <div class="flex min-w-0 shrink-0 items-center gap-2 md:gap-0">
      <NuxtLink
        to="/"
        class="mm-app-header-brand shrink-0 md:hidden"
        aria-label="Поток — на лендинг"
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
        type="button"
        variant="ghost"
        size="icon"
        class="size-9 shrink-0 text-muted-foreground hover:text-destructive sm:hidden"
        aria-label="Выйти"
        @click="logout"
      >
        <LogOut class="size-4" />
      </Button>
      <Button
        v-if="authStore.isAuthenticated"
        type="button"
        variant="outline"
        size="sm"
        class="hidden gap-1.5 sm:inline-flex"
        @click="logout"
      >
        <LogOut class="size-3.5 shrink-0" />
        Выйти
      </Button>
    </div>
  </header>
</template>
