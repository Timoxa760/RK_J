<script setup lang="ts">
import { useSidebar } from '~/components/ui/sidebar/utils'

const { bootstrap } = useAdvisorShellProvider()
const { setOpenMobile, isMobile } = useSidebar()
const { pending, clear } = usePendingMobileSidebar()

watch(
  pending,
  (open) => {
    if (!open || !isMobile.value) return
    setOpenMobile(true)
    clear()
  },
  { flush: 'post' }
)

onMounted(async () => {
  if (pending.value && isMobile.value) {
    setOpenMobile(true)
    clear()
  }
  await bootstrap()
})
</script>

<template>
  <span class="sr-only" aria-hidden="true" />
</template>
