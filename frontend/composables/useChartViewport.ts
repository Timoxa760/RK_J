import { useElementSize } from '@vueuse/core'

/** Ширина контейнера; до измерения — окно или mobile-first, не «999» (ломало compact-режим). */
function fallbackWidth(): number {
  if (import.meta.client && typeof window !== 'undefined') {
    return window.innerWidth
  }
  return 360
}

export function useChartViewport() {
  const containerRef = ref<HTMLElement | null>(null)
  const { width: measuredWidth } = useElementSize(containerRef)

  const width = computed(() =>
    measuredWidth.value > 0 ? measuredWidth.value : fallbackWidth()
  )

  const isCompact = computed(() => width.value < 480)
  const isMedium = computed(() => width.value < 768)

  return {
    containerRef,
    width,
    isCompact,
    isMedium
  }
}
