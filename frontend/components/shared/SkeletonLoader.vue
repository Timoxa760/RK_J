<script setup lang="ts">
defineProps<{
  class?: string
  height?: string
}>()
</script>

<template>
  <div
    class="animate-pulse rounded-xl bg-[color:var(--mm-bg-muted)]"
    :class="$props.class"
    :style="{ minHeight: height || '300px' }"
  />
</template>
