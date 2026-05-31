<script setup lang="ts">
import { Sparkles, TrendingUp } from 'lucide-vue-next'
import type { CategoriesResponse, FinancialProfile } from '~/types/api'
import { formatRub } from '~/constants/productCopy'
import {
  buildUserCategoryOptions,
  CATEGORY_EMOJI
} from '~/constants/expenseCategories'
import { buildScenarioPreview } from '~/utils/dashboardProjections'

const props = defineProps<{
  profile: FinancialProfile | null
  categories: CategoriesResponse | null
  embedded?: boolean
}>()

const selectedCategory = defineModel<string>('selectedCategory', { default: '' })
const percent = defineModel<number>('percent', { default: 20 })

const categoryOptions = computed(() => buildUserCategoryOptions(props.categories))

watch(
  categoryOptions,
  (options) => {
    if (!options.length) {
      selectedCategory.value = ''
      return
    }
    if (!options.some((row) => row.name === selectedCategory.value)) {
      selectedCategory.value = options[0]!.name
    }
  },
  { immediate: true }
)

watch(percent, (value) => {
  if (!Number.isFinite(value)) {
    percent.value = 20
    return
  }
  percent.value = Math.min(40, Math.max(10, Math.round(value / 5) * 5))
})

const activeOption = computed(() =>
  categoryOptions.value.find((row) => row.name === selectedCategory.value)
)

const preview = computed(() =>
  buildScenarioPreview({
    profile: props.profile,
    categories: props.categories,
    categoryName: selectedCategory.value,
    reductionPercent: percent.value,
    months: 12
  })
)

const savingBarWidth = computed(() => {
  const spend = preview.value.categorySpend
  const maxGain = Math.max(Math.round(spend * 0.4 * preview.value.months), 1)
  return Math.max(Math.round((preview.value.totalGain / maxGain) * 100), 8)
})

const scenarioBarLabel = computed(() => {
  const name = preview.value.categoryName
  return `С «${name}» −${percent.value}%`
})

const contextLine = computed(() => {
  if (!preview.value.hasData) return ''
  if (preview.value.incomeKnown) {
    return `Подушка за ${preview.value.months} мес.: ${formatRub(preview.value.baselineEnd)} → ${formatRub(preview.value.optimizedEnd)} при текущем свободном потоке ${preview.value.freeCashflow > 0 ? '+' : ''}${formatRub(preview.value.freeCashflow)}/мес.`
  }
  return 'Укажите доход в профиле — добавим прогноз роста подушки.'
})

function categoryHint(option: (typeof categoryOptions.value)[number]): string {
  const share = Math.round(option.share * 100)
  return `${share}% трат · ${formatRub(option.amount)}/мес`
}
</script>

<template>
  <component
    :is="embedded ? 'div' : 'Card'"
    data-demo="scenario-simulator"
    class="mm-scenario-simulator"
    :class="{ 'mm-scenario-simulator--embedded': embedded }"
  >
    <component
      :is="embedded ? 'div' : 'CardHeader'"
      class="mm-scenario-simulator__header"
    >
      <div class="flex items-start gap-3">
        <span
          class="flex size-9 shrink-0 items-center justify-center rounded-full bg-primary/15 text-primary"
          aria-hidden="true"
        >
          <Sparkles class="size-4" />
        </span>
        <div class="min-w-0 space-y-0.5">
          <CardTitle class="text-base leading-snug">А если чуть меньше тратить?</CardTitle>
          <CardDescription class="text-sm leading-relaxed">
            Категории из ваших покупок — сразу видно, сколько можно отложить
          </CardDescription>
        </div>
      </div>
    </component>

    <component :is="embedded ? 'div' : 'CardContent'" class="mm-scenario-simulator__body">
      <section class="mm-scenario-simulator__section">
        <p class="mm-scenario-simulator__section-label">Куда «резать»</p>
        <div v-if="categoryOptions.length" class="flex flex-wrap gap-2">
          <button
            v-for="opt in categoryOptions"
            :key="opt.name"
            type="button"
            class="mm-scenario-simulator__chip"
            :class="{ 'mm-scenario-simulator__chip--active': selectedCategory === opt.name }"
            @click="selectedCategory = opt.name"
          >
            <span aria-hidden="true">{{ CATEGORY_EMOJI[opt.name] ?? '📦' }}</span>
            <span>{{ opt.name }}</span>
          </button>
        </div>
        <p v-else class="text-sm leading-relaxed text-muted-foreground">
          Пока нет трат по категориям — добавьте покупку голосом или вручную.
        </p>
        <p v-if="activeOption" class="text-xs text-muted-foreground">
          {{ categoryHint(activeOption) }}
        </p>
      </section>

      <section class="mm-scenario-simulator__slider-wrap">
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
          min="10"
          max="40"
          step="5"
          class="mm-scenario-simulator__range mt-3 w-full"
          :disabled="!categoryOptions.length"
        />
        <div class="mt-1.5 flex justify-between text-xs text-muted-foreground">
          <span>10%</span>
          <span>40%</span>
        </div>
      </section>

      <Transition name="mm-scenario-fade" mode="out-in">
        <div
          v-if="preview.hasData"
          :key="`${selectedCategory}-${percent}`"
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

          <div class="mm-scenario-simulator__scale" :aria-label="scenarioBarLabel">
            <div class="mm-scenario-simulator__bar-track">
              <div
                class="mm-scenario-simulator__bar mm-scenario-simulator__bar--gain"
                :style="{ width: `${savingBarWidth}%` }"
              />
            </div>
            <p class="mm-scenario-simulator__scale-caption">
              {{ scenarioBarLabel }} · максимум на шкале — 40%
            </p>
          </div>

          <p v-if="contextLine" class="mm-scenario-simulator__context">
            {{ contextLine }}
          </p>

          <p class="mm-scenario-simulator__footnote">
            <TrendingUp class="size-4 shrink-0 text-primary" aria-hidden="true" />
            <span>
              «{{ preview.categoryName }}» — {{ formatRub(preview.categorySpend) }} за месяц.
            </span>
          </p>
        </div>

        <div v-else key="empty" class="mm-scenario-simulator__empty">
          <p>{{ preview.message }}</p>
        </div>
      </Transition>
    </component>
  </component>
</template>
