<script setup lang="ts">
import { isAppFeatureEnabled } from '~/constants/featureFlags'

const ACTS = [
  { path: '/dashboard', label: 'Картина недели', selector: '[data-demo="narrative"]' },
  { path: '/dashboard', label: 'Финансовый план', selector: '[data-demo="financial-plan"]' },
  { path: '/dashboard', label: 'Советник', selector: '[data-demo="advisor-chat"]' },
  { path: '/dashboard', label: 'Добавить трату', selector: '[data-demo="add-expense"]' },
  { path: '/receipts', label: 'Лента чеков', selector: '[data-demo="receipts"]' },
  ...(isAppFeatureEnabled('creditsNav')
    ? [
        { path: '/credits', label: 'Ипотека', selector: '[data-demo="mortgage-form"]' } as const,
        { path: '/credits', label: 'Сравнение банков', selector: '[data-demo="bank-compare"]' } as const
      ]
    : [])
] as const

const active = ref(false)
const step = ref(0)
let timer: ReturnType<typeof setTimeout> | null = null

async function goToStep(index: number) {
  step.value = index
  const act = ACTS[index]
  if (!act) return
  await navigateTo(act.path)
  await nextTick()
  setTimeout(() => highlight(act.selector), 400)
}

function highlight(selector: string) {
  document.querySelectorAll('.demo-highlight').forEach((el) => {
    el.classList.remove('demo-highlight')
  })
  const el = document.querySelector(selector)
  if (el) {
    el.classList.add('demo-highlight')
    el.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

function scheduleNext() {
  if (timer) clearTimeout(timer)
  if (!active.value) return
  timer = setTimeout(() => {
    if (step.value < ACTS.length - 1) {
      goToStep(step.value + 1)
      scheduleNext()
    } else {
      stop()
    }
  }, 5000)
}

async function start() {
  active.value = true
  step.value = 0
  await goToStep(0)
  scheduleNext()
}

function stop() {
  active.value = false
  step.value = 0
  if (timer) clearTimeout(timer)
  document.querySelectorAll('.demo-highlight').forEach((el) => {
    el.classList.remove('demo-highlight')
  })
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && active.value) stop()
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  if (timer) clearTimeout(timer)
})

defineExpose({ start, stop })
</script>

<template>
  <div
    v-if="active"
    class="fixed inset-x-0 bottom-0 z-50 border-t border-[color:var(--mm-border)] bg-[color:var(--mm-surface)]/95 px-4 py-3 backdrop-blur mm-safe-bottom"
  >
    <div class="mx-auto flex max-w-4xl flex-wrap items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <span class="text-xs font-medium uppercase tracking-wide text-[color:var(--mm-primary)]">Режим демо</span>
        <span class="text-sm text-[color:var(--mm-text)]">{{ ACTS[step]?.label }}</span>
      </div>
      <div class="flex gap-1">
        <button
          v-for="(_, i) in ACTS"
          :key="i"
          type="button"
          class="h-2 w-8 rounded-full transition-colors"
          :class="i === step ? 'bg-[color:var(--mm-primary)]' : i < step ? 'bg-[color:var(--mm-primary-muted)]' : 'bg-[color:var(--mm-border)]'"
          :aria-label="`Шаг ${i + 1}`"
          @click="goToStep(i)"
        />
      </div>
      <button
        type="button"
        class="text-xs text-[color:var(--mm-text-soft)] hover:text-[color:var(--mm-text)]"
        @click="stop"
      >
        Esc — выход
      </button>
    </div>
  </div>
</template>
