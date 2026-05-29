<script setup lang="ts">
import { useAuthStore } from '~/store/authStore'

const route = useRoute()
const authStore = useAuthStore()
const sidebarOpen = ref(false)

const isBarePage = computed(() => route.path === '/login' || route.path === '/')
const showShell = computed(() => !isBarePage.value)

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}

function closeSidebar() {
  sidebarOpen.value = false
}

watch(
  () => route.path,
  () => {
    sidebarOpen.value = false
  }
)

onMounted(() => {
  authStore.hydrate()
})
</script>

<template>
  <div class="relative min-h-screen overflow-x-hidden bg-[color:var(--mm-bg)] text-[color:var(--mm-text)]">
    <SharedBackgroundFlow />

    <SharedSidebar
      v-if="showShell"
      :mobile-open="sidebarOpen"
      @close="closeSidebar"
    />
    <div class="relative z-10" :class="showShell ? 'lg:pl-64' : ''">
      <SharedAppHeader v-if="showShell" />
      <main
        :class="
          showShell
            ? 'px-3 py-4 pb-24 sm:px-6 sm:py-6 sm:pb-6 lg:px-8 lg:pb-8'
            : 'flex min-h-[100dvh] items-center justify-center px-4 py-8'
        "
      >
        <NuxtPage />
      </main>
    </div>

    <SharedMobileTabBar v-if="showShell" @open-menu="toggleSidebar" />
    <SharedToast />
  </div>
</template>
