<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    value: string
    tone?: 'default' | 'accent' | 'warn' | 'primary' | 'good'
    valueUppercase?: boolean
  }>(),
  {
    tone: 'default',
    valueUppercase: false
  }
)

const toneClass = computed(() => `mm-key-strip--${props.tone}`)
</script>

<template>
  <div class="mm-key-strip" :class="toneClass">
    <span class="mm-key-strip__label">{{ label }}</span>
    <p class="mm-key-strip__value" :class="{ 'mm-key-strip__value--upper': valueUppercase }">
      {{ value }}
    </p>
  </div>
</template>
