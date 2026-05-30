<script setup lang="ts">
import { ArrowRight, Target, Wallet } from 'lucide-vue-next'
import type { OnboardingDraft } from '~/types/api'

defineProps<{
  draft: OnboardingDraft
  summary: {
    income: number
    monthlySaving: number
    fixedTotal: number
    freeCashflow: number
    goalForecast: string
    runwayMonths: number | null
  }
}>()

const emit = defineEmits<{
  back: []
  next: []
}>()
</script>

<template>
  <OnboardingStepShell
    title="Вот ваша картина"
    description="Коротко — и что сделать дальше."
    show-back
    next-label="Продолжить"
    @back="emit('back')"
    @next="emit('next')"
  >
    <div class="grid gap-3 sm:grid-cols-2">
      <OnboardingMetricCard label="Доход" accent :icon="Wallet">
        {{ summary.income.toLocaleString('ru-RU') }}
        <span class="text-sm font-normal text-[color:var(--mm-text-muted)]">₽/мес</span>
        <template v-if="draft.emergency_fund" #hint>
          Подушка {{ draft.emergency_fund.toLocaleString('ru-RU') }} ₽
        </template>
      </OnboardingMetricCard>

      <OnboardingMetricCard label="Цель" :icon="Target">
        <span class="block text-sm font-medium leading-snug">{{ draft.goal_title }}</span>
        <span class="mt-1 block text-lg font-semibold tabular-nums">
          {{ draft.goal_amount.toLocaleString('ru-RU') }} ₽
        </span>
      </OnboardingMetricCard>
    </div>

    <div class="mm-onb-callout space-y-2">
      <p v-if="summary.fixedTotal">
        Постоянные платежи ~{{ summary.fixedTotal.toLocaleString('ru-RU') }} ₽/мес, после расходов
        остаётся ~{{ summary.freeCashflow.toLocaleString('ru-RU') }} ₽.
      </p>
      <p v-else-if="draft.skipped_expenses" class="text-[color:var(--mm-text-muted)]">
        Постоянные платежи пока не указаны — учтём ваши обычные траты.
      </p>
      <p v-if="summary.runwayMonths">
        Запаса хватит примерно на {{ summary.runwayMonths }} мес. при текущих тратах.
      </p>
      <p class="font-medium">{{ summary.goalForecast }}</p>
    </div>

    <div
      class="flex items-start gap-3 rounded-2xl border border-[color:var(--mm-primary)]/25 bg-[color:var(--mm-primary-soft)]/50 px-4 py-3"
    >
      <ArrowRight class="mt-0.5 size-5 shrink-0 text-[color:var(--mm-primary)]" />
      <div class="text-sm">
        <p class="font-medium text-[color:var(--mm-primary)]">Первое действие</p>
        <p class="mt-0.5 text-[color:var(--mm-text-muted)]">
          Добавьте одну покупку голосом или вручную. Чеки с кассы — по желанию; зарплату ФНС не
          видит, доход вы уже назвали.
        </p>
      </div>
    </div>
  </OnboardingStepShell>
</template>
