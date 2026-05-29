<script setup lang="ts">
import { useAuthStore } from '~/store/authStore'

onMounted(() => {
  useAuthStore().hydrate()
})
</script>

<template>
  <div class="mm-landing relative min-h-screen overflow-x-hidden bg-[color:var(--mm-bg)] text-[color:var(--mm-text)]">
    <SharedBackgroundFlow />

    <LandingHeader />

    <main class="relative z-10">
      <NuxtPage />
    </main>
  </div>
</template>
