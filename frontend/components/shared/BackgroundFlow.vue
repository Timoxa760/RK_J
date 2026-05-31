<script setup lang="ts">
import FlowSvgLayer from '~/components/shared/background/FlowSvgLayer.vue'
import { FLOW_PULSES } from '~/constants/backgroundFlow'
</script>

<template>
  <ClientOnly>
    <div class="mm-flow" aria-hidden="true">
      <FlowSvgLayer tracks accent />
      <FlowSvgLayer :pulses="FLOW_PULSES" />
    </div>
    <template #fallback>
      <div class="mm-flow mm-flow--fallback" aria-hidden="true" />
    </template>
  </ClientOnly>
</template>
