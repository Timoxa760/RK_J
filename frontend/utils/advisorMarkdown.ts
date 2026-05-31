/**
 * Явные склейки от LLM (без пробела). Не трогаем нормальные русские слова с теми же слогами.
 */
const GLUED_REPAIRS: ReadonlyArray<readonly [RegExp, string]> = [
  [/притекущих\s*данных/gi, 'при текущих данных'],
  [/притекущих/gi, 'при текущих'],
  [/свободных\s*денег/gi, 'свободных денег'],
  [/свободныхденег/gi, 'свободных денег'],
  [/откладывать\s*на\s*цель/gi, 'Откладывать на цель'],
  [/откладыватьнацель/gi, 'Откладывать на цель'],
  [/накредит/gi, 'на кредит'],
  [/главный\s*вопрос/gi, 'Главный вопрос'],
  [/главныйвопрос/gi, 'Главный вопрос'],
  [/илиесть/gi, 'или есть'],
  [/замесяц/gi, 'за месяц'],
  [/занеделю/gi, 'за неделю'],
  [/вмесяц/gi, 'в месяц'],
  [/намесяц/gi, 'на месяц'],
  [/безкатегори/gi, 'без категори'],
  [/безпробел/gi, 'без пробел'],
  [/([а-яё]{3,})(данных|денег)\b/gi, '$1 $2'],
  [/([а-яё]+)(на)(цель|кредит)\b/gi, '$1 $2 $3'],
  [/(\d)(₽)/g, '$1 $2'],
  [/(\d)(%)/g, '$1 $2']
]

const EXPLICIT_SPLIT_FIXES: ReadonlyArray<readonly [RegExp, string]> = [
  [/за\s+по\s+лнении/gi, 'заполнении'],
  [/за\s+по\s+лнен/gi, 'заполнен'],
  [/по\s+это\s+му/gi, 'поэтому'],
  [/по\s+это\s+м([^а-яё]|$)/gi, 'поэтом$1'],
  [/де\s+по\s+зит/gi, 'депозит'],
  [/го\s+до\s+вых/gi, 'годовых'],
  [/кр\s+из\s+ис/gi, 'кризис'],
  [/пред\s+по\s+чтительнее/gi, 'предпочтительнее'],
  [/рабо\s+тают/gi, 'работают'],
  [/не\s+рабают/gi, 'не работают'],
  [/без\s+укания/gi, 'без уточнения'],
  [/ставки\s*цб/gi, 'ставки ЦБ'],
  [/ис\s+по\s+льз/gi, 'использ'],
  [/не\s+до\s+стающ/gi, 'недостающ'],
  [/единственныйфи\s+на\s+нсовый/gi, 'единственный финансовый'],
  [/фи\s+на\s+нсовый/gi, 'финансовый'],
  [/месяцона/gi, 'месяц она'],
  [/увасуже/gi, 'у вас уже'],
  [/заменын[а-яёa-z]+/gi, 'замены нет']
]

const PREP_INSIDE_WORD_FALSE = new Set([
  'идти', 'идём', 'идем', 'иду', 'шли', 'шёл', 'шел',
  'смотреть', 'смотрите', 'смотрим', 'ориентир', 'опираться'
])

const PREP_INSIDE_WORD_NA_FALSE = new Set([
  'он', 'она', 'они', 'мы', 'вы', 'то', 'ту', 'ему', 'ей', 'им'
])

const RE_SHORT_FRAGMENT_CHAIN =
  /(^|[\s(«"'—-])([а-яё]{1,3}(?: [а-яё]{1,3}){2,})([\s,.!?;:»)\]"'—-]|$)/gi

const RE_MIXED_FRAGMENT_CHAIN =
  /(^|[\s(«"'—-])((?:[а-яё]{1,3}\s+){2,}[а-яё]{4,12})([\s,.!?;:»)\]"'—-]|$)/gi

const RE_PREP_INSIDE_WORD = /([а-яё]{3,6}) по ([б-джзклмнпрстфхцчшщ][а-яё]{4,})/gi

const RE_PREP_INSIDE_WORD_NA = /([а-яё]{1,3}) на ([б-джзклмнпрстфхцчшщ][а-яё]{3,})/gi

const SPLIT_MERGE_STOP = new Set([
  'за', 'по', 'от', 'до', 'на', 'в', 'и', 'но', 'или', 'при', 'для', 'из', 'с', 'к', 'у', 'о', 'об',
  'не', 'ни', 'вы', 'мы', 'он', 'я', 'ты', 'её', 'ее', 'их', 'то', 'как', 'что', 'это', 'там', 'тут',
  'уже', 'ещё', 'eще', 'все', 'всё', 'без', 'про', 'под', 'над', 'мне', 'вас', 'нас',   'ему', 'ей', 'им',
  'бы', 'ли', 'же', 'она', 'они', 'нет', 'есть'
])

const RE_EMBEDDED_PREP_SPLIT =
  /([а-яё]{3,}) (?:за|по|от|или) ([а-яё]{1,3})([\s,.!?;:«»\)\]"'\-—]|$)/gi
const RE_SUFFIX_SPLIT = /([а-яё]{4,}) ([а-яё]{1,3})([\s,.!?;:«»\)\]"'\-—]|$)/gi

const RUSSIAN_VOWELS = new Set(['а', 'е', 'ё', 'и', 'о', 'у', 'ы', 'э', 'ю', 'я'])

function isRussianVowel(ch: string): boolean {
  return RUSSIAN_VOWELS.has(ch.toLowerCase())
}

function endsWithRussianVowel(word: string): boolean {
  const ch = word.slice(-1)
  return ch ? isRussianVowel(ch) : false
}

function shouldMergeSplitFragment(stem: string, frag: string): boolean {
  const lower = frag.toLowerCase()
  if (SPLIT_MERGE_STOP.has(lower)) return false
  if (lower.length === 1) {
    if (lower === 'и') return false
    if (isRussianVowel(lower) && endsWithRussianVowel(stem)) return false
  }
  return true
}

function mergeSplitMatch(full: string, stem: string, frag: string, tail: string): string {
  if (!shouldMergeSplitFragment(stem, frag)) return full
  return stem + frag + tail
}

function repairShortFragmentChains(text: string): string {
  return text.replace(RE_SHORT_FRAGMENT_CHAIN, (full, prefix, chain, suffix) => {
    const parts = chain.trim().split(/\s+/)
    if (parts.length < 3 || parts.some((p) => p.length > 3)) return full
    if (parts.every((p) => SPLIT_MERGE_STOP.has(p.toLowerCase()))) return full
    const merged = parts.join('')
    if (merged.length < 5 || merged.length > 24) return full
    return `${prefix}${merged}${suffix}`
  })
}

function repairMixedFragmentChains(text: string): string {
  return text.replace(RE_MIXED_FRAGMENT_CHAIN, (full, prefix, chain, suffix) => {
    const parts = chain.trim().split(/\s+/)
    if (parts.length < 3) return full
    if (parts.slice(0, -1).some((p) => p.length > 3 || p.length < 2)) return full
    const last = parts[parts.length - 1]!
    if (last.length < 4 || last.length > 12) return full
    const merged = parts.join('')
    if (merged.length < 6 || merged.length > 32) return full
    return `${prefix}${merged}${suffix}`
  })
}

function repairPrepInsideWordNa(text: string): string {
  return text.replace(RE_PREP_INSIDE_WORD_NA, (full, stem, cont) => {
    if (PREP_INSIDE_WORD_NA_FALSE.has(String(stem).toLowerCase())) return full
    return `${stem}на${cont}`
  })
}

function repairPrepInsideWord(text: string): string {
  return text.replace(RE_PREP_INSIDE_WORD, (full, stem, cont) => {
    if (PREP_INSIDE_WORD_FALSE.has(String(stem).toLowerCase())) return full
    return `${stem}по${cont}`
  })
}

/** Убирает лишние пробелы внутри русских слов (артефакт LLM). */
export function repairSplitRussianWords(text: string): string {
  if (!text) return text
  let out = text
  for (const [pattern, replacement] of EXPLICIT_SPLIT_FIXES) {
    out = out.replace(pattern, replacement)
  }
  out = repairShortFragmentChains(out)
  out = repairMixedFragmentChains(out)
  out = repairPrepInsideWord(out)
  out = repairPrepInsideWordNa(out)
  out = out.replace(RE_EMBEDDED_PREP_SPLIT, (full, stem, frag, tail) =>
    mergeSplitMatch(full, stem, frag, tail)
  )
  for (let i = 0; i < 6; i++) {
    const next = out.replace(RE_SUFFIX_SPLIT, (full, stem, frag, tail) =>
      mergeSplitMatch(full, stem, frag, tail)
    )
    if (next === out) break
    out = next
  }
  return out
}

/** Приводит ответ советника к читаемому виду (пробелы, абзацы, без таблиц). */
export function formatAdvisorReplyText(text: string): string {
  let out = text.replace(/\r\n/g, '\n').trim()
  if (!out) return ''

  out = out.replace(/\|{2,}/g, '\n').replace(/\|---\|/g, '\n')
  out = convertPipeBlocks(out)
  out = convertPipeRows(out)
  out = out.replace(/([.!?])\s*([А-ЯЁ][^\n:]{0,40}:)/g, '$1\n\n$2')
  out = out.replace(/^(#{1,6})([^\s#\n])/gm, '$1 $2')
  out = out.replace(/([^\n])(#{1,6}\s)/g, '$1\n\n$2')
  out = out.replace(/([^\n])\s*---\s*([^\n])/g, '$1\n\n---\n\n$2')

  for (const [pattern, replacement] of GLUED_REPAIRS) {
    out = out.replace(pattern, replacement)
  }

  out = repairAdvisorSpacing(out)
  out = out.replace(/^[ \t|—-]+$/gm, '')

  out = out
    .split('\n')
    .map((line) => repairAdvisorSpacing(line.trim()))
    .filter((line) => line && line !== '|' && line !== '---')
    .join('\n')
    .replace(/[ \t]{2,}/g, ' ')
    .trim()

  return repairSplitRussianWords(out)
}

function convertPipeBlocks(text: string): string {
  return text.replace(/(?:\|[^|\n]{1,120}){2,}/g, (block) => pipeCellsToList(block))
}

function convertPipeRows(text: string): string {
  return text.replace(/^[^|\n]*\|[^|\n]+$/gm, (row) => pipeCellsToList(row))
}

function pipeCellsToList(block: string): string {
  const rows = block
    .split('|')
    .map((cell) => cell.trim())
    .filter((cell) => cell && cell !== '---' && !cell.startsWith('---'))

  if (!rows.length) {
    return block.replace(/\|/g, ' ').replace(/[ \t]{2,}/g, ' ').trim()
  }
  return rows.map((row) => `- ${row}`).join('\n')
}

function repairAdvisorSpacing(s: string): string {
  return s
    .replace(/([,.:;!?])([^\s\d])/g, '$1 $2')
    .replace(/—/g, ' — ')
    .replace(/—\s*(\d)/g, '— $1')
    .replace(/(\d)\s*₽/g, '$1 ₽')
    .replace(/(на|по|за|от|до)\s*(\d)/gi, '$1 $2')
    .replace(/₽\s*([А-ЯA-Z])/g, '₽ $1')
    .replace(/([а-яё])([А-ЯЁ])/g, '$1 $2')
    .replace(/[ \t]{2,}/g, ' ')
    .trim()
}

/** @deprecated use formatAdvisorReplyText */
export function normalizeAdvisorMarkdown(text: string): string {
  return formatAdvisorReplyText(text)
}

export function formatChatMessageTime(ts: number): string {
  const date = new Date(ts)
  const now = new Date()
  const sameDay =
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate()

  const time = date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
  if (sameDay) return time

  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export function formatChatDayLabel(ts: number): string | null {
  const date = new Date(ts)
  const now = new Date()
  const startOf = (d: Date) => new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
  const diffDays = Math.round((startOf(now) - startOf(date)) / 86_400_000)

  if (diffDays === 0) return 'Сегодня'
  if (diffDays === 1) return 'Вчера'
  if (diffDays < 7) {
    return date.toLocaleDateString('ru-RU', { weekday: 'long' })
  }
  return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' })
}

export function chatDayKey(ts: number): string {
  const d = new Date(ts)
  return `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`
}
