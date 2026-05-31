<script setup lang="ts">
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import type { AdvisorReplyBlock } from '~/types/api'
import { formatAdvisorReplyText, repairSplitRussianWords } from '~/utils/advisorMarkdown'
import { hasStructuredBlocks, normalizeAdvisorBlocks, parseAdvisorStoredContent } from '~/utils/advisorStructured'

const props = defineProps<{
  content: string
  title?: string
  blocks?: AdvisorReplyBlock[]
  streaming?: boolean
}>()

marked.setOptions({ breaks: true, gfm: true })

const html = ref('')

const resolved = computed(() => {
  if (hasStructuredBlocks(props.blocks)) {
    const normalized = normalizeAdvisorBlocks(props.blocks, props.title)
    return {
      title: normalized.title,
      blocks: normalized.blocks!,
      plain: repairSplitRussianWords(props.content)
    }
  }
  return parseAdvisorStoredContent(props.content)
})

const useStructured = computed(() => hasStructuredBlocks(resolved.value.blocks))

function renderMarkdown() {
  const normalized = formatAdvisorReplyText(resolved.value.plain)
  if (!normalized) {
    html.value = ''
    return
  }
  const raw = marked.parse(normalized, { async: false }) as string
  html.value = DOMPurify.sanitize(raw, {
    ADD_TAGS: ['table', 'thead', 'tbody', 'tr', 'th', 'td'],
    ADD_ATTR: ['colspan', 'rowspan']
  })
}

watch(
  () => [props.content, props.blocks, props.title],
  () => {
    if (!useStructured.value) renderMarkdown()
  },
  { immediate: true }
)
</script>

<template>
  <AdvisorStructuredContent
    v-if="useStructured"
    :title="resolved.title"
    :blocks="resolved.blocks!"
    :streaming="streaming"
  />
  <div v-else-if="streaming" class="flex items-center gap-2 text-sm text-muted-foreground">
    <span class="advisor-markdown-cursor" aria-hidden="true" />
    <span>Готовлю ответ…</span>
  </div>
  <div
    v-else-if="html"
    class="advisor-markdown text-base leading-relaxed"
    v-html="html"
  />
  <p
    v-else-if="content && !streaming"
    class="whitespace-pre-wrap text-base leading-relaxed"
  >
    {{ content }}
  </p>
  <span v-else-if="streaming" class="advisor-markdown-cursor" aria-hidden="true" />
</template>

<style scoped>
.advisor-markdown-cursor {
  display: inline-block;
  width: 2px;
  height: 1em;
  margin-left: 2px;
  vertical-align: text-bottom;
  background: currentColor;
  opacity: 0.45;
  animation: advisor-cursor-blink 1s step-end infinite;
}

@keyframes advisor-cursor-blink {
  50% {
    opacity: 0;
  }
}
</style>
