<script setup lang="ts">
import type { AdvisorReplyBlock } from '~/types/api'

defineProps<{
  title?: string
  blocks: AdvisorReplyBlock[]
  streaming?: boolean
}>()
</script>

<template>
  <article class="advisor-structured">
    <h3 v-if="title" class="advisor-structured__title">
      {{ title }}
    </h3>

    <template v-for="(block, index) in blocks" :key="`${block.type}-${index}`">
      <p v-if="block.type === 'lead'" class="advisor-structured__lead">
        {{ block.text }}
      </p>

      <h4 v-else-if="block.type === 'heading'" class="advisor-structured__heading">
        {{ block.text }}
      </h4>

      <p v-else-if="block.type === 'paragraph'" class="advisor-structured__paragraph">
        {{ block.text }}
      </p>

      <ul v-else-if="block.type === 'list'" class="advisor-structured__list">
        <li v-for="(item, i) in block.items" :key="`${index}-${i}`">
          {{ item }}
        </li>
      </ul>

      <div
        v-else-if="block.type === 'callout'"
        class="advisor-structured__callout"
        :class="
          block.tone === 'info'
            ? 'advisor-structured__callout--info'
            : 'advisor-structured__callout--action'
        "
      >
        {{ block.text }}
      </div>
    </template>

    <span v-if="streaming" class="advisor-markdown-cursor" aria-hidden="true" />
  </article>
</template>

<style scoped>
.advisor-structured__title {
  @apply mb-3 text-lg font-bold leading-snug tracking-tight;
}

.advisor-structured__lead {
  @apply text-base font-medium leading-relaxed;
}

.advisor-structured__heading {
  @apply mt-4 text-base font-semibold leading-snug first:mt-0;
}

.advisor-structured__paragraph {
  @apply mt-2 text-base leading-relaxed text-foreground first:mt-0;
}

.advisor-structured__list {
  @apply mt-2 list-disc space-y-1.5 pl-5 text-base leading-relaxed first:mt-0;
}

.advisor-structured__callout {
  @apply mt-3 rounded-xl border px-3.5 py-3 text-sm leading-relaxed first:mt-0;
}

.advisor-structured__callout--action {
  @apply font-medium;
  border-color: color-mix(in srgb, var(--mm-primary) 30%, transparent);
  background-color: color-mix(in srgb, var(--mm-primary-soft) 70%, transparent);
}

.advisor-structured__callout--info {
  @apply border-border bg-muted/40;
}

.advisor-markdown-cursor {
  display: inline-block;
  width: 2px;
  height: 1em;
  margin-left: 2px;
  vertical-align: text-bottom;
  background: currentColor;
  opacity: 0.45;
  animation: advisor-cursor-blink 1s step-end infinite;
}

@keyframes advisor-cursor-blink {
  50% {
    opacity: 0;
  }
}
</style>
