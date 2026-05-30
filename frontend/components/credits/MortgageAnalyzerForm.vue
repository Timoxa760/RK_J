<script setup lang="ts">
defineProps<{
  mortgageAmount: number
  monthlyIncome?: number
  savings?: number
  loading?: boolean
}>()

const emit = defineEmits<{
  analyze: [amount: number]
}>()

const amount = defineModel<number>('amount', { default: 12_000_000 })

function submit() {
  if (amount.value > 0) emit('analyze', amount.value)
}
</script>

<template>
  <Card data-demo="mortgage-form">
    <CardHeader>
      <CardTitle class="text-base">Ипотечный разбор</CardTitle>
      <CardDescription>
        Скажите сумму — Поток прикинет платёж, запас и сравнит банки. Один ответ вместо таблиц.
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form class="flex flex-col gap-3 sm:flex-row sm:items-end" @submit.prevent="submit">
        <div class="min-w-0 flex-1 space-y-2">
          <Label for="mortgage-amount">Сумма ипотеки, ₽</Label>
          <Input
            id="mortgage-amount"
            v-model.number="amount"
            type="number"
            min="500_000"
            step="100_000"
            :disabled="loading"
          />
          <p v-if="monthlyIncome" class="text-xs text-muted-foreground">
            Доход {{ monthlyIncome.toLocaleString('ru-RU') }} ₽/мес
            <span v-if="savings"> · подушка {{ savings.toLocaleString('ru-RU') }} ₽</span>
          </p>
        </div>
        <Button type="submit" class="w-full shrink-0 sm:w-auto" :disabled="loading">
          {{ loading ? 'Считаем…' : 'Проверить' }}
        </Button>
      </form>
    </CardContent>
  </Card>
</template>
