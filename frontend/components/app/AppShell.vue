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
      <div class="flex min-h-svh w-full bg-background">
        <aside
          class="hidden w-64 shrink-0 border-r border-border bg-background md:block"
          aria-hidden="true"
        />
        <div class="flex min-w-0 flex-1 flex-col">
          <header class="flex h-14 shrink-0 items-center border-b px-4 sm:px-6">
            <div class="h-4 w-32 animate-pulse rounded bg-muted" />
          </header>
          <main class="flex-1 px-3 py-4 pb-24 sm:px-6 sm:py-6 lg:px-8">
            <slot />
          </main>
        </div>
      </div>
    </template>
  </ClientOnly>
</template>
