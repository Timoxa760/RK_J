<script setup lang="ts">
import { FLOW_LINE_PATHS } from '~/constants/flowLinePaths'
import type { FlowPulseAnim } from '~/constants/backgroundFlow'
import {
  FLOW_ACCENT_INDICES,
  FLOW_LINE_INDICES,
  FLOW_STAGE_GROUP,
  FLOW_STAGE_VIEWBOX,
  flowAnimVars
} from '~/constants/backgroundFlow'

const props = defineProps<{
  tracks?: boolean
  accent?: boolean
  pulses?: FlowPulseAnim[]
}>()

const hasTracks = computed(() => Boolean(props.tracks))
const hasAnim = computed(() => Boolean(props.pulses?.length))
const isAnimOnly = computed(() => hasAnim.value && !hasTracks.value)
</script>

<template>
  <svg
    v-if="hasTracks || hasAnim"
    class="mm-flow__stage"
    :class="{ 'mm-flow__stage--anim': isAnimOnly }"
    :viewBox="FLOW_STAGE_VIEWBOX"
    preserveAspectRatio="xMidYMid slice"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <g :transform="FLOW_STAGE_GROUP">
      <template v-if="hasTracks">
        <path
          v-for="i in FLOW_LINE_INDICES"
          :key="`track-${i}`"
          :d="FLOW_LINE_PATHS[i]"
          class="mm-flow__track"
        />
        <path
          v-for="i in accent ? FLOW_ACCENT_INDICES : []"
          :key="`accent-${i}`"
          :d="FLOW_LINE_PATHS[i]"
          class="mm-flow__track mm-flow__track--accent"
        />
      </template>

      <template v-if="hasAnim">
        <path
          v-for="pulse in pulses"
          :key="`pulse-${pulse.index}`"
          :d="FLOW_LINE_PATHS[pulse.index]"
          class="mm-flow__pulse"
          :style="flowAnimVars(pulse.duration, pulse.delay)"
        />
      </template>
    </g>
  </svg>
</template>
