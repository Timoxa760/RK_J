<script setup lang="ts">
import { Send } from 'lucide-vue-next'
import { ADVISOR } from '~/constants/productCopy'
import { QUICK_PROMPTS } from '~/utils/advisorChat'

const props = withDefaults(
  defineProps<{
    messages: Array<{ id: string; role: 'user' | 'assistant'; content: string }>
    typing?: boolean
    error?: string | null
    /** Компактный чат в левом сайдбаре */
    sidebar?: boolean
  }>(),
  { sidebar: false }
)

const emit = defineEmits<{
  send: [text: string]
}>()

const draft = ref('')
const listRef = ref<HTMLElement | null>(null)

watch(
  () => [props.messages.length, props.typing],
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
</script>

<template>
  <Card
    id="advisor-chat"
    data-demo="advisor-chat"
    class="flex flex-col overflow-hidden"
    :class="
      sidebar
        ? 'mm-sidebar-advisor-card h-full min-h-0 border-0 bg-transparent shadow-none'
        : 'h-full'
    "
  >
    <CardHeader
      class="gap-1 shrink-0"
      :class="sidebar ? 'p-0 pb-2' : 'gap-1.5 p-4 pb-2 sm:p-5 sm:pb-3'"
    >
      <CardTitle :class="sidebar ? 'text-sm font-semibold' : 'text-lg font-semibold'">
        {{ ADVISOR.chatTitle }}
      </CardTitle>
      <CardDescription v-if="!sidebar" class="text-base">{{ ADVISOR.chatHint }}</CardDescription>
      <p v-else class="text-xs leading-snug text-muted-foreground">{{ ADVISOR.chatHintSidebar }}</p>
    </CardHeader>
    <CardContent
      class="flex min-h-0 flex-1 flex-col space-y-2 p-0"
      :class="sidebar ? 'pb-0' : 'space-y-3 pb-4 sm:pb-5'"
    >
      <div
        ref="listRef"
        class="min-h-0 flex-1 space-y-2 overflow-y-auto py-1"
        :class="sidebar ? 'px-0' : 'max-h-[min(480px,52vh)] space-y-3 px-4 sm:px-5'"
        aria-live="polite"
      >
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="flex"
          :class="msg.role === 'user' ? 'justify-end' : 'justify-start'"
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
            {{ msg.content }}
          </div>
        </div>

        <div v-if="typing" class="flex justify-start">
          <div
            class="rounded-2xl rounded-bl-md border bg-muted/50 text-muted-foreground"
            :class="sidebar ? 'px-3 py-2 text-sm' : 'px-4 py-3 text-base'"
          >
            Печатаю…
          </div>
        </div>
      </div>

      <div class="flex flex-wrap gap-1.5" :class="sidebar ? '' : 'px-4'">
        <Button
          v-for="prompt in QUICK_PROMPTS"
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

      <form class="flex gap-2" :class="sidebar ? '' : 'px-4'" @submit.prevent="submit">
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

      <p v-if="error" class="text-xs text-destructive" :class="sidebar ? '' : 'px-4 text-sm'">{{ error }}</p>
    </CardContent>
  </Card>
</template>
