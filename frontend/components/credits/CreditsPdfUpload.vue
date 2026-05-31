<script setup lang="ts">
import { FileUp } from 'lucide-vue-next'

const { scanResult, scanLoading, scanError, scanContract } = useCredits()

const inputRef = ref<HTMLInputElement | null>(null)

function openPicker() {
  inputRef.value?.click()
}

async function onChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  try {
    await scanContract(file)
  } catch {
    // scanContract записывает текст в scanError
  }
}
</script>

<template>
  <div class="space-y-4">
    <input
      ref="inputRef"
      type="file"
      accept="application/pdf"
      class="sr-only"
      @change="onChange"
    />
    <p class="text-xs text-muted-foreground">
      Кредит, займ или кредит наличными — загрузите ваш договор или «Индивидуальные условия», не общие условия продукта.
    </p>
    <Button type="button" class="gap-2" data-demo="credit-pdf-upload" :disabled="scanLoading" @click="openPicker">
      <FileUp class="size-4" />
      {{ scanLoading ? 'Распознаём…' : 'Выбрать PDF договора' }}
    </Button>

    <Alert v-if="scanError" variant="destructive">
      <AlertDescription>{{ scanError }}</AlertDescription>
    </Alert>

    <CreditsScanReport v-if="scanResult" :scan="scanResult" />
  </div>
</template>
