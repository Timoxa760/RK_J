<script setup lang="ts">
import { Check, Mic, PenLine } from 'lucide-vue-next'
import type { ReceiptManualResponse, ReceiptVoiceResponse } from '~/types/api'

withDefaults(
  defineProps<{
    /** Кнопка «Готово» после успеха — в модалке дашборда */
    dismissible?: boolean
  }>(),
  { dismissible: false }
)

const emit = defineEmits<{
  added: []
  confirmed: []
}>()

const baseTabs = ['voice', 'manual'] as const
type PurchaseTab = (typeof baseTabs)[number]
type Phase = 'input' | 'success'

const inputTab = ref<PurchaseTab>('voice')
const phase = ref<Phase>('input')
const successHint = ref('')

const activeTabIndex = computed(() => baseTabs.indexOf(inputTab.value))

const { submitting, submitManual, lastResult } = useReceiptSubmit()

function formatPurchaseResult(
  result: ReceiptManualResponse | ReceiptVoiceResponse | null
): string | null {
  if (!result) return null
  if ('total' in result) {
    return `${result.store} — ${result.total.toLocaleString('ru-RU')} ₽ · ${result.category}`
  }
  return `${result.store} — ${result.amount.toLocaleString('ru-RU')} ₽ · ${result.category}`
}

function showSuccess(fallback?: string) {
  successHint.value =
    formatPurchaseResult(lastResult.value) ?? fallback ?? 'Покупка отправлена — Поток обновит картину денег.'
  phase.value = 'success'
  emit('added')
}

function selectTab(id: PurchaseTab) {
  if (phase.value === 'success') return
  inputTab.value = id
}

function onVoiceDone() {
  showSuccess()
}

async function onManual(payload: {
  store: string
  amount: number
  category: string
  date: string
}) {
  await submitManual(payload)
  showSuccess()
}

function confirmSuccess() {
  phase.value = 'input'
  successHint.value = ''
  emit('confirmed')
}
</script>

<template>
  <div>
    <div
      v-if="phase === 'success'"
      class="mm-purchase-success"
      role="status"
      aria-live="polite"
    >
      <div class="mm-purchase-success__icon" aria-hidden="true">
        <Check class="size-8 stroke-[2.5]" />
      </div>
      <p class="mm-purchase-success__title">Записали</p>
      <p class="mm-purchase-success__hint">{{ successHint }}</p>
      <Button
        v-if="dismissible"
        type="button"
        class="mm-purchase-success__btn mt-2 w-full"
        @click="confirmSuccess"
      >
        Готово
      </Button>
    </div>

    <template v-else>
      <div class="mm-onb-tabs mb-5" role="tablist" aria-label="Способ ввода покупки">
        <div
          class="mm-onb-tabs__indicator"
          :style="{
            width: `${100 / baseTabs.length}%`,
            transform: `translateX(${activeTabIndex * 100}%)`
          }"
          aria-hidden="true"
        />
        <button
          type="button"
          role="tab"
          class="mm-onb-tabs__btn inline-flex items-center justify-center gap-1.5"
          :class="{ 'mm-onb-tabs__btn--active': inputTab === 'voice' }"
          :aria-selected="inputTab === 'voice'"
          @click="selectTab('voice')"
        >
          <Mic class="size-3.5" aria-hidden="true" />
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
      </div>

      <div v-show="inputTab === 'voice'" role="tabpanel">
        <OnboardingPurchaseVoice @done="onVoiceDone" />
      </div>

      <div v-show="inputTab === 'manual'" role="tabpanel">
        <OnboardingPurchaseManual :busy="submitting" @submit="onManual" />
      </div>
    </template>
  </div>
</template>
