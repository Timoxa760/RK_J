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
const shellReady = useShellReady()

const pageTitles: Record<string, { title: string; subtitle: string }> = {
  '/dashboard': { title: NAV.dashboard, subtitle: NAV.dashboardSubtitle },
  '/advisor': { title: NAV.advisor, subtitle: NAV.advisorSubtitle },
  '/receipts': { title: NAV.receipts, subtitle: 'Куда уходят деньги' },
  '/credits': { title: NAV.creditsTitle, subtitle: NAV.creditsSubtitle },
  '/profile': { title: NAV.profile, subtitle: 'Доход, запас и цель' }
}

const meta = computed(() => pageTitles[route.path] ?? { title: 'Поток', subtitle: '' })
const profileActive = computed(() => route.path === '/profile')
const displaySubtitle = computed(() => meta.value.subtitle)
</script>

<template>
  <header
    class="mm-app-shell-header z-20 flex h-14 shrink-0 items-center gap-2 border-b bg-background/95 px-3 backdrop-blur-sm sm:h-16 sm:gap-3 sm:px-6 mm-safe-top md:border-b-0"
  >
    <div class="mm-app-header-title-block">
      <h1 class="truncate text-base font-semibold leading-tight sm:text-lg">
        {{ title || meta.title }}
      </h1>
      <p
        v-if="subtitle || displaySubtitle"
        class="mm-app-header-subtitle hidden sm:block"
      >
        {{ subtitle || displaySubtitle }}
      </p>
    </div>

    <div class="relative z-10 flex shrink-0 items-center gap-1.5 sm:gap-2">
      <slot name="actions" />
      <NuxtLink
        v-if="shellReady && authStore.isAuthenticated"
        to="/profile"
        class="mm-app-header-profile-link md:hidden"
        :class="{ 'mm-app-header-profile-link--active': profileActive }"
        :aria-current="profileActive ? 'page' : undefined"
      >
        <User class="size-[1.125rem] shrink-0" stroke-width="2" />
        <span>{{ NAV.profile }}</span>
      </NuxtLink>
    </div>
  </header>
</template>
