<script setup lang="ts">
import { FileCheck, RefreshCw } from 'lucide-vue-next'

const emit = defineEmits<{
  synced: []
}>()

const { loading: qrLoading, submitQr } = useFns()
const { loading: mcoLoading, authStep, phone, startAuth, verifyAuth, syncReceipts } = useFnsMco()

const qr = ref('')
const smsCode = ref('')

const busy = computed(() => qrLoading.value || mcoLoading.value)

async function onTicket() {
  if (!qr.value.trim()) return
  await submitQr(qr.value)
  qr.value = ''
  emit('synced')
}

async function onVerify() {
  if (!smsCode.value.trim()) return
  await verifyAuth(smsCode.value)
}

async function onSync() {
  const count = await syncReceipts()
  if (count > 0) emit('synced')
}
</script>

<template>
  <div class="space-y-4">
    <p class="text-sm text-muted-foreground">
      <strong class="font-medium text-foreground">ФНС — автоматический слой.</strong>
      Чеки с онлайн-касс по QR или через «Мои чеки».
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
          :disabled="busy"
        />
        <Button :disabled="busy || !qr.trim()" @click="onTicket">
          <FileCheck class="mr-2 size-4" />
          Проверить чек
        </Button>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="pb-2">
        <CardTitle class="text-base">«Мои чеки» (MCO)</CardTitle>
        <CardDescription>
          Синхронизация покупок с lkdr.nalog.ru · телефон {{ phone || '—' }}
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-3">
        <template v-if="authStep === 'idle'">
          <Button variant="secondary" :disabled="busy || !phone" @click="startAuth">
            Подключить «Мои чеки»
          </Button>
        </template>

        <template v-else-if="authStep === 'code_sent'">
          <Input
            v-model="smsCode"
            inputmode="numeric"
            maxlength="6"
            placeholder="Код из SMS (demo: 0000)"
            :disabled="busy"
          />
          <Button :disabled="busy || !smsCode.trim()" @click="onVerify">
            Подтвердить
          </Button>
        </template>

        <template v-else>
          <p class="text-sm text-muted-foreground">Аккаунт ФНС подключён.</p>
          <Button :disabled="busy" @click="onSync">
            <RefreshCw class="mr-2 size-4" :class="{ 'animate-spin': mcoLoading }" />
            Синхронизировать чеки
          </Button>
        </template>
      </CardContent>
    </Card>
  </div>
</template>
