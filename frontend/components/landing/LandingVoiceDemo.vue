<script setup lang="ts">
import { Mic } from 'lucide-vue-next'
import { demoScenarios } from '~/constants/landingContent'

type Phase = 'listening' | 'user' | 'typing' | 'assistant' | 'done'

const activeScenario = ref(0)
const phase = ref<Phase>('listening')
const reducedMotion = ref(false)

let timers: ReturnType<typeof setTimeout>[] = []

function clearTimers() {
  timers.forEach(clearTimeout)
  timers = []
}

function schedule(fn: () => void, ms: number) {
  timers.push(setTimeout(fn, ms))
}

function runSequence() {
  clearTimers()
  phase.value = reducedMotion.value ? 'done' : 'listening'

  if (reducedMotion.value) return

  schedule(() => {
    phase.value = 'user'
  }, 900)

  schedule(() => {
    phase.value = 'typing'
  }, 2200)

  schedule(() => {
    phase.value = 'assistant'
  }, 4000)
}

function selectScenario(index: number) {
  if (index === activeScenario.value) return
  activeScenario.value = index
  runSequence()
}

const scenario = computed(() => demoScenarios[activeScenario.value]!)
const userLine = computed(() => scenario.value.lines.find((l) => l.role === 'user'))
const assistantLine = computed(() => scenario.value.lines.find((l) => l.role === 'assistant'))

const isUserPreview = computed(() => phase.value === 'listening')
const showUserClear = computed(
  () => phase.value === 'user' || phase.value === 'typing' || phase.value === 'assistant' || phase.value === 'done'
)
const showTyping = computed(() => phase.value === 'typing')
const showAssistant = computed(
  () => phase.value === 'assistant' || phase.value === 'done'
)
const isListening = computed(() => phase.value === 'listening')

const statusLabel = computed(() => {
  if (isListening.value) return 'Поток слушает'
  if (showTyping.value) return 'Готовит ответ…'
  return 'Поток отвечает'
})

onMounted(() => {
  reducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  runSequence()
})

onUnmounted(clearTimers)
</script>

<template>
  <div
    id="demo"
    class="mm-landing-demo mm-voice-demo mm-landing-glass mm-landing-card mm-landing-card--glow relative w-full p-4 sm:p-5 lg:p-6"
  >
    <div class="mm-landing-demo-tabs mb-3 sm:mb-4">
      <div
        class="mm-landing-demo-tabs__indicator"
        :style="{
          width: `${100 / demoScenarios.length}%`,
          transform: `translateX(${activeScenario * 100}%)`
        }"
        aria-hidden="true"
      />
      <button
        v-for="(item, index) in demoScenarios"
        :key="item.id"
        type="button"
        class="mm-landing-demo-tabs__btn"
        :class="{ 'mm-landing-demo-tabs__btn--active': index === activeScenario }"
        @click="selectScenario(index)"
      >
        {{ item.label }}
      </button>
    </div>

    <div class="mm-landing-demo-chat mt-3 sm:mt-4">
      <div class="mm-landing-demo-status mb-2 flex items-center gap-2">
        <span
          class="mm-landing-demo-status-dot"
          :class="{ 'mm-landing-demo-status-dot--listen': isListening }"
        />
        <p class="mm-landing-demo-status__label text-xs font-medium text-[color:var(--mm-text)] sm:text-sm">
          {{ statusLabel }}
        </p>
        <span class="text-xs text-[color:var(--mm-text-soft)]">· {{ scenario.context }}</span>
      </div>

      <div class="mm-landing-demo-messages">
        <div
          v-if="userLine"
          class="mm-landing-demo-user-slot"
        >
          <div
            class="mm-landing-demo-layer mm-landing-chat-bubble mm-landing-chat-bubble--user mm-landing-chat-bubble--voice"
            :class="{
              'mm-landing-demo-layer--preview': isUserPreview,
              'mm-landing-demo-layer--visible': showUserClear
            }"
            :aria-hidden="!isUserPreview && !showUserClear"
          >
            <span class="mm-landing-chat-bubble__voice-tag">
              <Mic class="h-3 w-3" stroke-width="2" />
              Голосовое
            </span>
            <p class="mt-1.5">{{ userLine.text }}</p>
          </div>
        </div>

        <div class="mm-landing-assistant-slot">
          <div
            class="mm-landing-demo-layer mm-landing-chat-bubble mm-landing-chat-bubble--typing"
            :class="{ 'mm-landing-demo-layer--visible': showTyping }"
            :aria-hidden="!showTyping"
            aria-label="Поток печатает"
          >
            <span /><span /><span />
          </div>
          <div
            v-if="assistantLine"
            class="mm-landing-demo-layer mm-landing-chat-bubble mm-landing-chat-bubble--assistant"
            :class="{ 'mm-landing-demo-layer--visible': showAssistant }"
            :aria-hidden="!showAssistant"
          >
            {{ assistantLine.text }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mm-landing-demo-chat {
  overflow: visible;
}

.mm-landing-demo-status__label {
  transition: opacity 0.5s cubic-bezier(0.22, 1, 0.36, 1);
}

.mm-landing-demo-messages {
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
  min-height: 13.5rem;
  overflow: visible;
}

.mm-landing-demo-user-slot {
  position: relative;
  min-height: 5.25rem;
}

.mm-landing-assistant-slot {
  position: relative;
  flex: 1;
  min-height: 7.75rem;
}

.mm-landing-demo-layer {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  opacity: 0;
  transform: translate3d(0, 6px, 0);
  visibility: hidden;
  filter: blur(0);
  transition:
    opacity 0.55s cubic-bezier(0.16, 1, 0.3, 1),
    transform 0.55s cubic-bezier(0.16, 1, 0.3, 1),
    filter 0.65s cubic-bezier(0.16, 1, 0.3, 1),
    visibility 0.55s;
  will-change: opacity, transform, filter;
}

.mm-landing-demo-layer--preview {
  opacity: 0.82;
  visibility: visible;
  transform: translate3d(0, 0, 0);
  filter: blur(7px);
  transition:
    opacity 0.35s ease,
    filter 0.35s ease,
    transform 0.35s ease,
    visibility 0s;
}

.mm-landing-demo-layer--visible {
  opacity: 1;
  transform: translate3d(0, 0, 0);
  visibility: visible;
  filter: blur(0);
}

.mm-landing-chat-bubble--typing {
  width: fit-content;
  max-width: 100%;
}

.mm-landing-chat-bubble--assistant {
  height: auto;
}

@media (prefers-reduced-motion: reduce) {
  .mm-landing-demo-layer {
    transition: none;
  }

  .mm-landing-demo-layer:not(.mm-landing-demo-layer--visible):not(.mm-landing-demo-layer--preview) {
    display: none;
  }
}
</style>
