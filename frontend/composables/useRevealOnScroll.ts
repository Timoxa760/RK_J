export function useRevealOnScroll(options?: { threshold?: number }) {
  const root = ref<HTMLElement | null>(null)
  const visible = ref(false)
  let observer: IntersectionObserver | null = null

  onMounted(() => {
    if (!root.value) return

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reduced) {
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
        threshold: options?.threshold ?? 0.12,
        rootMargin: '0px 0px -6% 0px'
      }
    )

    observer.observe(root.value)
  })

  onUnmounted(() => {
    observer?.disconnect()
  })

  return { root, visible }
}
