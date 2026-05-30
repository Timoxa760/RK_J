<script setup lang="ts">
import { useMotion } from '@vueuse/motion'
import { Mic } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    listening?: boolean
    parsing?: boolean
    compact?: boolean
    mini?: boolean
    dock?: boolean
    gentle?: boolean
    ambient?: boolean
  }>(),
  {
    listening: false,
    parsing: false,
    compact: false,
    mini: false,
    dock: false,
    gentle: false,
    ambient: false
  }
)

const orbRef = ref<HTMLElement | null>(null)

const { apply } = useMotion(orbRef, {
  initial: { scale: 1 },
  enter: { scale: 1 }
})

watch(
  () => [props.listening, props.parsing, props.ambient, props.gentle, props.dock] as const,
  ([listening, parsing, ambient, gentle, dock]) => {
    if (dock) return

    const soft = gentle
      ? { type: 'spring' as const, stiffness: 110, damping: 24, mass: 0.9 }
      : { type: 'spring' as const, stiffness: 220, damping: 18 }

    if (parsing) {
      apply({
        scale: soft ? 0.98 : 0.96,
        transition: { duration: soft ? 650 : 400, ease: 'easeOut' }
      })
      return
    }
    if (listening) {
      apply({ scale: soft ? 1.05 : 1.08, transition: soft })
      return
    }
    if (ambient) {
      apply({
        scale: soft ? 1.015 : 1.03,
        transition: { duration: soft ? 900 : 600, ease: 'easeInOut' }
      })
      return
    }
    apply({ scale: 1, transition: { duration: soft ? 700 : 400, ease: 'easeOut' } })
  },
  { immediate: true }
)
</script>

<template>
  <div
    ref="orbRef"
    class="mm-onb-mic-orb-visual"
    :class="{
      'mm-onb-mic-orb-visual--compact': compact && !mini && !dock,
      'mm-onb-mic-orb-visual--mini': mini,
      'mm-onb-mic-orb-visual--dock': dock,
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
