<script setup lang="ts">
import { FLOW_LINE_PATHS } from '~/constants/flowLinePaths'

const props = withDefaults(
  defineProps<{
    variant?: 'hero' | 'brand'
    brandAlign?: 'stretch' | 'center'
  }>(),
  { variant: 'hero', brandAlign: 'stretch' }
)

const preserveAspect = computed(() => {
  if (props.variant !== 'brand') return 'xMinYMin meet'
  return props.brandAlign === 'center' ? 'xMidYMid meet' : 'none'
})

const patternId = computed(() =>
  props.variant === 'brand' ? 'mm-hero-flow-pattern-brand' : 'mm-hero-flow-pattern'
)

const pulseIndices = new Set([1, 7, 13, 19, 25])
const viewBox = '0 296 530 134'
</script>

<template>
  <svg
    class="mm-hero-flow-word"
    :class="{ 'mm-hero-flow-word--brand': variant === 'brand' }"
    :viewBox="viewBox"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    :preserveAspectRatio="preserveAspect"
    role="img"
    aria-label="Поток"
  >
    <defs>
      <pattern :id="patternId" patternUnits="userSpaceOnUse" width="422" height="596">
        <rect
          class="mm-hero-flow-word__fill"
          width="422"
          height="596"
          fill="var(--mm-hero-flow-fill, #96c8f9)"
        />
        <g class="mm-hero-flow-word__pattern">
          <path
            v-for="(d, i) in FLOW_LINE_PATHS"
            :key="`pat-${i}`"
            :d="d"
            class="mm-hero-flow-word__line"
            :class="{ 'mm-hero-flow-word__line--pulse': pulseIndices.has(i) }"
          />
        </g>
      </pattern>
    </defs>

    <text
      class="mm-hero-flow-word__text"
      x="0"
      y="404"
      font-size="106"
      font-weight="700"
      letter-spacing="-0.02em"
      :fill="`url(#${patternId})`"
    >
      ПОТОК
    </text>
  </svg>
</template>
