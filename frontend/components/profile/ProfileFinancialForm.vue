<script setup lang="ts">
import type { FinancialProfile } from '~/types/api'

const { profile, loadProfile, saveProfile, syncProfileToApi } = useFinancialProfile()
const saving = ref(false)

const draft = ref<FinancialProfile>({
  active_income: 0,
  passive_income: 0,
  emergency_fund: 0
})

const saved = ref(false)

const draftTotalIncome = computed(
  () => Math.max(0, draft.value.active_income) + Math.max(0, draft.value.passive_income)
)

onMounted(() => {
  loadProfile()
  draft.value = { ...profile.value }
})

watch(profile, (value) => {
  draft.value = { ...value }
})

async function submit() {
  saving.value = true
  try {
    saveProfile({
      active_income: Math.max(0, draft.value.active_income),
      passive_income: Math.max(0, draft.value.passive_income),
      emergency_fund: Math.max(0, draft.value.emergency_fund),
      goal_title: draft.value.goal_title?.trim() ?? '',
      goal_amount: Math.max(0, draft.value.goal_amount ?? 0),
      skipped_goal: (draft.value.goal_amount ?? 0) < 1000
    })
    await syncProfileToApi(profile.value)
    saved.value = true
    setTimeout(() => {
      saved.value = false
    }, 2000)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <Card data-demo="profile-financial">
    <CardHeader>
      <CardTitle class="text-base">Ваши цифры</CardTitle>
      <CardDescription>
        Доход и запас нужны для прогноза на главной и в разделе «Прогноз»
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form class="grid gap-4 sm:grid-cols-2" @submit.prevent="submit">
        <div class="space-y-2">
          <Label for="active-income">Активный доход, ₽/мес</Label>
          <Input
            id="active-income"
            v-model.number="draft.active_income"
            type="number"
            min="0"
            step="1000"
          />
        </div>
        <div class="space-y-2">
          <Label for="passive-income">Пассивный доход, ₽/мес</Label>
          <Input
            id="passive-income"
            v-model.number="draft.passive_income"
            type="number"
            min="0"
            step="1000"
          />
        </div>
        <div class="space-y-2 sm:col-span-2">
          <Label for="emergency-fund">Подушка безопасности, ₽</Label>
          <Input
            id="emergency-fund"
            v-model.number="draft.emergency_fund"
            type="number"
            min="0"
            step="5000"
          />
        </div>

        <div class="space-y-2 sm:col-span-2">
          <Label for="goal-title">Финансовая цель</Label>
          <Input id="goal-title" v-model="draft.goal_title" type="text" placeholder="Отпуск, подушка…" />
        </div>
        <div class="space-y-2">
          <Label for="goal-amount">Сумма цели, ₽</Label>
          <Input
            id="goal-amount"
            v-model.number="draft.goal_amount"
            type="number"
            min="0"
            step="1000"
          />
        </div>

        <div class="flex flex-col gap-2 sm:col-span-2 sm:flex-row sm:items-center sm:justify-between">
          <p class="text-sm text-muted-foreground">
            Суммарный доход: {{ draftTotalIncome.toLocaleString('ru-RU') }} ₽/мес
          </p>
          <Button type="submit" class="w-full sm:w-auto" :disabled="saving">
            {{ saving ? 'Сохраняем…' : 'Сохранить' }}
          </Button>
        </div>

        <p v-if="saved" class="text-sm text-primary sm:col-span-2">Сохранено</p>
      </form>
    </CardContent>
  </Card>
</template>
