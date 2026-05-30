<script setup lang="ts">
import { ArrowRight } from 'lucide-vue-next'
import { ctaContent } from '~/constants/landingContent'
</script>

<template>
  <LandingReveal>
    <section class="mm-landing-section mm-landing-section--panel">
      <div class="mm-landing-section__inner">
        <div class="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between sm:gap-8">
          <div class="max-w-xl">
            <p class="mm-landing-section__eyebrow">
              Следующий шаг
            </p>
            <h2 class="mm-landing-section__title mt-3">
              {{ ctaContent.title }}
            </h2>
            <p class="mm-landing-section__lead mt-3">
              {{ ctaContent.description }}
            </p>
          </div>

          <NuxtLink
            to="/login"
            class="mm-btn-primary mm-landing-cta group min-h-12 w-full shrink-0 sm:w-auto"
          >
            {{ ctaContent.button }}
            <ArrowRight class="h-4 w-4 transition group-hover:translate-x-0.5" stroke-width="2" />
          </NuxtLink>
        </div>
      </div>
    </section>
  </LandingReveal>
</template>
