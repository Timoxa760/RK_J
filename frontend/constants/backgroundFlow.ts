import { FLOW_LINE_PATHS } from '~/constants/flowLinePaths'

export const FLOW_VIEWBOX = '0 0 422 596'

/** Кроп + масштаб: линии расходятся на весь экран */
export const FLOW_STAGE_VIEWBOX = '-80 210 580 420'
export const FLOW_STAGE_GROUP = 'translate(-40, -60) scale(1.48)'
export const FLOW_LINE_COUNT = FLOW_LINE_PATHS.length

const PULSE_EVERY = 3

export interface FlowPulseAnim {
  index: number
  duration: number
  delay: number
}

function slotOf(index: number) {
  return Math.floor(index / PULSE_EVERY)
}

function pulseTiming(index: number): FlowPulseAnim {
  const slot = slotOf(index)
  return {
    index,
    duration: 16 + (slot % 5) * 2,
    delay: -(slot * 3.2)
  }
}

export const FLOW_PULSES: FlowPulseAnim[] = FLOW_LINE_PATHS.map((_, i) => i)
  .filter((i) => i % PULSE_EVERY === 1)
  .map((index) => pulseTiming(index))

export const FLOW_LINE_INDICES = FLOW_LINE_PATHS.map((_, i) => i)

export const FLOW_PULSE_LINE_INDICES = new Set(
  FLOW_LINE_INDICES.filter((i) => i % PULSE_EVERY === 1)
)

export const FLOW_ACCENT_INDICES = FLOW_LINE_INDICES.filter(
  (i) => !FLOW_PULSE_LINE_INDICES.has(i)
)

export function flowAnimVars(duration: number, delay: number) {
  return {
    '--flow-duration': `${duration}s`,
    '--flow-delay': `${delay}s`
  } as Record<string, string>
}
