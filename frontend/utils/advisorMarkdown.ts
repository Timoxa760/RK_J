type TextFix = readonly [RegExp, string]

const EXPLICIT_SPLIT_FIXES: TextFix[] = [
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
  [/меся\s+-(\s|$)/gi, 'месяце$1'],
  [/увасуже/gi, 'у вас уже'],
  [/заменын[a-яёa-z]+/gi, 'замены нет']
]

const TOKEN_GLUED_MAP = new Map<string, string>(
  Object.entries({
    расходынапродукты: 'расходы на продукты',
    бюджетнаследующий: 'бюджет на следующий',
    откладыватьнацель: 'откладывать на цель',
    откладыватьнакредит: 'откладывать на кредит',
    свободныхденег: 'свободных денег',
    притекущихданных: 'при текущих данных',
    притекущих: 'при текущих',
    главныйвопрос: 'главный вопрос',
    илиесть: 'или есть',
    замесяц: 'за месяц',
    занеделю: 'за неделю',
    вмесяц: 'в месяц',
    намесяц: 'на месяц',
    накредит: 'на кредит',
    безкатегори: 'без категори',
    безпробел: 'без пробел',
    сокративэти: 'сократив эти',
    сокративэтот: 'сократив этот',
    сокративэту: 'сократив эту'
  })
)

const KNOWN_SHORT_TAIL = [
  'эти',
  'этот',
  'этой',
  'этом',
  'эту',
  'это',
  'или',
  'уже',
  'ещё',
  'eще',
  'нет'
]

const SAFE_LEFT_STEMS = new Set([
  'расходы',
  'бюджет',
  'сократив',
  'откладывать',
  'отложите',
  'сократите',
  'переведите',
  'направьте',
  'главный',
  'свободных',
  'свободные',
  'свободный'
])

const RE_VERBISH_LEFT = /(?:ив|ивш|ая|ые|ий|ую|ите|ьте|ать|ить|еть)$/i

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

function applyLeadingCase(source: string, replacement: string): string {
  if (!source.length || !replacement.length) return replacement
  const first = source[0]!
  if (first === first.toUpperCase() && first !== first.toLowerCase()) {
    return first + replacement.slice(1)
  }
  return replacement
}

function repairGluedToken(token: string): string {
  const lower = token.toLowerCase()
  const mapped = TOKEN_GLUED_MAP.get(lower)
  if (mapped) return applyLeadingCase(token, mapped)

  for (const tail of KNOWN_SHORT_TAIL) {
    if (!lower.endsWith(tail)) continue
    const left = lower.slice(0, -tail.length)
    if (left.length < 4) continue
    if (!SAFE_LEFT_STEMS.has(left) && !RE_VERBISH_LEFT.test(left)) continue
    return applyLeadingCase(token, `${left} ${tail}`)
  }

  return token
}

function repairTokensInText(text: string): string {
  return text.replace(/[а-яё]+/gi, (token) => repairGluedToken(token))
}

/** Убирает известные артефакты LLM без разрезания нормальных русских слов. */
export function repairSplitRussianWords(text: string): string {
  if (!text) return text
  let out = text
  for (const [pattern, replacement] of EXPLICIT_SPLIT_FIXES) {
    out = out.replace(pattern, replacement)
  }
  return repairTokensInText(out)
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
