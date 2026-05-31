<script setup lang="ts">
import { ArrowRight } from 'lucide-vue-next'
import { heroContent, AUTH_REGISTER_PATH } from '~/constants/landingContent'
</script>

<template>
  <LandingReveal immediate>
    <section class="mm-landing-section mm-landing-section--hero mm-landing-hero">
      <div class="grid gap-8 lg:grid-cols-2 lg:items-center lg:gap-10 xl:gap-12">
        <div class="mm-landing-hero__copy mx-auto w-full max-w-xl min-w-0 lg:mx-0 mm-landing-hero-cascade">
        <p class="mm-landing-section__eyebrow">
          {{ heroContent.eyebrow }}
        </p>

        <h1 class="mm-hero-title mt-4 text-[color:var(--mm-text)] sm:mt-5">
          <span class="block">Расходы текут —</span>
          <span class="block mm-hero-title__flow-line">
            <SharedHeroFlowWord />
          </span>
          <span class="block mm-hero-title__accent">их разберёт</span>
        </h1>

        <div class="mm-landing-hero__lead-box">
          <p class="mm-landing-hero__lead">
            {{ heroContent.subtitle }}
          </p>
          <p class="mm-landing-hero__lead-note">
            {{ heroContent.subtitleNote }}
          </p>
        </div>

        <div class="mt-8">
          <NuxtLink :to="AUTH_REGISTER_PATH" class="mm-btn-primary mm-landing-cta mm-landing-cta--pulse group min-h-12 w-full sm:w-auto">
            Начать бесплатно
            <ArrowRight class="h-4 w-4 transition group-hover:translate-x-0.5" stroke-width="2" />
          </NuxtLink>
        </div>
      </div>

        <div class="mm-landing-hero__demo mx-auto w-full max-w-xl min-w-0 lg:mx-0 lg:max-w-none">
          <ClientOnly>
            <LandingVoiceDemo />
            <template #fallback>
              <div
                class="mm-landing-demo mm-landing-glass mm-landing-card min-h-[26.5rem] w-full animate-pulse rounded-2xl bg-white/50"
                aria-hidden="true"
              />
            </template>
          </ClientOnly>
        </div>
      </div>
    </section>
  </LandingReveal>
</template>
