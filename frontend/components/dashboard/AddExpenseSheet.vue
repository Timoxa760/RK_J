<script setup lang="ts">
const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
  added: []
}>()

function onAdded() {
  emit('added')
}

function closeSheet() {
  open.value = false
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent
      class="mm-purchase-sheet max-h-[min(90dvh,720px)] w-[calc(100%-2rem)] max-w-md gap-0 overflow-y-auto rounded-2xl p-0 sm:max-w-lg"
    >
      <DialogHeader class="space-y-2 px-5 pb-0 pt-5 text-left">
        <DialogTitle class="text-xl font-semibold tracking-tight">
          Добавить покупку
        </DialogTitle>
        <DialogDescription class="text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
          Одна покупка — и Поток обновит картину денег. Как в онбординге: голос, вручную или чек.
        </DialogDescription>
      </DialogHeader>

      <div class="px-5 pb-5 pt-4">
        <PurchaseInputTabs
          v-if="open"
          show-photo
          dismissible
          @added="onAdded"
          @confirmed="closeSheet"
        />
      </div>
    </DialogContent>
  </Dialog>
</template>
