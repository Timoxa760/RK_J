<script setup lang="ts">
import type { ProviderId } from '~/types/api'
import { useAuthStore } from '~/store/authStore'

const authStore = useAuthStore()
const { providers, connected, loading, error, success, connect, sync } = useProviders()

const connectProvider = ref<ProviderId | null>(null)
const credPhone = ref('')
const credPassword = ref('')

function openConnect(id: ProviderId) {
  connectProvider.value = id
  credPhone.value = authStore.user?.phone ?? ''
  credPassword.value = ''
}

async function submitConnect() {
  if (!connectProvider.value) return
  await connect(connectProvider.value, {
    phone: credPhone.value,
    password: credPassword.value
  })
  if (!error.value) connectProvider.value = null
}
</script>

<template>
  <div class="space-y-6">
    <Card>
      <CardHeader>
        <CardTitle class="text-base">Аккаунт</CardTitle>
      </CardHeader>
      <CardContent>
        <dl v-if="authStore.user" class="space-y-3 text-sm">
          <div class="flex flex-col gap-0.5 sm:flex-row sm:gap-2">
            <dt class="text-muted-foreground">Телефон</dt>
            <dd class="font-medium">{{ authStore.user.phone }}</dd>
          </div>
          <div v-if="authStore.user.name" class="flex flex-col gap-0.5 sm:flex-row sm:gap-2">
            <dt class="text-muted-foreground">Имя</dt>
            <dd class="font-medium">{{ authStore.user.name }}</dd>
          </div>
        </dl>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-base">Мои магазины</CardTitle>
        <CardDescription>Подключите личный кабинет — чеки подтянутся автоматически</CardDescription>
      </CardHeader>
      <CardContent>
        <Alert v-if="error" variant="destructive" class="mb-3">
          <AlertDescription>{{ error }}</AlertDescription>
        </Alert>
        <p v-if="success" class="mb-3 text-sm text-primary">{{ success }}</p>

        <ul class="space-y-2">
          <li
            v-for="p in providers"
            :key="p.id"
            class="flex flex-col gap-2 rounded-lg border px-4 py-3 sm:flex-row sm:items-center sm:justify-between"
          >
            <div>
              <p class="font-medium">{{ p.label }}</p>
              <p class="text-xs text-muted-foreground">
                {{ connected[p.id] ? `Статус: ${connected[p.id]}` : 'Не подключён' }}
              </p>
            </div>
            <div class="flex gap-2">
              <Button
                variant="secondary"
                size="sm"
                :disabled="loading === p.id"
                @click="openConnect(p.id)"
              >
                Подключить
              </Button>
              <Button
                v-if="connected[p.id]"
                variant="outline"
                size="sm"
                :disabled="loading === p.id"
                @click="sync(p.id)"
              >
                Синхр.
              </Button>
            </div>
          </li>
        </ul>

        <p class="mt-4 text-xs text-muted-foreground">
          После подключения откройте
          <NuxtLink to="/receipts" class="text-primary underline">ленту расходов</NuxtLink>.
        </p>
      </CardContent>
    </Card>

    <Dialog :open="Boolean(connectProvider)" @update:open="(v) => !v && (connectProvider = null)">
      <DialogContent v-if="connectProvider" class="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>
            Подключить {{ providers.find((x) => x.id === connectProvider)?.label }}
          </DialogTitle>
        </DialogHeader>
        <form class="space-y-4" @submit.prevent="submitConnect">
          <div class="space-y-2">
            <Label for="cred-phone">Телефон</Label>
            <Input id="cred-phone" v-model="credPhone" type="tel" placeholder="Телефон" />
          </div>
          <div class="space-y-2">
            <Label for="cred-pass">Пароль от ЛК</Label>
            <Input id="cred-pass" v-model="credPassword" type="password" placeholder="Пароль" />
          </div>
          <DialogFooter class="gap-2 sm:gap-0">
            <Button type="button" variant="outline" @click="connectProvider = null">Отмена</Button>
            <Button type="submit" :disabled="!!loading">
              {{ loading ? 'Подключение…' : 'Подключить' }}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>
