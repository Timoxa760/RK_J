<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  showBack?: boolean
  hideNext?: boolean
  nextLabel?: string
  nextDisabled?: boolean
  loading?: boolean
  secondaryAction?: { label: string } | null
  skipHint?: string
  /** Кнопки действий на всю ширину карточки (экран приветствия). */
  stretchActions?: boolean
}>()

const emit = defineEmits<{
  back: []
  next: []
  secondary: []
}>()
</script>

<template>
  <article class="mm-onboarding-card">
    <header class="space-y-1.5 lg:space-y-2">
      <h2 class="text-xl font-semibold tracking-tight text-[color:var(--mm-text)] sm:text-2xl">
        {{ title }}
      </h2>
      <p v-if="description" class="text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
        {{ description }}
      </p>
    </header>

    <div class="mt-4 space-y-4 sm:mt-5 sm:space-y-5 lg:mt-4 lg:space-y-4">
      <slot />

      <p
        v-if="skipHint"
        class="text-sm leading-relaxed text-[color:var(--mm-text-muted)]"
      >
        {{ skipHint }}
      </p>

      <div
        class="mm-onboarding-step-footer mm-safe-bottom flex gap-2 border-t border-border/60 pt-4"
        :class="
          stretchActions
            ? 'flex-col'
            : 'flex-col-reverse sm:flex-row sm:items-center sm:justify-between'
        "
      >
        <Button
          v-if="showBack"
          type="button"
          variant="ghost"
          class="min-h-11 w-full sm:w-auto"
          @click="emit('back')"
        >
          Назад
        </Button>
        <div v-else-if="!stretchActions" class="hidden sm:block" />

        <div
          class="flex w-full flex-col gap-2"
          :class="stretchActions ? '' : 'sm:w-auto sm:flex-row'"
        >
          <Button
            v-if="secondaryAction"
            type="button"
            variant="outline"
            class="min-h-11 w-full"
            :class="stretchActions ? '' : 'sm:w-auto'"
            @click="emit('secondary')"
          >
            {{ secondaryAction.label }}
          </Button>
          <Button
            v-if="!hideNext"
            type="button"
            class="min-h-11 w-full"
            :class="stretchActions ? '' : 'sm:w-auto'"
            :disabled="nextDisabled || loading"
            @click="emit('next')"
          >
            {{ loading ? 'Сохраняем…' : nextLabel ?? 'Далее' }}
          </Button>
        </div>
      </div>
    </div>
  </article>
</template>
