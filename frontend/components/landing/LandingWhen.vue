<script setup lang="ts">
import { moments } from '~/constants/landingContent'
import { landingIcons } from '~/components/landing/landingIcons'
</script>

<template>
  <LandingReveal>
    <section id="when" class="mm-landing-section mm-landing-section--panel">
      <div class="mm-landing-section__inner">
        <p class="mm-landing-section__eyebrow">Когда включать</p>
        <h2 class="mm-landing-section__title mt-3">В вашем ритме — не по расписанию бухгалтера</h2>
        <p class="mm-landing-section__lead mt-3 max-w-xl">
          Там, где поток трат уже пошёл, а разобраться некогда.
        </p>

        <ol class="mm-landing-stagger mm-landing-timeline mt-6 space-y-0 sm:mt-8 sm:hidden">
          <li
            v-for="(item, index) in moments"
            :key="item.title"
            class="mm-landing-timeline__item relative flex gap-4 pb-4 pl-1"
          >
            <div class="flex flex-col items-center">
              <span class="mm-landing-timeline__dot mm-landing-timeline__dot--pulse">
                <component :is="landingIcons[item.icon]" class="h-3.5 w-3.5" stroke-width="2" />
              </span>
              <span
                v-if="index < moments.length - 1"
                class="mm-landing-timeline__line mt-2 flex-1"
                aria-hidden="true"
              />
            </div>
            <div class="mm-landing-card mm-landing-card--interactive min-w-0 flex-1 p-4">
              <h3 class="font-bold text-[color:var(--mm-text)]">{{ item.title }}</h3>
              <p class="mt-1.5 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
                {{ item.description }}
              </p>
            </div>
          </li>
        </ol>

        <div class="mt-6 hidden sm:mt-8 sm:block">
          <div class="mm-landing-timeline-h-track grid grid-cols-3 gap-4">
            <div
              v-for="item in moments"
              :key="`${item.title}-dot`"
              class="mm-landing-timeline-h-track__cell relative flex justify-center"
            >
              <span class="mm-landing-timeline__dot mm-landing-timeline__dot--pulse relative z-[1]">
                <component :is="landingIcons[item.icon]" class="h-3.5 w-3.5" stroke-width="2" />
              </span>
            </div>
          </div>

          <ol class="mm-landing-stagger mt-4 grid grid-cols-3 gap-4">
            <li
              v-for="item in moments"
              :key="item.title"
              class="mm-landing-card mm-landing-card--interactive p-5 sm:p-6"
            >
              <h3 class="font-bold text-[color:var(--mm-text)]">{{ item.title }}</h3>
              <p class="mt-2 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
                {{ item.description }}
              </p>
            </li>
          </ol>
        </div>
      </div>
    </section>
  </LandingReveal>
</template>
