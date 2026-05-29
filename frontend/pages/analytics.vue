<script setup lang="ts">
const {
  timeMachine,
  forecast,
  insights,
  loading,
  scenarioResult,
  loadAll,
  simulateScenario
} = useAnalytics()

const cutCategory = ref('Кафе')
const percent = ref(20)

onMounted(() => {
  loadAll()
})

async function runSimulation() {
  await simulateScenario({ cutCategory: cutCategory.value, percent: percent.value })
}
</script>

<template>
  <div class="mm-page-shell">
    <SharedSkeletonLoader v-if="loading" height="480px" />

    <template v-else>
      <article class="mm-card p-3 sm:p-4" data-demo="timemachine">
        <h2 class="mb-2 text-sm font-semibold text-[color:var(--mm-text)] sm:mb-3">Машина времени</h2>
        <p v-if="timeMachine" class="mb-2 text-xs text-[color:var(--mm-primary)]">
          Разница с оптимистичным сценарием: +{{ timeMachine.delta.toLocaleString('ru-RU') }} ₽
        </p>
        <div class="mm-chart-wrap">
          <ChartsTimeMachineChart :data="timeMachine" />
        </div>
      </article>

      <article class="mm-card p-3 sm:p-4">
        <h2 class="mb-2 text-sm font-semibold text-[color:var(--mm-text)] sm:mb-3">Прогноз трат</h2>
        <div class="mm-chart-wrap">
          <ChartsForecastChart :data="forecast" />
        </div>
      </article>

      <article class="mm-card p-4 sm:p-5">
        <h2 class="mb-3 text-sm font-semibold text-[color:var(--mm-text)]">Симулятор «что если»</h2>
        <div class="flex flex-col gap-3 sm:flex-row sm:flex-wrap">
          <select
            v-model="cutCategory"
            class="w-full rounded-lg border border-[color:var(--mm-border)] bg-white px-3 py-2.5 text-sm sm:w-auto"
          >
            <option>Кафе</option>
            <option>Продукты</option>
            <option>Развлечения</option>
          </select>
          <div class="flex items-center gap-2">
            <input
              v-model.number="percent"
              type="number"
              min="5"
              max="50"
              class="w-full rounded-lg border border-[color:var(--mm-border)] px-3 py-2.5 text-sm sm:w-24"
            />
            <span class="shrink-0 text-sm text-[color:var(--mm-text-soft)]">%</span>
          </div>
          <button
            type="button"
            class="mm-btn-primary w-full !py-2.5 sm:w-auto sm:!px-4"
            @click="runSimulation"
          >
            Симулировать
          </button>
        </div>
        <p v-if="scenarioResult" class="mt-3 text-sm text-[color:var(--mm-text-muted)]">{{ scenarioResult }}</p>
      </article>

      <article class="mm-card p-4 sm:p-6">
        <h2 class="mb-4 text-sm font-semibold text-[color:var(--mm-text)]">Инсайты</h2>
        <InsightsInsightsPanel v-if="insights" :insights="insights.insights" />
      </article>
    </template>
  </div>
</template>
