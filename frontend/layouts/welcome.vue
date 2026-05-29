<script setup lang="ts">
import { Menu, User, X } from 'lucide-vue-next'
import { useAuthStore } from '~/store/authStore'

const route = useRoute()
const authStore = useAuthStore()
const mobileNavOpen = ref(false)

const appNavLinks = [
  { label: 'Диагноз', to: '/dashboard' },
  { label: 'Расходы', to: '/receipts' },
  { label: 'Кредиты', to: '/credits' },
  { label: 'Прогноз', to: '/analytics' }
]

function closeMobileNav() {
  mobileNavOpen.value = false
}

function toggleMobileNav() {
  mobileNavOpen.value = !mobileNavOpen.value
}

watch(
  () => route.path,
  () => {
    mobileNavOpen.value = false
  }
)

onMounted(() => {
  authStore.hydrate()
})
</script>

<template>
  <div class="relative min-h-screen overflow-x-hidden bg-[color:var(--mm-bg)] text-[color:var(--mm-text)]">
    <SharedBackgroundFlow />

    <header class="relative z-30 bg-white/90 backdrop-blur-md mm-safe-top">
      <div class="mx-auto flex h-14 max-w-[1400px] items-center gap-2 px-4 sm:h-[var(--mm-header-height)] sm:gap-6 sm:px-8 lg:px-12">
        <NuxtLink
          to="/"
          class="shrink-0 text-lg font-semibold tracking-tight text-[color:var(--mm-text)] sm:text-xl"
          @click="closeMobileNav"
        >
          Поток
        </NuxtLink>

        <nav
          v-if="authStore.isAuthenticated"
          class="hidden flex-1 items-center justify-center gap-6 md:flex lg:gap-8"
        >
          <NuxtLink
            v-for="link in appNavLinks"
            :key="link.label"
            :to="link.to"
            class="text-sm text-[color:var(--mm-text-muted)] transition hover:text-[color:var(--mm-text)]"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>

        <div class="ml-auto flex items-center gap-2">
          <button
            v-if="authStore.isAuthenticated"
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-full border border-[color:var(--mm-border)] text-[color:var(--mm-text-muted)] transition hover:bg-[color:var(--mm-bg-muted)] md:hidden"
            :aria-expanded="mobileNavOpen"
            aria-controls="welcome-mobile-nav"
            :aria-label="mobileNavOpen ? 'Закрыть меню' : 'Открыть меню'"
            @click="toggleMobileNav"
          >
            <X v-if="mobileNavOpen" class="h-5 w-5" stroke-width="1.75" />
            <Menu v-else class="h-5 w-5" stroke-width="1.75" />
          </button>

          <NuxtLink
            v-if="!authStore.isAuthenticated"
            to="/login"
            class="mm-pill shrink-0 !px-3 !py-1.5 text-xs sm:!px-4 sm:!py-2 sm:text-sm transition hover:border-[color:var(--mm-text-soft)]"
          >
            <User class="h-4 w-4" stroke-width="1.75" />
            <span class="hidden sm:inline">Войти</span>
          </NuxtLink>

          <NuxtLink
            v-else
            to="/dashboard"
            class="mm-btn-primary shrink-0 !px-3 !py-1.5 text-xs sm:!px-4 sm:!py-2 sm:text-sm"
            @click="closeMobileNav"
          >
            К приложению
          </NuxtLink>
        </div>
      </div>

      <Transition name="welcome-nav">
        <div
          v-if="mobileNavOpen && authStore.isAuthenticated"
          id="welcome-mobile-nav"
          class="border-t border-[color:var(--mm-border)] bg-white md:hidden"
        >
          <nav class="mx-auto max-w-[1400px] px-4 py-3">
            <NuxtLink
              v-for="link in appNavLinks"
              :key="link.label"
              :to="link.to"
              class="flex min-h-11 items-center rounded-xl px-3 text-sm font-medium text-[color:var(--mm-text-muted)] transition active:bg-[color:var(--mm-bg-muted)]"
              @click="closeMobileNav"
            >
              {{ link.label }}
            </NuxtLink>
          </nav>
        </div>
      </Transition>
    </header>

    <div
      v-if="mobileNavOpen"
      class="fixed inset-0 z-20 bg-black/20 md:hidden"
      aria-hidden="true"
      @click="closeMobileNav"
    />

    <main class="relative z-10">
      <NuxtPage />
    </main>
  </div>
</template>

<style scoped>
.welcome-nav-enter-active,
.welcome-nav-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.welcome-nav-enter-from,
.welcome-nav-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
