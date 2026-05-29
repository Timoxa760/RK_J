<script setup lang="ts">
import { proofStats } from '~/constants/landingContent'
</script>

<template>
  <LandingReveal>
    <section class="mm-landing-section mm-landing-section--tight">
      <div class="mm-landing-stagger mm-landing-proof flex flex-col gap-3 sm:flex-row sm:items-stretch sm:justify-center sm:gap-4">
        <div
          v-for="(stat, index) in proofStats"
          :key="stat.value"
          class="mm-landing-proof__pill group flex flex-1 flex-col items-center rounded-2xl px-5 py-4 text-center sm:max-w-xs"
          :style="{ '--stagger-i': index }"
        >
          <span class="mm-landing-proof__value">{{ stat.value }}</span>
          <span class="mt-1 text-xs text-[color:var(--mm-text-muted)] sm:text-sm">{{ stat.label }}</span>
        </div>
      </div>
    </section>
  </LandingReveal>
</template>
