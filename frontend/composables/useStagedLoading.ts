import type { FinancialReportLoadingStage } from '~/constants/financialReportLoading'

export function useStagedLoading(
  stages: FinancialReportLoadingStage[],
  active: Ref<boolean>
) {
  const currentIndex = ref(0)
  const finishing = ref(false)
  let timers: ReturnType<typeof setTimeout>[] = []

  function clearTimers() {
    for (const timer of timers) clearTimeout(timer)
    timers = []
  }

  function scheduleStages() {
    clearTimers()
    currentIndex.value = 0
    finishing.value = false

    let offset = 0
    for (let i = 0; i < stages.length - 1; i++) {
      const durationMs = stages[i]?.durationMs ?? 2400
      offset += durationMs
      const nextIndex = i + 1
      timers.push(
        setTimeout(() => {
          if (active.value && !finishing.value) {
            currentIndex.value = nextIndex
          }
        }, offset)
      )
    }
  }

  watch(
    active,
    (isActive) => {
      if (isActive) {
        scheduleStages()
        return
      }

      clearTimers()
      if (currentIndex.value === 0 && !finishing.value) return

      finishing.value = true
      currentIndex.value = stages.length - 1
      timers.push(
        setTimeout(() => {
          currentIndex.value = 0
          finishing.value = false
        }, 320)
      )
    },
    { immediate: true }
  )

  onScopeDispose(clearTimers)

  const currentStage = computed(() => stages[currentIndex.value] ?? stages[0])
  const isLastStage = computed(() => currentIndex.value >= stages.length - 1)

  const progress = computed(() => {
    const total = Math.max(stages.length, 1)
    const base = (currentIndex.value + 1) / total
    return Math.min(100, Math.round(base * 100))
  })

  return {
    currentIndex,
    currentStage,
    progress,
    isLastStage,
    finishing
  }
}
