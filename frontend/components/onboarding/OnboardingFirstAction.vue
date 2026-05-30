<script setup lang="ts">
defineProps<{
  finishing: boolean
  expenseAdded: boolean
}>()

const emit = defineEmits<{
  back: []
  finish: []
  added: []
}>()
</script>

<template>
  <OnboardingStepShell
    title="Первая покупка"
    description="Одна покупка — и Поток увидит ваши траты. Доход и цель вы уже назвали; чеки с кассы — по желанию."
    show-back
    :next-label="expenseAdded ? 'Открыть дашборд' : 'Пропустить и открыть дашборд'"
    :loading="finishing"
    @back="emit('back')"
    @next="emit('finish')"
  >
    <PurchaseInputTabs @added="emit('added')" />
  </OnboardingStepShell>
</template>
