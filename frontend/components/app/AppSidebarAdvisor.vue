<script setup lang="ts">
import type { AdvisorContext } from '~/utils/advisorChat'
import type { AdvisorChatAction } from '~/types/api'
import type { ChatTurn } from '~/composables/useAdvisorChat'

defineProps<{
  messages: ChatTurn[]
  typing?: boolean
  error?: string | null
  context?: AdvisorContext | null
}>()

const emit = defineEmits<{
  send: [text: string]
  action: [action: AdvisorChatAction]
}>()
</script>

<template>
  <ClientOnly>
    <div class="mm-sidebar-advisor flex min-h-0 flex-1 flex-col">
      <AdvisorChat
        sidebar
        :messages="messages"
        :typing="typing"
        :error="error"
        :context="context"
        @send="emit('send', $event)"
        @action="emit('action', $event)"
      />
      <NuxtLink
        to="/advisor"
        class="mt-2 block text-center text-xs text-primary underline-offset-2 hover:underline md:hidden"
      >
        Открыть полный чат
      </NuxtLink>
    </div>
    <template #fallback>
      <div class="mm-sidebar-advisor flex min-h-[12rem] flex-1 flex-col gap-2 p-1">
        <div class="h-4 w-20 animate-pulse rounded bg-muted" />
        <div class="h-20 flex-1 animate-pulse rounded-xl bg-muted/50" />
        <div class="h-9 animate-pulse rounded-lg bg-muted/50" />
      </div>
    </template>
  </ClientOnly>
</template>
