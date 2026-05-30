<script setup lang="ts">
const emit = defineEmits<{
  submit: [payload: { store: string; amount: number; category: string; date: string }]
}>()

defineProps<{ busy?: boolean }>()

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
  <form class="space-y-4" @submit.prevent="send">
    <p class="text-sm text-muted-foreground">
      Быстрый ручной ввод — одна покупка без голоса. Подходит, если не хотите подключать ФНС.
    </p>

    <div class="space-y-2">
      <Label for="amount">Сумма, ₽</Label>
      <Input id="amount" v-model.number="amount" type="number" min="1" required :disabled="busy" />
    </div>

    <div class="space-y-2">
      <Label for="store">Магазин</Label>
      <Input
        id="store"
        v-model="store"
        placeholder="Пятёрочка, Ozon…"
        :disabled="busy"
      />
    </div>

    <div class="space-y-2">
      <Label for="category">Категория</Label>
      <Select v-model="category" :disabled="busy">
        <SelectTrigger id="category">
          <SelectValue placeholder="Категория" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem v-for="c in categories" :key="c" :value="c">{{ c }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-2">
      <Label for="date">Дата</Label>
      <Input id="date" v-model="date" type="date" :disabled="busy" />
    </div>

    <Button type="submit" :disabled="busy">
      {{ busy ? 'Сохраняем…' : 'Добавить расход' }}
    </Button>
  </form>
</template>
