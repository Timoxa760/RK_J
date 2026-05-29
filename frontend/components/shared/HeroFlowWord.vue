<script setup lang="ts">
import { FLOW_LINE_PATHS } from '~/constants/flowLinePaths'

const patternId = `mm-hero-line1-${useId()}`
const pulseIndices = new Set([1, 7, 13, 19, 25])

/** Кроп вокруг «ПОТОК» в координатах line1.svg */
const viewBox = '0 298 358 114'
</script>

<template>
  <svg
    class="mm-hero-flow-word"
    :viewBox="viewBox"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    preserveAspectRatio="xMinYMin meet"
    role="img"
    aria-label="Поток"
  >
    <defs>
      <pattern
        :id="patternId"
        patternUnits="userSpaceOnUse"
        width="422"
        height="596"
      >
        <rect width="422" height="596" fill="#e8955f" />
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
