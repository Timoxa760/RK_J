<script setup lang="ts">
import { RotateCcw, Send } from 'lucide-vue-next'
import type { AdvisorChatAction } from '~/types/api'
import { ADVISOR } from '~/constants/productCopy'
import { buildDynamicQuickPrompts, type AdvisorContext } from '~/utils/advisorChat'
import {
  chatDayKey,
  formatChatDayLabel,
  formatChatMessageTime
} from '~/utils/advisorMarkdown'
import type { ChatTurn } from '~/composables/useAdvisorChat'

const props = withDefaults(
  defineProps<{
    messages: ChatTurn[]
    typing?: boolean
    error?: string | null
    context?: AdvisorContext | null
    /** Полноэкранная страница /advisor */
    fullPage?: boolean
    showReset?: boolean
  }>(),
  { fullPage: false, showReset: false }
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

function sourceLabel(msg: ChatTurn): string | null {
  if (msg.role !== 'assistant') return null
  if (msg.source === 'heuristic') return ADVISOR.chatSourceHeuristic
  if (msg.source === 'gemini') return ADVISOR.chatSourceAi
  return null
}

function showDayDivider(index: number): string | null {
  const msg = props.messages[index]
  if (!msg) return null
  const prev = props.messages[index - 1]
  if (prev && chatDayKey(prev.createdAt) === chatDayKey(msg.createdAt)) return null
  return formatChatDayLabel(msg.createdAt)
}
</script>

<template>
  <Card
    :id="fullPage ? 'advisor-chat' : undefined"
    data-demo="advisor-chat"
    class="flex flex-col overflow-hidden"
    :class="fullPage ? 'h-full min-h-0 border-0 bg-transparent shadow-none' : 'h-full'"
  >
    <CardHeader class="shrink-0 gap-1.5 p-4 pb-2 sm:p-5 sm:pb-3">
      <div class="flex items-start justify-between gap-2">
        <div>
          <CardTitle class="text-lg font-semibold">
            {{ ADVISOR.chatTitle }}
          </CardTitle>
          <CardDescription class="text-base">{{ ADVISOR.chatHint }}</CardDescription>
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
      class="flex min-h-0 flex-1 flex-col space-y-3 p-0 pb-4 sm:pb-5"
    >
      <div
        ref="listRef"
        class="min-h-0 flex-1 space-y-3 overflow-y-auto py-1"
        :class="fullPage ? 'px-3 sm:px-4' : 'max-h-[min(480px,52vh)] px-4 sm:px-5'"
        aria-live="polite"
      >
        <div
          v-if="!messages.length && !typing"
          class="flex min-h-[12rem] flex-col items-center justify-center px-2 py-6 text-center sm:min-h-[14rem]"
        >
          <p class="max-w-md text-sm leading-relaxed text-muted-foreground">
            {{ ADVISOR.chatEmptyHint }}
          </p>
        </div>

        <template v-for="(msg, index) in messages" :key="msg.id">
          <div
            v-if="showDayDivider(index)"
            class="flex justify-center py-1"
          >
            <span class="rounded-full bg-muted/60 px-3 py-0.5 text-xs text-muted-foreground">
              {{ showDayDivider(index) }}
            </span>
          </div>
          <div
            class="flex flex-col"
            :class="msg.role === 'user' ? 'items-end' : 'items-start'"
          >
            <div
              class="max-w-[92%] rounded-2xl px-4 py-3 shadow-sm sm:max-w-[88%]"
              :class="
                msg.role === 'user'
                  ? 'rounded-br-md bg-primary text-[color:var(--mm-text)]'
                  : 'rounded-bl-md border bg-card text-foreground'
              "
            >
              <template v-if="msg.role === 'assistant'">
                <span v-if="msg.streaming" class="text-muted-foreground">
                  Печатаю…
                </span>
                <AdvisorMessageContent
                  v-else
                  :content="msg.content"
                  :title="msg.title"
                  :blocks="msg.blocks"
                />
              </template>
              <p v-else class="whitespace-pre-wrap text-base leading-relaxed">
                {{ msg.content }}
              </p>
            </div>
            <div
              class="mt-1 flex max-w-[92%] flex-wrap items-center gap-x-2 gap-y-0.5 px-1 sm:max-w-[88%]"
              :class="msg.role === 'user' ? 'justify-end' : 'justify-start'"
            >
              <time
                :datetime="new Date(msg.createdAt).toISOString()"
                class="text-[10px] text-muted-foreground"
              >
                {{ formatChatMessageTime(msg.createdAt) }}
              </time>
              <span
                v-if="sourceLabel(msg)"
                class="text-[10px] text-muted-foreground"
              >
                · {{ sourceLabel(msg) }}
              </span>
            </div>
            <AdvisorChatActions
              v-if="msg.role === 'assistant' && msg.actions?.length"
              :actions="msg.actions"
              :disabled="typing"
              @action="emit('action', $event)"
            />
          </div>
        </template>

        <div v-if="typing && !messages.some((m) => m.streaming)" class="flex justify-start">
          <div class="rounded-2xl rounded-bl-md border bg-muted/50 px-4 py-3 text-base text-muted-foreground">
            Печатаю…
          </div>
        </div>
      </div>

      <div
        v-if="quickPrompts.length"
        class="flex flex-wrap gap-1.5 px-3 sm:px-4"
      >
        <p
          v-if="!messages.length"
          class="mb-1 w-full text-xs font-medium text-muted-foreground"
        >
          Подсказки по вашей картине
        </p>
        <Button
          v-for="prompt in quickPrompts"
          :key="prompt"
          type="button"
          variant="outline"
          size="sm"
          class="rounded-full text-sm"
          :disabled="typing"
          @click="onPrompt(prompt)"
        >
          {{ prompt }}
        </Button>
      </div>

      <form class="flex gap-2 px-3 sm:px-4" @submit.prevent="submit">
        <Input
          v-model="draft"
          :placeholder="ADVISOR.chatPlaceholder"
          class="min-h-12 flex-1 text-base"
          :disabled="typing"
          autocomplete="off"
        />
        <Button
          type="submit"
          size="icon"
          class="size-11 shrink-0"
          :disabled="typing || !draft.trim()"
        >
          <Send class="size-4" />
          <span class="sr-only">Отправить</span>
        </Button>
      </form>

      <p v-if="error" class="px-3 text-sm text-destructive sm:px-4">
        {{ error }}
      </p>
    </CardContent>
  </Card>
</template>
