<script setup lang="ts">
import { Goal, LineChart, Zap } from 'lucide-vue-next'
import { ONBOARDING } from '~/constants/productCopy'

const emit = defineEmits<{
  start: []
}>()

const highlights = [
  { icon: LineChart, text: ONBOARDING.welcomeBullet1 },
  { icon: Goal, text: ONBOARDING.welcomeBullet2 },
  { icon: Zap, text: ONBOARDING.welcomeBullet3 }
] as const
</script>

<template>
  <OnboardingStepShell
    title="Поток покажет, как ваши деньги работают на цели"
    description="Несколько простых шагов — и вы увидите свою картину и следующий шаг."
    hide-back
    next-label="Начать"
    @next="emit('start')"
  >
    <ul class="mm-onb-welcome-highlights">
      <li
        v-for="(item, index) in highlights"
        :key="index"
        class="mm-onb-welcome-highlights__item"
      >
        <span class="mm-onb-welcome-highlights__icon" aria-hidden="true">
          <component :is="item.icon" class="size-[1.125rem]" />
        </span>
        <p class="mm-onb-welcome-highlights__text">
          {{ item.text }}
        </p>
      </li>
    </ul>
  </OnboardingStepShell>
</template>
