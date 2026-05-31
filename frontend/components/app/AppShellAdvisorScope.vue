<script setup lang="ts">
const { bootstrap } = useAdvisorShellProvider()

onMounted(async () => {
  await bootstrap()
})
</script>

<template>
  <slot />
</template>
