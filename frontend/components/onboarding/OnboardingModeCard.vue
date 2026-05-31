<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  title: string
  description: string
  recommended?: boolean
  badge?: string
  icon: Component
}>()

const emit = defineEmits<{
  select: []
}>()
</script>

<template>
  <button
    type="button"
    class="mm-onb-mode-card flex h-full w-full flex-col text-left sm:min-h-[17.5rem]"
    :class="{ 'mm-onb-mode-card--recommended': recommended }"
    @click="emit('select')"
  >
    <span v-if="badge" class="mm-onb-mode-card__badge">{{ badge }}</span>
    <div class="mm-onb-mode-card__icon">
      <component :is="icon" class="size-6" aria-hidden="true" />
    </div>
    <p class="mt-3 text-base font-semibold text-[color:var(--mm-text)]">{{ title }}</p>
    <p class="mt-1.5 flex-1 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
      {{ description }}
    </p>
    <div class="mm-onb-mode-card__visual">
      <slot />
    </div>
  </button>
</template>
