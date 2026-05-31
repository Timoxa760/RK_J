<script setup lang="ts">
const route = useRoute()

const isBarePage = computed(() => route.path === '/login' || route.path === '/')
</script>

<template>
  <div class="font-sans antialiased relative min-h-screen overflow-x-clip bg-background text-foreground">
    <SharedBackgroundFlow v-if="isBarePage" />
    <div
      v-if="isBarePage"
      class="pointer-events-none absolute inset-0 z-[1] bg-[#f9fcff]/80"
      aria-hidden="true"
    />

    <AppShell v-if="!isBarePage">
      <NuxtPage />
    </AppShell>

    <main
      v-else
      class="relative z-10 flex min-h-[100dvh] items-center justify-center px-4 py-8"
    >
      <NuxtPage />
    </main>
  </div>
</template>
