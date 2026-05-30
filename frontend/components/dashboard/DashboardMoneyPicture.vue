<script setup lang="ts">
import type { CategoriesResponse, ForecastResponse, TimeMachineResponse } from '~/types/api'
import { ADVISOR } from '~/constants/productCopy'

defineProps<{
  categories: CategoriesResponse | null
  forecast: ForecastResponse | null
  timemachine: TimeMachineResponse | null
  categoriesSummary: string
  currentSavings?: number | null
  loading?: boolean
  embedded?: boolean
}>()
</script>

<template>
  <component
    :is="embedded ? 'div' : 'Card'"
    class="overflow-hidden"
    :class="embedded ? 'space-y-5' : undefined"
    data-demo="categories"
  >
    <component :is="embedded ? 'div' : 'CardHeader'" class="gap-2" :class="embedded ? 'space-y-1' : 'p-4 pb-2 sm:p-5 sm:pb-3'">
      <component :is="embedded ? 'h3' : 'CardTitle'" class="text-base font-semibold leading-snug">
        {{ ADVISOR.categoriesTitle }}
      </component>
      <p class="text-sm leading-relaxed text-muted-foreground">
        {{ categoriesSummary }}
      </p>
    </component>

    <component :is="embedded ? 'div' : 'CardContent'" :class="embedded ? 'space-y-5' : 'space-y-5 p-4 pt-0 sm:p-5 sm:pt-0'">
      <Skeleton v-if="loading" class="h-48 w-full rounded-xl" />

      <template v-else>
        <section aria-label="Категории трат">
          <ChartsSimpleCategoryBars :data="categories" />
        </section>

        <Separator />

        <div class="grid gap-5 lg:grid-cols-2 lg:gap-6">
          <section
            class="rounded-xl border bg-muted/30 p-4"
            aria-label="Накопления"
            data-demo="timemachine"
          >
            <h3 class="text-base font-semibold text-foreground">
              {{ ADVISOR.savingsTitle }}
            </h3>
            <ChartsSimpleSavingsChart
              :data="timemachine"
              :current-savings="currentSavings"
              class="mt-3"
            />
          </section>

          <section class="rounded-xl border bg-muted/30 p-4" aria-label="Прогноз трат">
            <h3 class="text-base font-semibold text-foreground">
              {{ ADVISOR.forecastTitle }}
            </h3>
            <ChartsSimpleForecastChart :data="forecast" class="mt-3" />
          </section>
        </div>
      </template>
    </component>
  </component>
</template>
