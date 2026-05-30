<script setup lang="ts">
import { Sparkles, TrendingUp } from 'lucide-vue-next'
import type { CategoriesResponse, FinancialProfile } from '~/types/api'
import { SCENARIO_OPTIONS } from '~/types/api'
import { formatRub } from '~/constants/productCopy'
import { buildScenarioPreview } from '~/utils/dashboardProjections'

const props = defineProps<{
  profile: FinancialProfile | null
  categories: CategoriesResponse | null
  embedded?: boolean
}>()

const scenario = defineModel<'reduce_delivery' | 'reduce_cafe' | 'reduce_entertainment' | 'custom'>(
  'scenario',
  { default: 'reduce_cafe' }
)
const percent = defineModel<number>('percent', { default: 20 })

const scenarioMeta: Record<
  (typeof SCENARIO_OPTIONS)[number]['value'],
  { emoji: string; hint: string }
> = {
  reduce_delivery: { emoji: '🛵', hint: 'Сервисы доставки еды' },
  reduce_cafe: { emoji: '☕', hint: 'Кофе, обеды, рестораны' },
  reduce_entertainment: { emoji: '🎬', hint: 'Кино, подписки, досуг' },
  custom: { emoji: '✂️', hint: 'Любая категория трат' }
}

const preview = computed(() =>
  buildScenarioPreview({
    profile: props.profile,
    categories: props.categories,
    scenario: scenario.value,
    reductionPercent: percent.value,
    months: 12
  })
)

const compareMax = computed(() =>
  Math.max(preview.value.baselineEnd, preview.value.optimizedEnd, 1)
)

const baselineWidth = computed(() =>
  Math.round((preview.value.baselineEnd / compareMax.value) * 100)
)

const optimizedWidth = computed(() =>
  Math.round((preview.value.optimizedEnd / compareMax.value) * 100)
)
</script>

<template>
  <component
    :is="embedded ? 'div' : 'Card'"
    data-demo="scenario-simulator"
    class="mm-scenario-simulator"
    :class="embedded ? 'space-y-5' : undefined"
  >
    <component :is="embedded ? 'div' : 'CardHeader'" class="space-y-2">
      <div class="flex items-center gap-2">
        <span
          class="flex size-9 items-center justify-center rounded-full bg-primary/15 text-primary"
          aria-hidden="true"
        >
          <Sparkles class="size-4" />
        </span>
        <div>
          <CardTitle class="text-base">А если чуть меньше тратить?</CardTitle>
          <CardDescription class="text-sm">
            Подберите категорию — сразу увидите, сколько можно отложить
          </CardDescription>
        </div>
      </div>
    </component>

    <component :is="embedded ? 'div' : 'CardContent'" class="space-y-5">
      <div class="space-y-2">
        <p class="text-sm font-medium text-foreground">Куда «резать»</p>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="opt in SCENARIO_OPTIONS"
            :key="opt.value"
            type="button"
            class="mm-scenario-simulator__chip"
            :class="{ 'mm-scenario-simulator__chip--active': scenario === opt.value }"
            @click="scenario = opt.value"
          >
            <span aria-hidden="true">{{ scenarioMeta[opt.value].emoji }}</span>
            <span>{{ opt.label }}</span>
          </button>
        </div>
        <p class="text-xs text-muted-foreground">
          {{ scenarioMeta[scenario].hint }}
        </p>
      </div>

      <div class="mm-scenario-simulator__slider-wrap">
        <div class="flex items-end justify-between gap-3">
          <Label for="scenario-percent" class="text-sm font-medium">На сколько меньше</Label>
          <p class="mm-scenario-simulator__percent" aria-live="polite">
            {{ percent }}%
          </p>
        </div>
        <input
          id="scenario-percent"
          v-model.number="percent"
          type="range"
          min="5"
          max="50"
          step="5"
          class="mm-scenario-simulator__range mt-3 w-full"
        />
        <div class="mt-1 flex justify-between text-xs text-muted-foreground">
          <span>5%</span>
          <span>50%</span>
        </div>
      </div>

      <Transition name="mm-scenario-fade" mode="out-in">
        <div
          v-if="preview.hasData"
          :key="`${scenario}-${percent}`"
          class="mm-scenario-simulator__result"
        >
          <div class="mm-scenario-simulator__metrics">
            <div class="mm-scenario-simulator__metric">
              <p class="mm-scenario-simulator__metric-label">Экономия в месяц</p>
              <p class="mm-scenario-simulator__metric-value">
                +{{ formatRub(preview.monthlySaving) }}
              </p>
            </div>
            <div class="mm-scenario-simulator__metric mm-scenario-simulator__metric--accent">
              <p class="mm-scenario-simulator__metric-label">За 12 месяцев</p>
              <p class="mm-scenario-simulator__metric-value">
                +{{ formatRub(preview.totalGain) }}
              </p>
            </div>
          </div>

          <div class="space-y-3">
            <p class="text-sm font-medium text-foreground">Накопления через год</p>
            <div class="space-y-2">
              <div class="mm-scenario-simulator__bar-row">
                <span class="mm-scenario-simulator__bar-label">Как сейчас</span>
                <div class="mm-scenario-simulator__bar-track">
                  <div
                    class="mm-scenario-simulator__bar mm-scenario-simulator__bar--base"
                    :style="{ width: `${baselineWidth}%` }"
                  />
                </div>
                <span class="mm-scenario-simulator__bar-value">
                  {{ formatRub(preview.baselineEnd) }}
                </span>
              </div>
              <div class="mm-scenario-simulator__bar-row">
                <span class="mm-scenario-simulator__bar-label">С «{{ preview.scenarioLabel }}» −{{ percent }}%</span>
                <div class="mm-scenario-simulator__bar-track">
                  <div
                    class="mm-scenario-simulator__bar mm-scenario-simulator__bar--gain"
                    :style="{ width: `${optimizedWidth}%` }"
                  />
                </div>
                <span class="mm-scenario-simulator__bar-value mm-scenario-simulator__bar-value--accent">
                  {{ formatRub(preview.optimizedEnd) }}
                </span>
              </div>
            </div>
          </div>

          <p class="mm-scenario-simulator__footnote">
            <TrendingUp class="size-4 shrink-0 text-primary" aria-hidden="true" />
            <span>Разница {{ formatRub(preview.totalGain) }} — из вашего дохода и трат в профиле.</span>
          </p>
        </div>

        <div v-else key="empty" class="mm-scenario-simulator__empty">
          <p>{{ preview.message }}</p>
        </div>
      </Transition>
    </component>
  </component>
</template>
