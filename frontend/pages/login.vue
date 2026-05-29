<script setup lang="ts">

const { register, login } = useAuth()

const phone = ref('')
const code = ref('')
const step = ref<'phone' | 'code'>('phone')
const loading = ref(false)
const error = ref('')

function normalizePhone(value: string) {
  const digits = value.replace(/\D/g, '')
  if (digits.startsWith('8') && digits.length === 11) {
    return `+7${digits.slice(1)}`
  }
  if (digits.startsWith('7') && digits.length === 11) {
    return `+${digits}`
  }
  if (digits.length === 10) {
    return `+7${digits}`
  }
  return value.startsWith('+') ? value : `+${digits}`
}

function formatPhoneInput(e: Event) {
  const input = e.target as HTMLInputElement
  let digits = input.value.replace(/\D/g, '')
  if (!digits.length) {
    phone.value = ''
    return
  }
  if (digits[0] === '8') digits = `7${digits.slice(1)}`
  if (digits[0] !== '7') digits = `7${digits}`
  digits = digits.slice(0, 11)
  const parts = [
    '+7',
    digits.length > 1 ? ` (${digits.slice(1, 4)}` : '',
    digits.length >= 4 ? `) ${digits.slice(4, 7)}` : '',
    digits.length >= 7 ? `-${digits.slice(7, 9)}` : '',
    digits.length >= 9 ? `-${digits.slice(9, 11)}` : ''
  ]
  phone.value = parts.join('')
}

async function submitPhone() {
  error.value = ''
  const normalized = normalizePhone(phone.value)
  if (normalized.replace(/\D/g, '').length < 11) {
    error.value = 'Введите корректный номер телефона'
    return
  }
  loading.value = true
  try {
    await register(normalized)
    phone.value = normalized
    step.value = 'code'
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}

async function submitCode() {
  error.value = ''
  if (code.value.length < 4) {
    error.value = 'Введите код из SMS'
    return
  }
  loading.value = true
  try {
    await login(phone.value, code.value)
    await navigateTo('/dashboard')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка входа'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="relative z-10 w-full max-w-md mm-card rounded-2xl p-6 sm:rounded-3xl sm:p-8">
    <h1 class="text-2xl font-medium tracking-tight text-[color:var(--mm-text)]">Вход в Поток</h1>
    <p class="mt-2 text-sm text-[color:var(--mm-text-muted)]">
      {{
        step === 'phone'
          ? 'Войдите, чтобы сохранить историю разговоров с помощником'
          : 'Введите код из SMS'
      }}
    </p>

    <form
      v-if="step === 'phone'"
      class="mt-6 space-y-4"
      @submit.prevent="submitPhone"
    >
      <div>
        <label class="block text-sm font-medium text-[color:var(--mm-text-muted)]" for="phone">Телефон</label>
        <input
          id="phone"
          :value="phone"
          type="tel"
          placeholder="+7 (999) 123-45-67"
          class="mt-1 w-full rounded-xl border border-[color:var(--mm-border)] bg-[color:var(--mm-bg-elevated)] px-3 py-2 text-[color:var(--mm-text)] focus:border-[color:var(--mm-primary)] focus:outline-none focus:ring-1 focus:ring-[color:var(--mm-primary-muted)]"
          @input="formatPhoneInput"
        />
      </div>
      <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
      <button type="submit" class="mm-btn-primary w-full disabled:opacity-50" :disabled="loading">
        {{ loading ? 'Отправка…' : 'Получить код' }}
      </button>
    </form>

    <form
      v-else
      class="mt-6 space-y-4"
      @submit.prevent="submitCode"
    >
      <p class="text-sm text-[color:var(--mm-text-soft)]">{{ phone }}</p>
      <div>
        <label class="block text-sm font-medium text-[color:var(--mm-text-muted)]" for="code">Код</label>
        <input
          id="code"
          v-model="code"
          type="text"
          inputmode="numeric"
          maxlength="4"
          placeholder="Код из SMS"
          class="mt-1 w-full rounded-xl border border-[color:var(--mm-border)] bg-[color:var(--mm-bg-elevated)] px-3 py-2 text-[color:var(--mm-text)] focus:border-[color:var(--mm-primary)] focus:outline-none focus:ring-1 focus:ring-[color:var(--mm-primary-muted)]"
        />
      </div>
      <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
      <button type="submit" class="mm-btn-primary w-full disabled:opacity-50" :disabled="loading">
        {{ loading ? 'Вход…' : 'Войти' }}
      </button>
      <button
        type="button"
        class="w-full text-sm text-[color:var(--mm-text-soft)] hover:text-[color:var(--mm-text-muted)]"
        @click="step = 'phone'; code = ''; error = ''"
      >
        Изменить номер
      </button>
    </form>
  </div>
</template>
