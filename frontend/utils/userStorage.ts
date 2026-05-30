import { useAuthStore } from '~/store/authStore'

/** Стабильный ключ пользователя для localStorage (телефон или id). */
export function normalizeUserKey(phone?: string | null, userId?: string | null): string {
  const digits = (phone ?? '').replace(/\D/g, '')
  if (digits.length >= 10) {
    if (digits.length === 11 && digits.startsWith('8')) return `7${digits.slice(1)}`
    if (digits.length === 11 && digits.startsWith('7')) return digits
    if (digits.length === 10) return `7${digits}`
    return digits
  }
  if (userId) return userId.replace(/[^a-zA-Z0-9._-]/g, '') || 'anonymous'
  return 'anonymous'
}

export function userStorageKey(prefix: string, phone?: string | null, userId?: string | null): string {
  return `${prefix}:${normalizeUserKey(phone, userId)}`
}

export function currentUserStorageKey(prefix: string): string {
  const authStore = useAuthStore()
  return userStorageKey(prefix, authStore.user?.phone, authStore.user?.id)
}
