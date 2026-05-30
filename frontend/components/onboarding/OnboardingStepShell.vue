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
}>()

const emit = defineEmits<{
  back: []
  next: []
  secondary: []
}>()
</script>

<template>
  <article class="mm-onboarding-card w-full p-6 sm:p-8">
    <header class="space-y-2">
      <h2 class="text-xl font-semibold tracking-tight text-[color:var(--mm-text)] sm:text-2xl">
        {{ title }}
      </h2>
      <p v-if="description" class="text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
        {{ description }}
      </p>
    </header>

    <div class="mt-6 space-y-6">
      <slot />

      <div class="flex flex-col-reverse gap-2 border-t border-border/60 pt-5 sm:flex-row sm:items-center sm:justify-between">
        <Button
          v-if="showBack"
          type="button"
          variant="ghost"
          class="w-full sm:w-auto"
          @click="emit('back')"
        >
          Назад
        </Button>
        <div v-else class="hidden sm:block" />

        <div class="flex w-full flex-col gap-2 sm:w-auto sm:flex-row">
          <Button
            v-if="secondaryAction"
            type="button"
            variant="outline"
            class="w-full sm:w-auto"
            @click="emit('secondary')"
          >
            {{ secondaryAction.label }}
          </Button>
          <Button
            v-if="!hideNext"
            type="button"
            class="w-full sm:w-auto"
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
