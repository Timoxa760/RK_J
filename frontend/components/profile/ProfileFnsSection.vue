<script setup lang="ts">
import { RefreshCw, ShieldCheck } from 'lucide-vue-next'
import { FNS } from '~/constants/productCopy'
import { useAuthStore } from '~/store/authStore'

const authStore = useAuthStore()
const {
  connection,
  syncing,
  connectBusy,
  connectLoading,
  connectStep,
  connectPhone,
  error,
  success,
  importedCount,
  pendingCount,
  sendConnectCode,
  verifyConnectCode,
  resetConnectFlow,
  backToPhoneStep,
  sync,
  disconnect
} = useFns()

const phone = ref('')
const smsCode = ref('')
const connectOpen = ref(false)

onMounted(() => {
  phone.value = connection.value.phone || authStore.user?.phone || ''
})

function formatSyncTime(iso?: string): string | null {
  if (!iso) return null
  return new Date(iso).toLocaleString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const lastSyncLabel = computed(() => formatSyncTime(connection.value.last_sync_at))

function openConnectDialog() {
  smsCode.value = ''
  resetConnectFlow()
  connectOpen.value = true
}

function closeConnectDialog() {
  connectOpen.value = false
  smsCode.value = ''
  resetConnectFlow()
}

async function submitPhoneStep() {
  await sendConnectCode(phone.value)
}

async function submitCodeStep() {
  const result = await verifyConnectCode(smsCode.value)
  if (result && !error.value) {
    closeConnectDialog()
  }
}

async function resendCode() {
  smsCode.value = ''
  await sendConnectCode(phone.value)
}
</script>

<template>
  <Card data-demo="profile-fns">
    <CardHeader>
      <CardTitle class="text-base">{{ FNS.title }}</CardTitle>
      <CardDescription>{{ FNS.hint }}</CardDescription>
    </CardHeader>

    <CardContent class="space-y-4">
      <Alert v-if="error && !connectOpen" variant="destructive">
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>
      <p v-if="success && !connectOpen" class="text-sm text-primary">{{ success }}</p>

      <div class="rounded-xl border px-4 py-3">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div class="space-y-1">
            <div class="flex items-center gap-2">
              <ShieldCheck
                class="size-4"
                :class="connection.connected ? 'text-primary' : 'text-muted-foreground'"
              />
              <p class="font-medium">
                {{ connection.connected ? FNS.statusConnected : FNS.statusDisconnected }}
              </p>
            </div>
            <p v-if="connection.connected && connection.phone" class="text-xs text-muted-foreground">
              {{ connection.phone }}
            </p>
            <p v-if="connection.connected && importedCount" class="text-xs text-muted-foreground">
              {{ FNS.importedLine(importedCount) }}
            </p>
            <p v-if="lastSyncLabel" class="text-xs text-muted-foreground">
              {{ FNS.lastSync(lastSyncLabel) }}
            </p>
          </div>

          <div class="flex flex-wrap gap-2">
            <Button
              v-if="!connection.connected"
              type="button"
              size="sm"
              :disabled="connectLoading"
              @click="openConnectDialog"
            >
              {{ connectLoading ? FNS.connecting : FNS.connect }}
            </Button>
            <template v-else>
              <Button
                type="button"
                size="sm"
                variant="secondary"
                class="gap-2"
                :disabled="syncing"
                @click="sync()"
              >
                <RefreshCw class="size-4" :class="{ 'animate-spin': syncing }" />
                {{ syncing ? FNS.syncing : FNS.sync }}
              </Button>
              <Button type="button" size="sm" variant="outline" :disabled="syncing" @click="disconnect">
                {{ FNS.disconnect }}
              </Button>
            </template>
          </div>
        </div>

        <p
          v-if="connection.connected && pendingCount > 0"
          class="mt-3 text-xs text-muted-foreground"
        >
          {{ FNS.pendingLine(pendingCount) }}
        </p>
      </div>

      <Button v-if="connection.connected" variant="link" class="h-auto p-0 text-sm" as-child>
        <NuxtLink to="/receipts">{{ FNS.openReceipts }}</NuxtLink>
      </Button>
    </CardContent>
  </Card>

  <Dialog
    :open="connectOpen"
    @update:open="(value) => (value ? (connectOpen = true) : closeConnectDialog())"
  >
    <DialogContent class="min-w-0 overflow-x-hidden sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ FNS.connect }}</DialogTitle>
        <DialogDescription>
          {{
            connectStep === 'phone'
              ? FNS.connectDialog
              : FNS.codeSent(connectPhone)
          }}
        </DialogDescription>
      </DialogHeader>

      <form
        v-if="connectStep === 'phone'"
        class="min-w-0 space-y-4"
        @submit.prevent="submitPhoneStep"
      >
        <div class="min-w-0 space-y-2">
          <Label for="fns-phone">{{ FNS.phoneLabel }}</Label>
          <Input id="fns-phone" v-model="phone" type="tel" autocomplete="tel" class="w-full" />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <DialogFooter class="flex-col gap-2 sm:flex-col">
          <Button type="button" variant="outline" class="w-full sm:w-auto" @click="closeConnectDialog">
            Отмена
          </Button>
          <Button type="submit" class="w-full sm:w-auto" :disabled="connectBusy">
            {{ connectBusy ? FNS.sendingCode : FNS.sendCode }}
          </Button>
        </DialogFooter>
      </form>

      <form v-else class="min-w-0 space-y-4" @submit.prevent="submitCodeStep">
        <div class="min-w-0 space-y-2">
          <Label for="fns-code">{{ FNS.codeLabel }}</Label>
          <Input
            id="fns-code"
            v-model="smsCode"
            type="text"
            inputmode="numeric"
            autocomplete="one-time-code"
            maxlength="8"
            class="w-full"
            :placeholder="FNS.codePlaceholder"
          />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <div class="flex flex-wrap gap-x-4 gap-y-1">
          <Button
            type="button"
            variant="link"
            class="h-auto p-0 text-sm"
            :disabled="connectLoading"
            @click="backToPhoneStep"
          >
            {{ FNS.changePhone }}
          </Button>
          <Button
            type="button"
            variant="link"
            class="h-auto p-0 text-sm"
            :disabled="connectBusy"
            @click="resendCode"
          >
            {{ connectBusy ? FNS.sendingCode : FNS.resendCode }}
          </Button>
        </div>
        <DialogFooter class="flex-col gap-2 sm:flex-col">
          <Button
            type="button"
            variant="outline"
            class="w-full sm:w-auto"
            :disabled="connectLoading"
            @click="closeConnectDialog"
          >
            Отмена
          </Button>
          <Button type="submit" class="w-full sm:w-auto" :disabled="connectLoading">
            {{
              syncing
                ? FNS.syncing
                : connectBusy
                  ? FNS.verifying
                  : FNS.confirmCode
            }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
