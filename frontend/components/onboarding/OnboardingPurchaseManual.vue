<script setup lang="ts">
import { STANDARD_EXPENSE_CATEGORIES } from '~/constants/expenseCategories'

const emit = defineEmits<{
  submit: [payload: { store: string; amount: number; category: string; date: string }]
}>()

defineProps<{
  busy?: boolean
}>()

const amount = ref<number | ''>('')
const category = ref<string>(STANDARD_EXPENSE_CATEGORIES[0])
const date = ref(new Date().toISOString().slice(0, 10))

function send() {
  const n = typeof amount.value === 'number' ? amount.value : Number(amount.value)
  if (!n || n <= 0) return
  emit('submit', {
    store: category.value,
    amount: n,
    category: category.value,
    date: date.value
  })
}
</script>

<template>
  <form class="mm-onb-form-panel" @submit.prevent="send">
    <div class="mm-onb-field">
      <Label for="onb-purchase-category" class="mm-onb-field-label">Тип траты</Label>
      <Select v-model="category" :disabled="busy">
        <SelectTrigger id="onb-purchase-category">
          <SelectValue placeholder="Выберите категорию" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem v-for="c in STANDARD_EXPENSE_CATEGORIES" :key="c" :value="c">
            {{ c }}
          </SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="mm-onb-field">
      <Label for="onb-purchase-amount" class="mm-onb-field-label">Сумма, ₽</Label>
      <Input
        id="onb-purchase-amount"
        v-model.number="amount"
        type="number"
        min="1"
        step="1"
        class="text-lg"
        placeholder="3200"
        required
        :disabled="busy"
      />
    </div>

    <div class="mm-onb-field">
      <Label for="onb-purchase-date" class="mm-onb-field-label">Дата</Label>
      <Input id="onb-purchase-date" v-model="date" type="date" :disabled="busy" />
    </div>

    <Button type="submit" class="w-full" :disabled="busy">
      {{ busy ? 'Сохраняем…' : 'Добавить покупку' }}
    </Button>
  </form>
</template>
