<script setup lang="ts">
import { benefits } from '~/constants/landingContent'
import { landingIcons } from '~/components/landing/landingIcons'
</script>

<template>
  <LandingReveal>
    <section id="value" class="mm-landing-section mm-landing-section--panel">
      <div class="mm-landing-section__inner">
      <p class="mm-landing-section__eyebrow">Что даёт Поток</p>
      <h2 class="mm-landing-section__title mt-3">Прозрачность без перегруза</h2>
      <p class="mm-landing-section__lead mt-3 max-w-xl">
        Не учёт ради учёта — понятный разбор и одно действие, которое имеет смысл сейчас.
      </p>

      <div class="mm-landing-stagger mt-6 grid gap-4 sm:mt-8 sm:grid-cols-3 sm:gap-5">
        <article
          v-for="item in benefits"
          :key="item.title"
          class="mm-landing-card mm-landing-card--accent mm-landing-card--interactive group flex flex-col p-5 sm:p-6"
        >
          <div
            class="flex h-11 w-11 items-center justify-center rounded-2xl bg-gradient-to-br from-[color:var(--mm-primary-soft)] to-white text-[color:var(--mm-primary)] shadow-sm transition-transform duration-300 group-hover:scale-110 group-hover:rotate-3"
          >
            <component :is="landingIcons[item.icon]" class="h-5 w-5" stroke-width="1.75" />
          </div>
          <h3 class="mt-4 font-bold text-[color:var(--mm-text)]">{{ item.title }}</h3>
          <p class="mt-2 flex-1 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
            {{ item.description }}
          </p>
        </article>
      </div>
      </div>
    </section>
  </LandingReveal>
</template>
