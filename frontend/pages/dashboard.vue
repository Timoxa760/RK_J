<script setup lang="ts">
const { sankey, stores, categories, compare, timemachine, loading, error, loadAll, retry } =
  useDashboard()

onMounted(() => {
  loadAll()
})
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <div
      v-if="error"
      class="flex flex-col gap-3 rounded-xl border border-stone-300/80 bg-stone-100 px-4 py-3 text-sm text-stone-700 sm:flex-row sm:items-center sm:justify-between"
    >
      <span>{{ error }}</span>
      <button
        type="button"
        class="shrink-0 rounded-lg border border-red-300 bg-white px-3 py-1.5 text-xs font-medium hover:bg-red-50"
        @click="retry"
      >
        Повторить
      </button>
    </div>

    <div class="grid gap-4 sm:gap-6 lg:grid-cols-2">
      <article class="mm-card p-3 sm:p-4 lg:col-span-2">
        <h2 class="mb-2 text-sm font-medium text-[color:var(--mm-text)] sm:mb-3">Откуда деньги и куда уходят</h2>
        <SharedSkeletonLoader v-if="loading && !sankey" height="260px" />
        <div v-else class="mm-chart-wrap">
          <ChartsSankeyChart :data="sankey" />
        </div>
      </article>

      <article class="mm-card p-3 sm:p-4">
        <h2 class="mb-2 text-sm font-medium text-[color:var(--mm-text)] sm:mb-3">Куда уходят деньги — по магазинам</h2>
        <SharedSkeletonLoader v-if="loading && !stores" height="260px" />
        <div v-else class="mm-chart-wrap">
          <ChartsBubbleChart :data="stores" />
        </div>
      </article>

      <article class="mm-card p-3 sm:p-4" data-demo="categories">
        <h2 class="mb-2 text-sm font-medium text-[color:var(--mm-text)] sm:mb-3">Категории</h2>
        <p class="mb-2 text-xs text-[color:var(--mm-text-soft)]">Клик по сектору — детализация</p>
        <SharedSkeletonLoader v-if="loading && !categories" height="260px" />
        <div v-else class="mm-chart-wrap">
          <ChartsCategoryPie :data="categories" />
        </div>
      </article>

      <article class="mm-card p-3 sm:p-4">
        <h2 class="mb-2 text-sm font-medium text-[color:var(--mm-text)] sm:mb-3">Сравнение месяцев</h2>
        <SharedSkeletonLoader v-if="loading && !compare" height="260px" />
        <div v-else class="mm-chart-wrap">
          <ChartsDonutCompare :data="compare" />
        </div>
      </article>

      <article class="mm-card p-3 sm:p-4 lg:col-span-2">
        <h2 class="mb-2 text-sm font-medium text-[color:var(--mm-text)] sm:mb-3">Машина времени</h2>
        <SharedSkeletonLoader v-if="loading && !timemachine" height="240px" />
        <div v-else class="mm-chart-wrap">
          <ChartsTimeMachineChart :data="timemachine" />
        </div>
      </article>
    </div>
  </div>
</template>
