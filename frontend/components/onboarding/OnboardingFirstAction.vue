<script setup lang="ts">
import { PenLine, Receipt } from 'lucide-vue-next'

defineProps<{
  finishing: boolean
  expenseAdded: boolean
}>()

const emit = defineEmits<{
  back: []
  finish: []
  added: []
}>()

const inputTab = ref<'voice' | 'manual' | 'fns'>('voice')
const tabOrder = ['voice', 'manual', 'fns'] as const

const activeTabIndex = computed(() => tabOrder.indexOf(inputTab.value))

const { submitting, submitManual, submitVoiceAudio } = useReceiptSubmit()

async function onVoice(audio: Blob) {
  await submitVoiceAudio(audio)
  emit('added')
}

async function onManual(payload: {
  store: string
  amount: number
  category: string
  date: string
}) {
  await submitManual(payload)
  emit('added')
}

function onFnsSynced() {
  emit('added')
}

function selectTab(id: (typeof tabOrder)[number]) {
  inputTab.value = id
}
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
    <div class="mm-onb-tabs mb-5" role="tablist" aria-label="Способ ввода покупки">
      <div
        class="mm-onb-tabs__indicator"
        :style="{
          width: `${100 / tabOrder.length}%`,
          transform: `translateX(${activeTabIndex * 100}%)`
        }"
        aria-hidden="true"
      />
      <button
        type="button"
        role="tab"
        class="mm-onb-tabs__btn"
        :class="{ 'mm-onb-tabs__btn--active': inputTab === 'voice' }"
        :aria-selected="inputTab === 'voice'"
        @click="selectTab('voice')"
      >
        Голос
      </button>
      <button
        type="button"
        role="tab"
        class="mm-onb-tabs__btn inline-flex items-center justify-center gap-1.5"
        :class="{ 'mm-onb-tabs__btn--active': inputTab === 'manual' }"
        :aria-selected="inputTab === 'manual'"
        @click="selectTab('manual')"
      >
        <PenLine class="size-3.5" aria-hidden="true" />
        Вручную
      </button>
      <button
        type="button"
        role="tab"
        class="mm-onb-tabs__btn inline-flex items-center justify-center gap-1.5"
        :class="{ 'mm-onb-tabs__btn--active': inputTab === 'fns' }"
        :aria-selected="inputTab === 'fns'"
        @click="selectTab('fns')"
      >
        <Receipt class="size-3.5" aria-hidden="true" />
        Чеки
      </button>
    </div>

    <div v-show="inputTab === 'voice'" role="tabpanel">
      <OnboardingPurchaseVoice :busy="submitting" @submit="onVoice" />
    </div>

    <div v-show="inputTab === 'manual'" role="tabpanel">
      <OnboardingPurchaseManual :busy="submitting" @submit="onManual" />
    </div>

    <div v-show="inputTab === 'fns'" class="space-y-3" role="tabpanel">
      <div class="mm-onb-form-panel border-dashed bg-[color:var(--mm-primary-soft)]/30">
        <p class="text-sm font-medium text-[color:var(--mm-text)]">
          Только расходы по чекам
        </p>
        <p class="mt-1 text-xs leading-relaxed text-[color:var(--mm-text-muted)]">
          Чеки с кассы — по желанию. Зарплату и накопления сервис не видит, только покупки.
        </p>
      </div>
      <DashboardFnsExpensePanel @synced="onFnsSynced" />
    </div>

    <p
      v-if="expenseAdded"
      class="mt-4 text-center text-sm font-medium text-[color:var(--mm-primary)]"
    >
      Покупка добавлена — можно открыть дашборд
    </p>
  </OnboardingStepShell>
</template>
