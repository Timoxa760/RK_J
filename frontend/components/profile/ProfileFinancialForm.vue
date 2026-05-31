<script setup lang="ts">
import type { FinancialProfile, FixedExpense } from '~/types/api'
import { defaultFinancialProfile } from '~/composables/useFinancialProfile'
import { PROFILE } from '~/constants/productCopy'

const { profile, loadProfile, saveProfile, syncProfileToApi } = useFinancialProfile()
const saving = ref(false)

const draft = ref<FinancialProfile>({ ...defaultFinancialProfile })

const saved = ref(false)

const draftTotalIncome = computed(
  () => Math.max(0, draft.value.active_income) + Math.max(0, draft.value.passive_income)
)

const draftFixedTotal = computed(() =>
  (draft.value.fixed_expenses ?? []).reduce((sum, row) => sum + Math.max(0, row.amount), 0)
)

function syncDraftFromProfile(value: FinancialProfile) {
  draft.value = {
    ...value,
    fixed_expenses: [...(value.fixed_expenses ?? [])]
  }
}

onMounted(() => {
  loadProfile()
  syncDraftFromProfile(profile.value)
})

watch(profile, syncDraftFromProfile)

function addFixedExpense() {
  draft.value.fixed_expenses = [...(draft.value.fixed_expenses ?? []), { title: '', amount: 0 }]
}

function updateFixedExpense(index: number, patch: Partial<FixedExpense>) {
  draft.value.fixed_expenses = (draft.value.fixed_expenses ?? []).map((item, i) =>
    i === index ? { ...item, ...patch } : item
  )
}

function removeFixedExpense(index: number) {
  draft.value.fixed_expenses = (draft.value.fixed_expenses ?? []).filter((_, i) => i !== index)
}

function normalizeFixedExpenses(rows: FixedExpense[]): FixedExpense[] {
  return rows
    .filter((row) => row.title.trim() || row.amount > 0)
    .map((row) => ({
      title: row.title.trim(),
      amount: Math.max(0, row.amount)
    }))
}

async function submit() {
  saving.value = true
  try {
    const fixedExpenses = normalizeFixedExpenses(draft.value.fixed_expenses ?? [])
    const hasFixed = fixedExpenses.some((row) => row.amount > 0)

    saveProfile({
      active_income: Math.max(0, draft.value.active_income),
      passive_income: Math.max(0, draft.value.passive_income),
      emergency_fund: Math.max(0, draft.value.emergency_fund),
      goal_title: draft.value.goal_title?.trim() ?? '',
      goal_amount: Math.max(0, draft.value.goal_amount ?? 0),
      skipped_goal: (draft.value.goal_amount ?? 0) < 1000,
      fixed_expenses: fixedExpenses,
      skipped_expenses: !hasFixed
    })
    await syncProfileToApi(profile.value)
    syncDraftFromProfile(profile.value)
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
      <CardTitle class="text-base">{{ PROFILE.formTitle }}</CardTitle>
      <CardDescription>
        {{ PROFILE.formHint }}
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form class="grid gap-4 sm:grid-cols-2" @submit.prevent="submit">
        <div class="space-y-2">
          <Label for="active-income">Зарплата и работа, ₽/мес</Label>
          <Input
            id="active-income"
            v-model.number="draft.active_income"
            type="number"
            min="0"
            step="1000"
          />
        </div>
        <div class="space-y-2">
          <Label for="passive-income">Другое поступление, ₽/мес</Label>
          <Input
            id="passive-income"
            v-model.number="draft.passive_income"
            type="number"
            min="0"
            step="1000"
          />
        </div>
        <div class="space-y-2 sm:col-span-2">
          <Label for="emergency-fund">Накопления на чёрный день, ₽</Label>
          <Input
            id="emergency-fund"
            v-model.number="draft.emergency_fund"
            type="number"
            min="0"
            step="5000"
          />
        </div>

        <div class="space-y-2 sm:col-span-2">
          <Label for="goal-title">На что копите</Label>
          <Input id="goal-title" v-model="draft.goal_title" type="text" placeholder="Отпуск, ремонт…" />
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

        <div class="space-y-3 sm:col-span-2">
          <div>
            <p class="text-sm font-medium">{{ PROFILE.fixedExpensesTitle }}</p>
            <p class="mt-0.5 text-sm text-muted-foreground">
              {{ PROFILE.fixedExpensesHint }}
            </p>
            <p
              v-if="draft.skipped_expenses && !draftFixedTotal"
              class="mt-1 text-xs text-muted-foreground"
            >
              На опросе вы пропустили этот блок — добавьте платежи ниже.
            </p>
          </div>

          <ul v-if="draft.fixed_expenses?.length" class="space-y-3">
            <li
              v-for="(item, index) in draft.fixed_expenses"
              :key="index"
              class="grid gap-2 sm:grid-cols-[1fr_120px_auto]"
            >
              <Input
                :model-value="item.title"
                placeholder="Название"
                @update:model-value="updateFixedExpense(index, { title: String($event) })"
              />
              <Input
                :model-value="item.amount"
                type="number"
                min="0"
                step="500"
                placeholder="₽/мес"
                @update:model-value="updateFixedExpense(index, { amount: Number($event) })"
              />
              <Button type="button" variant="ghost" size="sm" @click="removeFixedExpense(index)">
                Удалить
              </Button>
            </li>
          </ul>

          <Button type="button" variant="secondary" class="w-full sm:w-auto" @click="addFixedExpense">
            {{ PROFILE.addFixedExpense }}
          </Button>

          <p v-if="draftFixedTotal" class="text-sm text-muted-foreground">
            Платежи каждый месяц: {{ draftFixedTotal.toLocaleString('ru-RU') }} ₽/мес
          </p>
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
