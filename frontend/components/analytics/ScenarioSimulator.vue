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

const compareMax = computed(() =>
  Math.max(preview.value.currentBalance, preview.value.optimizedEnd, preview.value.baselineEnd, 1)
)

const baselineWidth = computed(() =>
  Math.round((preview.value.currentBalance / compareMax.value) * 100)
)

const optimizedWidth = computed(() =>
  Math.round((preview.value.optimizedEnd / compareMax.value) * 100)
)

const scenarioBarLabel = computed(() => {
  const name = preview.value.categoryName
  return `С «${name}» −${percent.value}%`
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
            Категории из ваших покупок — сразу видно, сколько можно отложить
          </CardDescription>
        </div>
      </div>
    </component>

    <component :is="embedded ? 'div' : 'CardContent'" class="space-y-5">
      <div class="space-y-2">
        <p class="text-sm font-medium text-foreground">Куда «резать»</p>
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
        <p v-else class="text-sm text-muted-foreground">
          Пока нет трат по категориям — добавьте покупку голосом или вручную.
        </p>
        <p v-if="activeOption" class="text-xs text-muted-foreground">
          {{ categoryHint(activeOption) }}
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
          :disabled="!categoryOptions.length"
        />
        <div class="mt-1 flex justify-between text-xs text-muted-foreground">
          <span>5%</span>
          <span>50%</span>
        </div>
      </div>

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

          <div class="space-y-3">
            <p class="text-sm font-medium text-foreground">Подушка через год</p>
            <div class="space-y-2">
              <div class="mm-scenario-simulator__bar-row">
                <span class="mm-scenario-simulator__bar-label">Как сейчас</span>
                <div class="mm-scenario-simulator__bar-track">
                  <div
                    class="mm-scenario-simulator__bar mm-scenario-simulator__bar--base"
                    :style="{ width: `${Math.max(baselineWidth, 8)}%` }"
                  />
                </div>
                <span class="mm-scenario-simulator__bar-value">
                  {{ formatRub(preview.currentBalance) }}
                </span>
              </div>
              <div
                v-if="preview.baselineEnd !== preview.currentBalance"
                class="mm-scenario-simulator__bar-row"
              >
                <span class="mm-scenario-simulator__bar-label">Без изменений</span>
                <div class="mm-scenario-simulator__bar-track">
                  <div
                    class="mm-scenario-simulator__bar mm-scenario-simulator__bar--muted"
                    :style="{
                      width: `${Math.max(Math.round((preview.baselineEnd / compareMax) * 100), 8)}%`
                    }"
                  />
                </div>
                <span class="mm-scenario-simulator__bar-value">
                  {{ formatRub(preview.baselineEnd) }}
                </span>
              </div>
              <div class="mm-scenario-simulator__bar-row">
                <span class="mm-scenario-simulator__bar-label">{{ scenarioBarLabel }}</span>
                <div class="mm-scenario-simulator__bar-track">
                  <div
                    class="mm-scenario-simulator__bar mm-scenario-simulator__bar--gain"
                    :style="{ width: `${Math.max(optimizedWidth, 8)}%` }"
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
            <span>
              «{{ preview.categoryName }}» — {{ formatRub(preview.categorySpend) }} за месяц.
              <template v-if="preview.freeCashflow !== 0">
                Свободно {{ preview.freeCashflow > 0 ? '+' : '' }}{{ formatRub(preview.freeCashflow) }}/мес. из профиля и трат.
              </template>
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
