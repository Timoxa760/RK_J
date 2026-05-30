<script setup lang="ts">
import { allowDigitKeydown, digitsOnly, onDigitPaste } from '~/utils/numericInput'

const { register, login } = useAuth()

const phoneDigits = ref('')
const codeDigits = ref('')
const step = ref<'phone' | 'code'>('phone')
const loading = ref(false)
const error = ref('')

const phone = computed({
  get(): string {
    const d = phoneDigits.value
    if (!d.length) return '+7'
    return [
      '+7',
      d.length > 0 ? ` (${d.slice(0, 3)}` : '',
      d.length >= 3 ? `) ${d.slice(3, 6)}` : '',
      d.length >= 6 ? `-${d.slice(6, 8)}` : '',
      d.length >= 8 ? `-${d.slice(8, 10)}` : ''
    ].join('')
  },
  set(raw: string | number) {
    let digits = digitsOnly(String(raw), 11)
    if (digits.startsWith('8')) digits = digits.slice(1)
    if (digits.startsWith('7')) digits = digits.slice(1)
    phoneDigits.value = digits.slice(0, 10)
  }
})

const code = computed({
  get: () => codeDigits.value,
  set(raw: string | number) {
    codeDigits.value = digitsOnly(String(raw), 4)
  }
})

const phoneComplete = computed(() => phoneDigits.value.length === 10)

function onPhonePaste(event: ClipboardEvent) {
  onDigitPaste(event, (digits) => {
    phone.value = digits
  }, 11)
}

function onCodePaste(event: ClipboardEvent) {
  onDigitPaste(event, (digits) => {
    codeDigits.value = digits
  }, 4)
}

function normalizePhone(display: string) {
  const digits = display.replace(/\D/g, '')
  if (digits.length === 11 && digits.startsWith('7')) return `+${digits}`
  if (digits.length === 10) return `+7${digits}`
  return display
}

async function submitPhone() {
  error.value = ''
  if (!phoneComplete.value) {
    error.value = 'Введите номер полностью — 10 цифр после +7'
    return
  }
  loading.value = true
  try {
    const normalized = normalizePhone(phone.value)
    await register(normalized)
    phoneDigits.value = normalized.replace(/\D/g, '').replace(/^7/, '').slice(-10)
    step.value = 'code'
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}

function backToPhone() {
  step.value = 'phone'
  codeDigits.value = ''
  error.value = ''
}

async function submitCode() {
  error.value = ''
  if (codeDigits.value.length < 4) {
    error.value = 'Введите 4 цифры кода'
    return
  }
  loading.value = true
  try {
    await login(normalizePhone(phone.value), codeDigits.value)
    useFinancialProfile().loadProfile()
    useGoals().fetchGoals()
    const { isComplete } = useOnboarding()
    await navigateTo(isComplete() ? '/dashboard' : '/onboarding')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка входа'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Card
    class="relative z-10 w-full max-w-md border-[color:var(--mm-border)] bg-white/95 shadow-[0_24px_64px_-28px_rgb(0_0_0_/_0.18)] backdrop-blur-md"
  >
    <CardHeader>
      <CardTitle class="text-2xl">Вход в Поток</CardTitle>
      <CardDescription>
        {{
          step === 'phone'
            ? 'Войдите, чтобы сохранить историю разговоров с помощником'
            : 'Введите код из SMS'
        }}
      </CardDescription>
    </CardHeader>
    <CardContent>
      <form v-if="step === 'phone'" class="space-y-4" @submit.prevent="submitPhone">
        <div class="space-y-2">
          <Label for="phone">Телефон</Label>
          <Input
            id="phone"
            v-model="phone"
            type="tel"
            inputmode="numeric"
            pattern="[0-9+() -]*"
            autocomplete="tel"
            maxlength="18"
            placeholder="+7 (999) 123-45-67"
            @keydown="allowDigitKeydown"
            @paste="onPhonePaste"
          />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <Button type="submit" class="w-full" :disabled="loading || !phoneComplete">
          {{ loading ? 'Отправка…' : 'Получить код' }}
        </Button>
      </form>

      <form v-else class="space-y-4" @submit.prevent="submitCode">
        <p class="text-sm text-muted-foreground">{{ phone }}</p>
        <div class="space-y-2">
          <Label for="code">Код</Label>
          <Input
            id="code"
            v-model="code"
            type="text"
            inputmode="numeric"
            pattern="[0-9]*"
            autocomplete="one-time-code"
            maxlength="4"
            placeholder="0000"
            @keydown="allowDigitKeydown"
            @paste="onCodePaste"
          />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <Button type="submit" class="w-full" :disabled="loading || codeDigits.length < 4">
          {{ loading ? 'Вход…' : 'Войти' }}
        </Button>
        <Button type="button" variant="ghost" class="w-full" @click="backToPhone">
          Изменить номер
        </Button>
      </form>
    </CardContent>
  </Card>
</template>
