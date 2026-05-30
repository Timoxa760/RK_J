<script setup lang="ts">
import { buildDigestPageNarrative, buildWeeklyAction } from '~/utils/pageNarrative'

const { digest, loading, error, loadDigest } = useDigest()
const { primaryGoal, fetchGoals } = useGoals()

const pageLoading = computed(() => loading.value)

onMounted(async () => {
  await Promise.all([loadDigest(), fetchGoals()])
})

const weeklyAction = computed(() =>
  buildWeeklyAction(digest.value, null, null)
)

const pageNarrative = computed(() => {
  const block = buildDigestPageNarrative({
    digest: digest.value,
    primaryGoal: primaryGoal.value,
    weeklyAction: weeklyAction.value
  })
  return {
    ...block,
    weeklyAction: undefined
  }
})

const periodLabel = computed(() => {
  if (!digest.value) return ''
  return `${digest.value.period.from} — ${digest.value.period.to}`
})
</script>

<template>
  <div class="mx-auto w-full max-w-4xl space-y-6">
    <Alert v-if="error" variant="destructive">
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <SharedPageNarrative :narrative="pageNarrative" :loading="pageLoading && !digest" />

    <Skeleton v-if="pageLoading && !digest" class="h-48 w-full" />

    <template v-else-if="digest">
      <DigestGoalProgress :goal="primaryGoal" :monthly-saving="digest.saved" />

      <DigestWeeklyAction :action="weeklyAction" :period-label="periodLabel" />

      <section class="space-y-4" aria-label="Динамика за период">
        <div>
          <h2 class="text-sm font-medium text-muted-foreground">Динамика за период</h2>
        </div>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-base">Цифры месяца</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
              <div class="rounded-lg bg-muted p-3 text-center">
                <p class="text-xs text-muted-foreground">Потрачено</p>
                <p class="mt-1 font-semibold">{{ digest.total_spent.toLocaleString('ru-RU') }} ₽</p>
              </div>
              <div class="rounded-lg bg-muted p-3 text-center">
                <p class="text-xs text-muted-foreground">Доход</p>
                <p class="mt-1 font-semibold">{{ digest.total_income.toLocaleString('ru-RU') }} ₽</p>
              </div>
              <div class="rounded-lg bg-muted p-3 text-center">
                <p class="text-xs text-muted-foreground">Отложено</p>
                <p class="mt-1 font-semibold text-primary">{{ digest.saved.toLocaleString('ru-RU') }} ₽</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-base">Траты под контролем</CardTitle>
            <CardDescription>{{ digest.insights_summary }}</CardDescription>
          </CardHeader>
          <CardContent>
            <p class="text-3xl font-bold text-primary">{{ digest.mindfulness_rating }}/100</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-base">По категориям</CardTitle>
          </CardHeader>
          <CardContent>
            <ul class="space-y-2">
              <li
                v-for="cat in digest.by_category"
                :key="cat.name"
                class="flex justify-between gap-2 text-sm"
              >
                <span>{{ cat.name }}</span>
                <span class="shrink-0 text-muted-foreground">
                  {{ cat.total.toLocaleString('ru-RU') }} ₽ · {{ cat.percent }}%
                  <span v-if="cat.trend" class="text-primary">{{ cat.trend }}</span>
                </span>
              </li>
            </ul>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-base">Облако покупок</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="flex flex-wrap gap-2">
              <Badge v-for="word in digest.word_cloud" :key="word" variant="secondary">
                {{ word }}
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="pb-2">
            <CardTitle class="text-base">Топ магазинов</CardTitle>
          </CardHeader>
          <CardContent>
            <ul class="space-y-2 text-sm">
              <li
                v-for="store in digest.top_stores"
                :key="store.name"
                class="flex justify-between gap-2"
              >
                <span>{{ store.name }} ({{ store.visits }} визитов)</span>
                <span class="shrink-0 font-medium">{{ store.total.toLocaleString('ru-RU') }} ₽</span>
              </li>
            </ul>
          </CardContent>
        </Card>
      </section>
    </template>

    <Card v-else>
      <CardContent class="py-10 text-center">
        <p class="text-sm text-muted-foreground">Не удалось загрузить сводку за месяц.</p>
        <Button class="mt-4" variant="secondary" @click="loadDigest">Повторить</Button>
      </CardContent>
    </Card>
  </div>
</template>
