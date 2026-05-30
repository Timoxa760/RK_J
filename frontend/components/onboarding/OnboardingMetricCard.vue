<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  label: string
  icon?: Component
  accent?: boolean
}>()
</script>

<template>
  <div class="mm-onb-metric" :class="{ 'mm-onb-metric--accent': accent }">
    <div class="mm-onb-metric__label" :class="{ 'mm-onb-metric__label--accent': accent }">
      <component v-if="icon" :is="icon" class="size-4" aria-hidden="true" />
      <span>{{ label }}</span>
    </div>
    <div class="mm-onb-metric__value">
      <slot />
    </div>
    <p v-if="$slots.hint" class="mm-onb-metric__hint">
      <slot name="hint" />
    </p>
  </div>
</template>
