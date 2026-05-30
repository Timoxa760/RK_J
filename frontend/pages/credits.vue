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
  hasCredits,
  canAnalyzeMortgage,
  fetchDashboard,
  scanContract
} = useCredits()

const { breakdown, loading: mortgageLoading, error: mortgageError, analyze } = useMortgage()

const mortgageAmount = ref(12_000_000)
const fileInput = ref<HTMLInputElement | null>(null)

const pageNarrative = computed(() => buildCreditsPageNarrative(dashboard.value, hasCredits.value))

onMounted(async () => {
  await fetchDashboard()
  if (canAnalyzeMortgage.value && dashboard.value) {
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
  if (file) {
    void scanContract(file).then(() => {
      if (canAnalyzeMortgage.value && dashboard.value) {
        void runMortgageAnalysis(mortgageAmount.value)
      }
    })
  }
}

function openFilePicker() {
  fileInput.value?.click()
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

    <Card v-else-if="dashboard && !canAnalyzeMortgage" class="border-dashed">
      <CardHeader>
        <CardTitle class="text-base">Нужен доход для ипотечного разбора</CardTitle>
        <CardDescription>
          Заполните доход в онбординге или профиле — или загрузите PDF кредитного договора для DTI.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <input
          ref="fileInput"
          type="file"
          accept="application/pdf"
          class="sr-only"
          @change="onFileChange"
        />
        <Button type="button" :disabled="scanLoading" @click="openFilePicker">
          {{ scanLoading ? 'Распознаём…' : 'Выбрать PDF договора' }}
        </Button>
        <div v-if="scanResult" class="rounded-lg bg-muted p-4 text-sm">
          <p><strong>Банк:</strong> {{ scanResult.parsed.bank }}</p>
          <p><strong>Ставка:</strong> {{ scanResult.parsed.rate }}%</p>
          <p><strong>Платёж:</strong> {{ scanResult.parsed.monthly_payment.toLocaleString('ru-RU') }} ₽</p>
        </div>
      </CardContent>
    </Card>

    <template v-else-if="dashboard && canAnalyzeMortgage">
      <CreditsHealthCards v-if="hasCredits" :dashboard="dashboard" />

      <CreditsMortgageAnalyzerForm
        v-model:amount="mortgageAmount"
        data-demo="mortgage-form"
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

      <Card v-if="hasCredits">
        <CardHeader>
          <CardTitle class="text-base">Текущие кредиты</CardTitle>
          <CardDescription>Учитываются при расчёте доли дохода на кредиты</CardDescription>
        </CardHeader>
        <CardContent>
          <ul class="divide-y">
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
        </CardContent>
      </Card>

      <Card class="border-dashed">
        <CardHeader>
          <CardTitle class="text-base">
            {{ hasCredits ? 'Добавить договор (PDF)' : 'Кредитный договор (PDF, опционально)' }}
          </CardTitle>
          <CardDescription>
            {{
              hasCredits
                ? 'Загрузите ещё один PDF — обновим расчёт DTI'
                : 'PDF уточнит DTI и список обязательств; ипотечный разбор уже доступен по доходу из профиля'
            }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Input type="file" accept="application/pdf" @change="onFileChange" />
          <p v-if="scanLoading" class="mt-2 text-sm text-muted-foreground">Распознаём…</p>
          <div v-if="scanResult" class="mt-4 rounded-lg bg-muted p-4 text-sm">
            <p><strong>Банк:</strong> {{ scanResult.parsed.bank }}</p>
            <p><strong>Ставка:</strong> {{ scanResult.parsed.rate }}%</p>
            <p><strong>Платёж:</strong> {{ scanResult.parsed.monthly_payment.toLocaleString('ru-RU') }} ₽</p>
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
