<script setup lang="ts">
import { buildCreditsPageNarrative } from '~/utils/pageNarrative'
import { isAppFeatureEnabled } from '~/constants/featureFlags'
import { buildCreditsPortfolioSummary } from '~/utils/creditReport'

definePageMeta({
  middleware: () => {
    if (!isAppFeatureEnabled('creditsNav')) {
      return navigateTo('/dashboard')
    }
  }
})

const {
  dashboard,
  enrichedDashboard,
  loading: creditsLoading,
  fetchError,
  scanResult,
  hasCredits,
  fetchDashboard
} = useCredits()

const pageNarrative = computed(() => buildCreditsPageNarrative(enrichedDashboard.value ?? dashboard.value))
const portfolioSummary = computed(() =>
  enrichedDashboard.value ? buildCreditsPortfolioSummary(enrichedDashboard.value) : ''
)

onMounted(() => {
  fetchDashboard()
})
</script>

<template>
  <div class="mx-auto w-full max-w-6xl space-y-6">
    <SharedPageNarrative :narrative="pageNarrative" :loading="creditsLoading && !dashboard" />

    <Alert v-if="fetchError" variant="destructive">
      <AlertDescription>{{ fetchError }}</AlertDescription>
    </Alert>

    <Skeleton v-if="creditsLoading && !dashboard" class="h-40 w-full" />

    <template v-else-if="dashboard">
      <CreditsPdfUpload />

      <template v-if="hasCredits && enrichedDashboard">
        <p v-if="portfolioSummary" class="text-sm text-muted-foreground">
          {{ portfolioSummary }}
        </p>
        <CreditsHealthCards :dashboard="enrichedDashboard" />
        <CreditsContractsList :dashboard="enrichedDashboard" />
      </template>
    </template>

    <Card v-else>
      <CardContent class="py-10 text-center">
        <p class="text-sm text-muted-foreground">Не удалось загрузить данные по кредитам.</p>
        <Button class="mt-4" variant="secondary" @click="fetchDashboard">Повторить</Button>
      </CardContent>
    </Card>
  </div>
</template>
