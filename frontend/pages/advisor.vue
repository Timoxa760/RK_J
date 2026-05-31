<script setup lang="ts">
import { ADVISOR } from '~/constants/productCopy'
import { buildWeeklyActionPrompt } from '~/utils/advisorChat'

const {
  advisorContext,
  messages,
  typing,
  chatError,
  sendMessage,
  resetChat,
  runAction,
  bootstrap
} = useAdvisorShell()

onMounted(() => {
  bootstrap()
})

const weeklyHint = computed(() => {
  const insight = advisorContext.value?.topInsight
  const action = advisorContext.value?.diagnosis?.main_action
  if (insight?.title) return insight.title
  if (action?.title) return action.title
  return null
})

const weeklyPrompt = computed(() => {
  const insight = advisorContext.value?.topInsight
  const action = advisorContext.value?.diagnosis?.main_action
  if (insight?.title) {
    return buildWeeklyActionPrompt(insight.title, insight.body ?? insight.description)
  }
  if (action?.title) {
    return buildWeeklyActionPrompt(action.title, action.description)
  }
  return null
})

async function onWeeklyHintClick() {
  if (weeklyPrompt.value) await sendMessage(weeklyPrompt.value)
}
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col">
    <p
      v-if="weeklyHint && weeklyPrompt"
      class="mb-3 shrink-0 rounded-xl border bg-muted/40 px-3 py-2 text-sm text-muted-foreground"
    >
      <span class="font-medium text-foreground">{{ ADVISOR.chatWeeklyHint }}:</span>
      {{ weeklyHint }}
      <button
        type="button"
        class="ml-1 text-primary underline-offset-2 hover:underline"
        @click="onWeeklyHintClick"
      >
        Спросить
      </button>
    </p>

    <AdvisorChat
      full-page
      show-reset
      :messages="messages"
      :typing="typing"
      :error="chatError"
      :context="advisorContext"
      class="min-h-0 flex-1"
      @send="sendMessage"
      @reset="resetChat"
      @action="runAction"
    />
  </div>
</template>
