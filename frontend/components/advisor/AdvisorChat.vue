<script setup lang="ts">
import { RotateCcw, Send } from 'lucide-vue-next'
import type { AdvisorChatAction } from '~/types/api'
import { ADVISOR } from '~/constants/productCopy'
import { buildDynamicQuickPrompts, type AdvisorContext } from '~/utils/advisorChat'
import type { ChatTurn } from '~/composables/useAdvisorChat'

const props = withDefaults(
  defineProps<{
    messages: ChatTurn[]
    typing?: boolean
    error?: string | null
    context?: AdvisorContext | null
    /** Компактный чат в левом сайдбаре */
    sidebar?: boolean
    /** Полноэкранная страница /advisor */
    fullPage?: boolean
    showReset?: boolean
  }>(),
  { sidebar: false, fullPage: false, showReset: false }
)

const emit = defineEmits<{
  send: [text: string]
  reset: []
  action: [action: AdvisorChatAction]
}>()

const draft = ref('')
const listRef = ref<HTMLElement | null>(null)

const quickPrompts = computed(() =>
  props.context ? buildDynamicQuickPrompts(props.context) : []
)

watch(
  () => [props.messages.length, props.typing, props.messages.map((m) => m.content).join('')],
  async () => {
    await nextTick()
    const el = listRef.value
    if (el) el.scrollTop = el.scrollHeight
  }
)

function submit() {
  const text = draft.value.trim()
  if (!text) return
  draft.value = ''
  emit('send', text)
}

function onPrompt(text: string) {
  emit('send', text)
}

function isLocalSource(msg: ChatTurn) {
  return msg.source === 'local'
}
</script>

<template>
  <Card
    :id="fullPage ? 'advisor-chat' : undefined"
    data-demo="advisor-chat"
    class="flex flex-col overflow-hidden"
    :class="[
      sidebar
        ? 'mm-sidebar-advisor-card h-full min-h-0 border-0 bg-transparent shadow-none'
        : fullPage
          ? 'h-full min-h-0 border-0 bg-transparent shadow-none'
          : 'h-full'
    ]"
  >
    <CardHeader
      class="gap-1 shrink-0"
      :class="sidebar ? 'p-0 pb-2' : 'gap-1.5 p-4 pb-2 sm:p-5 sm:pb-3'"
    >
      <div class="flex items-start justify-between gap-2">
        <div>
          <CardTitle :class="sidebar ? 'text-sm font-semibold' : 'text-lg font-semibold'">
            {{ ADVISOR.chatTitle }}
          </CardTitle>
          <CardDescription v-if="!sidebar" class="text-base">{{ ADVISOR.chatHint }}</CardDescription>
          <p v-else class="text-xs leading-snug text-muted-foreground">{{ ADVISOR.chatHintSidebar }}</p>
        </div>
        <Button
          v-if="showReset"
          type="button"
          variant="ghost"
          size="sm"
          class="shrink-0 gap-1 text-xs"
          :disabled="typing"
          @click="emit('reset')"
        >
          <RotateCcw class="size-3.5" />
          {{ ADVISOR.chatReset }}
        </Button>
      </div>
    </CardHeader>
    <CardContent
      class="flex min-h-0 flex-1 flex-col space-y-2 p-0"
      :class="sidebar ? 'pb-0' : 'space-y-3 pb-4 sm:pb-5'"
    >
      <div
        ref="listRef"
        class="min-h-0 flex-1 space-y-2 overflow-y-auto py-1"
        :class="
          sidebar
            ? 'px-0'
            : fullPage
              ? 'space-y-3 px-3 sm:px-4'
              : 'max-h-[min(480px,52vh)] space-y-3 px-4 sm:px-5'
        "
        aria-live="polite"
      >
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="flex flex-col"
          :class="msg.role === 'user' ? 'items-end' : 'items-start'"
        >
          <div
            class="whitespace-pre-wrap rounded-2xl leading-relaxed shadow-sm"
            :class="[
              sidebar ? 'max-w-[92%] px-3 py-2 text-sm' : 'max-w-[88%] px-4 py-3 text-base',
              msg.role === 'user'
                ? 'rounded-br-md bg-primary text-[color:var(--mm-text)]'
                : 'rounded-bl-md border bg-card text-foreground'
            ]"
          >
            <span v-if="msg.streaming && !msg.content" class="text-muted-foreground">Печатаю…</span>
            <span v-else>{{ msg.content }}</span>
          </div>
          <p
            v-if="msg.role === 'assistant' && isLocalSource(msg)"
            class="mt-1 text-[10px] text-muted-foreground"
          >
            {{ ADVISOR.chatLocalReply }}
          </p>
          <AdvisorChatActions
            v-if="msg.role === 'assistant' && msg.actions?.length"
            :actions="msg.actions"
            :disabled="typing"
            @action="emit('action', $event)"
          />
        </div>

        <div v-if="typing && !messages.some((m) => m.streaming)" class="flex justify-start">
          <div
            class="rounded-2xl rounded-bl-md border bg-muted/50 text-muted-foreground"
            :class="sidebar ? 'px-3 py-2 text-sm' : 'px-4 py-3 text-base'"
          >
            Печатаю…
          </div>
        </div>
      </div>

      <div class="flex flex-wrap gap-1.5" :class="sidebar ? '' : 'px-3 sm:px-4'">
        <Button
          v-for="prompt in quickPrompts"
          :key="prompt"
          type="button"
          variant="outline"
          size="sm"
          class="rounded-full text-xs"
          :class="sidebar ? 'h-7 px-2' : 'text-sm'"
          :disabled="typing"
          @click="onPrompt(prompt)"
        >
          {{ prompt }}
        </Button>
      </div>

      <form
        class="flex gap-2"
        :class="sidebar ? '' : 'px-3 sm:px-4'"
        @submit.prevent="submit"
      >
        <Input
          v-model="draft"
          :placeholder="ADVISOR.chatPlaceholder"
          class="flex-1"
          :class="sidebar ? 'min-h-9 text-sm' : 'min-h-12 text-base'"
          :disabled="typing"
          autocomplete="off"
        />
        <Button
          type="submit"
          size="icon"
          class="shrink-0"
          :class="sidebar ? 'size-9' : 'size-11'"
          :disabled="typing || !draft.trim()"
        >
          <Send class="size-4" />
          <span class="sr-only">Отправить</span>
        </Button>
      </form>

      <p v-if="error" class="text-xs text-destructive" :class="sidebar ? '' : 'px-3 text-sm sm:px-4'">
        {{ error }}
      </p>
    </CardContent>
  </Card>
</template>
