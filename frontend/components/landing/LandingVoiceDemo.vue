<script setup lang="ts">
import { demoScenarios } from '~/constants/landingContent'

type Phase = 'listening' | 'user' | 'typing' | 'assistant' | 'done'

const USER_HOLD_MS = 800
const THINK_MS = 3000

const activeScenario = ref(0)
const phase = ref<Phase>('listening')
const reducedMotion = ref(false)
const typedLength = ref(0)
const typewriterDone = ref(false)

let timers: ReturnType<typeof setTimeout>[] = []
let typeFrame = 0

function easeOutCubic(t: number) {
  return 1 - (1 - t) ** 3
}

function clearTimers() {
  timers.forEach(clearTimeout)
  timers = []
  if (typeFrame) {
    cancelAnimationFrame(typeFrame)
    typeFrame = 0
  }
}

function schedule(fn: () => void, ms: number) {
  timers.push(setTimeout(fn, ms))
}

function typewriterDuration(text: string) {
  return Math.max(2400, Math.min(4200, text.length * 46))
}

function formatMoneyInText(text: string) {
  return text.replace(/(\d[\d\s]*)\s*₽/g, (_, amount: string) => {
    const normalized = amount.trim().replace(/\s+/g, '\u00a0')
    return `${normalized}\u00a0₽`
  })
}

function startTypewriter(onComplete: () => void) {
  if (typeFrame) {
    cancelAnimationFrame(typeFrame)
    typeFrame = 0
  }

  const text = userLine.value?.text ?? ''
  typedLength.value = 0
  typewriterDone.value = false

  if (reducedMotion.value || !text) {
    typedLength.value = text.length
    typewriterDone.value = true
    onComplete()
    return
  }

  const duration = typewriterDuration(text)
  const started = performance.now()

  const tick = (now: number) => {
    const progress = Math.min(1, (now - started) / duration)
    typedLength.value = Math.max(0, Math.round(text.length * easeOutCubic(progress)))
    if (progress < 1) {
      typeFrame = requestAnimationFrame(tick)
      return
    }

    typedLength.value = text.length
    typewriterDone.value = true
    onComplete()
  }

  typeFrame = requestAnimationFrame(tick)
}

function runSequence() {
  clearTimers()
  typewriterDone.value = false
  phase.value = reducedMotion.value ? 'done' : 'listening'

  if (reducedMotion.value) {
    typedLength.value = userLine.value?.text.length ?? 0
    return
  }

  startTypewriter(() => {
    schedule(() => {
      phase.value = 'user'

      schedule(() => {
        phase.value = 'typing'

        schedule(() => {
          phase.value = 'assistant'
        }, THINK_MS)
      }, USER_HOLD_MS)
    }, 200)
  })
}

function selectScenario(index: number) {
  if (index === activeScenario.value) return
  clearTimers()
  activeScenario.value = index
  typedLength.value = 0
  typewriterDone.value = false
  phase.value = 'listening'
  nextTick(() => runSequence())
}

const scenario = computed(() => demoScenarios[activeScenario.value]!)
const userLine = computed(() => scenario.value.lines.find((l) => l.role === 'user'))
const assistantLine = computed(() => scenario.value.lines.find((l) => l.role === 'assistant'))

const responseMode = computed(() => {
  if (phase.value === 'listening' || phase.value === 'user') return 'idle'
  if (phase.value === 'typing') return 'typing'
  if (phase.value === 'assistant' || phase.value === 'done') return 'assistant'
  return 'idle'
})

const isListening = computed(() => phase.value === 'listening')
const isParsing = computed(() => phase.value === 'typing')
const micAmbient = computed(() => phase.value === 'user' || phase.value === 'assistant' || phase.value === 'done')

const dockTranscriptFull = computed(() => formatMoneyInText(userLine.value?.text ?? ''))

const dockTranscript = computed(() => {
  const text = userLine.value?.text ?? ''
  if (!text) return ''
  let slice = text
  if (phase.value === 'listening' && !reducedMotion.value) {
    slice = text.slice(0, typedLength.value)
  }
  if (phase.value === 'listening' || phase.value === 'user' || phase.value === 'typing' || phase.value === 'assistant' || phase.value === 'done') {
    return formatMoneyInText(slice)
  }
  return ''
})

const assistantText = computed(() => formatMoneyInText(assistantLine.value?.text ?? ''))

const dockStatus = computed(() => {
  if (phase.value === 'listening') return 'Поток слушает…'
  if (phase.value === 'user') return 'Расход записан'
  if (isParsing.value) return 'Готовит ответ…'
  if (responseMode.value === 'assistant') return 'Поток ответил'
  return 'Скажите, что потратили'
})

const chartMax = computed(() => Math.max(...scenario.value.chartBars, 1))

onMounted(() => {
  reducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  runSequence()
})

onUnmounted(clearTimers)
</script>

<template>
  <div
    id="demo"
    class="mm-landing-demo mm-landing-demo--chat mm-voice-demo mm-landing-glass mm-landing-card mm-landing-card--glow relative flex w-full flex-col p-4 sm:p-5 lg:p-6"
  >
    <div class="mm-landing-demo-tabs shrink-0">
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

    <div class="mm-landing-demo-body mt-3 flex min-h-0 flex-1 flex-col sm:mt-4">
      <div class="mm-landing-demo-response flex min-h-0 flex-1 flex-col">
        <p class="mm-landing-demo-response__context mb-3 shrink-0 text-xs text-[color:var(--mm-text-soft)]">
          {{ scenario.context }}
        </p>

        <div
          class="mm-landing-demo-response__shell"
          aria-live="polite"
          aria-relevant="additions text"
        >
          <div class="mm-landing-demo-response__panel">
            <Transition name="mm-landing-demo-panel" mode="out-in">
              <div
                v-if="responseMode === 'idle'"
                :key="`idle-${activeScenario}`"
                class="mm-landing-demo-response__idle flex h-full min-h-full items-center justify-center rounded-2xl border border-dashed border-[color:var(--mm-border-subtle)] bg-white/40 px-4 py-6 text-center"
              >
                <p class="mm-landing-demo-response__idle-text text-sm leading-relaxed text-[color:var(--mm-text-soft)]">
                  <Transition name="mm-landing-demo-idle" mode="out-in">
                    <span v-if="phase === 'listening' && !typewriterDone" key="listen">Слушаем вашу фразу…</span>
                    <span v-else key="recorded">Расход записан — Поток готовит разбор</span>
                  </Transition>
                </p>
              </div>

              <div
                v-else-if="responseMode === 'typing'"
                :key="`typing-${activeScenario}`"
                class="mm-landing-demo-response__thinking flex h-full min-h-full flex-col justify-end gap-2 pt-1"
              >
                <p class="text-xs text-[color:var(--mm-text-soft)]">Поток думает…</p>
                <div class="mm-landing-demo-msg mm-landing-demo-msg--assistant">
                  <div
                    class="mm-landing-chat-bubble mm-landing-chat-bubble--typing"
                    aria-label="Поток печатает"
                  >
                    <span></span><span></span><span></span>
                  </div>
                </div>
              </div>

              <div
                v-else
                :key="`assistant-${activeScenario}`"
                class="mm-landing-demo-response__answer flex h-full min-h-full flex-col gap-3 pt-1"
              >
                <div class="mm-landing-demo-msg mm-landing-demo-msg--assistant">
                  <div class="mm-landing-chat-bubble mm-landing-chat-bubble--assistant">
                    {{ assistantText }}
                  </div>
                </div>

                <div class="mm-landing-demo-chart">
                  <div class="mm-landing-demo-chart__head">
                    <span class="mm-landing-demo-chart__label">{{ scenario.chartLabel }}</span>
                    <span class="mm-landing-demo-chart__metric">{{ scenario.metric }}</span>
                  </div>
                  <div class="mm-landing-demo-chart__bars" role="img" :aria-label="`${scenario.chartLabel}: ${scenario.metric}`">
                    <div
                      v-for="(bar, index) in scenario.chartBars"
                      :key="`${scenario.id}-${index}`"
                      class="mm-landing-demo-chart__bar"
                      :class="{ 'mm-landing-demo-chart__bar--active': index === scenario.highlightIndex }"
                      :style="{ height: `${(bar / chartMax) * 100}%`, animationDelay: `${0.08 + index * 0.07}s` }"
                    />
                  </div>
                </div>
              </div>
            </Transition>
          </div>
        </div>
      </div>

      <div
        class="mm-landing-demo-voice-dock shrink-0"
        :class="{
          'mm-landing-demo-voice-dock--listen': isListening || isParsing,
          'mm-landing-demo-voice-dock--done': responseMode === 'assistant'
        }"
      >
        <p class="mm-landing-demo-voice-dock__transcript">
          <span v-if="dockTranscriptFull" class="mm-landing-demo-voice-dock__transcript-ghost" aria-hidden="true">
            «{{ dockTranscriptFull }}»
          </span>
          <span class="mm-landing-demo-voice-dock__transcript-live">
            <span v-if="dockTranscript">«{{ dockTranscript }}»</span>
            <span v-else class="text-[color:var(--mm-text-soft)]">…</span>
          </span>
        </p>

        <div class="mm-landing-demo-voice-dock__orb" aria-hidden="true">
          <OnboardingMicOrbVisual
            dock
            gentle
            :listening="isListening"
            :parsing="isParsing"
            :ambient="micAmbient"
          />
        </div>

        <p class="mm-landing-demo-voice-dock__status">
          <span
            class="mm-landing-demo-status-dot"
            :class="{ 'mm-landing-demo-status-dot--listen': isListening || isParsing }"
          />
          <span class="mm-landing-demo-voice-dock__status-text">{{ dockStatus }}</span>
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mm-landing-demo--chat {
  min-height: 26.5rem;
}

.mm-landing-demo-response__idle-text {
  min-height: 2.75rem;
}

.mm-landing-demo-idle-enter-active,
.mm-landing-demo-idle-leave-active {
  transition: opacity 0.45s cubic-bezier(0.22, 1, 0.36, 1);
}

.mm-landing-demo-idle-enter-from,
.mm-landing-demo-idle-leave-to {
  opacity: 0;
}

.mm-landing-demo-body {
  transition: opacity 0.35s ease;
}

.mm-landing-demo-response__shell {
  position: relative;
  height: 13.75rem;
  flex-shrink: 0;
}

.mm-landing-demo-response__panel {
  position: relative;
  height: 100%;
  overflow-x: hidden;
  overflow-y: auto;
  overscroll-behavior: contain;
}

.mm-landing-demo-panel-enter-active,
.mm-landing-demo-panel-leave-active {
  transition:
    opacity 0.55s cubic-bezier(0.22, 1, 0.36, 1),
    transform 0.55s cubic-bezier(0.22, 1, 0.36, 1);
}

.mm-landing-demo-panel-leave-active {
  position: absolute;
  inset: 0;
  width: 100%;
  pointer-events: none;
}

.mm-landing-demo-panel-enter-active {
  position: relative;
  width: 100%;
}

.mm-landing-demo-text-enter-active,
.mm-landing-demo-text-leave-active {
  transition:
    opacity 0.45s cubic-bezier(0.22, 1, 0.36, 1),
    transform 0.45s cubic-bezier(0.22, 1, 0.36, 1);
}

.mm-landing-demo-panel-enter-from,
.mm-landing-demo-panel-leave-to,
.mm-landing-demo-text-enter-from,
.mm-landing-demo-text-leave-to {
  opacity: 0;
  transform: translate3d(0, 6px, 0);
}

.mm-landing-demo-msg {
  display: flex;
  width: 100%;
}

.mm-landing-demo-msg--assistant {
  justify-content: flex-start;
}

.mm-landing-demo-msg--assistant .mm-landing-chat-bubble--assistant,
.mm-landing-demo-msg--assistant .mm-landing-chat-bubble--typing {
  margin-right: 0;
}

.mm-landing-demo-chart {
  border-radius: 1rem;
  border: 1px solid var(--mm-border-subtle);
  background: rgb(255 255 255 / 0.72);
  padding: 0.875rem 1rem 1rem;
  animation: mm-landing-demo-chart-in 0.7s cubic-bezier(0.22, 1, 0.36, 1) both;
}

@keyframes mm-landing-demo-chart-in {
  from {
    opacity: 0;
    transform: translate3d(0, 8px, 0);
  }
  to {
    opacity: 1;
    transform: translate3d(0, 0, 0);
  }
}

.mm-landing-demo-chart__head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.mm-landing-demo-chart__label {
  min-width: 0;
  font-size: 0.75rem;
  line-height: 1.3;
  color: var(--mm-text-soft);
}

.mm-landing-demo-chart__metric {
  flex-shrink: 0;
  font-size: 0.75rem;
  font-weight: 600;
  line-height: 1.3;
  color: var(--mm-landing-brand);
  white-space: nowrap;
}

.mm-landing-demo-chart__bars {
  display: flex;
  align-items: flex-end;
  gap: 0.35rem;
  height: 4.5rem;
}

.mm-landing-demo-chart__bar {
  flex: 1;
  min-width: 0;
  border-radius: 0.375rem 0.375rem 0.125rem 0.125rem;
  background: rgb(var(--mm-landing-brand-rgb) / 0.18);
  transform-origin: bottom;
  animation: mm-landing-demo-bar-in 0.85s cubic-bezier(0.22, 1, 0.36, 1) both;
}

.mm-landing-demo-chart__bar--active {
  background: linear-gradient(180deg, var(--mm-landing-brand-light) 0%, var(--mm-landing-brand) 100%);
}

@keyframes mm-landing-demo-bar-in {
  from {
    transform: scaleY(0.2);
    opacity: 0.2;
  }
  to {
    transform: scaleY(1);
    opacity: 1;
  }
}

.mm-landing-demo-voice-dock {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  grid-template-rows: auto 5.5rem 1.125rem;
  justify-items: center;
  align-items: start;
  gap: 0.25rem;
  margin-top: 0.75rem;
  border-top: 1px solid var(--mm-border-subtle);
  padding: 0.625rem 0 0.25rem;
}

.mm-landing-demo-voice-dock__transcript {
  position: relative;
  display: block;
  width: 100%;
  max-height: 2.5rem;
  overflow: hidden;
  justify-self: stretch;
  align-self: end;
  text-align: center;
  font-size: 0.75rem;
  line-height: 1.35;
  color: var(--mm-text);
  transition: color 0.6s ease;
  text-wrap: pretty;
  word-break: normal;
  overflow-wrap: normal;
  hyphens: none;
}

.mm-landing-demo-voice-dock__transcript-ghost {
  display: block;
  visibility: hidden;
  pointer-events: none;
  user-select: none;
}

.mm-landing-demo-voice-dock__transcript-live {
  position: absolute;
  inset: 0;
  display: -webkit-box;
  overflow: hidden;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  align-content: center;
}

@media (min-width: 640px) {
  .mm-landing-demo-response__shell {
    height: 14.5rem;
  }

  .mm-landing-demo--chat {
    min-height: 27.5rem;
  }

  .mm-landing-demo-voice-dock__transcript {
    font-size: 0.8125rem;
    line-height: 1.4;
  }

  .mm-landing-demo-chart__label,
  .mm-landing-demo-chart__metric {
    font-size: 0.8125rem;
  }
}

.mm-landing-demo-voice-dock--done .mm-landing-demo-voice-dock__transcript {
  color: var(--mm-text-soft);
}

.mm-landing-demo-voice-dock__orb {
  display: flex;
  width: 5.5rem;
  height: 5.5rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  justify-self: center;
  align-self: start;
  margin-inline: auto;
  overflow: hidden;
  border-radius: 9999px;
  isolation: isolate;
}

.mm-landing-demo-voice-dock__orb :deep(.mm-onb-mic-orb-visual) {
  transform: none !important;
}

.mm-landing-demo-voice-dock__status {
  display: flex;
  height: 1.125rem;
  align-items: center;
  justify-content: center;
  justify-self: center;
  gap: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--mm-text-soft);
}

.mm-landing-demo-voice-dock__status-text {
  transition: opacity 0.4s ease, color 0.5s ease;
}

.mm-landing-demo-voice-dock--listen .mm-landing-demo-voice-dock__status {
  color: var(--mm-landing-brand);
}

@media (prefers-reduced-motion: reduce) {
  .mm-landing-demo-panel-enter-active,
  .mm-landing-demo-panel-leave-active,
  .mm-landing-demo-text-enter-active,
  .mm-landing-demo-text-leave-active {
    transition: none;
  }

  .mm-landing-demo-chart {
    animation: none;
  }

  .mm-landing-demo-chart__bar {
    animation: none;
  }

  .mm-landing-demo-idle-enter-active,
  .mm-landing-demo-idle-leave-active {
    transition: none;
  }
}
</style>
