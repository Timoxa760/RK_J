<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{
  error: NuxtError
}>()

const message = computed(() => {
  const raw = props.error?.message ?? 'Неизвестная ошибка'
  if (raw.includes('__vrv_devtools')) {
    return 'Ошибка инициализации страницы. Перезапустите dev-сервер (rm -rf .nuxt node_modules/.cache && npm run dev).'
  }
  return raw
})

function retry() {
  clearError({ redirect: '/' })
}
</script>

<template>
  <div class="flex min-h-[100dvh] items-center justify-center bg-[#f9fcff] px-4">
    <div class="w-full max-w-md rounded-2xl border border-[#e0d6cc] bg-white p-6 shadow-lg">
      <h1 class="text-lg font-semibold text-[#1c1c1a]">
        {{ error.statusCode === 404 ? 'Страница не найдена' : 'Ошибка' }}
      </h1>
      <p class="mt-2 text-sm text-[#6d6760]">{{ message }}</p>
      <Button class="mt-4 w-full" type="button" @click="retry">
        На главную
      </Button>
    </div>
  </div>
</template>
