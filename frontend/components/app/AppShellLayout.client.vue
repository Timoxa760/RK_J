<script setup lang="ts">
import { useAuthStore } from '~/store/authStore'

const authStore = useAuthStore()
const demoModeRef = ref<{ start: () => void; stop: () => void } | null>(null)
const config = useRuntimeConfig()
const route = useRoute()

const showDemoTour = computed(
  () => config.public.demoMode && route.query.tour === '1'
)

const { open: addExpenseOpen, notifyAdded } = useAddExpenseSheet()

onMounted(() => {
  authStore.hydrate()
  if (showDemoTour.value) {
    nextTick(() => demoModeRef.value?.start())
  }
})

watch(showDemoTour, (enabled) => {
  if (enabled) nextTick(() => demoModeRef.value?.start())
})
</script>

<template>
  <SidebarProvider
    :default-open="true"
    class="mm-app-shell mm-app-shell--with-advisor"
    :style="{ '--sidebar-width': '20rem' }"
  >
    <AppShellAdvisorScope>
      <AppSidebar />
      <SidebarInset class="mm-app-shell-inset">
        <SharedAppHeader />
        <div class="mm-app-shell-main flex-1 px-3 py-4 pb-24 sm:px-6 sm:py-6 sm:pb-6 lg:px-8 lg:pb-8">
          <slot />
        </div>
        <SharedMobileTabBar />
      </SidebarInset>
    </AppShellAdvisorScope>
    <ClientOnly>
      <DashboardAddExpenseSheet v-model:open="addExpenseOpen" @added="notifyAdded" />
    </ClientOnly>
    <DemoMode v-if="config.public.demoMode" ref="demoModeRef" />
    <Sonner />
  </SidebarProvider>
</template>
