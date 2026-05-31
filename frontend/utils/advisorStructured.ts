import type { AdvisorReplyBlock } from '~/types/api'
import { repairSplitRussianWords } from '~/utils/advisorMarkdown'

export interface ParsedAdvisorContent {
  title?: string
  blocks?: AdvisorReplyBlock[]
  plain: string
}

interface StoredEnvelope {
  v?: number
  title?: string
  blocks?: AdvisorReplyBlock[]
  plain?: string
}

function repairBlock(block: AdvisorReplyBlock): AdvisorReplyBlock {
  if (block.type === 'list') {
    return {
      ...block,
      items: block.items?.map((item) => repairSplitRussianWords(item))
    }
  }
  return {
    ...block,
    text: block.text ? repairSplitRussianWords(block.text) : block.text
  }
}

function repairBlocks(blocks?: AdvisorReplyBlock[]): AdvisorReplyBlock[] | undefined {
  return blocks?.map(repairBlock)
}

/** Разбирает content из API/истории: JSON v1 или plain text. */
export function parseAdvisorStoredContent(content: string): ParsedAdvisorContent {
  const trimmed = content.trim()
  if (!trimmed) return { plain: '' }

  if (trimmed.startsWith('{')) {
    try {
      const env = JSON.parse(trimmed) as StoredEnvelope
      if (env.v === 1 && Array.isArray(env.blocks) && env.blocks.length > 0) {
        const blocks = repairBlocks(env.blocks)
        const title = env.title ? repairSplitRussianWords(env.title) : env.title
        const plain = repairSplitRussianWords(
          env.plain?.trim() || structuredToPlain(title, blocks)
        )
        return { title, blocks, plain }
      }
    } catch {
      /* plain text */
    }
  }

  return { plain: repairSplitRussianWords(trimmed) }
}

export function normalizeAdvisorBlocks(
  blocks?: AdvisorReplyBlock[] | null,
  title?: string
): { title?: string; blocks?: AdvisorReplyBlock[] } {
  return {
    title: title ? repairSplitRussianWords(title) : title,
    blocks: repairBlocks(blocks ?? undefined)
  }
}

export function structuredToPlain(title?: string, blocks?: AdvisorReplyBlock[]): string {
  const parts: string[] = []
  if (title?.trim()) parts.push(title.trim())
  for (const block of blocks ?? []) {
    if (block.type === 'list') {
      for (const item of block.items ?? []) {
        if (item.trim()) parts.push(`• ${item.trim()}`)
      }
    } else if (block.text?.trim()) {
      parts.push(block.text.trim())
    }
  }
  return parts.join('\n\n')
}

export function hasStructuredBlocks(blocks?: AdvisorReplyBlock[] | null): boolean {
  return Boolean(blocks?.length)
}
