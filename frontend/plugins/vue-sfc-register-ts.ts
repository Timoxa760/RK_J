import { createRequire } from 'node:module'

const require = createRequire(import.meta.url)
let registered = false

function registerTsForVueSfc() {
  if (registered) return
  registered = true
  require('vue/compiler-sfc')
}

/** Ensures @vue/compiler-sfc can resolve defineProps types in Vite worker threads. */
export function vueSfcRegisterTs() {
  return {
    name: 'vue-sfc-register-ts',
    enforce: 'pre' as const,
    buildStart() {
      registerTsForVueSfc()
    },
    transform(_code: string, id: string) {
      if (!registered && id.includes('.vue')) {
        registerTsForVueSfc()
      }
      return null
    },
  }
}
