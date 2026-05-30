<script setup lang="ts">
import { Keyboard, Mic } from 'lucide-vue-next'

const emit = defineEmits<{
  back: []
  selectText: []
  selectVoice: []
}>()
</script>

<template>
  <OnboardingStepShell
    title="Как пройти опрос?"
    description="Сначала соберём доходы, запас и цель. Выберите способ один раз — дальше все вопросы будут в этом формате."
    show-back
    hide-next
    @back="emit('back')"
  >
    <div class="grid items-stretch gap-4 sm:grid-cols-2">
      <OnboardingModeCard
        recommended
        badge="Рекомендуем"
        title="Голосом"
        description="Один вопрос за раз — ответ вслух, без полей для ввода."
        :icon="Mic"
        @select="emit('selectVoice')"
      >
        <div class="mm-onb-mode-preview">
          <div class="mm-onb-mode-preview__media">
            <OnboardingMicOrbVisual compact ambient />
          </div>
          <p class="mm-onb-mode-preview__caption">Голосом — быстрее</p>
        </div>
      </OnboardingModeCard>

      <OnboardingModeCard
        title="Текстом"
        description="Пошаговые поля — если хотите ввести цифры вручную."
        :icon="Keyboard"
        @select="emit('selectText')"
      >
        <div class="mm-onb-mode-preview">
          <div class="mm-onb-mode-preview__media">
            <OnboardingTextVisual compact />
          </div>
          <p class="mm-onb-mode-preview__caption">Текстом — по полям</p>
        </div>
      </OnboardingModeCard>
    </div>
  </OnboardingStepShell>
</template>
