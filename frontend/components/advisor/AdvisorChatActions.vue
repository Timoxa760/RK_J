<script setup lang="ts">
import type { AdvisorChatAction } from '~/types/api'
import { ADVISOR } from '~/constants/productCopy'

defineProps<{
  actions: AdvisorChatAction[]
  disabled?: boolean
}>()

const emit = defineEmits<{
  action: [action: AdvisorChatAction]
}>()
</script>

<template>
  <div v-if="actions.length" class="mt-2 flex flex-wrap gap-1.5">
    <Button
      v-for="(action, index) in actions"
      :key="`${action.type}-${action.label}-${index}`"
      type="button"
      variant="outline"
      size="sm"
      class="h-7 rounded-full px-2.5 text-xs"
      :disabled="disabled"
      @click="emit('action', action)"
    >
      {{ action.label }}
    </Button>
  </div>
</template>
