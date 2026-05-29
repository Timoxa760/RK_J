<script setup lang="ts">
import { howSteps } from '~/constants/landingContent'
import { landingIcons } from '~/components/landing/landingIcons'
</script>

<template>
  <LandingReveal>
    <section id="how" class="mm-landing-section mm-landing-section--panel">
      <div class="mm-landing-section__inner">
      <p class="mm-landing-section__eyebrow">Как это работает</p>
      <h2 class="mm-landing-section__title mt-3">Три шага — и поток под контролем</h2>
      <p class="mm-landing-section__lead mt-3 max-w-xl">
        Без таблиц и бухгалтерии. Вы говорите — Поток разбирает и предлагает одну правку.
      </p>

      <div class="mm-landing-stagger mm-landing-steps mt-6 grid gap-4 sm:mt-8 sm:grid-cols-3 sm:gap-5">
        <article
          v-for="step in howSteps"
          :key="step.num"
          class="mm-landing-card mm-landing-card--step mm-landing-card--interactive group relative overflow-hidden p-5 sm:p-6"
        >
          <div class="relative z-[1]">
            <div class="flex items-center justify-between gap-3">
              <span class="mm-landing-step-num">{{ step.num }}</span>
              <component
                :is="landingIcons[step.icon]"
                class="h-5 w-5 text-[color:var(--mm-primary)] transition-transform duration-300 group-hover:scale-110"
                stroke-width="1.75"
              />
            </div>
            <h3 class="mt-4 font-bold text-[color:var(--mm-text)]">{{ step.title }}</h3>
            <p class="mt-2 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
              {{ step.description }}
            </p>
          </div>
          <span class="mm-landing-step-watermark" aria-hidden="true">{{ step.num }}</span>
        </article>
      </div>
      </div>
    </section>
  </LandingReveal>
</template>
