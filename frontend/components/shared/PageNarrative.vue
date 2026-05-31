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

const slots = useSlots()

const resolved = computed(() => {
  const n = props.narrative
  return {
    headline: props.headline ?? n?.headline ?? '',
    paragraphs: props.paragraphs ?? n?.paragraphs ?? [],
    goalOpportunityThousands: n?.goalOpportunityThousands ?? null,
    healthEmoji: props.healthEmoji ?? n?.healthEmoji,
    healthTone: props.healthTone ?? n?.healthTone,
    badgeLabel: props.badgeLabel ?? n?.badgeLabel,
    weeklyAction: props.weeklyAction ?? n?.weeklyAction,
    adviceHint: props.adviceHint ?? n?.adviceHint ?? ADVISOR.weeklyAdviceHintShort,
    callout: props.callout ?? n?.callout,
    incomeDisplay: n?.incomeDisplay ?? null,
    expensesDisplay: n?.expensesDisplay ?? null,
    expensesWarn: n?.expensesWarn ?? false
  }
})

const isHero = computed(
  () =>
    Boolean(
      resolved.value.weeklyAction ||
        slots.aside ||
        resolved.value.incomeDisplay ||
        resolved.value.expensesDisplay
    )
)

const showMoneyRow = computed(
  () => Boolean(resolved.value.incomeDisplay || resolved.value.expensesDisplay)
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
      class="mm-narrative-hero"
    >
      <template v-if="isHero">
        <div class="mm-narrative-hero__top">
          <aside
            v-if="resolved.weeklyAction"
            class="mm-narrative-hero__advice"
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

          <div v-if="$slots.aside" class="mm-narrative-hero__aside">
            <slot name="aside" />
          </div>
        </div>

        <div v-if="showMoneyRow" class="mm-narrative-hero__money-row" aria-label="Доход и траты">
          <div v-if="resolved.incomeDisplay" class="mm-narrative-hero__money-card">
            <span class="mm-narrative-hero__money-label">Доход</span>
            <span class="mm-narrative-hero__money-value">{{ resolved.incomeDisplay }}</span>
          </div>
          <div v-if="resolved.expensesDisplay" class="mm-narrative-hero__money-card">
            <span class="mm-narrative-hero__money-label">Траты</span>
            <span
              class="mm-narrative-hero__money-value"
              :class="{ 'text-amber-800': resolved.expensesWarn }"
            >
              {{ resolved.expensesDisplay }}
            </span>
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
