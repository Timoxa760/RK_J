<script setup lang="ts">
import type { PageNarrativeBlock } from '~/utils/pageNarrative'
import type { HealthTone } from '~/utils/dashboardSummary'
import { ADVISOR, GOALS } from '~/constants/productCopy'

const props = defineProps<{
  narrative?: PageNarrativeBlock | null
  headline?: string
  paragraphs?: string[]
  healthEmoji?: PageNarrativeBlock['healthEmoji']
  healthTone?: HealthTone
  badgeLabel?: string
  weeklyAction?: string
  adviceHint?: string
  callout?: string
  loading?: boolean
}>()

const resolved = computed(() => {
  const n = props.narrative
  return {
    headline: props.headline ?? n?.headline ?? '',
    paragraphs: props.paragraphs ?? n?.paragraphs ?? [],
    contextFacts: n?.contextFacts ?? [],
    goalOpportunityThousands: n?.goalOpportunityThousands ?? null,
    healthEmoji: props.healthEmoji ?? n?.healthEmoji,
    healthTone: props.healthTone ?? n?.healthTone,
    badgeLabel: props.badgeLabel ?? n?.badgeLabel,
    weeklyAction: props.weeklyAction ?? n?.weeklyAction,
    adviceHint: props.adviceHint ?? n?.adviceHint ?? ADVISOR.weeklyAdviceHintShort,
    callout: props.callout ?? n?.callout
  }
})

const isHero = computed(
  () =>
    Boolean(
      resolved.value.weeklyAction ||
        resolved.value.contextFacts.length
    )
)
</script>

<template>
  <section class="w-full" aria-label="Кратко о странице" data-demo="narrative">
    <div v-if="$slots.actions" class="mb-3 flex justify-center">
      <slot name="actions" />
    </div>

    <Skeleton v-if="loading" class="h-32 w-full rounded-2xl" />

    <section
      v-else-if="resolved.headline || resolved.paragraphs.length || isHero"
      class="mm-narrative-hero space-y-3"
    >
      <template v-if="isHero">
        <aside
          v-if="resolved.weeklyAction"
          class="mm-tier-1 mm-narrative-hero__advice"
          aria-label="Совет недели"
        >
          <div class="mm-narrative-hero__advice-meta">
            <span class="mm-narrative-hero__advice-badge">{{ ADVISOR.weeklyAdviceTitle }}</span>
            <span
              v-if="resolved.goalOpportunityThousands"
              class="mm-narrative-hero__advice-hook"
            >
              {{ GOALS.opportunityAmount(resolved.goalOpportunityThousands) }}
            </span>
          </div>
          <p class="mm-narrative-hero__advice-text">{{ resolved.weeklyAction }}</p>
          <p class="mm-narrative-hero__advice-hint">{{ resolved.adviceHint }}</p>
        </aside>

        <div
          v-if="resolved.contextFacts.length"
          class="mm-narrative-hero__grid"
          aria-label="Ключевые цифры"
        >
          <div
            v-for="fact in resolved.contextFacts"
            :key="fact.id"
            class="mm-narrative-hero__tile"
            :class="{
              'mm-narrative-hero__tile--accent': fact.tone === 'accent',
              'mm-narrative-hero__tile--warn': fact.tone === 'warn'
            }"
          >
            <span class="mm-narrative-hero__tile-label">{{ fact.label }}</span>
            <span class="mm-narrative-hero__tile-value">{{ fact.value }}</span>
          </div>
        </div>
      </template>

      <template v-else>
        <div class="mm-tier-2 rounded-xl border bg-card p-4 sm:p-5">
          <p v-if="resolved.headline" class="text-xl font-semibold leading-snug">
            {{ resolved.headline }}
          </p>
          <p
            v-for="(paragraph, i) in resolved.paragraphs"
            :key="i"
            class="mt-2 text-base leading-relaxed text-muted-foreground"
          >
            {{ paragraph }}
          </p>
        </div>
      </template>

      <slot />
    </section>
  </section>
</template>
