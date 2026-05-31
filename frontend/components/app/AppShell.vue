<script setup lang="ts">
import { useAuthStore } from '~/store/authStore'

const authStore = useAuthStore()

onMounted(() => {
  authStore.hydrate()
})
</script>

<template>
  <ClientOnly>
    <AppShellLayout>
      <slot />
    </AppShellLayout>
    <template #fallback>
      <div class="mm-app-shell">
        <aside class="mm-app-sidebar hidden md:flex" aria-hidden="true" />
        <div class="mm-app-shell-content">
          <header
            class="mm-app-shell-header z-20 flex h-14 shrink-0 items-center border-b bg-background/95 px-3 mm-safe-top md:border-b-0"
          >
            <div class="h-4 w-32 animate-pulse rounded bg-muted" />
          </header>
          <main class="mm-app-shell-main">
            <div class="mx-auto max-w-6xl space-y-4">
              <div class="h-24 animate-pulse rounded-2xl bg-muted/60" />
              <div class="h-96 animate-pulse rounded-2xl bg-muted/40" />
            </div>
          </main>
          <nav
            class="mm-mobile-tab-bar fixed inset-x-0 bottom-0 z-30 border-t bg-background/95 md:hidden mm-safe-bottom"
            aria-hidden="true"
          >
            <div class="mm-mobile-tab-bar__grid mm-mobile-tab-bar__grid--compact h-14" />
          </nav>
        </div>
      </div>
    </template>
  </ClientOnly>
</template>
