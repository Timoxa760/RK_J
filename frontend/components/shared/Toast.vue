<script setup lang="ts">
const toast = useState<{ message: string; type: 'error' | 'success' | 'info' } | null>(
  'app-toast',
  () => null
)

let hideTimer: ReturnType<typeof setTimeout> | null = null

watch(toast, (value) => {
  if (hideTimer) clearTimeout(hideTimer)
  if (value) {
    hideTimer = setTimeout(() => {
      toast.value = null
    }, 4000)
  }
})

const typeClass = {
  error: 'border-red-200 bg-red-50 text-red-800',
  success: 'border-[color:var(--mm-primary-muted)] bg-[color:var(--mm-primary-soft)] text-[color:var(--mm-primary-hover)]',
  info: 'border-[color:var(--mm-primary-muted)] bg-[color:var(--mm-primary-soft)] text-[color:var(--mm-primary-hover)]'
}
</script>

<template>
  <Teleport to="body">
    <Transition name="toast">
      <div
        v-if="toast"
        class="fixed bottom-4 left-4 right-4 z-[100] mx-auto max-w-sm rounded-xl border px-4 py-3 text-sm font-medium shadow-lg sm:bottom-6 sm:left-auto sm:right-6 mm-safe-bottom"
        :class="typeClass[toast.type]"
        role="status"
      >
        {{ toast.message }}
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(0.5rem);
}
</style>
