<script setup lang="ts">
const authStore = useAuthStore()
</script>

<template>
  <section class="mm-card p-4 sm:p-6">
    <h2 class="text-lg font-semibold text-[color:var(--mm-text)]">Профиль</h2>
    <dl v-if="authStore.user" class="mt-4 space-y-3 text-sm">
      <div class="flex flex-col gap-0.5 sm:flex-row sm:gap-2">
        <dt class="text-[color:var(--mm-text-soft)]">Телефон</dt>
        <dd class="font-medium text-[color:var(--mm-text)]">{{ authStore.user.phone }}</dd>
      </div>
      <div v-if="authStore.user.name" class="flex flex-col gap-0.5 sm:flex-row sm:gap-2">
        <dt class="text-[color:var(--mm-text-soft)]">Имя</dt>
        <dd class="font-medium text-[color:var(--mm-text)]">{{ authStore.user.name }}</dd>
      </div>
    </dl>
    <p class="mt-6 text-sm text-[color:var(--mm-text-soft)]">Тёмная тема — в следующих версиях</p>
  </section>
</template>
