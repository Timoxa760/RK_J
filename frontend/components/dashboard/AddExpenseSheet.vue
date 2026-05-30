<script setup lang="ts">
import { Camera, Mic, PenLine, Receipt } from 'lucide-vue-next'

const open = defineModel<boolean>('open', { default: false })

const tab = ref('voice')
const { submitting, submitManual } = useReceiptSubmit()

const emit = defineEmits<{
  added: []
}>()

async function onManual(payload: {
  store: string
  amount: number
  category: string
  date: string
}) {
  await submitManual(payload)
  finish()
}

function finish() {
  emit('added')
  open.value = false
}

function onFnsSynced() {
  emit('added')
}

watch(open, (v) => {
  if (v) tab.value = 'voice'
})
</script>

<template>
  <Sheet v-model:open="open">
    <SheetContent side="bottom" class="max-h-[90dvh] overflow-y-auto rounded-t-2xl sm:max-w-lg sm:rounded-2xl">
      <SheetHeader>
        <SheetTitle>Добавить</SheetTitle>
        <SheetDescription>Один способ — Поток разберёт и обновит картину</SheetDescription>
      </SheetHeader>

      <Tabs v-model="tab" class="mt-4">
        <TabsList class="grid w-full grid-cols-4">
          <TabsTrigger value="voice" class="gap-1 text-xs">
            <Mic class="size-3.5" />
            Голос
          </TabsTrigger>
          <TabsTrigger value="manual" class="gap-1 text-xs">
            <PenLine class="size-3.5" />
            Вручную
          </TabsTrigger>
          <TabsTrigger value="fns" class="gap-1 text-xs">
            <Receipt class="size-3.5" />
            ФНС
          </TabsTrigger>
          <TabsTrigger value="photo" class="gap-1 text-xs">
            <Camera class="size-3.5" />
            Чек
          </TabsTrigger>
        </TabsList>

        <TabsContent value="voice" class="mt-4">
          <DashboardVoiceExpenseInput @done="finish" />
        </TabsContent>
        <TabsContent value="manual" class="mt-4">
          <DashboardManualExpenseForm :busy="submitting" @submit="onManual" />
        </TabsContent>
        <TabsContent value="fns" class="mt-4">
          <DashboardFnsExpensePanel @synced="onFnsSynced" />
        </TabsContent>
        <TabsContent value="photo" class="mt-4 space-y-3 text-center">
          <p class="text-sm text-muted-foreground">
            Скан фото чека скоро появится. Сейчас — QR ФНС, голос или ручной ввод.
          </p>
          <Button variant="secondary" as-child>
            <NuxtLink to="/dashboard" @click="open = false">Смотреть категории на главной</NuxtLink>
          </Button>
        </TabsContent>
      </Tabs>
    </SheetContent>
  </Sheet>
</template>
