<script setup lang="ts">
import { ctaContent } from '~/constants/landingContent'
</script>

<template>
  <LandingReveal>
    <section class="mm-landing-section mm-landing-section--tight">
      <div class="mm-landing-cta-band mm-landing-card--glow rounded-3xl px-5 py-7 sm:px-8 sm:py-9">
        <div class="flex flex-col items-start gap-5 sm:flex-row sm:items-center sm:justify-between">
          <div class="max-w-lg">
            <h2 class="text-xl font-bold tracking-tight text-[color:var(--mm-text)] sm:text-2xl">
              {{ ctaContent.title }}
            </h2>
            <p class="mt-2 text-sm leading-relaxed text-[color:var(--mm-text-muted)] sm:text-base">
              {{ ctaContent.description }}
            </p>
          </div>
          <NuxtLink
            to="/login"
            class="mm-btn-primary mm-landing-cta mm-landing-cta--pulse min-h-12 w-full shrink-0 sm:w-auto"
          >
            {{ ctaContent.button }}
          </NuxtLink>
        </div>
      </div>
    </section>
  </LandingReveal>
</template>
