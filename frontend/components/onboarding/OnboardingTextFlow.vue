<script setup lang="ts">
import { GOAL_KIND_LABELS } from '~/constants/onboardingGoals'
import type { GoalKind, OnboardingDraft } from '~/types/api'

const props = defineProps<{
  draft: OnboardingDraft
  step: number
}>()

const emit = defineEmits<{
  back: []
  next: []
  patch: [partial: Partial<OnboardingDraft>]
  setGoalKind: [kind: GoalKind]
  patchBreakdown: [partial: Partial<OnboardingDraft['emergency_breakdown']>]
  addFixedExpense: []
  updateFixedExpense: [index: number, partial: Partial<{ title: string; amount: number }>]
  removeFixedExpense: [index: number]
  skipExpenses: []
}>()

const inputMode = ref<'text' | 'voice'>('text')

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

const canProceed = computed(() => {
  switch (props.step) {
    case 1:
      return props.draft.active_income > 0 || props.draft.passive_income > 0
    case 2:
      return true
    case 3:
      return props.draft.goal_amount >= 1000
    case 4:
      return true
    default:
      return false
  }
})

function onVoicePatch(patch: Partial<OnboardingDraft>) {
  emit('patch', patch)
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
    :secondary-action="step === 4 ? { label: 'Пропустить' } : null"
    @back="emit('back')"
    @next="emit('next')"
    @secondary="emit('skipExpenses')"
  >
    <div class="mm-onb-input-mode mb-5" role="tablist" aria-label="Способ ответа">
      <button
        type="button"
        role="tab"
        class="mm-onb-input-mode__btn"
        :class="{ 'mm-onb-input-mode__btn--active': inputMode === 'text' }"
        :aria-selected="inputMode === 'text'"
        @click="inputMode = 'text'"
      >
        Вручную
      </button>
      <button
        type="button"
        role="tab"
        class="mm-onb-input-mode__btn"
        :class="{ 'mm-onb-input-mode__btn--active': inputMode === 'voice' }"
        :aria-selected="inputMode === 'voice'"
        @click="inputMode = 'voice'"
      >
        Голосом
      </button>
    </div>

    <OnboardingSurveyVoiceBlock
      v-if="inputMode === 'voice'"
      :survey-step="step"
      :draft="draft"
      @patch="onVoicePatch"
    />

    <template v-else>
      <div v-if="step === 1" class="mm-onb-form-panel">
        <div class="space-y-2">
          <Label for="onb-active">Основной доход, ₽/мес</Label>
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
        <div class="space-y-2">
          <Label for="onb-passive">Дополнительный доход, ₽/мес</Label>
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
          <div class="mt-3 space-y-3">
            <div class="grid gap-3 sm:grid-cols-3">
              <div class="space-y-2">
                <Label for="onb-cash">Наличные</Label>
                <Input
                  id="onb-cash"
                  :model-value="draft.emergency_breakdown.cash"
                  type="number"
                  min="0"
                  @update:model-value="emit('patchBreakdown', { cash: Number($event) })"
                />
              </div>
              <div class="space-y-2">
                <Label for="onb-deposit">Вклад</Label>
                <Input
                  id="onb-deposit"
                  :model-value="draft.emergency_breakdown.deposit"
                  type="number"
                  min="0"
                  @update:model-value="emit('patchBreakdown', { deposit: Number($event) })"
                />
              </div>
              <div class="space-y-2">
                <Label for="onb-invest">Инвестиции</Label>
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
      </div>

      <div v-else-if="step === 3" class="mm-onb-form-panel">
        <div class="space-y-2">
          <Label>Тип цели</Label>
          <div class="grid grid-cols-2 gap-2">
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

      <div v-else-if="step === 4" class="mm-onb-form-panel space-y-4">
        <ul v-if="draft.fixed_expenses.length" class="space-y-3">
          <li
            v-for="(item, index) in draft.fixed_expenses"
            :key="index"
            class="grid gap-2 sm:grid-cols-[1fr_120px_auto]"
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

        <p v-if="draft.skipped_expenses" class="text-sm text-muted-foreground">
          Учтём ваши обычные траты, когда добавите покупки.
        </p>
      </div>
    </template>
  </OnboardingStepShell>
</template>
