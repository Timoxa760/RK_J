export function useRevealOnScroll(options?: {
  threshold?: number
  immediate?: boolean
  delay?: number
}) {
  const root = ref<HTMLElement | null>(null)
  const visible = ref(false)
  const animate = ref(false)
  let observer: IntersectionObserver | null = null
  let delayTimer: ReturnType<typeof setTimeout> | null = null

  onMounted(() => {
    if (!root.value) return

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reduced) {
      visible.value = true
      return
    }

    animate.value = true

    if (options?.immediate) {
      const delay = options.delay ?? 0
      if (delay > 0) {
        delayTimer = setTimeout(() => {
          visible.value = true
        }, delay)
      } else {
        visible.value = true
      }
      return
    }

    const rect = root.value.getBoundingClientRect()
    if (rect.top < window.innerHeight * 0.92) {
      visible.value = true
      return
    }

    observer = new IntersectionObserver(
      ([entry]) => {
        if (entry?.isIntersecting) {
          visible.value = true
          observer?.disconnect()
        }
      },
      {
        threshold: options?.threshold ?? 0.1,
        rootMargin: '0px 0px -4% 0px'
      }
    )

    observer.observe(root.value)
  })

  onUnmounted(() => {
    observer?.disconnect()
    if (delayTimer) clearTimeout(delayTimer)
  })

  return { root, visible, animate }
}
