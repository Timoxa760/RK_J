<script setup lang="ts">
import { Keyboard, Mic, Sparkles } from 'lucide-vue-next'

const emit = defineEmits<{
  selectText: []
  selectVoice: []
}>()
</script>

<template>
  <section class="space-y-8">
    <header class="flex items-start gap-3">
      <div class="mm-onb-mode-card__icon shrink-0">
        <Sparkles class="size-6" aria-hidden="true" />
      </div>
      <div>
        <h2 class="text-xl font-semibold tracking-tight text-[color:var(--mm-text)] sm:text-2xl">
          Как пройти опрос?
        </h2>
        <p class="mt-2 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
          Сначала соберём доходы, подушку и цель — вы сами расскажете текстом или голосом. Чеки из
          ФНС подключим позже: они дают только покупки, не зарплату.
        </p>
      </div>
    </header>

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
          <p class="mm-onb-mode-preview__caption">
            Голосом — быстрее
          </p>
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
          <p class="mm-onb-mode-preview__caption">
            Текстом — по полям
          </p>
        </div>
      </OnboardingModeCard>
    </div>

    <p class="text-center text-xs text-[color:var(--mm-text-muted)]">
      Можно сменить способ на следующем шаге — данные сохранятся локально.
    </p>
  </section>
</template>
