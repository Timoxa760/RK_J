<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'hero' | 'brand'
    brandAlign?: 'stretch' | 'center'
  }>(),
  { variant: 'hero', brandAlign: 'stretch' }
)
</script>

<template>
  <ClientOnly>
    <SharedHeroFlowWordGraphic :variant="variant" :brand-align="brandAlign" />
    <template #fallback>
      <span
        class="mm-hero-flow-word mm-hero-flow-word--fallback"
        :class="{ 'mm-hero-flow-word--brand': variant === 'brand' }"
        aria-label="Поток"
      >
        ПОТОК
      </span>
    </template>
  </ClientOnly>
</template>
