<script setup lang="ts">
const { root, visible } = useRevealOnScroll()
</script>

<template>
  <div
    ref="root"
    class="mm-landing-reveal"
    :class="{ 'mm-landing-reveal--visible': visible }"
  >
    <slot />
  </div>
</template>
