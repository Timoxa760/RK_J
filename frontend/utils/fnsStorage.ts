import type { FnsConnectionState } from '~/types/api'
import { currentUserStorageKey } from '~/utils/userStorage'

const FNS_PREFIX = 'potok:fns'

const defaultState = (): FnsConnectionState => ({
  connected: false,
  phone: '',
  auto_sync: true,
  imported_ids: []
})

export function readFnsConnection(): FnsConnectionState {
  if (!import.meta.client) return defaultState()
  try {
    const raw = localStorage.getItem(currentUserStorageKey(FNS_PREFIX))
    if (!raw) return defaultState()
    return { ...defaultState(), ...JSON.parse(raw) }
  } catch {
    return defaultState()
  }
}

export function writeFnsConnection(state: FnsConnectionState) {
  if (!import.meta.client) return
  localStorage.setItem(currentUserStorageKey(FNS_PREFIX), JSON.stringify(state))
}
