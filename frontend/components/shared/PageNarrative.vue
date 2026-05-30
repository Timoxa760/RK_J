<script setup lang="ts">
import type { PageNarrativeBlock } from '~/utils/pageNarrative'
import type { HealthTone } from '~/utils/dashboardSummary'

const props = defineProps<{
  narrative?: PageNarrativeBlock | null
  headline?: string
  paragraphs?: string[]
  healthEmoji?: PageNarrativeBlock['healthEmoji']
  healthTone?: HealthTone
  badgeLabel?: string
  weeklyAction?: string
  callout?: string
  loading?: boolean
}>()

const toneVariant: Record<HealthTone, 'default' | 'secondary' | 'destructive'> = {
  good: 'default',
  warn: 'secondary',
  risk: 'destructive'
}

const resolved = computed(() => {
  const n = props.narrative
  return {
    headline: props.headline ?? n?.headline ?? '',
    paragraphs: props.paragraphs ?? n?.paragraphs ?? [],
    healthEmoji: props.healthEmoji ?? n?.healthEmoji,
    healthTone: props.healthTone ?? n?.healthTone,
    badgeLabel: props.badgeLabel ?? n?.badgeLabel,
    weeklyAction: props.weeklyAction ?? n?.weeklyAction,
    callout: props.callout ?? n?.callout
  }
})

const borderClass = computed(() => {
  const tone = resolved.value.healthTone
  if (tone === 'risk') return 'border-l-red-500 bg-red-50/40'
  if (tone === 'warn') return 'border-l-amber-500 bg-amber-50/40'
  if (tone === 'good') return 'border-l-emerald-500 bg-emerald-50/40'
  return 'border-l-border bg-muted/30'
})
</script>

<template>
  <section class="w-full space-y-3" aria-label="Кратко о странице">
    <div v-if="$slots.actions" class="flex justify-end">
      <slot name="actions" />
    </div>

    <Skeleton v-if="loading" class="h-24 w-full" />

    <Card
      v-else-if="resolved.headline || resolved.paragraphs.length"
      class="border-l-4"
      :class="borderClass"
    >
      <CardContent class="space-y-3 py-4">
        <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
          <div class="min-w-0 space-y-2">
            <p v-if="resolved.healthEmoji" class="flex items-center gap-2 text-base font-medium">
              <span aria-hidden="true">{{ resolved.healthEmoji }}</span>
              <span>{{ resolved.headline }}</span>
            </p>
            <p v-else-if="resolved.headline" class="text-base font-medium">
              {{ resolved.headline }}
            </p>
            <p
              v-for="(paragraph, i) in resolved.paragraphs"
              :key="i"
              class="text-sm leading-relaxed text-muted-foreground"
            >
              {{ paragraph }}
            </p>
          </div>
          <Badge
            v-if="resolved.badgeLabel"
            :variant="resolved.healthTone ? toneVariant[resolved.healthTone] : 'secondary'"
            class="w-fit shrink-0"
          >
            {{ resolved.badgeLabel }}
          </Badge>
        </div>

        <p v-if="resolved.callout" class="text-xs text-muted-foreground">
          {{ resolved.callout }}
        </p>

        <div
          v-if="resolved.weeklyAction"
          class="rounded-lg border border-primary/20 bg-primary/5 px-3 py-2 text-sm"
        >
          <p class="text-xs font-medium text-primary">Что сделать на этой неделе</p>
          <p class="mt-0.5 leading-snug">{{ resolved.weeklyAction }}</p>
        </div>

        <slot />
      </CardContent>
    </Card>
  </section>
</template>
