<script setup lang="ts">
const props = defineProps<{
  messages: Array<{ id: string; role: 'user' | 'assistant'; content: string }>
  typing?: boolean
  error?: string | null
}>()

const emit = defineEmits<{
  send: [text: string]
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
        @send="emit('send', $event)"
      />
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
