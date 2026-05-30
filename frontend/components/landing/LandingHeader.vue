<script setup lang="ts">
import { ArrowRight, Menu, User, X } from 'lucide-vue-next'
import { LANDING_ANCHORS } from '~/constants/landingContent'
import { scrollToLandingSection } from '~/composables/useLandingAnchor'
import { needsOnboarding } from '~/composables/useOnboarding'
import { useAuthStore } from '~/store/authStore'
import { isAppFeatureEnabled } from '~/constants/featureFlags'

const route = useRoute()
const authStore = useAuthStore()
const mobileNavOpen = ref(false)
const headerScrolled = ref(false)
const headerAnimate = ref(false)
const headerEntered = ref(false)
/** После mount — иначе SSR (гость) ≠ клиент с token из localStorage → hydration mismatch. */
const headerAuthReady = ref(false)
const headerAuthenticated = ref(false)

const isLanding = computed(() => route.path === '/')

const appEntryPath = computed(() =>
  authStore.isAuthenticated && needsOnboarding()
    ? '/onboarding'
    : '/dashboard'
)

const appNavLinks = computed(() =>
  [
    { label: 'Диагноз', to: '/dashboard' },
    { label: 'Расходы', to: '/receipts' },
    isAppFeatureEnabled('creditsNav') ? { label: 'Кредиты', to: '/credits' } : null,
    { label: 'Советник', to: '/dashboard' }
  ].filter(Boolean) as Array<{ label: string; to: string }>
)

function closeMobileNav() {
  mobileNavOpen.value = false
}

function toggleMobileNav() {
  mobileNavOpen.value = !mobileNavOpen.value
}

function goToAnchor(id: string) {
  scrollToLandingSection(id)
  closeMobileNav()
}

function onScroll() {
  headerScrolled.value = window.scrollY > 8
}

watch(
  () => route.path,
  () => {
    mobileNavOpen.value = false
  }
)

watch(mobileNavOpen, (open) => {
  if (!import.meta.client) return
  document.body.style.overflow = open ? 'hidden' : ''
})

onMounted(() => {
  if (!authStore.token) {
    authStore.hydrate()
  }
  headerAuthenticated.value = authStore.isAuthenticated
  headerAuthReady.value = true
  onScroll()
  window.addEventListener('scroll', onScroll, { passive: true })

  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (!reduced) {
    headerAnimate.value = true
    requestAnimationFrame(() => {
      headerEntered.value = true
    })
  } else {
    headerEntered.value = true
  }
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  closeMobileNav()
})
</script>

<template>
  <header
    class="mm-landing-header"
    :class="{
      'mm-landing-header--scrolled': headerScrolled,
      'mm-landing-header--open': mobileNavOpen,
      'mm-landing-header--animate': headerAnimate,
      'mm-landing-header--entered': headerEntered
    }"
  >
    <div class="mm-landing-header__bar mx-auto flex h-14 max-w-[1400px] items-center gap-3 px-4 sm:h-16 sm:px-8 lg:px-12">
      <NuxtLink to="/" class="mm-landing-header__brand group shrink-0" @click="closeMobileNav">
        <span class="text-lg font-bold tracking-tight text-[color:var(--mm-text)] sm:text-xl">Поток</span>
        <span class="hidden text-[10px] font-medium uppercase tracking-[0.2em] text-[color:var(--mm-text-soft)] sm:block">
          голос · траты
        </span>
      </NuxtLink>

      <nav
        v-if="headerAuthReady && headerAuthenticated"
        class="mm-landing-header__nav hidden flex-1 items-center justify-center gap-1 md:flex lg:gap-2"
        aria-label="Приложение"
      >
        <NuxtLink
          v-for="link in appNavLinks"
          :key="link.label"
          :to="link.to"
          class="mm-landing-header__link rounded-full px-3 py-2 text-sm"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <nav
        v-else-if="isLanding"
        class="mm-landing-header__nav hidden flex-1 items-center justify-center gap-1 sm:flex lg:gap-2"
        aria-label="Разделы лендинга"
      >
        <button
          v-for="link in LANDING_ANCHORS"
          :key="link.id"
          type="button"
          class="mm-landing-header__link rounded-full px-3 py-2 text-sm"
          @click="goToAnchor(link.id)"
        >
          {{ link.label }}
        </button>
      </nav>

      <div class="ml-auto flex items-center gap-2">
        <NuxtLink
          v-if="!headerAuthReady || !headerAuthenticated"
          to="/login"
          class="mm-btn-primary mm-landing-cta mm-landing-header__login hidden min-h-10 !px-4 !py-2 text-sm sm:inline-flex"
        >
          Войти
        </NuxtLink>

        <NuxtLink
          v-else
          :to="appEntryPath"
          class="mm-btn-primary mm-landing-cta hidden min-h-10 !px-4 !py-2 text-sm sm:inline-flex"
          @click="closeMobileNav"
        >
          {{ needsOnboarding() ? 'Пройти опрос' : 'К приложению' }}
        </NuxtLink>

        <NuxtLink
          v-if="!headerAuthReady || !headerAuthenticated"
          to="/login"
          class="mm-landing-header__login-icon mm-landing-brand-text flex h-10 w-10 items-center justify-center rounded-full border border-[color:var(--mm-border)] transition hover:bg-[color:var(--mm-landing-brand-soft)] sm:hidden"
          aria-label="Войти"
        >
          <User class="h-4 w-4" stroke-width="1.75" />
        </NuxtLink>

        <button
          v-if="(headerAuthReady && headerAuthenticated) || isLanding"
          type="button"
          class="mm-landing-header__menu-btn flex h-10 w-10 items-center justify-center rounded-full border border-[color:var(--mm-border)] text-[color:var(--mm-text-muted)] transition hover:bg-[color:var(--mm-bg-muted)] sm:hidden"
          :aria-expanded="mobileNavOpen"
          aria-controls="landing-mobile-nav"
          :aria-label="mobileNavOpen ? 'Закрыть меню' : 'Открыть меню'"
          @click="toggleMobileNav"
        >
          <X v-if="mobileNavOpen" class="h-5 w-5" stroke-width="1.75" />
          <Menu v-else class="h-5 w-5" stroke-width="1.75" />
        </button>
      </div>
    </div>

    <Transition name="landing-mobile-nav">
      <div
        v-if="mobileNavOpen"
        id="landing-mobile-nav"
        class="mm-landing-header__drawer sm:hidden"
      >
        <nav class="mm-landing-header__drawer-inner mx-auto max-w-[1400px] px-4 pb-5 pt-2">
          <template v-if="headerAuthReady && headerAuthenticated">
            <p class="mm-landing-header__drawer-label">Приложение</p>
            <NuxtLink
              v-for="link in appNavLinks"
              :key="link.label"
              :to="link.to"
              class="mm-landing-header__drawer-link"
              @click="closeMobileNav"
            >
              {{ link.label }}
            </NuxtLink>
          </template>

          <template v-else-if="isLanding">
            <p class="mm-landing-header__drawer-label">Разделы</p>
            <button
              v-for="link in LANDING_ANCHORS"
              :key="link.id"
              type="button"
              class="mm-landing-header__drawer-link w-full text-left"
              @click="goToAnchor(link.id)"
            >
              {{ link.label }}
            </button>

            <div class="mt-4 border-t border-[color:var(--mm-border-subtle)] pt-4">
              <NuxtLink
                to="/login"
                class="mm-btn-primary mm-landing-cta flex min-h-12 w-full items-center justify-center gap-2 text-sm"
                @click="closeMobileNav"
              >
                Войти в Поток
                <ArrowRight class="h-4 w-4" stroke-width="2" />
              </NuxtLink>
            </div>
          </template>
        </nav>
      </div>
    </Transition>
  </header>

  <div
    v-if="mobileNavOpen"
    class="mm-landing-header__backdrop fixed inset-0 z-20 sm:hidden"
    aria-hidden="true"
    @click="closeMobileNav"
  />
</template>
