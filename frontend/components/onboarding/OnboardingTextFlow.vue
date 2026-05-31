<script setup lang="ts">
import { GOAL_KIND_LABELS } from '~/constants/onboardingGoals'
import {
  ONBOARDING_SKIP_HINTS,
  ONBOARDING_SKIP_LABEL,
  surveyStepToSkipKey
} from '~/constants/onboardingSkips'
import type { GoalKind, OnboardingDraft } from '~/types/api'

const props = defineProps<{
  draft: OnboardingDraft
  step: number
}>()

const emit = defineEmits<{
  back: []
  next: []
  skip: []
  patch: [partial: Partial<OnboardingDraft>]
  setGoalKind: [kind: GoalKind]
  patchBreakdown: [partial: Partial<OnboardingDraft['emergency_breakdown']>]
  addFixedExpense: []
  updateFixedExpense: [index: number, partial: Partial<{ title: string; amount: number }>]
  removeFixedExpense: [index: number]
}>()

const goalKinds = Object.entries(GOAL_KIND_LABELS) as [GoalKind, string][]

const titles: Record<number, { title: string; description: string }> = {
  1: {
    title: 'Доходы',
    description: 'Зарплата, пенсия, сдача квартиры — всё, что приходит каждый месяц.'
  },
  2: {
    title: 'Запас',
    description: 'Сколько отложено на чёрный день. Разбивку можно не заполнять.'
  },
  3: {
    title: 'Цель',
    description: 'На что копите и какая сумма — покажем, когда примерно дойдёте.'
  },
  4: {
    title: 'Постоянные платежи',
    description: 'Что платите каждый месяц: аренда, кредит, детский сад.'
  }
}

const skipKey = computed(() => surveyStepToSkipKey(props.step))
const skipHint = computed(() => (skipKey.value ? ONBOARDING_SKIP_HINTS[skipKey.value] : undefined))

const canProceed = computed(() => {
  switch (props.step) {
    case 1:
      return props.draft.active_income > 0 || props.draft.passive_income > 0
    case 2:
      return props.draft.emergency_fund > 0
    case 3:
      return props.draft.goal_amount >= 1000
    case 4:
      return props.draft.fixed_expenses.some((item) => item.title.trim() && item.amount > 0)
    default:
      return false
  }
})

function onVoicePatch(patch: Partial<OnboardingDraft>) {
  if (patch.emergency_breakdown) {
    const b = patch.emergency_breakdown
    emit('patch', {
      ...patch,
      emergency_fund:
        patch.emergency_fund ?? b.cash + b.deposit + b.investments
    })
  } else {
    emit('patch', patch)
  }
  if (patch.goal_kind) {
    emit('setGoalKind', patch.goal_kind)
  }
}
</script>

<template>
  <OnboardingStepShell
    v-if="titles[step]"
    :title="titles[step].title"
    :description="titles[step].description"
    show-back
    :next-disabled="!canProceed"
    :next-label="step === 4 ? 'К моей картине' : 'Далее'"
    :secondary-action="{ label: ONBOARDING_SKIP_LABEL }"
    :skip-hint="skipHint"
    @back="emit('back')"
    @next="emit('next')"
    @secondary="emit('skip')"
  >
    <div v-if="step === 1" class="mm-onb-form-panel grid gap-4 md:grid-cols-2 md:items-end">
      <div class="mm-onb-field">
        <Label for="onb-active" class="mm-onb-field-label">Основной доход, ₽/мес</Label>
        <Input
          id="onb-active"
          :model-value="draft.active_income"
          type="number"
          min="0"
          step="1000"
          class="text-lg"
          @update:model-value="emit('patch', { active_income: Number($event) })"
        />
      </div>
      <div class="mm-onb-field">
        <Label for="onb-passive" class="mm-onb-field-label">Дополнительный доход, ₽/мес</Label>
        <Input
          id="onb-passive"
          :model-value="draft.passive_income"
          type="number"
          min="0"
          step="1000"
          class="text-lg"
          @update:model-value="emit('patch', { passive_income: Number($event) })"
        />
      </div>
    </div>

    <div v-else-if="step === 2" class="space-y-4">
      <div class="mm-onb-form-panel">
        <div class="space-y-2">
          <Label for="onb-fund">Общая сумма, ₽</Label>
          <Input
            id="onb-fund"
            :model-value="draft.emergency_fund"
            type="number"
            min="0"
            step="5000"
            class="text-lg"
            @update:model-value="emit('patch', { emergency_fund: Number($event) })"
          />
        </div>
      </div>
      <div class="mm-onb-form-panel border-dashed">
        <p class="text-xs font-medium text-[color:var(--mm-text-muted)]">Разбивка (необязательно)</p>
        <div class="mt-3 grid gap-3 sm:grid-cols-3 sm:items-end">
          <div class="mm-onb-field">
            <Label for="onb-cash" class="mm-onb-field-label">Наличные</Label>
            <Input
              id="onb-cash"
              :model-value="draft.emergency_breakdown.cash"
              type="number"
              min="0"
              @update:model-value="emit('patchBreakdown', { cash: Number($event) })"
            />
          </div>
          <div class="mm-onb-field">
            <Label for="onb-deposit" class="mm-onb-field-label">Вклад</Label>
            <Input
              id="onb-deposit"
              :model-value="draft.emergency_breakdown.deposit"
              type="number"
              min="0"
              @update:model-value="emit('patchBreakdown', { deposit: Number($event) })"
            />
          </div>
          <div class="mm-onb-field">
            <Label for="onb-invest" class="mm-onb-field-label">Инвестиции</Label>
            <Input
              id="onb-invest"
              :model-value="draft.emergency_breakdown.investments"
              type="number"
              min="0"
              @update:model-value="emit('patchBreakdown', { investments: Number($event) })"
            />
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="step === 3" class="mm-onb-form-panel">
      <div class="space-y-2">
        <Label>Тип цели</Label>
        <div class="grid grid-cols-2 gap-2 md:grid-cols-3 lg:grid-cols-4">
          <button
            v-for="[kind, label] in goalKinds"
            :key="kind"
            type="button"
            class="mm-onb-goal-pill"
            :class="{ 'mm-onb-goal-pill--active': draft.goal_kind === kind }"
            @click="emit('setGoalKind', kind)"
          >
            {{ label }}
          </button>
        </div>
      </div>
      <div v-if="draft.goal_kind === 'other'" class="space-y-2">
        <Label for="onb-goal-title">Название</Label>
        <Input
          id="onb-goal-title"
          :model-value="draft.goal_title"
          @update:model-value="emit('patch', { goal_title: String($event) })"
        />
      </div>
      <div class="space-y-2">
        <Label for="onb-goal-amount">Сумма цели, ₽</Label>
        <Input
          id="onb-goal-amount"
          :model-value="draft.goal_amount"
          type="number"
          min="1000"
          step="1000"
          class="text-lg"
          @update:model-value="emit('patch', { goal_amount: Number($event) })"
        />
      </div>
    </div>

    <div v-else-if="step === 4" class="mm-onb-form-panel">
      <ul v-if="draft.fixed_expenses.length" class="space-y-3">
        <li
          v-for="(item, index) in draft.fixed_expenses"
          :key="index"
          class="grid items-center gap-2 sm:grid-cols-[1fr_120px_auto]"
        >
          <Input
            :model-value="item.title"
            placeholder="Название"
            @update:model-value="emit('updateFixedExpense', index, { title: String($event) })"
          />
          <Input
            :model-value="item.amount"
            type="number"
            min="0"
            step="500"
            placeholder="₽/мес"
            @update:model-value="emit('updateFixedExpense', index, { amount: Number($event) })"
          />
          <Button
            type="button"
            variant="ghost"
            size="sm"
            @click="emit('removeFixedExpense', index)"
          >
            Удалить
          </Button>
        </li>
      </ul>

      <Button type="button" variant="secondary" class="w-full" @click="emit('addFixedExpense')">
        Добавить строку
      </Button>
    </div>
  </OnboardingStepShell>
</template>
