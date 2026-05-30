<script setup lang="ts">
import { allowDigitKeydown, digitsOnly, onDigitPaste } from '~/utils/numericInput'
import { clearAnonymousOnboardingKeys, isOnboardingComplete } from '~/composables/useOnboarding'
import { useAuthStore } from '~/store/authStore'

type Mode = 'login' | 'register' | 'forgot' | 'reset'

const { register, login, requestPasswordReset, resetPassword } = useAuth()
const authStore = useAuthStore()

const mode = ref<Mode>('login')
const phoneDigits = ref('')
const password = ref('')
const passwordConfirm = ref('')
const resetCode = ref('')
const loading = ref(false)
const error = ref('')
const info = ref('')

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

const phoneComplete = computed(() => phoneDigits.value.length === 10)
const passwordOk = computed(() => password.value.length >= 8)
const passwordsMatch = computed(() => password.value === passwordConfirm.value)

const titles: Record<Mode, string> = {
  login: 'Вход в Поток',
  register: 'Регистрация',
  forgot: 'Восстановление пароля',
  reset: 'Новый пароль'
}

const descriptions: Record<Mode, string> = {
  login: 'Войдите, чтобы сохранить историю разговоров с помощником',
  register: 'Создайте аккаунт — телефон и пароль от 8 символов',
  forgot: 'Укажите номер — мы отправим код для сброса',
  reset: 'Введите код и новый пароль'
}

function onPhonePaste(event: ClipboardEvent) {
  onDigitPaste(event, (digits) => {
    phone.value = digits
  }, 11)
}

function normalizePhone(display: string) {
  const digits = display.replace(/\D/g, '')
  if (digits.length === 11 && digits.startsWith('7')) return `+${digits}`
  if (digits.length === 10) return `+7${digits}`
  return display
}

function switchMode(next: Mode) {
  mode.value = next
  error.value = ''
  info.value = ''
  password.value = ''
  passwordConfirm.value = ''
  resetCode.value = ''
}

async function afterAuthSuccess() {
  clearAnonymousOnboardingKeys()
  useFinancialProfile().loadProfile()
  useGoals().fetchGoals()
  await navigateTo(
    isOnboardingComplete(authStore.user?.phone, authStore.user?.id)
      ? '/dashboard'
      : '/onboarding'
  )
}

async function submitLogin() {
  error.value = ''
  if (!phoneComplete.value || !passwordOk.value) {
    error.value = 'Введите телефон и пароль (мин. 8 символов)'
    return
  }
  loading.value = true
  try {
    await login(normalizePhone(phone.value), password.value)
    await afterAuthSuccess()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Неверный телефон или пароль'
  } finally {
    loading.value = false
  }
}

async function submitRegister() {
  error.value = ''
  if (!phoneComplete.value || !passwordOk.value) {
    error.value = 'Телефон и пароль (мин. 8 символов) обязательны'
    return
  }
  if (!passwordsMatch.value) {
    error.value = 'Пароли не совпадают'
    return
  }
  loading.value = true
  try {
    await register(normalizePhone(phone.value), password.value)
    await login(normalizePhone(phone.value), password.value)
    await afterAuthSuccess()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Не удалось зарегистрироваться'
  } finally {
    loading.value = false
  }
}

async function submitForgot() {
  error.value = ''
  info.value = ''
  if (!phoneComplete.value) {
    error.value = 'Введите номер полностью'
    return
  }
  loading.value = true
  try {
    await requestPasswordReset(normalizePhone(phone.value))
    info.value = 'Если аккаунт существует, код отправлен. Проверьте SMS.'
    mode.value = 'reset'
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Не удалось отправить код'
  } finally {
    loading.value = false
  }
}

async function submitReset() {
  error.value = ''
  if (!phoneComplete.value || !resetCode.value.trim() || !passwordOk.value) {
    error.value = 'Заполните все поля; пароль — мин. 8 символов'
    return
  }
  if (!passwordsMatch.value) {
    error.value = 'Пароли не совпадают'
    return
  }
  loading.value = true
  try {
    await resetPassword(
      normalizePhone(phone.value),
      resetCode.value.trim(),
      password.value
    )
    await login(normalizePhone(phone.value), password.value)
    await afterAuthSuccess()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Неверный код или ошибка сброса'
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
      <CardTitle class="text-2xl">{{ titles[mode] }}</CardTitle>
      <CardDescription>{{ descriptions[mode] }}</CardDescription>
    </CardHeader>
    <CardContent>
      <form
        v-if="mode === 'login'"
        class="space-y-4"
        @submit.prevent="submitLogin"
      >
        <div class="space-y-2">
          <Label for="phone">Телефон</Label>
          <Input
            id="phone"
            v-model="phone"
            type="tel"
            inputmode="numeric"
            autocomplete="tel"
            maxlength="18"
            placeholder="+7 (999) 123-45-67"
            @keydown="allowDigitKeydown"
            @paste="onPhonePaste"
          />
        </div>
        <div class="space-y-2">
          <Label for="password">Пароль</Label>
          <Input
            id="password"
            v-model="password"
            type="password"
            autocomplete="current-password"
            minlength="8"
            placeholder="Минимум 8 символов"
          />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <Button type="submit" class="w-full" :disabled="loading">
          {{ loading ? 'Вход…' : 'Войти' }}
        </Button>
        <div class="flex flex-col gap-2 text-center text-sm">
          <button type="button" class="text-primary underline-offset-4 hover:underline" @click="switchMode('register')">
            Создать аккаунт
          </button>
          <button type="button" class="text-muted-foreground underline-offset-4 hover:underline" @click="switchMode('forgot')">
            Забыли пароль?
          </button>
        </div>
      </form>

      <form
        v-else-if="mode === 'register'"
        class="space-y-4"
        @submit.prevent="submitRegister"
      >
        <div class="space-y-2">
          <Label for="reg-phone">Телефон</Label>
          <Input
            id="reg-phone"
            v-model="phone"
            type="tel"
            inputmode="numeric"
            autocomplete="tel"
            maxlength="18"
            placeholder="+7 (999) 123-45-67"
            @keydown="allowDigitKeydown"
            @paste="onPhonePaste"
          />
        </div>
        <div class="space-y-2">
          <Label for="reg-password">Пароль</Label>
          <Input
            id="reg-password"
            v-model="password"
            type="password"
            autocomplete="new-password"
            minlength="8"
          />
        </div>
        <div class="space-y-2">
          <Label for="reg-password2">Повторите пароль</Label>
          <Input
            id="reg-password2"
            v-model="passwordConfirm"
            type="password"
            autocomplete="new-password"
            minlength="8"
          />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <Button type="submit" class="w-full" :disabled="loading">
          {{ loading ? 'Регистрация…' : 'Зарегистрироваться' }}
        </Button>
        <Button type="button" variant="ghost" class="w-full" @click="switchMode('login')">
          Уже есть аккаунт — войти
        </Button>
      </form>

      <form
        v-else-if="mode === 'forgot'"
        class="space-y-4"
        @submit.prevent="submitForgot"
      >
        <div class="space-y-2">
          <Label for="forgot-phone">Телефон</Label>
          <Input
            id="forgot-phone"
            v-model="phone"
            type="tel"
            inputmode="numeric"
            autocomplete="tel"
            maxlength="18"
            placeholder="+7 (999) 123-45-67"
            @keydown="allowDigitKeydown"
            @paste="onPhonePaste"
          />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <p v-if="info" class="text-sm text-muted-foreground">{{ info }}</p>
        <Button type="submit" class="w-full" :disabled="loading || !phoneComplete">
          {{ loading ? 'Отправка…' : 'Отправить код' }}
        </Button>
        <Button type="button" variant="ghost" class="w-full" @click="switchMode('login')">
          Назад ко входу
        </Button>
      </form>

      <form
        v-else
        class="space-y-4"
        @submit.prevent="submitReset"
      >
        <p class="text-sm text-muted-foreground">{{ phone }}</p>
        <div class="space-y-2">
          <Label for="reset-code">Код из SMS</Label>
          <Input
            id="reset-code"
            v-model="resetCode"
            type="text"
            inputmode="numeric"
            autocomplete="one-time-code"
            placeholder="Код подтверждения"
          />
        </div>
        <div class="space-y-2">
          <Label for="reset-password">Новый пароль</Label>
          <Input
            id="reset-password"
            v-model="password"
            type="password"
            autocomplete="new-password"
            minlength="8"
          />
        </div>
        <div class="space-y-2">
          <Label for="reset-password2">Повторите пароль</Label>
          <Input
            id="reset-password2"
            v-model="passwordConfirm"
            type="password"
            autocomplete="new-password"
            minlength="8"
          />
        </div>
        <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
        <Button type="submit" class="w-full" :disabled="loading">
          {{ loading ? 'Сохранение…' : 'Сохранить пароль' }}
        </Button>
        <Button type="button" variant="ghost" class="w-full" @click="switchMode('login')">
          Назад ко входу
        </Button>
      </form>
    </CardContent>
  </Card>
</template>
