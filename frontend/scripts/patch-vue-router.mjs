import { readFileSync, writeFileSync, existsSync } from 'node:fs'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..')
const targets = [
  'node_modules/vue-router/dist/vue-router.mjs',
  'node_modules/vue-router/dist/vue-router.esm-browser.js',
  'node_modules/vue-router/dist/vue-router.cjs'
]

const needle = 'instance.__vrv_devtools = info;'
const replacement = 'if (instance) instance.__vrv_devtools = info;'

let patched = 0

for (const rel of targets) {
  const file = resolve(root, rel)
  if (!existsSync(file)) continue

  const code = readFileSync(file, 'utf8')
  if (!code.includes(needle) || code.includes(replacement)) continue

  writeFileSync(file, code.replaceAll(needle, replacement), 'utf8')
  patched += 1
  console.log(`[patch-vue-router] patched ${rel}`)
}

if (patched === 0) {
  console.log('[patch-vue-router] already patched or vue-router not found')
}
