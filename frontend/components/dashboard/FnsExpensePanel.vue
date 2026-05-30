<script setup lang="ts">
import { FileCheck } from 'lucide-vue-next'

const emit = defineEmits<{
  synced: []
}>()

const { loading, submitQr } = useFns()
const qr = ref('')

async function onTicket() {
  if (!qr.value.trim()) return
  await submitQr(qr.value)
  qr.value = ''
  emit('synced')
}
</script>

<template>
  <div class="space-y-4">
    <p class="text-sm text-muted-foreground">
      <strong class="font-medium text-foreground">ФНС — автоматический слой.</strong>
      Чеки с онлайн-касс по QR. Не обязательно для старта.
    </p>

    <Card>
      <CardHeader class="pb-2">
        <CardDescription>QR чека</CardDescription>
      </CardHeader>
      <CardContent class="space-y-3">
        <Textarea
          v-model="qr"
          rows="3"
          class="font-mono text-xs"
          placeholder="t=20260530T1200&s=5000.00&fn=...&i=...&fp=..."
          :disabled="loading"
        />
        <Button :disabled="loading || !qr.trim()" @click="onTicket">
          <FileCheck class="mr-2 size-4" />
          Проверить чек
        </Button>
      </CardContent>
    </Card>

    <Card class="border-dashed opacity-70">
      <CardHeader class="pb-2">
        <CardTitle class="text-base">«Мои чеки» (MCO)</CardTitle>
        <CardDescription>Подключение через OAuth — скоро</CardDescription>
      </CardHeader>
      <CardContent>
        <Button variant="secondary" disabled>
          Подключить «Мои чеки»
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
