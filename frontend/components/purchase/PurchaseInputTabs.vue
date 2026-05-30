<script setup lang="ts">
import { Camera, Check, PenLine, Receipt } from 'lucide-vue-next'
import type { ReceiptManualResponse, ReceiptVoiceResponse } from '~/types/api'

const props = withDefaults(
  defineProps<{
    /** Четвёртая вкладка «Фото чека» — на дашборде; в онбординге скрыта */
    showPhoto?: boolean
    /** Кнопка «Готово» после успеха — в модалке дашборда */
    dismissible?: boolean
  }>(),
  { showPhoto: false, dismissible: false }
)

const emit = defineEmits<{
  added: []
  confirmed: []
}>()

const baseTabs = ['voice', 'manual', 'fns'] as const
type PurchaseTab = (typeof baseTabs)[number] | 'photo'
type Phase = 'input' | 'success'

const tabOrder = computed((): PurchaseTab[] =>
  props.showPhoto ? [...baseTabs, 'photo'] : [...baseTabs]
)

const inputTab = ref<PurchaseTab>('voice')
const phase = ref<Phase>('input')
const successHint = ref('')

const activeTabIndex = computed(() => tabOrder.value.indexOf(inputTab.value))

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

function onFnsSynced() {
  showSuccess('Чек добавлен — расход учтён в ленте.')
}

function onPhotoSynced() {
  showSuccess('Чек по фото добавлен — расход учтён в ленте.')
}

function confirmSuccess() {
  phase.value = 'input'
  successHint.value = ''
  emit('confirmed')
}

watch(
  () => props.showPhoto,
  () => {
    if (!tabOrder.value.includes(inputTab.value)) {
      inputTab.value = 'voice'
    }
  }
)
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
        <button
          v-if="showPhoto"
          type="button"
          role="tab"
          class="mm-onb-tabs__btn inline-flex items-center justify-center gap-1.5"
          :class="{ 'mm-onb-tabs__btn--active': inputTab === 'photo' }"
          :aria-selected="inputTab === 'photo'"
          @click="selectTab('photo')"
        >
          <Camera class="size-3.5" aria-hidden="true" />
          Фото
        </button>
      </div>

      <div v-show="inputTab === 'voice'" role="tabpanel">
        <OnboardingPurchaseVoice @done="onVoiceDone" />
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

      <div v-if="showPhoto" v-show="inputTab === 'photo'" class="space-y-3" role="tabpanel">
        <div class="mm-onb-form-panel border-dashed bg-[color:var(--mm-primary-soft)]/30">
          <p class="text-sm font-medium text-[color:var(--mm-text)]">
            QR на бумажном чеке
          </p>
          <p class="mt-1 text-xs leading-relaxed text-[color:var(--mm-text-muted)]">
            Сфотографируйте QR — Поток проверит чек через ФНС.
          </p>
        </div>
        <DashboardPhotoReceiptPanel embedded @synced="onPhotoSynced" />
      </div>
    </template>
  </div>
</template>
