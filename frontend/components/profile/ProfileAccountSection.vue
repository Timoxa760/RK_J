<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'
import { useAuthStore } from '~/store/authStore'

defineProps<{
  showSurveyPrompt?: boolean
}>()

const emit = defineEmits<{
  retakeSurvey: []
}>()

const authStore = useAuthStore()
const { logout: signOut } = useAuth()

function logout() {
  signOut()
  navigateTo('/login')
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-base">Аккаунт</CardTitle>
      <CardDescription>Вход по номеру телефона — имя и магазины не нужны</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <dl v-if="authStore.user" class="space-y-3 text-sm">
        <div class="flex flex-col gap-0.5 sm:flex-row sm:gap-2">
          <dt class="text-muted-foreground">Телефон</dt>
          <dd class="font-medium tabular-nums">{{ authStore.user.phone }}</dd>
        </div>
      </dl>

      <div class="flex flex-wrap gap-2">
        <Button
          v-if="showSurveyPrompt"
          type="button"
          variant="secondary"
          size="sm"
          @click="emit('retakeSurvey')"
        >
          Пройти опрос заново
        </Button>
        <Button type="button" variant="outline" size="sm" @click="logout">
          Выйти
        </Button>
      </div>
    </CardContent>
  </Card>
</template>
