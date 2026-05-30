export const chartThemeLight = {
  backgroundColor: 'transparent',
  textStyle: { color: '#6d6760' },
  axisColor: '#f5dcc8',
  splitLine: '#f2ebe3',
  colors: ['#e8955f', '#f0a66b', '#f5c4a0', '#d4824a', '#f5dcc8', '#c9773f']
} as const

export function baseGrid(compact = false) {
  return compact
    ? { left: 40, right: 12, top: 36, bottom: 48, containLabel: true }
    : { left: 48, right: 24, top: 40, bottom: 40, containLabel: true }
}

export function formatAxisMoney(value: number): string {
  if (Math.abs(value) >= 1_000_000) return `${Math.round(value / 100_000) / 10}M`
  if (Math.abs(value) >= 1_000) return `${Math.round(value / 1000)}k`
  return String(Math.round(value))
}

export function formatMonthLabel(month: string): string {
  const match = month.match(/^(\d{4})-(\d{2})/)
  if (!match) return month
  return `${match[2]}/${match[1]!.slice(2)}`
}

export function sparseLabelInterval(count: number, compact: boolean): number | 'auto' {
  if (count <= 6) return 0
  if (compact) return Math.max(0, Math.ceil(count / 6) - 1)
  if (count <= 12) return 1
  return Math.ceil(count / 8) - 1
}

export function labelLayoutHideOverlap() {
  return { hideOverlap: true }
}

export function chartHeightClass(size: 'sm' | 'md' | 'lg' | 'full'): string {
  const map = {
    sm: 'h-[260px] lg:h-[300px]',
    md: 'h-[340px] sm:h-[300px] lg:h-[340px]',
    lg: 'h-[340px] lg:h-[400px]',
    full: 'h-[400px] sm:h-[360px] lg:h-[440px]'
  }
  return map[size]
}
