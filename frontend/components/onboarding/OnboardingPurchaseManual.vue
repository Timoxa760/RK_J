<script setup lang="ts">
const emit = defineEmits<{
  submit: [payload: { store: string; amount: number; category: string; date: string }]
}>()

defineProps<{
  busy?: boolean
}>()

const amount = ref<number | ''>('')
const store = ref('')
const category = ref('Продукты')
const date = ref(new Date().toISOString().slice(0, 10))

const categories = [
  'Продукты',
  'Кафе и рестораны',
  'Транспорт',
  'Одежда',
  'Развлечения',
  'Здоровье',
  'Прочее'
]

function send() {
  const n = typeof amount.value === 'number' ? amount.value : Number(amount.value)
  if (!n || n <= 0) return
  emit('submit', {
    store: store.value.trim() || 'Не указан',
    amount: n,
    category: category.value,
    date: date.value
  })
}
</script>

<template>
  <form class="mm-onb-form-panel space-y-4" @submit.prevent="send">
    <div class="space-y-2">
      <Label for="onb-purchase-amount">Сумма покупки, ₽</Label>
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

    <div class="space-y-2">
      <Label for="onb-purchase-store">Магазин</Label>
      <Input
        id="onb-purchase-store"
        v-model="store"
        placeholder="Пятёрочка, такси…"
        :disabled="busy"
      />
    </div>

    <div class="space-y-2">
      <Label for="onb-purchase-category">Категория</Label>
      <Select v-model="category" :disabled="busy">
        <SelectTrigger id="onb-purchase-category">
          <SelectValue placeholder="Категория" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem v-for="c in categories" :key="c" :value="c">
            {{ c }}
          </SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-2">
      <Label for="onb-purchase-date">Дата</Label>
      <Input id="onb-purchase-date" v-model="date" type="date" :disabled="busy" />
    </div>

    <Button type="submit" class="w-full" :disabled="busy">
      {{ busy ? 'Сохраняем…' : 'Добавить покупку' }}
    </Button>
  </form>
</template>
