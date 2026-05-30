<script setup lang="ts">
import { MessageCircle } from 'lucide-vue-next'
import type { Goal } from '~/types/api'
import { ADVISOR } from '~/constants/productCopy'
import {
  buildGettingStartedPrompt,
  buildGoalCloserLabel,
  buildGoalCloserPrompt,
  buildWeeklyActionPrompt,
  goalProgressPercent
} from '~/utils/advisorChat'

const props = withDefaults(
  defineProps<{
    label?: string
    prompt?: string
    goal?: Goal | null
    insightTitle?: string
    insightDescription?: string
    variant?: 'outline' | 'secondary' | 'default'
    size?: 'sm' | 'default'
    showIcon?: boolean
    class?: string
  }>(),
  {
    variant: 'outline',
    size: 'sm',
    showIcon: true
  }
)

const { openAdvisorChat } = useOpenAdvisorChat()

const resolvedPrompt = computed(() => {
  if (props.prompt) return props.prompt
  if (props.goal && props.goal.target_amount > 0) {
    return buildGoalCloserPrompt(props.goal)
  }
  if (props.insightTitle) {
    return buildWeeklyActionPrompt(props.insightTitle, props.insightDescription)
  }
  return buildGettingStartedPrompt()
})

const resolvedLabel = computed(() => {
  if (props.label) return props.label
  if (props.goal?.title) return buildGoalCloserLabel(props.goal.title)
  if (props.insightTitle) return ADVISOR.askAboutAction
  return ADVISOR.askGettingStarted
})

const showForGoal = computed(() => {
  if (!props.goal) return true
  return goalProgressPercent(props.goal) < 100
})

async function onClick() {
  await openAdvisorChat(resolvedPrompt.value)
}
</script>

<template>
  <Button
    v-if="showForGoal"
    :variant="variant"
    :size="size"
    :class="props.class"
    class="gap-1.5"
    type="button"
    @click="onClick"
  >
    <MessageCircle v-if="showIcon" class="size-3.5 shrink-0" />
    {{ resolvedLabel }}
  </Button>
</template>
