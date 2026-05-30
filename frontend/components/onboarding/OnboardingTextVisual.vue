<script setup lang="ts">
defineProps<{
  compact?: boolean
}>()
</script>

<template>
  <div
    class="mm-onb-text-visual"
    :class="{ 'mm-onb-text-visual--compact': compact }"
    aria-hidden="true"
  >
    <div class="mm-onb-text-visual__line mm-onb-text-visual__line--a" />
    <div class="mm-onb-text-visual__line mm-onb-text-visual__line--b" />
    <div class="mm-onb-text-visual__line mm-onb-text-visual__line--c" />
    <div class="mm-onb-text-visual__cursor" />
  </div>
</template>
