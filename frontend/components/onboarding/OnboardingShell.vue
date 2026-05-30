<script setup lang="ts">
import { ONBOARDING } from '~/constants/productCopy'

defineProps<{
  progress?: { current: number; total: number; label?: string } | null
  showProgress?: boolean
  stepLabels?: string[]
  /** На шаге «Старт» — только бейдж, без второго заголовка в шапке */
  minimalHeader?: boolean
}>()
</script>

<template>
  <div class="mm-onboarding relative mx-auto w-full max-w-xl">
    <div class="mm-onboarding__orb mm-onboarding__orb--a" aria-hidden="true" />
    <div class="mm-onboarding__orb mm-onboarding__orb--b" aria-hidden="true" />

    <header class="relative z-10 mb-6 text-center">
      <p
        class="inline-flex items-center gap-2 rounded-full border border-[color:var(--mm-primary)]/20 bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] text-[color:var(--mm-primary)] backdrop-blur-sm"
      >
        <span class="size-1.5 animate-pulse rounded-full bg-[color:var(--mm-primary)]" />
        Поток
      </p>
      <template v-if="!minimalHeader">
        <h1 class="mt-3 font-[family-name:var(--mm-font-sf-display)] text-2xl font-semibold tracking-tight sm:text-3xl">
          {{ ONBOARDING.shellTitle }}
        </h1>
        <p class="mt-2 text-sm text-[color:var(--mm-text-muted)]">
          {{ ONBOARDING.shellSubtitle }}
        </p>
      </template>
    </header>

    <OnboardingProgress
      v-if="showProgress && progress"
      class="relative z-10 mb-5"
      :current="progress.current"
      :total="progress.total"
      :label="progress.label"
      :step-labels="stepLabels"
    />

    <div class="relative z-20 w-full">
      <slot />
    </div>
  </div>
</template>
