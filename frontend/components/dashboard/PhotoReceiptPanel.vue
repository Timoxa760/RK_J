<script setup lang="ts">
import { Camera, Loader2 } from 'lucide-vue-next'

defineProps<{
  embedded?: boolean
}>()

const emit = defineEmits<{
  synced: []
}>()

const { loading, submitQr } = useFns()
const fileInput = ref<HTMLInputElement | null>(null)
const scanError = ref('')
const previewUrl = ref<string | null>(null)

async function decodeQrFromFile(file: File) {
  scanError.value = ''
  previewUrl.value = URL.createObjectURL(file)

  if (!('BarcodeDetector' in window)) {
    scanError.value =
      'Браузер не читает QR с фото — вставьте строку QR на вкладке «ФНС» или используйте голос.'
    return
  }

  try {
    const bitmap = await createImageBitmap(file)
    const detector = new BarcodeDetector({ formats: ['qr_code'] })
    const codes = await detector.detect(bitmap)
    bitmap.close()

    const raw = codes[0]?.rawValue?.trim()
    if (!raw) {
      scanError.value = 'QR на фото не найден — попробуйте другой ракурс или вкладку «ФНС».'
      return
    }

    await submitQr(raw)
    emit('synced')
  } catch (e) {
    scanError.value =
      e instanceof Error ? e.message : 'Не удалось прочитать QR — попробуйте вкладку «ФНС».'
  }
}

function onPickClick() {
  fileInput.value?.click()
}

function onFileChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) void decodeQrFromFile(file)
}

onBeforeUnmount(() => {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
})
</script>

<template>
  <div class="space-y-4">
    <p v-if="!embedded" class="text-sm text-muted-foreground">
      Сфотографируйте QR на бумажном чеке — Поток проверит его через ФНС.
    </p>

    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      capture="environment"
      class="sr-only"
      @change="onFileChange"
    />

    <Button type="button" class="w-full" :disabled="loading" @click="onPickClick">
      <Loader2 v-if="loading" class="mr-2 size-4 animate-spin" />
      <Camera v-else class="mr-2 size-4" />
      {{ loading ? 'Проверяем…' : 'Снять или выбрать фото чека' }}
    </Button>

    <img
      v-if="previewUrl"
      :src="previewUrl"
      alt="Предпросмотр чека"
      class="mx-auto max-h-40 rounded-lg border object-contain"
    />

    <p v-if="scanError" class="text-sm text-destructive">{{ scanError }}</p>
  </div>
</template>
