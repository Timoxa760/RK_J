<script setup lang="ts">
const { dashboard, loading, error, fetchDashboard } = useCredits()

onMounted(() => {
  fetchDashboard()
})

function dtiColor(value: number) {
  if (value < 35) return 'bg-[color:var(--mm-primary)]'
  if (value < 50) return 'bg-[#c4a574]'
  return 'bg-[#c4847a]'
}
</script>

<template>
  <div class="mm-page-shell">
    <div v-if="error" class="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ error }}
    </div>

    <SharedSkeletonLoader v-if="loading && !dashboard" height="360px" />

    <template v-else-if="dashboard">
      <div class="grid gap-4 sm:gap-6 lg:grid-cols-2">
        <article class="mm-card p-4 sm:p-6" data-demo="credits-dti">
          <h2 class="text-sm font-semibold text-[color:var(--mm-text)]">Кредитный светофор</h2>
          <p class="mt-1 text-xs text-[color:var(--mm-text-muted)]">
            Доля дохода на погашение долгов. Норма — ниже 35%.
          </p>
          <p class="mt-4 text-2xl font-bold text-[color:var(--mm-text)] sm:text-3xl">{{ dashboard.dti }}%</p>
          <div class="mt-3 h-2 overflow-hidden rounded-full bg-[color:var(--mm-bg-muted)]">
            <div
              class="h-full rounded-full transition-all"
              :class="dtiColor(dashboard.dti)"
              :style="{ width: `${dashboard.dti}%` }"
            />
          </div>
        </article>

        <article class="mm-card p-4 sm:p-6">
          <h2 class="text-sm font-semibold text-[color:var(--mm-text)]">Стресс-тест</h2>
          <p class="mt-4 text-2xl font-bold text-[color:var(--mm-text)] sm:text-3xl">{{ dashboard.stress_test_dti }}%</p>
          <div class="mt-3 h-2 overflow-hidden rounded-full bg-[color:var(--mm-bg-muted)]">
            <div
              class="h-full rounded-full"
              :class="dtiColor(dashboard.stress_test_dti)"
              :style="{ width: `${dashboard.stress_test_dti}%` }"
            />
          </div>
          <p class="mt-4 text-sm text-[color:var(--mm-text-muted)]">
            Доход: {{ dashboard.monthly_income.toLocaleString('ru-RU') }} ₽/мес
          </p>
        </article>
      </div>

      <section class="mm-card p-4 sm:p-6">
        <h2 class="text-sm font-semibold text-[color:var(--mm-text)]">Кредиты</h2>
        <ul class="mt-4 divide-y divide-[color:var(--mm-border-subtle)]">
          <li
            v-for="credit in dashboard.credits"
            :key="credit.id"
            class="flex flex-col gap-1 py-4 sm:flex-row sm:flex-wrap sm:justify-between sm:gap-2"
          >
            <div class="min-w-0">
              <p class="font-medium text-[color:var(--mm-text)]">{{ credit.name }}</p>
              <p class="text-xs text-[color:var(--mm-text-soft)]">{{ credit.rate }}%</p>
            </div>
            <div class="text-sm sm:text-right">
              <p class="font-medium">{{ credit.payment.toLocaleString('ru-RU') }} ₽/мес</p>
              <p class="text-[color:var(--mm-text-soft)]">{{ credit.balance.toLocaleString('ru-RU') }} ₽</p>
            </div>
          </li>
        </ul>
      </section>
    </template>
  </div>
</template>
