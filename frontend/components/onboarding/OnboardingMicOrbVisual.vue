<script setup lang="ts">
import { useMotion } from '@vueuse/motion'
import { Mic } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    listening?: boolean
    parsing?: boolean
    compact?: boolean
    ambient?: boolean
  }>(),
  {
    listening: false,
    parsing: false,
    compact: false,
    ambient: false
  }
)

const orbRef = ref<HTMLElement | null>(null)

const { apply } = useMotion(orbRef, {
  initial: { scale: 1 },
  enter: { scale: 1 }
})

watch(
  () => [props.listening, props.parsing, props.ambient] as const,
  ([listening, parsing, ambient]) => {
    if (parsing) {
      apply({ scale: 0.96, transition: { duration: 400 } })
      return
    }
    if (listening) {
      apply({ scale: 1.08, transition: { type: 'spring', stiffness: 220, damping: 18 } })
      return
    }
    if (ambient) {
      apply({ scale: 1.03, transition: { duration: 600 } })
      return
    }
    apply({ scale: 1, transition: { duration: 400 } })
  },
  { immediate: true }
)
</script>

<template>
  <div
    ref="orbRef"
    class="mm-onb-mic-orb-visual"
    :class="{
      'mm-onb-mic-orb-visual--compact': compact,
      'mm-onb-mic-orb-visual--listen': listening && !parsing,
      'mm-onb-mic-orb-visual--parse': parsing,
      'mm-onb-mic-orb-visual--idle': ambient || (!listening && !parsing)
    }"
    aria-hidden="true"
  >
    <span class="mm-onb-mic-orb-visual__blob mm-onb-mic-orb-visual__blob--a" />
    <span class="mm-onb-mic-orb-visual__blob mm-onb-mic-orb-visual__blob--b" />
    <span class="mm-onb-mic-orb-visual__blob mm-onb-mic-orb-visual__blob--c" />
    <span class="mm-onb-mic-orb-visual__core">
      <Mic class="mm-onb-mic-orb-visual__icon" :stroke-width="2.25" />
    </span>
  </div>
</template>
