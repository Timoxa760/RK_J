<script setup lang="ts">
import type { Goal } from '~/types/api'
import { buildGoalProgressText } from '~/utils/pageNarrative'

const { goals, loading, saving, error, fetchGoals, createGoal } = useGoals()
const { totalIncome } = useFinancialProfile()

const title = ref('')
const targetAmount = ref<number | null>(null)
const targetDate = ref('')

onMounted(() => {
  fetchGoals()
})

const monthlySavingEstimate = computed(() =>
  totalIncome.value > 0 ? Math.round(totalIncome.value * 0.1) : undefined
)

async function submitGoal() {
  if (!title.value.trim() || !targetAmount.value || targetAmount.value <= 0) return

  await createGoal({
    title: title.value.trim(),
    target_amount: targetAmount.value,
    target_date: targetDate.value || undefined,
    auto_save_percent: 10
  })

  title.value = ''
  targetAmount.value = null
  targetDate.value = ''
}

function progressPercent(goal: Goal) {
  return Math.min(
    100,
    goal.progress_percent ??
      Math.round((goal.current_amount / goal.target_amount) * 100)
  )
}
</script>

<template>
  <Card data-demo="profile-goals">
    <CardHeader>
      <CardTitle class="text-base">Цели</CardTitle>
      <CardDescription>Отпуск, подушка, крупная покупка — срок при текущем поведении</CardDescription>
    </CardHeader>
    <CardContent class="space-y-6">
      <Alert v-if="error" variant="destructive">
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>

      <Skeleton v-if="loading && !goals.length" class="h-24 w-full" />

      <ul v-else-if="goals.length" class="space-y-4">
        <li
          v-for="goal in goals"
          :key="goal.id"
          class="rounded-lg border px-4 py-3"
        >
          <div class="flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between">
            <p class="font-medium">{{ goal.title }}</p>
            <p class="text-sm text-muted-foreground">
              {{ goal.current_amount.toLocaleString('ru-RU') }} /
              {{ goal.target_amount.toLocaleString('ru-RU') }} ₽
            </p>
          </div>
          <div class="mt-3 h-2 overflow-hidden rounded-full bg-muted">
            <div
              class="h-full rounded-full bg-primary transition-all"
              :style="{ width: `${progressPercent(goal)}%` }"
            />
          </div>
          <p class="mt-2 text-xs text-muted-foreground">
            {{ buildGoalProgressText(goal, monthlySavingEstimate) }}
          </p>
        </li>
      </ul>

      <p v-else class="text-sm text-muted-foreground">
        Целей пока нет — создайте первую ниже.
      </p>

      <form class="grid gap-3 border-t pt-4 sm:grid-cols-2" @submit.prevent="submitGoal">
        <div class="space-y-2 sm:col-span-2">
          <Label for="goal-title">Название цели</Label>
          <Input id="goal-title" v-model="title" placeholder="Например, отпуск" required />
        </div>
        <div class="space-y-2">
          <Label for="goal-amount">Сумма, ₽</Label>
          <Input
            id="goal-amount"
            v-model.number="targetAmount"
            type="number"
            min="1000"
            step="1000"
            required
          />
        </div>
        <div class="space-y-2">
          <Label for="goal-date">Желаемая дата</Label>
          <Input id="goal-date" v-model="targetDate" type="date" />
        </div>
        <div class="sm:col-span-2">
          <Button type="submit" class="w-full sm:w-auto" :disabled="saving">
            {{ saving ? 'Сохраняем…' : 'Добавить цель' }}
          </Button>
        </div>
      </form>
    </CardContent>
  </Card>
</template>
