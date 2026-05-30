<script setup lang="ts">
import { buildCreditsPageNarrative } from '~/utils/pageNarrative'
import { isAppFeatureEnabled } from '~/constants/featureFlags'

definePageMeta({
  middleware: () => {
    if (!isAppFeatureEnabled('creditsNav')) {
      return navigateTo('/dashboard')
    }
  }
})

const {
  dashboard,
  loading: creditsLoading,
  error,
  scanResult,
  scanLoading,
  fetchDashboard,
  scanContract
} = useCredits()

const { breakdown, loading: mortgageLoading, error: mortgageError, analyze } = useMortgage()

const mortgageAmount = ref(12_000_000)

const pageNarrative = computed(() => buildCreditsPageNarrative(dashboard.value))

onMounted(async () => {
  await fetchDashboard()
  if (dashboard.value) {
    await runMortgageAnalysis(mortgageAmount.value)
  }
})

async function runMortgageAnalysis(amount: number) {
  if (!dashboard.value) return
  mortgageAmount.value = amount
  await analyze({
    mortgage_amount: amount,
    monthly_income: dashboard.value.monthly_income,
    savings: dashboard.value.savings,
    existing_dti: dashboard.value.dti,
    stress_test_months: dashboard.value.stress_test_months
  })
}

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) scanContract(file)
}
</script>

<template>
  <div class="mx-auto w-full max-w-6xl space-y-6">
    <SharedPageNarrative :narrative="pageNarrative" :loading="creditsLoading && !dashboard" />

    <Alert v-if="error" variant="destructive">
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <Alert v-if="mortgageError" variant="destructive">
      <AlertDescription>{{ mortgageError }}</AlertDescription>
    </Alert>

    <Skeleton v-if="creditsLoading && !dashboard" class="h-[360px] w-full" />

    <template v-else-if="dashboard">
      <CreditsHealthCards :dashboard="dashboard" />

      <CreditsMortgageAnalyzerForm
        v-model:amount="mortgageAmount"
        :monthly-income="dashboard.monthly_income"
        :savings="dashboard.savings"
        :loading="mortgageLoading"
        @analyze="runMortgageAnalysis"
      />

      <template v-if="breakdown">
        <section class="space-y-4" aria-label="Ипотечный разбор">
          <CreditsMortgageBreakdownSummary :breakdown="breakdown" />
          <CreditsScenarioNowWait :breakdown="breakdown" />
          <CreditsBankCompareTable :breakdown="breakdown" />
        </section>
      </template>

      <Card>
        <CardHeader>
          <CardTitle class="text-base">Текущие кредиты</CardTitle>
          <CardDescription>Учитываются при расчёте доли дохода на кредиты</CardDescription>
        </CardHeader>
        <CardContent>
          <ul v-if="dashboard.credits.length" class="divide-y">
            <li
              v-for="credit in dashboard.credits"
              :key="credit.id"
              class="flex flex-col gap-1 py-4 sm:flex-row sm:flex-wrap sm:justify-between"
            >
              <div class="min-w-0">
                <p class="font-medium">{{ credit.name ?? credit.bank }}</p>
                <p class="text-xs text-muted-foreground">{{ credit.rate }}% годовых</p>
              </div>
              <div class="text-sm sm:text-right">
                <p class="font-medium">
                  {{ (credit.payment ?? credit.monthly_payment ?? 0).toLocaleString('ru-RU') }} ₽/мес
                </p>
                <p class="text-muted-foreground">
                  Остаток: {{ (credit.balance ?? credit.remaining ?? 0).toLocaleString('ru-RU') }} ₽
                </p>
                <p v-if="credit.next_payment" class="text-xs text-muted-foreground">
                  След. платёж: {{ credit.next_payment }}
                </p>
              </div>
            </li>
          </ul>
          <p v-else class="text-sm text-muted-foreground">Активных кредитов нет.</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-base">Загрузить договор (PDF)</CardTitle>
          <CardDescription>Загрузите PDF — подставим ставку и платёж</CardDescription>
        </CardHeader>
        <CardContent>
          <Input type="file" accept="application/pdf" @change="onFileChange" />
          <p v-if="scanLoading" class="mt-2 text-sm text-muted-foreground">Распознаём…</p>
          <div v-if="scanResult" class="mt-4 rounded-lg bg-muted p-4 text-sm">
            <p><strong>Банк:</strong> {{ scanResult.parsed.bank }}</p>
            <p><strong>Ставка:</strong> {{ scanResult.parsed.rate }}%</p>
            <p><strong>Платёж:</strong> {{ scanResult.parsed.monthly_payment.toLocaleString('ru-RU') }} ₽</p>
            <p class="text-xs text-muted-foreground">
              Уверенность: {{ Math.round(scanResult.confidence * 100) }}%
            </p>
          </div>
        </CardContent>
      </Card>
    </template>

    <Card v-else>
      <CardContent class="py-10 text-center">
        <p class="text-sm text-muted-foreground">Не удалось загрузить данные по кредитам.</p>
        <Button class="mt-4" variant="secondary" @click="fetchDashboard">Повторить</Button>
      </CardContent>
    </Card>
  </div>
</template>
