<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    immediate?: boolean
    delay?: number
    threshold?: number
  }>(),
  {
    immediate: false,
    delay: 0,
    threshold: undefined
  }
)

const { root, visible, animate } = useRevealOnScroll({
  immediate: props.immediate,
  delay: props.delay,
  threshold: props.threshold
})
</script>

<template>
  <div
    ref="root"
    class="mm-landing-reveal"
    :class="{
      'mm-landing-reveal--animate': animate,
      'mm-landing-reveal--visible': visible
    }"
  >
    <slot />
  </div>
</template>
